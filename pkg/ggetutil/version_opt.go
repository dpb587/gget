package ggetutil

import "github.com/dpb587/gget/pkg/cli/opt"

type VersionOpt struct {
	Constraint opt.Constraint
	IsLatest   bool
}

func (v *VersionOpt) UnmarshalFlag(data string) error {
	if data == "latest" {
		v.IsLatest = true

		return nil
	}

	return v.Constraint.UnmarshalFlag(data)
}
