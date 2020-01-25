package ghet

import (
	"context"
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

	CD            string   `long:"cd" description:"change to directory before downloading"`
	IgnoreMissing []string `long:"ignore-missing" description:"if an asset is not found, skip it rather than failing"`

	Args AssetArgs `positional-args:"true"`
}

type AssetArgs struct {
	Origin OriginOpt  `positional-arg-name:"OWNER/REPOSITORY[@TAG]" description:"release reference"`
	Assets []AssetOpt `positional-arg-name:"[LOCAL-PATH=]ASSET-NAME" description:"glob-friendly asset name to download"`
}

func (c *AssetCmd) applySettings() {
	if c.Args.Origin.Server == "" {
		c.Args.Origin.Server = c.Global.Server
	}

	if len(c.Args.Assets) == 0 {
		c.Args.Assets = []AssetOpt{
			{
				RemoteMatch: "*",
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

	for _, asset := range release.Assets {
		var matched bool
		var resolved AssetOpt

		for _, assetOpt := range c.Args.Assets {
			resolved, matched = assetOpt.Resolve(asset.GetName())
			if matched {
				break
			}
		}

		if !matched {
			continue
		}

		assetMap[resolved.LocalPath] = asset
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
