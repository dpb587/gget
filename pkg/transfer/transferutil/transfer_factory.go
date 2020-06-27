package transferutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/dpb587/gget/pkg/service"
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

	if opts.ChecksumMode != "none" { // verify checksum
		var cs checksum.Checksum
		var err error

		if csr, ok := origin.(service.ChecksumSupportedResolvedResource); ok {
			cs, err = csr.GetChecksum(ctx, opts.ChecksumAcceptableAlgorithms)
			if err != nil {
				return nil, errors.Wrap(err, "getting checksum")
			}
		}

		if cs == nil && opts.ChecksumMode == "required" {
			return nil, fmt.Errorf("checksum required but not found: %s", opts.ChecksumAcceptableAlgorithms.Join(", "))
		} else if cs != nil {
			verifier, err := cs.NewVerifier(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "getting %s verifier", cs.Algorithm())
			}

			steps = append(
				steps,
				&step.VerifyChecksum{
					Verifier: verifier,
				},
			)
		}
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
	Executable                   bool
	ChecksumMode                 string
	ChecksumAcceptableAlgorithms checksum.AlgorithmList
}
