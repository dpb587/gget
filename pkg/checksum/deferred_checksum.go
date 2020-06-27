package checksum

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type deferredChecksum struct {
	manager   Manager
	resource  string
	algorithm Algorithm
}

func NewDeferredChecksum(manager Manager, resource string, algorithm Algorithm) Checksum {
	return &deferredChecksum{
		manager:   manager,
		resource:  resource,
		algorithm: algorithm,
	}
}

func (c deferredChecksum) Algorithm() Algorithm {
	return c.algorithm
}

func (c deferredChecksum) NewVerifier(ctx context.Context) (*HashVerifier, error) {
	checksum, err := c.requireChecksum(ctx)
	if err != nil {
		return nil, err
	}

	return checksum.NewVerifier(ctx)
}

func (c *deferredChecksum) requireChecksum(ctx context.Context) (Checksum, error) {
	checksums, err := c.manager.GetChecksums(ctx, c.resource, AlgorithmList{c.algorithm})
	if err != nil {
		return nil, errors.Wrap(err, "getting deferred checksum")
	} else if len(checksums) != 1 {
		return nil, fmt.Errorf("expected deferred checksum: %s", c.algorithm)
	}

	return checksums[0], nil
}
