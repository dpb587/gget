package export

import (
	"context"
	"io"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type YAMLExporter struct {
	ChecksumVerification checksum.VerificationProfile
}

var _ Exporter = YAMLExporter{}

func (e YAMLExporter) Export(ctx context.Context, w io.Writer, data *Data) error {
	res, err := newMarshalData(ctx, data, e.ChecksumVerification)
	if err != nil {
		return errors.Wrap(err, "preparing export")
	}

	buf, err := yaml.Marshal(res)
	if err != nil {
		return errors.Wrap(err, "marshalling")
	}

	w.Write(buf)
	w.Write([]byte("\n"))

	return nil
}
