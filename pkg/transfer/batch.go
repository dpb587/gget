package transfer

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/limiter"
	"github.com/vbauerster/mpb/v4"
)

type Batch struct {
	log       logrus.FieldLogger
	transfers []*Transfer
	limiter   *limiter.Limiter
	output    io.Writer

	errs  []string
	errsM sync.Mutex
}

func NewBatch(log logrus.FieldLogger, transfers []*Transfer, parallel int, output io.Writer) *Batch {
	return &Batch{
		log:       log,
		transfers: transfers,
		limiter:   limiter.New(parallel),
		output:    output,
	}
}

func (b *Batch) Transfer(ctx context.Context, failFast bool) error {
	pb := mpb.New(mpb.WithWidth(1), mpb.WithOutput(b.output))

	for _, d := range b.transfers {
		d.Prepare(pb)
	}

	var cancel context.CancelFunc
	if failFast {
		ctx, cancel = context.WithCancel(ctx)
	}

	for idx := range b.transfers {
		go b.transfer(idx, ctx, cancel)
	}

	pb.Wait()

	if len(b.errs) > 0 {
		return fmt.Errorf("transfers failed: %s", strings.Join(b.errs, ", "))
	}

	return nil
}

func (b *Batch) transfer(idx int, ctx context.Context, cancel context.CancelFunc) {
	b.limiter.Begin()
	defer b.limiter.End()

	xfer := b.transfers[idx]

	b.errsM.Lock()
	errsLen := len(b.errs)
	b.errsM.Unlock()

	if errsLen > 0 && cancel != nil {
		xfer.finalize("!", "skipped (due to previous error)")

		return
	}

	if err := xfer.Execute(ctx); err != nil {
		// only warning since it should be handled/printed elsewhere
		b.log.Warn(errors.Wrapf(err, "downloading %s", xfer.GetSubject()))

		b.errsM.Lock()
		if cancel != nil && len(b.errs) > 0 {
			// assume context was canceled and ignore as root cause
			xfer.finalize("!", "aborted (due to previous error)")
		} else {
			b.errs = append(b.errs, xfer.GetSubject())
			xfer.finalize("X", fmt.Sprintf("errored (%s)", err))
		}
		b.errsM.Unlock()

		if cancel != nil {
			cancel()
		}
	}
}
