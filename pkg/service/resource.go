package service

import (
	"context"
	"io"
)

type Resource string

type ResourceType string

type ResourceResolver interface {
	ResolveResource(ctx context.Context, resourceType ResourceType, resource Resource) ([]ResolvedResource, error)
}

type ResolvedResource interface {
	GetName() string
	GetSize() int64
	// GetLocation(ctx context.Context) (string, error)
	Open(ctx context.Context) (io.ReadCloser, error)
}
