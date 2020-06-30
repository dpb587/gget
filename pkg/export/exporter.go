package export

import (
	"context"
	"io"
)

type Exporter interface {
	Export(ctx context.Context, w io.Writer, data *Data) error
}

type TemplatedExporter interface {
	ParseTemplate(text string) error
}
