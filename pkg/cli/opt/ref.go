package opt

import (
	"github.com/dpb587/gget/pkg/service"
	"github.com/pkg/errors"
)

type Ref service.Ref

func (o *Ref) UnmarshalFlag(data string) error {
	parsed, err := service.ParseRefString(data)
	if err != nil {
		return errors.Wrap(err, "parsing ref option")
	}

	*o = Ref(parsed)

	return nil
}
