package step

import (
	"context"
	"fmt"
	"io"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/dpb587/gget/pkg/transfer"
	"github.com/vbauerster/mpb/v4/decor"
)

type VerifyChecksum struct {
	Verifier *checksum.HashVerifier
}

var _ transfer.Step = &VerifyChecksum{}
var _ io.Writer = &VerifyChecksum{}

func (dhv *VerifyChecksum) GetProgressParams() (int64, decor.Decorator) {
	name := fmt.Sprintf("verifying (%s)", dhv.Verifier.Algorithm())

	return 1, decor.Name(name, decor.WC{W: len(name), C: decor.DidentRight})
}

func (dhv *VerifyChecksum) Write(in []byte) (n int, err error) {
	return dhv.Verifier.Write(in)
}

func (dhv *VerifyChecksum) Execute(_ context.Context, s *transfer.State) error {
	err := dhv.Verifier.Verify()
	if err != nil {
		return err
	}

	s.Results = append(s.Results, fmt.Sprintf("%s OK", dhv.Verifier.Algorithm()))

	return nil
}
