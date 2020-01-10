package main

// go run . {repository} {version}
// find the repo with the GitHub API
// find the release based on version with the GitHub API
// for all release assets, download them

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

func main() {
	ownerRepo := strings.Split(os.Args[1], "/")
	version := "latest"
	if len(os.Args) > 2 {
		version = os.Args[2]
	}

	ctx := context.Background()
	client := github.NewClient(nil)

	repo, _, err := client.Repositories.Get(ctx, ownerRepo[0], ownerRepo[1])
	if err != nil {
		panic(errors.Wrap(err, "getting repo"))
	}

	var release *github.RepositoryRelease

	if version != "latest" {
		release, _, err = client.Repositories.GetReleaseByTag(ctx, repo.GetOwner().GetLogin(), repo.GetName(), version)
	} else {
		release, _, err = client.Repositories.GetLatestRelease(ctx, repo.GetOwner().GetLogin(), repo.GetName())
	}
	if err != nil {
		panic(errors.Wrap(err, "getting release"))
	}

	for _, asset := range release.Assets {
		fmt.Printf("%+v\n", asset.GetBrowserDownloadURL())
		remoteHandle, redirectURL, err := client.Repositories.DownloadReleaseAsset(ctx, repo.GetOwner().GetLogin(), repo.GetName(), asset.GetID())
		if err != nil {
			panic(errors.Wrapf(err, "getting asset %f", asset.GetName()))
		}

		if redirectURL != "" {
			res, err := http.DefaultClient.Get(redirectURL)
			if err != nil {
				panic(errors.Wrapf(err, "getting download url %s", redirectURL))
			}

			if res.StatusCode != 200 {
				panic(errors.Wrapf(fmt.Errorf("expected status 200: got %d", res.StatusCode), "getting download url %s", redirectURL))
			}

			remoteHandle = res.Body
		}

		localHandle, err := os.OpenFile(asset.GetName(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
		if err != nil {
			panic(errors.Wrapf(err, "creating %s", asset.GetName()))
		}

		_, err = io.Copy(localHandle, remoteHandle)
		if err != nil {
			panic(errors.Wrapf(err, "downloading %s", asset.GetName()))
		}
	}
}
