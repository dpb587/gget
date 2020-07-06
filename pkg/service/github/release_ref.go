package github

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/github/asset"
	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type ReleaseRef struct {
	refResolver *refResolver
	release     *github.RepositoryRelease

	targetRef       service.ResolvedRef
	checksumManager checksum.Manager
}

var _ service.ResolvedRef = &ReleaseRef{}
var _ service.ResourceResolver = &ReleaseRef{}

func (r *ReleaseRef) CanonicalRef() service.Ref {
	return r.refResolver.canonicalRef
}

func (r *ReleaseRef) GetMetadata(ctx context.Context) (service.RefMetadata, error) {
	targetRef, err := r.requireTargetRef(ctx)
	if err != nil {
		return nil, err
	}

	tagMetadata, err := targetRef.GetMetadata(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "getting commit metadata")
	}

	res := append(
		service.RefMetadata{
			{
				Name:  "github-release-id",
				Value: fmt.Sprintf("%d", r.release.GetID()),
			},
			{
				Name:  "github-release-published-at",
				Value: r.release.GetPublishedAt().Format(time.RFC3339),
			},
			{
				Name:  "github-release-body",
				Value: r.release.GetBody(),
			},
		},
		tagMetadata...,
	)

	return res, nil
}

func (r *ReleaseRef) ResolveResource(ctx context.Context, resourceType service.ResourceType, resource service.ResourceName) ([]service.ResolvedResource, error) {
	if resourceType == service.AssetResourceType {
		return r.resolveAssetResource(ctx, resource)
	}

	targetRef, err := r.requireTargetRef(ctx)
	if err != nil {
		return nil, err
	}

	return targetRef.ResolveResource(ctx, resourceType, resource)
}

func (r *ReleaseRef) requireTargetRef(ctx context.Context) (service.ResolvedRef, error) {
	if r.targetRef == nil {
		ref, err := r.refResolver.resolveTagWithRelease(ctx, r.release)
		if err != nil {
			return nil, errors.Wrap(err, "resolving commit")
		}

		r.targetRef = ref
	}

	return r.targetRef, nil
}

func (r *ReleaseRef) resolveAssetResource(ctx context.Context, resource service.ResourceName) ([]service.ResolvedResource, error) {
	var res []service.ResolvedResource

	for _, candidate := range r.release.Assets {
		if match, _ := filepath.Match(string(resource), candidate.GetName()); !match {
			continue
		}

		res = append(
			res,
			asset.NewResource(r.refResolver.client, r.refResolver.canonicalRef.Owner, r.refResolver.canonicalRef.Repository, candidate, r.requireChecksumManager()),
		)
	}

	return res, nil
}

func (r *ReleaseRef) requireChecksumManager() checksum.Manager {
	if r.checksumManager == nil {
		r.checksumManager = NewReleaseChecksumManager(r.refResolver.client, r.refResolver.canonicalRef.Owner, r.refResolver.canonicalRef.Repository, r.release)
	}

	return r.checksumManager
}
