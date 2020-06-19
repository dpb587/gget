package service

import (
	"strings"

	"github.com/Masterminds/semver"
)

type LookupRef struct {
	Ref
	RefVersions  []*semver.Constraints
	RefStability []string
}

func (lr LookupRef) SatisfiesStability(actual string) bool {
	if len(lr.RefStability) == 0 {
		return true
	}

	for _, desired := range lr.RefStability {
		if desired == "any" {
			return true
		} else if desired == actual {
			return true
		}
	}

	return false
}

func (lr LookupRef) SatisfiesVersion(actual string) (bool, error) {
	if len(lr.RefVersions) == 0 {
		return true, nil
	}

	ver, err := semver.NewVersion(strings.TrimPrefix(actual, "v"))
	if err != nil {
		return false, err
	}

	for _, desired := range lr.RefVersions {
		if !desired.Check(ver) {
			return false, nil
		}
	}

	return true, nil
}

func (lr LookupRef) ComplexRefModes() []string {
	var res []string

	if len(lr.RefVersions) > 0 {
		res = append(res, "version")
	}

	if len(lr.RefStability) > 0 {
		res = append(res, "stability")
	}

	return res
}

func (lr LookupRef) IsComplexRef() bool {
	return len(lr.ComplexRefModes()) > 0
}
