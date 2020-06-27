package transferutil

import (
	"context"
	"fmt"
	"io"
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

	if len(opts.ChecksumVerification.Acceptable) > 0 {
		var csl checksum.ChecksumList

		if csr, ok := origin.(service.ChecksumSupportedResolvedResource); ok {
			avail, err := csr.GetChecksums(ctx, opts.ChecksumVerification.Acceptable)
			if err != nil {
				return nil, errors.Wrap(err, "getting checksum")
			}

			csl = opts.ChecksumVerification.Selector.SelectChecksums(avail)
		}

		if opts.ChecksumVerification.Required && len(csl) == 0 {
			return nil, fmt.Errorf("acceptable checksum required but not found: %s", opts.ChecksumVerification.Acceptable.Join(", "))
		}

		for _, cs := range csl {
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

	return transfer.NewTransfer(origin, steps, opts.FinalStatus), nil
}

type TransferOptions struct {
	Executable           bool
	ChecksumVerification checksum.VerificationProfile
	FinalStatus          io.Writer
}
