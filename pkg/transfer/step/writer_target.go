package step

import (
	"context"
	"io"

	"github.com/dpb587/gget/pkg/transfer"
	"github.com/vbauerster/mpb/v4/decor"
)

type WriterTarget struct {
	FilePath string
	Writer   io.Writer
}

var _ transfer.Step = &WriterTarget{}
var _ io.Writer = &WriterTarget{}

func (dpi *WriterTarget) GetProgressParams() (int64, decor.Decorator) {
	return 0, nil
}

func (dpi *WriterTarget) Write(p []byte) (int, error) {
	return dpi.Writer.Write(p)
}

func (dpi *WriterTarget) Execute(_ context.Context, s *transfer.State) error {
	if dpi.FilePath != "" {
		s.LocalFilePath = dpi.FilePath
	}

	return nil
}
