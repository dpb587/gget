package checksum

import (
	"context"
	"hash"
)

type Checksum interface {
	Algorithm() string
	NewVerifier(context.Context) (*HashVerifier, error)
}

type hashChecksum struct {
	algorithm string
	expected  []byte
	hasher    func() hash.Hash
}

var _ Checksum = hashChecksum{}

func NewHashChecksum(algorithm string, expected []byte, hasher func() hash.Hash) Checksum {
	return hashChecksum{
		algorithm: algorithm,
		expected:  expected,
		hasher:    hasher,
	}
}

func (c hashChecksum) Algorithm() string {
	return c.algorithm
}

func (c hashChecksum) NewVerifier(ctx context.Context) (*HashVerifier, error) {
	return NewHashVerifier(c.algorithm, c.expected, c.hasher()), nil
}
