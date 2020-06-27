package checksum

import (
	"context"
	"hash"
)

type hashChecksum struct {
	algorithm Algorithm
	expected  []byte
	hasher    func() hash.Hash
}

var _ Checksum = hashChecksum{}

func NewHashChecksum(algorithm Algorithm, expected []byte, hasher func() hash.Hash) Checksum {
	return hashChecksum{
		algorithm: algorithm,
		expected:  expected,
		hasher:    hasher,
	}
}

func (c hashChecksum) Algorithm() Algorithm {
	return c.algorithm
}

func (c hashChecksum) NewVerifier(ctx context.Context) (*HashVerifier, error) {
	return NewHashVerifier(c.algorithm, c.expected, c.hasher()), nil
}
