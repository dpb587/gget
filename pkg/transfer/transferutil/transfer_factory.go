package transferutil

import (
	"context"
	"os"
	"path/filepath"

	"github.com/dpb587/gget/pkg/transfer"
	"github.com/dpb587/gget/pkg/transfer/step"
	"github.com/pkg/errors"
)

func BuildTransfer(ctx context.Context, origin transfer.DownloadAsset, targetPath string, opts TransferOptions) (*transfer.Transfer, error) {
	var steps []transfer.Step

	if targetPath == "-" {
		steps = append(
			steps,
			&step.WriterTarget{
				Writer: os.Stdout,
			},
		)
	} else {
		steps = append(
			steps,
			&step.TempFileTarget{
				Tmpdir: filepath.Dir(targetPath),
			},
		)
	}

	if ds, ok := origin.(transfer.StepProvider); ok {
		extraSteps, err := ds.GetTransferSteps(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "getting transfer steps")
		}

		steps = append(steps, extraSteps...)
	}

	if targetPath != "-" {
		if opts.Executable {
			steps = append(
				steps,
				&step.Executable{},
			)
		}

		steps = append(
			steps,
			&step.Rename{
				Target: targetPath,
			},
		)
	}

	return transfer.NewTransfer(origin, steps), nil
}

type TransferOptions struct {
	Executable bool
}
