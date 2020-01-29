package ghet

import (
	"fmt"
	"strings"

	"github.com/dpb587/ghet/pkg/service"
)

type RefOpt service.Ref

func (o *RefOpt) UnmarshalFlag(data string) error {
	slugVersion := strings.SplitN(data, "@", 2)
	ownerRepo := strings.SplitN(slugVersion[0], "/", 3)

	if len(slugVersion) == 2 {
		o.Ref = slugVersion[1]
	} else {
		o.Ref = ""
	}

	if len(ownerRepo) == 3 {
		o.Server = ownerRepo[0]
		o.Owner = ownerRepo[1]
		o.Repository = ownerRepo[2]
	} else if len(ownerRepo) == 2 {
		o.Server = ""
		o.Owner = ownerRepo[0]
		o.Repository = ownerRepo[1]
	} else {
		return fmt.Errorf("expected format: [server/]owner/repository[@version]")
	}

	return nil
}
