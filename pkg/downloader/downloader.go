package downloader

import (
	"context"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

type Download struct {
	asset      DownloadAsset
	verifiers  []DownloadVerifier
	installers []DownloadInstaller

	pb  *mpb.Progress
	pbp *mpb.Bar

	result []string
}

func NewDownload(asset DownloadAsset) *Download {
	return &Download{
		asset: asset,
	}
}

func (d *Download) AddVerifier(v DownloadVerifier) {
	d.verifiers = append(d.verifiers, v)
}

func (d *Download) AddInstaller(i DownloadInstaller) {
	d.installers = append(d.installers, i)
}

func (d *Download) SetProgressBar(pb *mpb.Progress) {
	d.pb = pb

	d.pbp = pb.AddBar(
		1,
		mpb.PrependDecorators(
			decor.Name(
				" ",
				decor.WC{W: 1, C: decor.DSyncSpaceR},
			),
			decor.Name(
				d.GetName(),
				decor.WC{W: len(d.GetName()), C: decor.DSyncSpaceR},
			),
			decor.Name(
				"queued",
				decor.WC{W: 4, C: decor.DidentRight},
			),
		),
	)

	d.pbp.SetTotal(1, true)

	for true {
		if d.pbp.Completed() {
			break
		}
	}
}

func (d *Download) GetName() string {
	return d.asset.GetName()
}

func (d *Download) Download(ctx context.Context) error {
	tmp, err := ioutil.TempFile("", "ghet-asset-*")
	if err != nil {
		return errors.Wrap(err, "opening temporary file")
	}

	err = d.download(ctx, tmp)
	if err != nil {
		defer os.RemoveAll(tmp.Name())

		return errors.Wrap(err, "downloading")
	}

	err = d.verify(ctx, tmp)
	if err != nil {
		defer os.RemoveAll(tmp.Name())

		return errors.Wrap(err, "verifying")
	}

	err = d.install(ctx, tmp)
	if err != nil {
		defer os.RemoveAll(tmp.Name())

		return errors.Wrap(err, "installing")
	}

	{
		done := "done"

		if len(d.result) > 0 {
			done = fmt.Sprintf("%s (%s)", done, strings.Join(d.result, "; "))
		}

		pbp := d.pb.AddBar(
			1,
			// mpb.BarParkTo(d.pbp),
			mpb.PrependDecorators(
				decor.Name(
					"√",
					decor.WC{W: 1, C: decor.DSyncSpaceR},
				),
				decor.Name(
					d.GetName(),
					decor.WC{W: len(d.GetName()), C: decor.DSyncSpaceR},
				),
				decor.Name(
					done,
					decor.WC{W: len(done), C: decor.DidentRight},
				),
			),
		)

		pbp.SetTotal(1, true)

		for true {
			if pbp.Completed() {
				break
			}
		}
	}

	return nil
}

func (d *Download) verify(ctx context.Context, tmp *os.File) error {
	for _, verifier := range d.verifiers {
		name := verifier.GetName()

		b := d.pb.AddBar(
			1,
			// mpb.BarParkTo(d.pbp),
			mpb.PrependDecorators(
				decor.Spinner(
					mpb.DefaultSpinnerStyle,
					decor.WC{W: 1, C: decor.DSyncSpaceR},
				),
				decor.Name(
					d.GetName(),
					decor.WC{W: len(name), C: decor.DSyncSpaceR},
				),
				decor.OnComplete(
					decor.Name(
						fmt.Sprintf("verifying (%s)", name),
						decor.WC{W: len(fmt.Sprintf("verifying (%s)", name)), C: decor.DidentRight},
					),
					fmt.Sprintf("verified (%s)", name),
				),
			),
		)

		res, err := verifier.Verify(ctx, tmp.Name())
		if err != nil {
			b.Abort(false)

			return errors.Wrapf(err, "verifying %s", verifier.GetName())
		}

		if res != "" {
			d.result = append(d.result, res)
		}

		b.SetTotal(1, true)

		for true {
			if b.Completed() {
				break
			}
		}
	}

	return nil
}
func (d *Download) install(ctx context.Context, tmp *os.File) error {
	for _, installer := range d.installers {
		name := installer.GetName()

		b := d.pb.AddBar(
			1,
			// mpb.BarParkTo(d.pbp),
			mpb.PrependDecorators(
				decor.Spinner(
					mpb.DefaultSpinnerStyle,
					decor.WC{W: 1, C: decor.DSyncSpaceR},
				),
				decor.Name(
					d.GetName(),
					decor.WC{W: len(name), C: decor.DSyncSpaceR},
				),
				decor.OnComplete(
					decor.Name(
						"installing",
						decor.WC{W: 10, C: decor.DidentRight},
					),
					"installed",
				),
			),
		)

		res, err := installer.Install(ctx, tmp.Name())
		if err != nil {
			b.Abort(false)

			return errors.Wrap(err, "installing")
		}

		if res != "" {
			d.result = append(d.result, res)
		}

		b.SetTotal(1, true)

		for true {
			if b.Completed() {
				break
			}
		}
	}

	return nil
}

func (d *Download) download(ctx context.Context, tmp *os.File) error {
	b := d.pb.AddBar(
		int64(d.asset.GetSize()),
		// mpb.BarParkTo(d.pbp),
		mpb.PrependDecorators(
			decor.Spinner(
				mpb.DefaultSpinnerStyle,
				decor.WC{W: 1, C: decor.DSyncSpaceR},
			),
			decor.Name(
				d.GetName(),
				decor.WC{W: len(d.asset.GetName()), C: decor.DSyncSpaceR},
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

	assetHandle, err := d.asset.Open(ctx)
	if err != nil {
		b.Abort(false)

		return errors.Wrap(err, "preparing read")
	}

	defer assetHandle.Close()

	r := b.ProxyReader(assetHandle)
	defer r.Close()

	var w io.Writer = tmp

	{
		writers := []io.Writer{tmp}

		for _, verifier := range d.verifiers {
			vw := verifier.GetStreamWriter()
			if vw == nil {
				continue
			}

			writers = append(writers, vw)
		}

		if len(writers) > 1 {
			w = io.MultiWriter(writers...)
		}
	}

	_, err = io.Copy(w, r)
	if err != nil {
		b.Abort(false)

		return errors.Wrapf(err, "downloading %s", d.asset.GetName())
	}

	b.SetTotal(int64(d.asset.GetSize()), true)

	for true {
		if b.Completed() {
			break
		}
	}

	return nil
}

type DownloadAsset interface {
	GetName() string
	GetSize() int
	Open(ctx context.Context) (io.ReadCloser, error)
}

type DownloadVerifier interface {
	GetName() string
	GetStreamWriter() io.Writer
	Verify(ctx context.Context, tmpfile string) (string, error)
}

type DownloadHashVerifier struct {
	Algo     string
	Expected string
	Actual   hash.Hash
}

func (dhv DownloadHashVerifier) GetName() string {
	return dhv.Algo
}

func (dhv DownloadHashVerifier) GetStreamWriter() io.Writer {
	return dhv.Actual
}

func (dhv DownloadHashVerifier) Verify(_ context.Context, _ string) (string, error) {
	actual := fmt.Sprintf("%x", dhv.Actual.Sum(nil))

	if dhv.Expected != actual {
		return "", fmt.Errorf("expected hash %s: got %s", dhv.Expected, actual)
	}

	return fmt.Sprintf("%s OK", dhv.Algo), nil
}

type DownloadInstaller interface {
	GetName() string
	Install(ctx context.Context, tmpfile string) (string, error)
}

type DownloadPathInstaller struct {
	Target string
}

func (dpi DownloadPathInstaller) GetName() string {
	return "installing"
}

func (dpi DownloadPathInstaller) Install(_ context.Context, tmpfile string) (string, error) {
	err := os.Rename(tmpfile, dpi.Target)
	if err != nil {
		return "", errors.Wrap(err, "renaming")
	}

	return "", nil
}
