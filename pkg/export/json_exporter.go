package export

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

type JSONExporter struct{}

var _ Exporter = JSONExporter{}

func (e JSONExporter) Export(ctx context.Context, w io.Writer, data *Data) error {
	res, err := newMarshalData(ctx, data)
	if err != nil {
		return errors.Wrap(err, "preparing export")
	}

	buf, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshalling")
	}

	w.Write(buf)
	w.Write([]byte("\n"))

	return nil
}
