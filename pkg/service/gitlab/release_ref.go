package gitlab

import (
	"context"
	"path"
	"path/filepath"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/gitlab/asset"
	"github.com/xanzy/go-gitlab"
)

type ReleaseRef struct {
	client    *gitlab.Client
	ref       service.Ref
	release   *gitlab.Release
	targetRef service.ResolvedRef

	checksumManager checksum.Manager
}

var _ service.ResourceResolver = &ReleaseRef{}

func (r *ReleaseRef) CanonicalRef() service.Ref {
	return r.ref
}

func (r *ReleaseRef) GetMetadata() []service.RefMetadata {
	return r.targetRef.GetMetadata()
}

func (r *ReleaseRef) ResolveResource(ctx context.Context, resourceType service.ResourceType, resource service.Resource) ([]service.ResolvedResource, error) {
	if resourceType == service.AssetResourceType {
		return r.resolveAssetResource(ctx, resource)
	}

	return r.targetRef.ResolveResource(ctx, resourceType, resource)
}

func (r *ReleaseRef) resolveAssetResource(ctx context.Context, resource service.Resource) ([]service.ResolvedResource, error) {
	var res []service.ResolvedResource

	for _, candidate := range r.release.Assets.Links {
		// TODO kind of weird to extract from remote url "file" name, but
		//   Name field is more traditionally a label. So currently using
		//   the file name that a browser would typically produce. Doesn't
		//   cover more complex URLs though.
		if match, _ := filepath.Match(string(resource), path.Base(candidate.URL)); !match {
			continue
		}

		res = append(
			res,
			asset.NewResource(r.client, r.ref.Owner, r.ref.Repository, candidate),
		)
	}

	return res, nil
}
