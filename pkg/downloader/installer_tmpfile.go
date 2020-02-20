package downloader

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v4/decor"
)

type DownloadTmpfileInstaller struct {
	Tmpdir string

	tmpfile *os.File
}

var _ Step = &DownloadTmpfileInstaller{}
var _ StepWriter = &DownloadTmpfileInstaller{}

func (dpi *DownloadTmpfileInstaller) GetProgressParams() (int64, decor.Decorator) {
	return 0, nil
}

func (dpi *DownloadTmpfileInstaller) GetWriter() (io.Writer, error) {
	p, err := ioutil.TempFile(dpi.Tmpdir, ".gget-*")
	if err != nil {
		return nil, errors.Wrap(err, "creating tempfile")
	}

	dpi.tmpfile = p

	return dpi.tmpfile, nil
}

func (dpi *DownloadTmpfileInstaller) Execute(_ context.Context, s *State) error {
	err := dpi.tmpfile.Close()
	if err != nil {
		return errors.Wrap(err, "closing file")
	}

	s.LocalFilePath = dpi.tmpfile.Name()

	return nil
}
