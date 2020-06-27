package checksum

import "context"

type Manager interface {
	GetChecksums(ctx context.Context, resource string, algos AlgorithmList) (ChecksumList, error)
}

type WriteableManager interface {
	Manager
	AddChecksum(string, Checksum)
}
