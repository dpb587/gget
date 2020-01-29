package service

import "context"

type Ref struct {
	Server     string
	Owner      string
	Repository string
	Ref        string
}

type RefResolver interface {
	ResolveRef(ctx context.Context, ref Ref) (ResolvedRef, error)
}

type ResolvedRef interface {
	ResourceResolver
}
