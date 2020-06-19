package gitlab

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/gitlab/archive"
	"github.com/dpb587/gget/pkg/service/gitlab/blob"
	"github.com/dpb587/gget/pkg/service/gitlab/gitlabutil"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type CommitRef struct {
	service.RefMetadataService

	client *gitlab.Client
	ref    service.Ref
	commit string

	archiveFileBase string
}

var _ service.ResourceResolver = &CommitRef{}

func (r *CommitRef) CanonicalRef() service.Ref {
	return r.ref
}

func (r *CommitRef) ResolveResource(ctx context.Context, resourceType service.ResourceType, resource service.ResourceName) ([]service.ResolvedResource, error) {
	switch resourceType {
	case service.ArchiveResourceType:
		return r.resolveArchiveResource(ctx, resource)
	case service.BlobResourceType:
		return r.resolveBlobResource(ctx, resource)
	}

	return nil, fmt.Errorf("unsupported resource type for commit ref: %s", resourceType)
}

func (r *CommitRef) resolveArchiveResource(ctx context.Context, resource service.ResourceName) ([]service.ResolvedResource, error) {
	// https://docs.gitlab.com/ce/api/repositories.html#get-file-archive
	candidates := []string{
		fmt.Sprintf("%s.bz2", r.archiveFileBase),
		fmt.Sprintf("%s.tar", r.archiveFileBase),
		fmt.Sprintf("%s.tar.bz2", r.archiveFileBase),
		fmt.Sprintf("%s.tar.gz", r.archiveFileBase),
		fmt.Sprintf("%s.tb2", r.archiveFileBase),
		fmt.Sprintf("%s.tbz", r.archiveFileBase),
		fmt.Sprintf("%s.tbz2", r.archiveFileBase),
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
				r.ref,
				r.commit,
				candidate,
				strings.TrimPrefix(candidate, fmt.Sprintf("%s.", r.archiveFileBase)),
			),
		)
	}

	return res, nil
}

func (r *CommitRef) resolveBlobResource(ctx context.Context, resource service.ResourceName) ([]service.ResolvedResource, error) {
	var res []service.ResolvedResource

	// get the full tree
	pt := true
	opts := &gitlab.ListTreeOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
		Ref:       &r.ref.Ref,
		Recursive: &pt,
	}

	for {
		nodes, resp, err := r.client.Repositories.ListTree(gitlabutil.GetRepositoryID(r.ref), opts)
		if err != nil {
			return nil, errors.Wrap(err, "getting commit tree")
		}

		for _, candidate := range nodes {
			if candidate.Type != "blob" {
				continue
			} else if match, _ := filepath.Match(string(resource), candidate.Path); !match {
				continue
			}

			res = append(res, blob.NewResource(r.client, r.ref, r.commit, candidate))
		}

		if resp.NextPage == 0 {
			break
		}

		opts.ListOptions.Page = resp.NextPage
	}

	return res, nil
}
