package checksum

import (
	"context"

	"github.com/pkg/errors"
)

type deferredChecksum struct {
	manager   Manager
	resource  string
	algorithm string
}

func NewDeferredChecksum(manager Manager, resource string, algorithm string) Checksum {
	return &deferredChecksum{
		manager:   manager,
		resource:  resource,
		algorithm: algorithm,
	}
}

func (c deferredChecksum) Algorithm() string {
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
	checksum, err := c.manager.GetChecksum(ctx, c.resource)
	if err != nil {
		return nil, errors.Wrap(err, "getting deferred checksum")
	} else if checksum == nil {
		return nil, errors.New("missing deferred checksum")
	}

	return checksum, nil
}
