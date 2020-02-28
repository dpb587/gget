package checksum

import "context"

type Manager interface {
	GetChecksum(ctx context.Context, resource string) (Checksum, bool, error)
}

type ManagerSetter interface {
	SetChecksum(string, Checksum)
}
