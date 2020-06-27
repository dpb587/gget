package checksum

import "fmt"

type HashVerificationError struct {
	algorithm Algorithm
	expected  []byte
	actual    []byte
}

var _ error = HashVerificationError{}

func (err HashVerificationError) Error() string {
	return fmt.Sprintf("expected %s checksum %x, but found %x", err.algorithm, err.expected, err.actual)
}
