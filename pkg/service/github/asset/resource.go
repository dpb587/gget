package asset

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type Resource struct {
	client            *github.Client
	releaseOwner      string
	releaseRepository string
	asset             github.ReleaseAsset
}

func NewResource(client *github.Client, releaseOwner, releaseRepository string, asset github.ReleaseAsset) *Resource {
	return &Resource{
		client:            client,
		releaseOwner:      releaseOwner,
		releaseRepository: releaseRepository,
		asset:             asset,
	}
}

func (r *Resource) GetName() string {
	return r.asset.GetName()
}

func (r *Resource) GetSize() int64 {
	return int64(r.asset.GetSize())
}

func (r *Resource) Open(ctx context.Context) (io.ReadCloser, error) {
	remoteHandle, redirectURL, err := r.client.Repositories.DownloadReleaseAsset(ctx, r.releaseOwner, r.releaseRepository, r.asset.GetID())
	if err != nil {
		return nil, errors.Wrap(err, "requesting asset")
	}

	if remoteHandle != nil {
		defer remoteHandle.Close()
	}

	if redirectURL != "" {
		res, err := http.DefaultClient.Get(redirectURL)
		if err != nil {
			return nil, errors.Wrapf(err, "getting download url %s", redirectURL)
		}

		if res.StatusCode != 200 {
			return nil, errors.Wrapf(fmt.Errorf("expected status 200: got %d", res.StatusCode), "getting download url %s", redirectURL)
		}

		remoteHandle = res.Body
	}

	return remoteHandle, nil
}
