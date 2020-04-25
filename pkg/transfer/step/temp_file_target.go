package step

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"github.com/dpb587/gget/pkg/transfer"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v4/decor"
)

type TempFileTarget struct {
	Tmpdir string

	tmpfile *os.File
}

var _ transfer.Step = &TempFileTarget{}
var _ io.Writer = &TempFileTarget{}

func (dpi *TempFileTarget) GetProgressParams() (int64, decor.Decorator) {
	return 0, nil
}

func (dpi *TempFileTarget) Prepare() error {

	return nil
}

func (dpi *TempFileTarget) Write(p []byte) (int, error) {
	if dpi.tmpfile == nil {
		p, err := ioutil.TempFile(dpi.Tmpdir, ".gget-*")
		if err != nil {
			return 0, errors.Wrap(err, "creating tempfile")
		}

		dpi.tmpfile = p
	}

	return dpi.tmpfile.Write(p)
}

func (dpi *TempFileTarget) Execute(_ context.Context, s *transfer.State) error {
	err := dpi.tmpfile.Close()
	if err != nil {
		return errors.Wrap(err, "closing file")
	}

	s.LocalFilePath = dpi.tmpfile.Name()

	return nil
}
