package transfer

import (
	"context"
	"io"

	"github.com/vbauerster/mpb/v4"
)

type State struct {
	Bar           *mpb.Bar
	LocalFilePath string
	Results       []string
}

type DownloadAsset interface {
	GetName() string
	GetSize() int64
	Open(ctx context.Context) (io.ReadCloser, error)
}
