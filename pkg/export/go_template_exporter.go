package export

import (
	"context"
	"io"
	"text/template"

	"github.com/pkg/errors"
)

type GoTemplateExporter struct {
	template *template.Template
}

var _ Exporter = &GoTemplateExporter{}
var _ TemplatedExporter = &GoTemplateExporter{}

func (e *GoTemplateExporter) ParseTemplate(text string) error {
	tmpl, err := template.New("export").Parse(text)
	if err != nil {
		return err
	}

	e.template = tmpl

	return nil
}

func (e *GoTemplateExporter) Export(ctx context.Context, w io.Writer, data *Data) error {
	// TODO lazy load and helper methods
	res, err := newMarshalData(ctx, data)
	if err != nil {
		return errors.Wrap(err, "preparing export")
	}

	err = e.template.Execute(w, res)
	if err != nil {
		return errors.Wrap(err, "executing template")
	}

	w.Write([]byte("\n"))

	return nil
}
