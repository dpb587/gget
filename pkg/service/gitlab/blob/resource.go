package blob

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/gitlab/gitlabutil"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type Resource struct {
	client *gitlab.Client
	ref    service.Ref
	target string
	node   *gitlab.TreeNode
}

var _ service.ResolvedResource = &Resource{}

func NewResource(client *gitlab.Client, ref service.Ref, target string, node *gitlab.TreeNode) *Resource {
	return &Resource{
		client: client,
		ref:    ref,
		target: target,
		node:   node,
	}
}

func (r *Resource) GetName() string {
	return r.node.Name
}

func (r *Resource) GetSize() int64 {
	return 0
}

func (r *Resource) Open(ctx context.Context) (io.ReadCloser, error) {
	// TODO switch to stream?
	bufRes, _, err := r.client.Repositories.Blob(gitlabutil.GetRepositoryID(r.ref), r.node.ID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting blob")
	}

	var apiRes apiResponse

	err = json.Unmarshal(bufRes, &apiRes)
	if err != nil {
		return nil, errors.Wrap(err, "parsing blob api")
	}

	var buf []byte

	if apiRes.Encoding == "base64" {
		buf, err = base64.StdEncoding.DecodeString(apiRes.Content)
		if err != nil {
			return nil, errors.Wrap(err, "decoding blob")
		}
	} else {
		return nil, errors.Wrap(err, "unsupported content encoding")
	}

	return ioutil.NopCloser(bytes.NewReader(buf)), nil
}

type apiResponse struct {
	Size     int    `json:"size"`
	Encoding string `json:"encoding"`
	Content  string `json:"content"`
	SHA      string `json:"sha"`
}
