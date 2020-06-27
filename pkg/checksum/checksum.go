package checksum

import (
	"context"
)

type Checksum interface {
	Algorithm() Algorithm
	NewVerifier(context.Context) (*HashVerifier, error)
}

type ChecksumList []Checksum

func (l ChecksumList) FilterAlgorithms(algorithms AlgorithmList) ChecksumList {
	var res ChecksumList

	for _, c := range l {
		for _, a := range algorithms {
			if c.Algorithm() == a {
				res = append(res, c)
			}
		}
	}

	return res
}
