package transfer

import (
	"context"
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/tidwall/limiter"
	"github.com/vbauerster/mpb/v4"
)

type Batch struct {
	transfers []*Transfer
	limiter   *limiter.Limiter
	output    io.Writer

	errs  []error
	errsM sync.Mutex
}

func NewBatch(transfers []*Transfer, parallel int, output io.Writer) *Batch {
	return &Batch{
		transfers: transfers,
		limiter:   limiter.New(parallel),
		output:    output,
	}
}

func (b *Batch) Transfer(ctx context.Context) error {
	pb := mpb.New(mpb.WithWidth(1), mpb.WithOutput(b.output))

	for _, d := range b.transfers {
		d.Prepare(pb)
	}

	for idx := range b.transfers {
		go b.transfer(idx, ctx)
	}

	pb.Wait()

	if len(b.errs) > 0 {
		// TODO multierr
		return b.errs[0]
	}

	return nil
}

func (b *Batch) transfer(idx int, ctx context.Context) {
	b.limiter.Begin()
	defer b.limiter.End()

	xfer := b.transfers[idx]

	if err := xfer.Execute(ctx); err != nil {
		b.errsM.Lock()
		b.errs = append(b.errs, errors.Wrapf(err, "downloading %s", xfer.GetSubject()))
		b.errsM.Unlock()
	}
}
