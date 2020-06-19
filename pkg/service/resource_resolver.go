package service

import (
	"context"
	"io"
)

type ResourceResolver interface {
	ResolveResource(ctx context.Context, resourceType ResourceType, resource ResourceName) ([]ResolvedResource, error)
}

type ResolvedResource interface {
	GetName() string
	GetSize() int64
	// GetLocation(ctx context.Context) (string, error)
	Open(ctx context.Context) (io.ReadCloser, error)
}
