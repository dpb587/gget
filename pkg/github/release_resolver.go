package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/dpb587/ghet/pkg/model"
	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

func ResolveRelease(ctx context.Context, client *github.Client, origin model.Origin) (*github.RepositoryRelease, error) {
	repo, _, err := client.Repositories.Get(ctx, origin.Owner, origin.Repository)
	if err != nil {
		return nil, errors.Wrap(err, "getting repo")
	}

	var release *github.RepositoryRelease

	if origin.Ref != "" {
		if strings.HasPrefix(origin.Ref, "v") {
			// dynamic matching

			constraint, err := semver.NewConstraint(strings.TrimPrefix(origin.Ref, "v"))
			if err != nil {
				return nil, errors.Wrap(err, "parsing version constraint")
			}

			opts := &github.ListOptions{}

			for {
				releases, res, err := client.Repositories.ListReleases(ctx, repo.GetOwner().GetLogin(), repo.GetName(), opts)
				if err != nil {
					return nil, errors.Wrap(err, "listing releases")
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
				return nil, fmt.Errorf("expected to find release version matching %s: no release found", origin.Ref)
			}
		} else {
			release, _, err = client.Repositories.GetReleaseByTag(ctx, repo.GetOwner().GetLogin(), repo.GetName(), origin.Ref)
		}
	} else {
		release, _, err = client.Repositories.GetLatestRelease(ctx, repo.GetOwner().GetLogin(), repo.GetName())
	}
	if err != nil {
		return nil, errors.Wrap(err, "getting release")
	}

	return release, nil
}
