package gget

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type ResourceMatcher string

func (o *ResourceMatcher) Match(remote string) bool {
	match, _ := filepath.Match(string(*o), remote)

	return match
}

func (o *ResourceMatcher) Validate() error {
	_, err := filepath.Match(string(*o), "test")
	if err != nil {
		return errors.Wrap(err, "expected valid Resource matcher")
	}

	return nil
}

type ResourceMatchers []ResourceMatcher

func (o ResourceMatchers) Match(remote string) ResourceMatchers {
	var res ResourceMatchers

	for _, m := range o {
		if !m.Match(remote) {
			continue
		}

		res = append(res, m)
	}

	return res
}

func (o ResourceMatchers) IsEmpty() bool {
	return len(o) == 0
}

type ResourceTransferSpec struct {
	RemoteMatch ResourceMatcher
	LocalPath   string
}

func (o *ResourceTransferSpec) Resolve(remote string) (ResourceTransferSpec, bool) {
	match := o.RemoteMatch.Match(remote)
	if !match {
		return ResourceTransferSpec{}, false
	}

	res := ResourceTransferSpec{
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

func (o *ResourceTransferSpec) UnmarshalFlag(data string) error {
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
