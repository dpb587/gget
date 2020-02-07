package checksum

import "context"

type InMemoryManager struct {
	checksums map[string]Checksum
}

var _ Manager = &InMemoryManager{}

func NewInMemoryManager() *InMemoryManager {
	return &InMemoryManager{
		checksums: map[string]Checksum{},
	}
}

func (m *InMemoryManager) GetChecksum(ctx context.Context, resource string) (Checksum, bool, error) {
	res, found := m.checksums[resource]
	if !found {
		return Checksum{}, false, nil
	}

	return res, true, nil
}

func (m *InMemoryManager) SetChecksum(resource string, checksum Checksum) {
	m.checksums[resource] = checksum
}
