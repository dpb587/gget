package opt

import (
	"path/filepath"
	"strings"
)

type ResourceTransfer struct {
	RemoteMatch ResourceMatcher
	LocalPath   string
}

type ResourceTransferList []ResourceTransfer

func (o *ResourceTransfer) Resolve(remote string) (ResourceTransfer, bool) {
	match := o.RemoteMatch.Match(remote)
	if !match {
		return ResourceTransfer{}, false
	}

	res := ResourceTransfer{
		RemoteMatch: ResourceMatcher(remote),
		LocalPath:   o.LocalPath,
	}

	if res.LocalPath == "" {
		res.LocalPath = remote
	} else if strings.HasSuffix(res.LocalPath, "/") {
		res.LocalPath = filepath.Join(res.LocalPath, string(res.RemoteMatch))
	}

	return res, true
}

func (o *ResourceTransfer) UnmarshalFlag(data string) error {
	dataSplit := strings.SplitN(data, "=", 2)

	if len(dataSplit) == 2 {
		o.RemoteMatch = ResourceMatcher(dataSplit[1])
		o.LocalPath = dataSplit[0]
	} else {
		o.RemoteMatch = ResourceMatcher(dataSplit[0])
		o.LocalPath = ""
	}

	if err := o.RemoteMatch.Validate(); err != nil {
		return err
	}

	return nil
}
