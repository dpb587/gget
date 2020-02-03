package gget

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type ResourceNameOpt string

func (o *ResourceNameOpt) Match(remote string) bool {
	match, _ := filepath.Match(string(*o), remote)

	return match
}

func (o *ResourceNameOpt) Validate() error {
	_, err := filepath.Match(string(*o), "test")
	if err != nil {
		return errors.Wrap(err, "expected valid Resource matcher")
	}

	return nil
}

type ResourcePathOpt struct {
	RemoteMatch ResourceNameOpt
	LocalPath   string
}

func (o *ResourcePathOpt) Resolve(remote string) (ResourcePathOpt, bool) {
	match := o.RemoteMatch.Match(remote)
	if !match {
		return ResourcePathOpt{}, false
	}

	res := ResourcePathOpt{
		RemoteMatch: ResourceNameOpt(remote),
		LocalPath:   o.LocalPath,
	}

	if res.LocalPath == "" {
		res.LocalPath = remote
	} else if strings.HasSuffix(res.LocalPath, "/") {
		res.LocalPath = filepath.Join(res.LocalPath, string(res.RemoteMatch))
	}

	return res, true
}

func (o *ResourcePathOpt) UnmarshalFlag(data string) error {
	dataSplit := strings.SplitN(data, "=", 2)

	if len(dataSplit) == 2 {
		o.RemoteMatch = ResourceNameOpt(dataSplit[1])
		o.LocalPath = dataSplit[0]
	} else {
		o.RemoteMatch = ResourceNameOpt(dataSplit[0])
		o.LocalPath = ""
	}

	if err := o.RemoteMatch.Validate(); err != nil {
		return err
	}

	return nil
}
