package github

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/dpb587/ghet/pkg/service"
	"github.com/dpb587/ghet/pkg/service/github/archive"
	"github.com/dpb587/ghet/pkg/service/github/asset"
	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type ReleaseRef struct {
	client  *github.Client
	repo    *github.Repository
	release *github.RepositoryRelease

	commitRef *CommitRef
}

func (r *ReleaseRef) ResolveResource(ctx context.Context, resourceType service.ResourceType, resource service.Resource) ([]service.ResolvedResource, error) {
	switch resourceType {
	case "archive":
		return r.resolveArchiveResource(ctx, resource)
	case "asset":
		return r.resolveAssetResource(ctx, resource)
	case "blob":
		return r.resolveBlobResource(ctx, resource)
	}

	return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
}

func (r *ReleaseRef) resolveArchiveResource(ctx context.Context, resource service.Resource) ([]service.ResolvedResource, error) {
	if resource == service.Resource("*") {
		// TODO glob-handling
		return nil, fmt.Errorf("TODO")
	}

	// TODO validation

	return []service.ResolvedResource{
		archive.NewResource(
			r.client,
			r.repo.GetOwner().GetLogin(),
			r.repo.GetName(),
			r.release.GetTargetCommitish(),
			string(resource),
		),
	}, nil
}

func (r *ReleaseRef) resolveAssetResource(ctx context.Context, resource service.Resource) ([]service.ResolvedResource, error) {
	var res []service.ResolvedResource

	for _, candidate := range r.release.Assets {
		if match, _ := filepath.Match(string(resource), candidate.GetName()); !match {
			continue
		}

		res = append(res, asset.NewResource(r.client, r.repo.GetOwner().GetLogin(), r.repo.GetName(), candidate))
	}

	return res, nil
}

func (r *ReleaseRef) resolveBlobResource(ctx context.Context, resource service.Resource) ([]service.ResolvedResource, error) {
	if r.commitRef == nil {
		// TODO GetTargetCommitish is incomplete
		commit, _, err := r.client.Repositories.GetCommit(ctx, r.repo.GetOwner().GetLogin(), r.repo.GetName(), r.release.GetTargetCommitish())
		if err != nil {
			return nil, errors.Wrap(err, "getting commit")
		}

		r.commitRef = &CommitRef{
			client: r.client,
			repo:   r.repo,
			commit: commit,
		}
	}

	return r.commitRef.resolveBlobResource(ctx, resource)
}
