package checksum

import "context"

type Manager interface {
	GetChecksum(ctx context.Context, resource string) (Checksum, error)
}

type WriteableManager interface {
	Manager
	AddChecksum(string, Checksum)
}
