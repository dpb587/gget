package opt

import (
	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

type Constraint struct {
	*semver.Constraints
	RawConstraint string
}

func (o *Constraint) UnmarshalFlag(data string) error {
	// having this early allows a hack where the global --version= value
	// can be non-compliant, but it does its own validation before the
	// rest of go-flags errors are emitted
	o.RawConstraint = data

	con, err := semver.NewConstraint(data)
	if err != nil {
		return errors.Wrap(err, "parsing version constraint")
	}

	o.Constraints = con

	return nil
}

type ConstraintList []*Constraint

func (cl ConstraintList) Constraints() []*semver.Constraints {
	var res []*semver.Constraints

	for _, v := range cl {
		res = append(res, v.Constraints)
	}

	return res
}
