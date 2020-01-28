package asset

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type Asset struct {
	client            *github.Client
	releaseOwner      string
	releaseRepository string
	asset             github.ReleaseAsset
}

func NewAsset(client *github.Client, releaseOwner, releaseRepository string, asset github.ReleaseAsset) *Asset {
	return &Asset{
		client:            client,
		releaseOwner:      releaseOwner,
		releaseRepository: releaseRepository,
		asset:             asset,
	}
}

func (a *Asset) GetName() string {
	return a.asset.GetName()
}

func (a *Asset) GetSize() int {
	return a.asset.GetSize()
}

func (a *Asset) GetLocation(ctx context.Context) (string, error) {
	return a.asset.GetBrowserDownloadURL(), nil
}

func (a *Asset) Open(ctx context.Context) (io.ReadCloser, error) {
	remoteHandle, redirectURL, err := a.client.Repositories.DownloadReleaseAsset(ctx, a.releaseOwner, a.releaseRepository, a.asset.GetID())
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
