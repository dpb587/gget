package asset

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/dpb587/gget/pkg/service"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type Resource struct {
	client            *gitlab.Client
	releaseOwner      string
	releaseRepository string
	asset             *gitlab.ReleaseLink
}

var _ service.ResolvedResource = &Resource{}

func NewResource(client *gitlab.Client, releaseOwner, releaseRepository string, asset *gitlab.ReleaseLink) *Resource {
	return &Resource{
		client:            client,
		releaseOwner:      releaseOwner,
		releaseRepository: releaseRepository,
		asset:             asset,
	}
}

func (r *Resource) GetName() string {
	return path.Base(r.asset.URL)
}

func (r *Resource) GetSize() int64 {
	return 0
}

func (r *Resource) Open(ctx context.Context) (io.ReadCloser, error) {
	res, err := http.DefaultClient.Get(r.asset.URL)
	if err != nil {
		return nil, errors.Wrapf(err, "getting %s", r.asset.URL)
	}

	if res.StatusCode != 200 {
		return nil, errors.Wrapf(fmt.Errorf("expected status 200: got %d", res.StatusCode), "getting %s", r.asset.URL)
	}

	return res.Body, nil
}
