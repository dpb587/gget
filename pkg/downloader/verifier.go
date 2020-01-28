package downloader

import (
	"context"
	"fmt"
	"hash"
	"io"

	"github.com/vbauerster/mpb/v4/decor"
)

type DownloadHashVerifier struct {
	Algo     string
	Expected string
	Actual   hash.Hash
}

var _ Step = &DownloadHashVerifier{}
var _ StepWriter = &DownloadHashVerifier{}

func (dhv *DownloadHashVerifier) GetProgressParams() (int64, decor.Decorator) {
	name := fmt.Sprintf("verifying (%s)", dhv.Algo)

	return 1, decor.Name(name, decor.WC{W: len(name), C: decor.DidentRight})
}

func (dhv *DownloadHashVerifier) GetWriter() (io.Writer, error) {
	return dhv.Actual, nil
}

func (dhv *DownloadHashVerifier) Execute(_ context.Context, s *State) error {
	actual := fmt.Sprintf("%x", dhv.Actual.Sum(nil))

	if dhv.Expected != actual {
		return fmt.Errorf("expected hash %s: hash is %s", dhv.Expected, actual)
	}

	s.Results = append(s.Results, fmt.Sprintf("%s OK", dhv.Algo))

	return nil
}
