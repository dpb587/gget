package service

import (
	"context"
	"io"

	"github.com/dpb587/gget/pkg/checksum"
)

type ResourceResolver interface {
	ResolveResource(ctx context.Context, resourceType ResourceType, resource ResourceName) ([]ResolvedResource, error)
}

type ResolvedResource interface {
	GetName() string
	GetSize() int64
	Open(ctx context.Context) (io.ReadCloser, error)
}

type ChecksumSupportedResolvedResource interface {
	GetChecksum(ctx context.Context, algos checksum.AlgorithmList) (checksum.Checksum, error)
}
