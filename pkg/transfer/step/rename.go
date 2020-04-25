package step

import (
	"context"
	"os"

	"github.com/dpb587/gget/pkg/transfer"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v4/decor"
)

type Rename struct {
	Target string
}

var _ transfer.Step = &Rename{}

func (dpi Rename) GetProgressParams() (int64, decor.Decorator) {
	return 0, nil
}

func (dpi Rename) Execute(_ context.Context, state *transfer.State) error {
	err := os.Rename(state.LocalFilePath, dpi.Target)
	if err != nil {
		return errors.Wrap(err, "renaming")
	}

	state.LocalFilePath = dpi.Target

	return nil
}
