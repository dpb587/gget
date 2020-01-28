package downloader

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v4/decor"
)

type DownloadRenameInstaller struct {
	Target string
}

var _ Step = &DownloadHashVerifier{}

func (dpi DownloadRenameInstaller) GetProgressParams() (int64, decor.Decorator) {
	return 0, nil
}

func (dpi DownloadRenameInstaller) Execute(_ context.Context, state *State) error {
	err := os.Rename(state.LocalFilePath, dpi.Target)
	if err != nil {
		return errors.Wrap(err, "renaming")
	}

	state.LocalFilePath = dpi.Target

	return nil
}
