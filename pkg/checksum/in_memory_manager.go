package checksum

import (
	"context"
)

type InMemoryManager struct {
	resourceChecksums map[string][]Checksum
}

var _ Manager = &InMemoryManager{}

func NewInMemoryManager() WriteableManager {
	return &InMemoryManager{
		resourceChecksums: map[string][]Checksum{},
	}
}

func (m *InMemoryManager) GetChecksum(ctx context.Context, resource string) (Checksum, error) {
	resourceChecksums, found := m.resourceChecksums[resource]
	if !found {
		return nil, nil
	}

	return StrongestChecksum(resourceChecksums), nil
}

func (m *InMemoryManager) Resources() []string {
	var res []string

	for resource := range m.resourceChecksums {
		res = append(res, resource)
	}

	return res
}

func (m *InMemoryManager) AddChecksum(resource string, checksum Checksum) {
	m.resourceChecksums[resource] = append(m.resourceChecksums[resource], checksum)
}
