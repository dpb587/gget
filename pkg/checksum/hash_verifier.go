package checksum

import (
	"bytes"
	"hash"
	"io"
)

type HashVerifier struct {
	algorithm string
	expected  []byte
	hasher    hash.Hash
}

var _ io.Writer = &HashVerifier{}

func NewHashVerifier(algorithm string, expected []byte, hasher hash.Hash) *HashVerifier {
	return &HashVerifier{
		algorithm: algorithm,
		expected:  expected,
		hasher:    hasher,
	}
}

func (hv *HashVerifier) Algorithm() string {
	return hv.algorithm
}

func (hv *HashVerifier) Write(p []byte) (int, error) {
	return hv.hasher.Write(p)
}

func (hv *HashVerifier) Verify() error {
	actual := hv.hasher.Sum(nil)

	if bytes.Compare(hv.expected, actual) == 0 {
		return nil
	}

	return &HashVerificationError{
		algorithm: hv.algorithm,
		expected:  hv.expected,
		actual:    actual,
	}
}
