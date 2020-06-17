package service

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/Masterminds/semver"
)

type Ref struct {
	Server     string
	Owner      string
	Repository string
	Ref        string
}

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

func (r Ref) String() string {
	str := r.Server

	if str == "" {
		str = "{:server-missing}"
	}

	str = path.Join(str, r.Owner, r.Repository)

	if r.Ref != "" {
		str = fmt.Sprintf("%s@%s", str, r.Ref)
	}

	return str
}

func ParseRefString(in string) (Ref, error) {
	slugVersion := strings.SplitN(in, "@", 2)
	ownerRepo := strings.SplitN(slugVersion[0], "/", 3)

	res := Ref{}

	if len(slugVersion) == 2 {
		res.Ref = slugVersion[1]
	} else {
		res.Ref = ""
	}

	if len(ownerRepo) == 3 {
		res.Server = ownerRepo[0]
		res.Owner = ownerRepo[1]
		res.Repository = ownerRepo[2]
	} else if len(ownerRepo) == 2 {
		res.Server = ""
		res.Owner = ownerRepo[0]
		res.Repository = ownerRepo[1]
	} else {
		return Ref{}, fmt.Errorf("input does not match expected format: [server/]owner/repository[@version]; received %s", in)
	}

	return res, nil
}

type RefMetadata struct {
	Name  string
	Value string
}

type RefResolver interface {
	ResolveRef(ctx context.Context, ref LookupRef) (ResolvedRef, error)
}

type ResolvedRef interface {
	ResourceResolver

	CanonicalRef() Ref
	GetMetadata() []RefMetadata
}

type RefMetadataService struct {
	Metadata []RefMetadata
}

func (rmf RefMetadataService) GetMetadata() []RefMetadata {
	return rmf.Metadata
}
