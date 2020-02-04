package app

import (
	"fmt"
	"time"
)

type Version struct {
	Name   string
	Semver string
	Commit string
	Built  time.Time
}

func MustVersion(name string, semver string, commit string, built_ string) Version {
	if name == "" {
		name = "unknown"
	}

	if semver == "" {
		semver = "0.0.0+dev"
	}

	if commit == "" {
		commit = "unknown"
	}

	var builtval time.Time
	var err error

	if built_ == "" {
		builtval = time.Now()
	} else {
		builtval, err = time.Parse(time.RFC3339, built_)
		if err != nil {
			panic(fmt.Errorf("cannot parse version time: %s", built_))
		}
	}

	return Version{name, semver, commit, builtval}
}

func (v Version) Version() string {
	return fmt.Sprintf("%s/%s", v.Name, v.Semver)
}

func (v Version) String() string {
	return fmt.Sprintf("%s/%s (commit %s; built %s)", v.Name, v.Semver, v.Commit, v.Built.Format(time.RFC3339))
}
