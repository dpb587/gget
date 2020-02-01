package service

import (
	"context"
	"io"
)

type Resource string

type ResourceType string

// ArchiveResourceType is a tar/zip export of the repository from the ref.
const ArchiveResourceType ResourceType = "archive"

// AssetResourceType is a user-provided file associated with the ref.
const AssetResourceType ResourceType = "asset"

// BlobResourceType is a blob of the repository at the ref.
const BlobResourceType ResourceType = "blob"

type ResourceResolver interface {
	ResolveResource(ctx context.Context, resourceType ResourceType, resource Resource) ([]ResolvedResource, error)
}

type ResolvedResource interface {
	GetName() string
	GetSize() int64
	// GetLocation(ctx context.Context) (string, error)
	Open(ctx context.Context) (io.ReadCloser, error)
}
