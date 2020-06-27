package service

import (
	"context"
)

type RefResolver interface {
	ResolveRef(ctx context.Context, ref LookupRef) (ResolvedRef, error)
}

type ResolvedRef interface {
	ResourceResolver

	CanonicalRef() Ref
	GetMetadata(ctx context.Context) (RefMetadata, error)
}
