package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

func main() {
	doneWg := new(sync.WaitGroup)

	pb := mpb.New(mpb.WithWidth(1), mpb.WithWaitGroup(doneWg))

	pbp1 := pb.AddBar(
		1,
		mpb.PrependDecorators(
			decor.Spinner(
				mpb.DefaultSpinnerStyle,
				decor.WC{W: 1, C: decor.DSyncSpaceR},
			),
			decor.Name(
				"one",
				decor.WC{W: len("one"), C: decor.DSyncSpaceR},
			),
			decor.Name(
				"queued",
				decor.WC{W: 4, C: decor.DidentRight},
			),
		),
	)

	pbp2 := pb.AddBar(
		1,
		mpb.BarParkTo(pbp1),
		mpb.PrependDecorators(
			decor.Spinner(
				mpb.DefaultSpinnerStyle,
				decor.WC{W: 1, C: decor.DSyncSpaceR},
			),
			decor.Name(
				"two",
				decor.WC{W: len("two"), C: decor.DSyncSpaceR},
			),
			decor.OnComplete(
				decor.Name(
					"downloading",
					decor.WC{W: 12, C: decor.DidentRight},
				),
				"downloaded",
			),
			decor.OnComplete(
				decor.NewPercentage("(%d)"),
				"",
			),
		),
	)

	time.Sleep(time.Second)

	pbp1.SetTotal(1, true)

	time.Sleep(time.Second)

	time.Sleep(time.Second)

	pbp2.SetTotal(1, false)

	time.Sleep(time.Second)

	var idx int
	pbp2.TraverseDecorators(func(d decor.Decorator) {
		idx = idx + 1

		if idx == 2 {
			fmt.Println("asdf")
			d = decor.Name(
				"installing",
				decor.WC{W: 12, C: decor.DidentRight},
			)
		}
	})

	time.Sleep(time.Second)

	pbp2.SetTotal(1, false)
}
