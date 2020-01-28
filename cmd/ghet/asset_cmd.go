package ghet

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/dpb587/ghet/pkg/downloader"
	"github.com/dpb587/ghet/pkg/github"
	githubasset "github.com/dpb587/ghet/pkg/github/asset"
	"github.com/dpb587/ghet/pkg/model"
	gogithub "github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v4"
)

type AssetCmd struct {
	*Global `no-flag:"true"`

	CD string `long:"cd" description:"change to directory before downloading"`

	IgnoreMissing []AssetNameOpt `long:"ignore-missing" description:"if an asset is not found, skip it rather than failing"`
	Exclude       []AssetNameOpt `long:"exclude" description:"asset name to exclude from downloads (glob-friendly)"`

	List bool `long:"list" description:"list matched assets instead of downloading"`

	Args AssetArgs `positional-args:"true"`
}

type AssetArgs struct {
	Origin OriginOpt      `positional-arg-name:"OWNER/REPOSITORY[@TAG]" description:"release reference"`
	Assets []AssetPathOpt `positional-arg-name:"[LOCAL-PATH=]ASSET-NAME" description:"asset name to download (glob-friendly)"`
}

func (c *AssetCmd) applySettings() {
	if c.Args.Origin.Server == "" {
		c.Args.Origin.Server = c.Global.Server
	}

	if len(c.Args.Assets) == 0 {
		c.Args.Assets = []AssetPathOpt{
			{
				RemoteMatch: AssetNameOpt("*"),
			},
		}
	}

	if c.CD != "" {
		for assetIdx, asset := range c.Args.Assets {
			c.Args.Assets[assetIdx].LocalPath = filepath.Join(c.CD, asset.LocalPath)
		}
	}
}

func (c *AssetCmd) Execute(_ []string) error {
	c.applySettings()

	ctx := context.Background()
	client := c.Global.GitHubClient(c.Args.Origin.Server)

	release, err := github.ResolveRelease(ctx, client, model.Origin(c.Args.Origin))
	if err != nil {
		return errors.Wrap(err, "resolving release reference")
	}

	checksums := githubasset.NewChecksumManager(release)

	assetMap := map[string]gogithub.ReleaseAsset{}
	assetMatches := make([]bool, len(c.Args.Assets))

	for _, asset := range release.Assets {
		{ // first check if its excluded
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

		assetMap[resolved.LocalPath] = asset
	}

	{ // did we find everything the user asked for?
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

	pb := mpb.New(mpb.WithWidth(1))

	var downloads []*downloader.Download

	for localPath, asset := range assetMap {
		d := downloader.NewDownload(githubasset.NewAsset(client, c.Args.Origin.Owner, c.Args.Origin.Repository, asset))

		cs, err := checksums.GetAssetChecksum(asset.GetName())
		if err != nil {
			return errors.Wrap(err, "getting asset checksum")
		}

		d.AddVerifier(&downloader.DownloadHashVerifier{
			Algo:     cs.Type,
			Expected: cs.Bytes,
			Actual:   cs.Hasher(),
		})

		d.AddInstaller(&downloader.DownloadPathInstaller{
			Target: localPath,
		})

		downloads = append(downloads, d)
	}

	sort.Slice(downloads, func(i, j int) bool {
		return downloads[i].GetName() < downloads[j].GetName()
	})

	for _, d := range downloads {
		d.SetProgressBar(pb)
	}

	for _, d := range downloads {
		err := d.Download(ctx)
		if err != nil {
			return errors.Wrapf(err, "downloading %s", d.GetName())
		}
	}

	pb.Wait()

	return nil
}
