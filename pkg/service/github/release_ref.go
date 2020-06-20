package github

import (
	"context"
	"path/filepath"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/github/asset"
	"github.com/google/go-github/v29/github"
)

type ReleaseRef struct {
	client    *github.Client
	ref       service.Ref
	release   *github.RepositoryRelease
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

func (r *ReleaseRef) ResolveResource(ctx context.Context, resourceType service.ResourceType, resource service.ResourceName) ([]service.ResolvedResource, error) {
	if resourceType == service.AssetResourceType {
		return r.resolveAssetResource(ctx, resource)
	}

	return r.targetRef.ResolveResource(ctx, resourceType, resource)
}

func (r *ReleaseRef) resolveAssetResource(ctx context.Context, resource service.ResourceName) ([]service.ResolvedResource, error) {
	var res []service.ResolvedResource

	for _, candidate := range r.release.Assets {
		if match, _ := filepath.Match(string(resource), candidate.GetName()); !match {
			continue
		}

		res = append(
			res,
			asset.NewResource(r.client, r.ref.Owner, r.ref.Repository, candidate, r.requireChecksumManager()),
		)
	}

	return res, nil
}

func (r *ReleaseRef) requireChecksumManager() checksum.Manager {
	if r.checksumManager == nil {
		r.checksumManager = NewReleaseChecksumManager(r.client, r.ref.Owner, r.ref.Repository, r.release)
	}

	return r.checksumManager
}
