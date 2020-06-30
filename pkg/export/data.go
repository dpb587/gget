package export

import (
	"context"
	"sort"
	"strings"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/dpb587/gget/pkg/service"
)

type metadataGetterFunc func(ctx context.Context) (service.RefMetadata, error)

type Data struct {
	origin               service.Ref
	metadataGetter       metadataGetterFunc
	resources            []service.ResolvedResource
	checksumVerification checksum.VerificationProfile

	metadata service.RefMetadata
}

func NewData(origin service.Ref, metadataGetter metadataGetterFunc, resources []service.ResolvedResource, checksumVerification checksum.VerificationProfile) *Data {
	return &Data{
		origin:               origin,
		metadataGetter:       metadataGetter,
		resources:            resources,
		checksumVerification: checksumVerification,
	}
}

func (d *Data) Origin() service.Ref {
	return d.origin
}

func (d *Data) Metadata(ctx context.Context) (service.RefMetadata, error) {
	if d.metadata == nil {
		metadata, err := d.metadataGetter(ctx)
		if err != nil {
			return nil, err
		}

		if metadata == nil {
			// avoid empty re-lookups
			metadata = service.RefMetadata{}
		}

		// deterministic
		sort.Slice(metadata, func(i, j int) bool {
			return strings.Compare(metadata[i].Name, metadata[j].Name) < 0
		})

		d.metadata = metadata
	}

	return d.metadata, nil
}

func (d *Data) Resources() []service.ResolvedResource {
	return d.resources
}
