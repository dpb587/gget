package export

import (
	"context"
	"io"
	"text/template"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/pkg/errors"
)

type GoTemplateExporter struct {
	Template             *template.Template
	ChecksumVerification checksum.VerificationProfile
}

func (e *GoTemplateExporter) Export(ctx context.Context, w io.Writer, data *Data) error {
	// TODO lazy load and helper methods
	res, err := newMarshalData(ctx, data, e.ChecksumVerification)
	if err != nil {
		return errors.Wrap(err, "preparing export")
	}

	return e.Template.Execute(w, res)
}
