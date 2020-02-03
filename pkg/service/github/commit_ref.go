package github

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/github/archive"
	"github.com/dpb587/gget/pkg/service/github/blob"
	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type CommitRef struct {
	service.RefMetadataService

	client *github.Client
	repo   *github.Repository
	commit string

	archiveFileBase string
}

var _ service.ResourceResolver = &CommitRef{}

func (r *CommitRef) ResolveResource(ctx context.Context, resourceType service.ResourceType, resource service.Resource) ([]service.ResolvedResource, error) {
	switch resourceType {
	case service.ArchiveResourceType:
		return r.resolveArchiveResource(ctx, resource)
	case service.BlobResourceType:
		return r.resolveBlobResource(ctx, resource)
	}

	return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
}

func (r *CommitRef) resolveArchiveResource(ctx context.Context, resource service.Resource) ([]service.ResolvedResource, error) {
	candidates := []string{
		fmt.Sprintf("%s.tar.gz", r.archiveFileBase),
		fmt.Sprintf("%s.zip", r.archiveFileBase),
	}

	var res []service.ResolvedResource

	for _, candidate := range candidates {
		if match, _ := filepath.Match(string(resource), candidate); !match {
			continue
		}

		res = append(
			res,
			archive.NewResource(
				r.client,
				r.repo.GetOwner().GetLogin(),
				r.repo.GetName(),
				r.commit,
				candidate,
			),
		)
	}

	return res, nil
}

func (r *CommitRef) resolveBlobResource(ctx context.Context, resource service.Resource) ([]service.ResolvedResource, error) {
	var res []service.ResolvedResource

	// get the full tree
	tree, _, err := r.client.Git.GetTree(ctx, r.repo.GetOwner().GetLogin(), r.repo.GetName(), r.commit, true)
	if err != nil {
		return nil, errors.Wrap(err, "getting commit tree")
	}

	for _, candidate := range tree.Entries {
		if match, _ := filepath.Match(string(resource), candidate.GetPath()); !match {
			continue
		}

		res = append(res, blob.NewResource(r.client, r.repo.GetOwner().GetLogin(), r.repo.GetName(), candidate))
	}

	return res, nil
}
