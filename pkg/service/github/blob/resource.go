package blob

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type Resource struct {
	client            *github.Client
	releaseOwner      string
	releaseRepository string
	asset             github.TreeEntry
}

func NewResource(client *github.Client, releaseOwner, releaseRepository string, asset github.TreeEntry) *Resource {
	return &Resource{
		client:            client,
		releaseOwner:      releaseOwner,
		releaseRepository: releaseRepository,
		asset:             asset,
	}
}

func (r *Resource) GetName() string {
	return r.asset.GetPath()
}

func (r *Resource) GetSize() int64 {
	return int64(r.asset.GetSize())
}

func (r *Resource) Open(ctx context.Context) (io.ReadCloser, error) {
	// TODO switch to stream?
	buf, _, err := r.client.Git.GetBlobRaw(ctx, r.releaseOwner, r.releaseRepository, r.asset.GetSHA())
	if err != nil {
		return nil, errors.Wrap(err, "getting blob")
	}

	return ioutil.NopCloser(bytes.NewReader(buf)), nil
}
