package checksum

import "context"

type InMemoryAliasManager struct {
	manager  WriteableManager
	resource string
}

var _ WriteableManager = &InMemoryAliasManager{}

func NewInMemoryAliasManager(resource string) WriteableManager {
	return &InMemoryAliasManager{
		manager:  NewInMemoryManager(),
		resource: resource,
	}
}

func (m *InMemoryAliasManager) GetChecksum(ctx context.Context, _ string) (Checksum, error) {
	return m.manager.GetChecksum(ctx, m.resource)
}

func (m *InMemoryAliasManager) AddChecksum(_ string, checksum Checksum) {
	m.manager.AddChecksum(m.resource, checksum)
}
