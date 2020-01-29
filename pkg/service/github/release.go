package github

import (
	"context"
	"path/filepath"

	"github.com/dpb587/ghet/pkg/service"
	"github.com/dpb587/ghet/pkg/service/github/asset"
	"github.com/google/go-github/v29/github"
)

type Release struct {
	client  *github.Client
	repo    *github.Repository
	release *github.RepositoryRelease
}

func (r *Release) ResolveResource(ctx context.Context, resource service.Resource) ([]service.ResolvedResource, error) {
	var res []service.ResolvedResource

	for _, candidate := range r.release.Assets {
		if match, _ := filepath.Match(string(resource), candidate.GetName()); !match {
			continue
		}

		res = append(res, asset.NewResource(r.client, r.repo.GetOwner().GetLogin(), r.repo.GetName(), candidate))
	}

	return res, nil
}
