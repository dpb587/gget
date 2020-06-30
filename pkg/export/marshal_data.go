package export

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/dpb587/gget/pkg/service"
	"github.com/pkg/errors"
)

func newMarshalData(ctx context.Context, data *Data) (marshalData, error) {
	origin := data.Origin()

	res := marshalData{
		Origin: marshalDataOrigin{
			String:     origin.String(),
			Service:    origin.Service,
			Server:     origin.Server,
			Owner:      origin.Owner,
			Repository: origin.Repository,
			Ref:        origin.Ref,
		},
	}

	metadata, err := data.Metadata(ctx)
	if err != nil {
		return marshalData{}, errors.Wrap(err, "getting metadata")
	}

	for _, metadatum := range metadata {
		res.Metadata = append(
			res.Metadata,
			marshalDataMetadatum{
				Key:   metadatum.Name,
				Value: metadatum.Value,
			},
		)
	}

	for _, resource := range data.Resources() {
		var size *int64

		rawSize := resource.GetSize()
		if rawSize > 0 {
			// TODO differentiate between 0 size and missing
			size = &rawSize
		}

		var jsonChecksums []marshalDataResourceChecksum

		if len(data.checksumVerification.Acceptable) > 0 {
			// TODO probably ought to refactor to respect .Required here as well
			if csr, ok := resource.(service.ChecksumSupportedResolvedResource); ok {
				checksums, err := csr.GetChecksums(ctx, data.checksumVerification.Acceptable)
				if err != nil {
					return marshalData{}, errors.Wrapf(err, "getting checksums of %s", resource.GetName())
				}

				for _, checksum := range data.checksumVerification.Selector.SelectChecksums(checksums) {
					v, err := checksum.NewVerifier(ctx)
					if err != nil {
						return marshalData{}, errors.Wrapf(err, "getting checksum %s of %s", checksum.Algorithm(), resource.GetName())
					}

					jsonChecksums = append(
						jsonChecksums,
						marshalDataResourceChecksum{
							Algo: string(v.Algorithm()),
							Data: fmt.Sprintf("%x", v.Expected()),
						},
					)
				}
			}
		}

		res.Resources = append(
			res.Resources,
			marshalDataResource{
				Name:      resource.GetName(),
				Size:      size,
				Checksums: jsonChecksums,
			},
		)
	}

	sort.Slice(res.Resources, func(i, j int) bool {
		return strings.Compare(res.Resources[i].Name, res.Resources[j].Name) < 0
	})

	return res, nil
}

type marshalData struct {
	Origin    marshalDataOrigin      `json:"origin"`
	Metadata  []marshalDataMetadatum `json:"metadata"`
	Resources []marshalDataResource  `json:"resources"`
}

type marshalDataOrigin struct {
	String     string `json:"string"`
	Service    string `json:"service"`
	Server     string `json:"server"`
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
}

type marshalDataMetadatum struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type marshalDataResource struct {
	Name      string                        `json:"name"`
	Size      *int64                        `json:"size,omitempty"`
	Checksums []marshalDataResourceChecksum `json:"checksums,omitempty"`
}

type marshalDataResourceChecksum struct {
	Algo string `json:"algo"`
	Data string `json:"data"`
}
