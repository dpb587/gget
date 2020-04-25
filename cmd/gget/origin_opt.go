package gget

import (
	"github.com/dpb587/gget/pkg/service"
	"github.com/pkg/errors"
)

type RefOpt service.Ref

func (o *RefOpt) UnmarshalFlag(data string) error {
	parsed, err := service.ParseRefString(data)
	if err != nil {
		return errors.Wrap(err, "parsing ref option")
	}

	*o = RefOpt(parsed)

	return nil
}
