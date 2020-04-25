package step

import (
	"context"
	"os"

	"github.com/dpb587/gget/pkg/transfer"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v4/decor"
)

type Executable struct{}

var _ transfer.Step = &Executable{}

func (dpi Executable) GetProgressParams() (int64, decor.Decorator) {
	return 0, nil
}

func (dpi Executable) Execute(_ context.Context, state *transfer.State) error {
	err := os.Chmod(state.LocalFilePath, 0755)
	if err != nil {
		return errors.Wrap(err, "chmod'ing")
	}

	state.Results = append(state.Results, "executable")

	return nil
}
