package main

// version-matching
// v# -> v#.latest
// v#.#.# -> assume exact match
// "~2.1, !~2.3.1"
//

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
	"path/filepath"
	"strings"

	"github.com/dpb587/ghet/pkg/checksum"
	"github.com/google/go-github/v29/github"
	"github.com/masterminds/semver"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func main() {
	slugVersion := strings.Split(os.Args[1], "@")
	ownerRepo := strings.Split(slugVersion[0], "/")

	version := "latest"
	if len(slugVersion) > 1 {
		version = slugVersion[1]
	}

	includeGlobs := []string{"*"}
	if len(os.Args) > 2 {
		includeGlobs = os.Args[2:]
	}

	var tc *http.Client
	ctx := context.Background()

	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		tc = oauth2.NewClient(ctx,
			oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: v},
			),
		)
	}

	client := github.NewClient(tc)

	repo, _, err := client.Repositories.Get(ctx, ownerRepo[0], ownerRepo[1])
	if err != nil {
		panic(errors.Wrap(err, "getting repo"))
	}

	var release *github.RepositoryRelease

	if version != "latest" {
		if strings.HasPrefix(version, "v") {
			// dynamic matching

			constraint, err := semver.NewConstraint(strings.TrimPrefix(version, "v"))
			if err != nil {
				panic(errors.Wrap(err, "parsing version constraint"))
			}

			opts := &github.ListOptions{}

			for {
				releases, res, err := client.Repositories.ListReleases(ctx, repo.GetOwner().GetLogin(), repo.GetName(), opts)
				if err != nil {
					panic(errors.Wrap(err, "listing releases"))
				}

				for _, candidateRelease := range releases {
					releaseVersion, err := semver.NewVersion(strings.TrimPrefix(candidateRelease.GetTagName(), "v"))
					if err != nil {
						// TODO log debug?
						continue
					}

					if constraint.Check(releaseVersion) {
						release = candidateRelease

						break
					}
				}

				if release != nil {
					break
				}

				opts.Page = res.NextPage

				if opts.Page == 0 {
					break
				} else if opts.Page > 5 {
					// just to be safe
					break // TODO customizable limit?
				}
			}

			if release == nil {
				panic(fmt.Errorf("expected to find release version matching %s: no release found", version))
			}
		} else {
			release, _, err = client.Repositories.GetReleaseByTag(ctx, repo.GetOwner().GetLogin(), repo.GetName(), version)
		}
	} else {
		release, _, err = client.Repositories.GetLatestRelease(ctx, repo.GetOwner().GetLogin(), repo.GetName())
	}
	if err != nil {
		panic(errors.Wrap(err, "getting release"))
	}

	parsedReleases := checksum.ParseReleaseNotes(release.GetBody())

	for _, asset := range release.Assets {
		var matched bool

		for _, includeGlob := range includeGlobs {
			if v, _ := filepath.Match(includeGlob, asset.GetName()); v {
				matched = true

				break
			}
		}

		if !matched {
			continue
		}

		checksum, found := parsedReleases.GetByName(asset.GetName())
		if !found {
			panic(errors.Wrapf(fmt.Errorf("no checksum found"), "downloading %s", asset.GetName()))
		}

		fmt.Printf("%s  %s\n", checksum.SHA, asset.GetName())

		verifierHash := checksum.NewHash()

		remoteHandle, redirectURL, err := client.Repositories.DownloadReleaseAsset(ctx, repo.GetOwner().GetLogin(), repo.GetName(), asset.GetID())
		if err != nil {
			panic(errors.Wrapf(err, "getting asset %f", asset.GetName()))
		}

		if remoteHandle != nil {
			defer remoteHandle.Close()
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

		defer localHandle.Close()

		tee := io.MultiWriter(localHandle, verifierHash)

		_, err = io.Copy(tee, remoteHandle)
		if err != nil {
			panic(errors.Wrapf(err, "downloading %s", asset.GetName()))
		}

		actualSHA := fmt.Sprintf("%x", verifierHash.Sum(nil))

		if actualSHA != checksum.SHA {
			panic(fmt.Errorf("expected checksum %s: got %s", checksum.SHA, actualSHA))
		}
	}
}
