package archive

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"

	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/gitlab/gitlabutil"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type Resource struct {
	client   *gitlab.Client
	ref      service.Ref
	target   string
	filename string
	format   string
}

var _ service.ResolvedResource = &Resource{}

func NewResource(client *gitlab.Client, ref service.Ref, target, filename, format string) *Resource {
	return &Resource{
		client:   client,
		ref:      ref,
		target:   target,
		filename: filename,
		format:   format,
	}
}

func (r *Resource) GetName() string {
	return r.filename
}

func (r *Resource) GetSize() int64 {
	return 0
}

func (r *Resource) Open(ctx context.Context) (io.ReadCloser, error) {
	buf, _, err := r.client.Repositories.Archive(gitlabutil.GetRepositoryID(r.ref), &gitlab.ArchiveOptions{
		Format: &r.format,
		SHA:    &r.target,
	})

	if err != nil {
		return nil, errors.Wrap(err, "getting archive url")
	}

	return ioutil.NopCloser(bytes.NewReader(buf)), nil
}
