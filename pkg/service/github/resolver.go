package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/dpb587/ghet/pkg/service"
	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type Service struct {
	client *github.Client
}

func NewService(client *github.Client) *Service {
	return &Service{
		client: client,
	}
}

var _ service.RefResolver = &Service{}

func (s Service) ResolveRef(ctx context.Context, ref service.Ref) (service.ResolvedRef, error) {
	repo, _, err := s.client.Repositories.Get(ctx, ref.Owner, ref.Repository)
	if err != nil {
		return nil, errors.Wrap(err, "getting repo")
	}

	var release *github.RepositoryRelease

	if ref.Ref != "" {
		if strings.HasPrefix(ref.Ref, "v") {
			// dynamic matching

			constraint, err := semver.NewConstraint(strings.TrimPrefix(ref.Ref, "v"))
			if err != nil {
				return nil, errors.Wrap(err, "parsing version constraint")
			}

			opts := &github.ListOptions{}

			for {
				releases, res, err := s.client.Repositories.ListReleases(ctx, repo.GetOwner().GetLogin(), repo.GetName(), opts)
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
				return nil, fmt.Errorf("expected to find release version matching %s: no release found", ref.Ref)
			}
		} else {
			release, _, err = s.client.Repositories.GetReleaseByTag(ctx, repo.GetOwner().GetLogin(), repo.GetName(), ref.Ref)
		}
	} else {
		release, _, err = s.client.Repositories.GetLatestRelease(ctx, repo.GetOwner().GetLogin(), repo.GetName())
	}
	if err != nil {
		return nil, errors.Wrap(err, "getting release")
	}

	return &Release{
		client:  s.client,
		repo:    repo,
		release: release,
	}, nil
}
