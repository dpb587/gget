package archive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/dpb587/gget/pkg/service"
	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type Resource struct {
	client   *github.Client
	ref      service.Ref
	target   string
	filename string
}

var _ service.ResolvedResource = &Resource{}

func NewResource(client *github.Client, ref service.Ref, target, filename string) *Resource {
	return &Resource{
		client:   client,
		ref:      ref,
		target:   target,
		filename: filename,
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
	if ext == ".gz" {
		ext = fmt.Sprintf("%s%s", filepath.Ext(strings.TrimSuffix(r.filename, ext)), ext)
	}

	switch ext {
	case ".tar.gz", ".tgz":
		archiveLink, _, err = r.client.Repositories.GetArchiveLink(ctx, r.ref.Owner, r.ref.Repository, github.Tarball, &github.RepositoryContentGetOptions{
			Ref: r.target,
		}, false)
	case ".zip":
		archiveLink, _, err = r.client.Repositories.GetArchiveLink(ctx, r.ref.Owner, r.ref.Repository, github.Zipball, &github.RepositoryContentGetOptions{
			Ref: r.target,
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
