package opt

import (
	"path/filepath"
	"strings"
)

type ResourceTransfer struct {
	RemoteMatch ResourceMatcher
	LocalDir string
	LocalName   string
}

type ResourceTransferList []ResourceTransfer

func (o *ResourceTransfer) LocalPath() string {
	return filepath.Join(o.LocalDir, o.LocalName)
}

func (o *ResourceTransfer) Resolve(remote string) (ResourceTransfer, bool) {
	match := o.RemoteMatch.Match(remote)
	if !match {
		return ResourceTransfer{}, false
	}

	res := ResourceTransfer{
		RemoteMatch: ResourceMatcher(remote),
		LocalDir: o.LocalDir,
		LocalName:   o.LocalName,
	}

	if res.LocalName == "" {
		res.LocalName = remote
	}

	return res, true
}

func (o *ResourceTransfer) UnmarshalFlag(data string) error {
	dataSplit := strings.SplitN(data, "=", 2)
	localPath := ""

	if len(dataSplit) == 2 {
		o.RemoteMatch = ResourceMatcher(dataSplit[1])
		localPath = dataSplit[0]
	} else {
		o.RemoteMatch = ResourceMatcher(dataSplit[0])
	}

	o.LocalDir, o.LocalName = filepath.Split(localPath)

	if err := o.RemoteMatch.Validate(); err != nil {
		return err
	}

	return nil
}
