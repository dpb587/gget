package downloader

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v4/decor"
)

type DownloadExecutableInstaller struct{}

var _ Step = &DownloadHashVerifier{}

func (dpi DownloadExecutableInstaller) GetProgressParams() (int64, decor.Decorator) {
	return 0, nil
}

func (dpi DownloadExecutableInstaller) Execute(_ context.Context, state *State) error {
	err := os.Chmod(state.LocalFilePath, 0755)
	if err != nil {
		return errors.Wrap(err, "chmod'ing")
	}

	state.Results = append(state.Results, "executable")

	return nil
}
