package downloader

import (
	"context"
	"io"

	"github.com/vbauerster/mpb/v4/decor"
)

type DownloadWriterInstaller struct {
	FilePath string
	Writer   io.Writer
}

var _ Step = &DownloadWriterInstaller{}
var _ StepWriter = &DownloadWriterInstaller{}

func (dpi *DownloadWriterInstaller) GetProgressParams() (int64, decor.Decorator) {
	return 0, nil
}

func (dpi *DownloadWriterInstaller) GetWriter() (io.Writer, error) {
	return dpi.Writer, nil
}

func (dpi *DownloadWriterInstaller) Execute(_ context.Context, s *State) error {
	if dpi.FilePath != "" {
		s.LocalFilePath = dpi.FilePath
	}

	return nil
}
