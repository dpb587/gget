package downloader

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

type Workflow struct {
	downloader DownloadAsset
	steps      []Step

	pb   *mpb.Progress
	bars []*mpb.Bar
}

func NewWorkflow(downloader DownloadAsset, steps ...Step) *Workflow {
	return &Workflow{
		downloader: downloader,
		steps:      steps,
	}
}

func (w *Workflow) GetSubject() string {
	return w.downloader.GetName()
}

func (w *Workflow) Prepare(pb *mpb.Progress) {
	w.pb = pb
	w.bars = make([]*mpb.Bar, len(w.steps)+4)

	w.bars[0] = w.newBar(pb, nil, " ", 1, decor.Name(
		"waiting",
		decor.WC{W: 7, C: decor.DidentRight},
	))

	w.bars[1] = w.newBar(pb, w.bars[0], "", 1, decor.Name(
		"connecting",
		decor.WC{W: 10, C: decor.DidentRight},
	))

	w.bars[2] = w.newBar(pb, w.bars[1], "", w.downloader.GetSize(), decor.OnComplete(
		decor.NewPercentage("downloading (%d)"),
		"downloaded",
	))

	lastBar := w.bars[2]

	for stepIdx, step := range w.steps {
		count, msg := step.GetProgressParams()
		if count == 0 {
			continue
		}

		w.bars[stepIdx+3] = w.newBar(pb, lastBar, "", count, msg)
		lastBar = w.bars[stepIdx+3]
	}

	w.bars[len(w.steps)+3] = w.newBar(pb, lastBar, "", 1, decor.Name(
		"finishing",
		decor.WC{W: 9, C: decor.DidentRight},
	))
}

func (w Workflow) Execute(ctx context.Context) error {
	{ // waiting
		w.bars[0].SetTotal(1, true)
	}

	var assetHandle io.ReadCloser

	{ // connecting
		var err error

		assetHandle, err = w.downloader.Open(ctx)
		if err != nil {
			w.bars[1].Abort(false)

			return errors.Wrap(err, "connecting")
		}

		w.bars[1].SetTotal(1, true)
	}

	defer assetHandle.Close()

	{ // downloading
		r := w.bars[2].ProxyReader(assetHandle)
		defer r.Close()

		var dw io.Writer

		{ // enumerate writers
			var writers []io.Writer

			for _, step := range w.steps {
				ssw, ok := step.(StepWriter)
				if !ok {
					continue
				}

				sw, err := ssw.GetWriter()
				if err != nil {
					w.bars[2].Abort(false)

					return errors.Wrap(err, "getting step writer")
				} else if sw == nil {
					continue
				}

				writers = append(writers, sw)
			}

			if len(writers) == 0 {
				w.bars[2].Abort(false)

				return fmt.Errorf("no download target found")
			}

			dw = io.MultiWriter(writers...)
		}

		_, err := io.Copy(dw, r)
		if err != nil {
			w.bars[2].Abort(false)

			return errors.Wrap(err, "downloading")
		}

		w.bars[2].SetTotal(int64(w.downloader.GetSize()), true)
	}

	var results []string

	{ // stepwise
		state := State{}

		for stepIdx, step := range w.steps {
			state.Bar = w.bars[stepIdx+3]

			err := step.Execute(ctx, &state)
			if err != nil {
				if state.Bar != nil {
					state.Bar.Abort(false)
				}

				return errors.Wrapf(err, "processing step %d", stepIdx)
			}

			if state.Bar != nil {
				count, _ := step.GetProgressParams()
				state.Bar.SetTotal(count, true)
			}
		}

		results = state.Results
	}

	{ // done
		name := fmt.Sprintf("done")
		if len(results) > 0 {
			name = fmt.Sprintf("%s (%s)", name, strings.Join(results, "; "))
		}

		pbb := w.newBar(w.pb, w.bars[len(w.steps)+3], "√", 1, decor.Name(
			name,
			decor.WC{W: len(name), C: decor.DidentRight},
		))

		w.bars[len(w.steps)+3].SetTotal(1, true)
		pbb.SetTotal(1, true)
	}

	return nil
}

func (w Workflow) newBar(pb *mpb.Progress, pbp *mpb.Bar, spinner string, count int64, msg decor.Decorator) *mpb.Bar {
	var spinnerd decor.Decorator

	switch spinner {
	case "":
		spinnerd = decor.Spinner(
			mpb.DefaultSpinnerStyle,
			decor.WC{W: 1, C: decor.DSyncSpaceR},
		)
	default:
		spinnerd = decor.Name(
			spinner,
			decor.WC{W: 1, C: decor.DSyncSpaceR},
		)
	}

	subject := w.downloader.GetName()

	return pb.AddBar(
		count,
		mpb.BarParkTo(pbp),
		mpb.PrependDecorators(
			spinnerd,
			decor.Name(
				subject,
				decor.WC{W: len(subject), C: decor.DSyncSpaceR},
			),
			msg,
		),
	)
}
