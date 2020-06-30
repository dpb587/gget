package export

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"k8s.io/client-go/util/jsonpath"
)

type JSONPathExporter struct {
	template *jsonpath.JSONPath
}

var _ Exporter = &JSONPathExporter{}
var _ TemplatedExporter = &JSONPathExporter{}

func (e *JSONPathExporter) ParseTemplate(text string) error {
	e.template = jsonpath.New("export")

	return e.template.Parse(text)
}

func (e *JSONPathExporter) Export(ctx context.Context, w io.Writer, data *Data) error {
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
