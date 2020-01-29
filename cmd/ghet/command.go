package ghet

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/dpb587/ghet/pkg/downloader"
	"github.com/dpb587/ghet/pkg/github"
	githubasset "github.com/dpb587/ghet/pkg/github/asset"
	"github.com/dpb587/ghet/pkg/model"
	gogithub "github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
	"github.com/tidwall/limiter"
	"github.com/vbauerster/mpb/v4"
)

type Command struct {
	*Runtime `group:"Runtime Options"`

	IgnoreMissing []AssetNameOpt `long:"ignore-missing" description:"if an asset is not found, skip it rather than failing" value-name:"[ASSET]" optional:"true" optional-value:"*"`
	Exclude       []AssetNameOpt `long:"exclude" description:"asset name to exclude from downloads (glob-friendly)" value-name:"ASSET"`

	List bool `long:"list" description:"list matched assets instead of downloading"`

	CD         string         `long:"cd" description:"change to directory before downloading"`
	Executable []AssetNameOpt `long:"exec" description:"apply executable permissions to downloads" value-name:"[ASSET]" optional:"true" optional-value:"*"`
	Stdout     bool           `long:"stdout" description:"write file contents to stdout rather than disk"`

	Args CommandArgs `positional-args:"true" required:"true"`
}

type CommandArgs struct {
	Origin OriginOpt      `positional-arg-name:"OWNER/REPOSITORY[@REF]" description:"release reference"`
	Assets []AssetPathOpt `positional-arg-name:"[LOCAL-PATH=]ASSET" description:"asset name(s) to download (glob-friendly)" optional:"true"`
}

func (c *Command) applySettings() {
	if c.Args.Origin.Server == "" {
		c.Args.Origin.Server = c.Runtime.Server
	}

	if len(c.Args.Assets) == 0 {
		c.Args.Assets = []AssetPathOpt{
			{
				RemoteMatch: AssetNameOpt("*"),
			},
		}
	}

	if c.Stdout {
		for assetIdx, asset := range c.Args.Assets {
			if asset.LocalPath != "" {
				continue
			}

			c.Args.Assets[assetIdx].LocalPath = "-"
		}
	}

	if c.CD != "" {
		for assetIdx, asset := range c.Args.Assets {
			if asset.LocalPath == "-" {
				continue
			}

			c.Args.Assets[assetIdx].LocalPath = filepath.Join(c.CD, asset.LocalPath)
		}
	}
}

func (c *Command) Execute(_ []string) error {
	c.applySettings()

	ctx := context.Background()
	client := c.Runtime.GitHubClient(c.Args.Origin.Server)

	release, err := github.ResolveRelease(ctx, client, model.Origin(c.Args.Origin))
	if err != nil {
		return errors.Wrap(err, "resolving release reference")
	}

	checksums := githubasset.NewChecksumManager(release)

	assetMap := map[string]gogithub.ReleaseAsset{}
	assetMatches := make([]bool, len(c.Args.Assets))

	for _, asset := range release.Assets {
		{ // first check if it is excluded
			var excluded bool

			for _, assetNameOpt := range c.Exclude {
				if assetNameOpt.Match(asset.GetName()) {
					excluded = true

					break
				}
			}

			if excluded {
				continue
			}
		}

		var resolved AssetPathOpt

		{ // now check if its a match
			var matched bool
			for assetPathOptIdx, assetPathOpt := range c.Args.Assets {
				resolved, matched = assetPathOpt.Resolve(asset.GetName())
				if !matched {
					continue
				}

				assetMatches[assetPathOptIdx] = true

				break
			}

			if !matched {
				continue
			}
		}

		if _, found := assetMap[resolved.LocalPath]; found {
			return fmt.Errorf("target file already specified: %s", resolved.LocalPath)
		}

		assetMap[resolved.LocalPath] = asset
	}

	{ // finally, did we find everything the user asked for?
		for assetIdx, assetMatched := range assetMatches {
			if assetMatched {
				continue
			}

			return errors.Wrap(fmt.Errorf("no asset matched: %s", c.Args.Assets[assetIdx].RemoteMatch), "expected matching assets")
		}
	}

	if c.List {
		for _, asset := range assetMap {
			fmt.Println(asset.GetName())
		}

		return nil
	}

	var downloads []*downloader.Workflow

	for localPath, asset := range assetMap {
		var steps []downloader.Step

		if localPath == "-" {
			steps = append(
				steps,
				&downloader.DownloadWriterInstaller{
					Writer: os.Stdout,
				},
			)
		} else {
			steps = append(
				steps,
				&downloader.DownloadTmpfileInstaller{},
			)
		}

		cs, csFound, err := checksums.GetAssetChecksum(asset.GetName())
		if err != nil {
			return errors.Wrap(err, "getting asset checksum")
		}

		if csFound {
			steps = append(
				steps,
				&downloader.DownloadHashVerifier{
					Algo:     cs.Type,
					Expected: cs.Bytes,
					Actual:   cs.Hasher(),
				},
			)
		}

		if localPath != "-" {
			for _, assetNameOpt := range c.Executable {
				if !assetNameOpt.Match(asset.GetName()) {
					continue
				}

				steps = append(
					steps,
					&downloader.DownloadExecutableInstaller{},
				)
			}

			steps = append(
				steps,
				&downloader.DownloadRenameInstaller{
					Target: localPath,
				},
			)
		}

		downloads = append(
			downloads,
			downloader.NewWorkflow(githubasset.NewAsset(client, c.Args.Origin.Owner, c.Args.Origin.Repository, asset), steps...),
		)
	}

	sort.Slice(downloads, func(i, j int) bool {
		return downloads[i].GetSubject() < downloads[j].GetSubject()
	})

	var pbO io.Writer = os.Stderr

	if c.Runtime.Quiet {
		pbO = ioutil.Discard
	}

	pb := mpb.New(mpb.WithWidth(1), mpb.WithOutput(pbO))

	for _, d := range downloads {
		d.Prepare(pb)
	}

	l := limiter.New(c.Runtime.Parallel)

	for _, d := range downloads {
		d := d
		go func() {
			l.Begin()
			defer l.End()

			err := d.Execute(ctx)
			if err != nil {
				// TODO concurrency
				panic(errors.Wrapf(err, "downloading %s", d.GetSubject()))
			}
		}()
	}

	pb.Wait()

	return nil
}
