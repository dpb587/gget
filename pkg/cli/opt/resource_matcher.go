package opt

import (
	"path/filepath"

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

type ResourceMatcherList []ResourceMatcher

func (o ResourceMatcherList) Match(remote string) ResourceMatcherList {
	var res ResourceMatcherList

	for _, m := range o {
		if !m.Match(remote) {
			continue
		}

		res = append(res, m)
	}

	return res
}

func (o ResourceMatcherList) IsEmpty() bool {
	return len(o) == 0
}
