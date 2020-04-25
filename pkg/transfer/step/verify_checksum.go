package step

import (
	"context"
	"fmt"
	"hash"
	"io"

	"github.com/dpb587/gget/pkg/transfer"
	"github.com/vbauerster/mpb/v4/decor"
)

type VerifyChecksum struct {
	Algo     string
	Expected string
	Actual   hash.Hash
}

var _ transfer.Step = &VerifyChecksum{}
var _ io.Writer = &VerifyChecksum{}

func (dhv *VerifyChecksum) GetProgressParams() (int64, decor.Decorator) {
	name := fmt.Sprintf("verifying (%s)", dhv.Algo)

	return 1, decor.Name(name, decor.WC{W: len(name), C: decor.DidentRight})
}

func (dhv *VerifyChecksum) Write(in []byte) (n int, err error) {
	return dhv.Actual.Write(in)
}

func (dhv *VerifyChecksum) Execute(_ context.Context, s *transfer.State) error {
	actual := fmt.Sprintf("%x", dhv.Actual.Sum(nil))

	if dhv.Expected != actual {
		return fmt.Errorf("expected hash %s: hash is %s", dhv.Expected, actual)
	}

	s.Results = append(s.Results, fmt.Sprintf("%s OK", dhv.Algo))

	return nil
}
