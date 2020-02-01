package service

import "context"

type Ref struct {
	Server     string
	Owner      string
	Repository string
	Ref        string
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
