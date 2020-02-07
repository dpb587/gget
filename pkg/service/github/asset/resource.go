package asset

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/dpb587/gget/pkg/downloader"
	"github.com/dpb587/gget/pkg/service"
	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type Resource struct {
	client            *github.Client
	releaseOwner      string
	releaseRepository string
	checksumManager   checksum.Manager
	asset             github.ReleaseAsset
}

var _ service.ResolvedResource = &Resource{}
var _ downloader.StepProvider = &Resource{}

func NewResource(client *github.Client, releaseOwner, releaseRepository string, asset github.ReleaseAsset, checksumManager checksum.Manager) *Resource {
	return &Resource{
		client:            client,
		releaseOwner:      releaseOwner,
		releaseRepository: releaseRepository,
		asset:             asset,
		checksumManager:   checksumManager,
	}
}

func (r *Resource) GetName() string {
	return r.asset.GetName()
}

func (r *Resource) GetSize() int64 {
	return int64(r.asset.GetSize())
}

func (r *Resource) GetDownloaderSteps(ctx context.Context) ([]downloader.Step, error) {
	cs, found, err := r.checksumManager.GetChecksum(ctx, r.asset.GetName())
	if err != nil {
		return nil, errors.Wrapf(err, "getting checksum of %s", r.asset.GetName())
	} else if !found {
		return nil, nil
	}

	res := []downloader.Step{
		&downloader.DownloadHashVerifier{
			Algo:     cs.Type,
			Expected: cs.Bytes,
			Actual:   cs.Hasher(),
		},
	}

	return res, nil
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
