package archive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type Resource struct {
	client            *github.Client
	releaseOwner      string
	releaseRepository string
	ref               string
	filename          string
}

func NewResource(client *github.Client, releaseOwner, releaseRepository, ref, filename string) *Resource {
	return &Resource{
		client:            client,
		releaseOwner:      releaseOwner,
		releaseRepository: releaseRepository,
		filename:          filename,
	}
}

func (r *Resource) GetName() string {
	return r.filename
}

func (r *Resource) GetSize() int64 {
	return 0
}

func (r *Resource) Open(ctx context.Context) (io.ReadCloser, error) {
	var archiveLink *url.URL
	var err error

	ext := filepath.Ext(r.filename)
	if ext == "gz" {
		ext = fmt.Sprintf("%s%s", filepath.Ext(strings.TrimSuffix(r.filename, ext)), ext)
	}

	switch ext {
	case ".tar.gz", ".tgz":
		archiveLink, _, err = r.client.Repositories.GetArchiveLink(ctx, r.releaseOwner, r.releaseRepository, github.Tarball, &github.RepositoryContentGetOptions{
			Ref: r.ref,
		}, false)
	case ".zip":
		archiveLink, _, err = r.client.Repositories.GetArchiveLink(ctx, r.releaseOwner, r.releaseRepository, github.Zipball, &github.RepositoryContentGetOptions{
			Ref: r.ref,
		}, false)
	default:
		return nil, fmt.Errorf("unrecognized extension: %s", ext)
	}

	if err != nil {
		return nil, errors.Wrap(err, "getting archive url")
	}

	res, err := http.DefaultClient.Get(archiveLink.String())
	if err != nil {
		return nil, errors.Wrap(err, "getting download url")
	}

	if res.StatusCode != 200 {
		return nil, errors.Wrapf(fmt.Errorf("expected status 200: got %d", res.StatusCode), "getting download url %s", archiveLink)
	}

	return res.Body, nil
}
