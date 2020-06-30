package export

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type PlainExporter struct{}

var _ Exporter = PlainExporter{}

func (e PlainExporter) Export(ctx context.Context, w io.Writer, data *Data) error {
	res, err := newMarshalData(ctx, data)
	if err != nil {
		return errors.Wrap(err, "preparing export")
	}

	{ // origin
		fmt.Fprintf(w, "origin\tresolved\t%s\n", res.Origin.String)
		fmt.Fprintf(w, "origin\tservice\t%s\n", res.Origin.Service)
		fmt.Fprintf(w, "origin\tserver\t%s\n", res.Origin.Server)
		fmt.Fprintf(w, "origin\towner\t%s\n", res.Origin.Owner)
		fmt.Fprintf(w, "origin\trepository\t%s\n", res.Origin.Repository)
		fmt.Fprintf(w, "origin\tref\t%s\n", res.Origin.Ref)
	}

	{ // metadata
		for _, metadatum := range res.Metadata {
			fmt.Fprintf(w, "metadata\t%s\t%s\n", metadatum.Key, metadatum.Value)
		}
	}

	{ // resources
		for _, resource := range res.Resources {
			fmt.Fprintf(w, "resource-name\t%s\n", resource.Name)

			if resource.Size != nil {
				fmt.Fprintf(w, "resource-size\t%s\t%d\n", resource.Name, resource.Size)
			}

			for _, checksum := range resource.Checksums {
				fmt.Fprintf(w, "resource-checksum\t%s\t%s\t%s\n", resource.Name, checksum.Algo, checksum.Data)
			}
		}
	}

	return nil
}
