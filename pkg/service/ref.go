package service

import (
	"context"
	"fmt"
	"strings"
)

type Ref struct {
	Server     string
	Owner      string
	Repository string
	Ref        string
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
	ResolveRef(ctx context.Context, ref Ref) (ResolvedRef, error)
}

type ResolvedRef interface {
	ResourceResolver

	GetMetadata() []RefMetadata
}

type RefMetadataService struct {
	Metadata []RefMetadata
}

func (rmf RefMetadataService) GetMetadata() []RefMetadata {
	return rmf.Metadata
}
