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

func (m *InMemoryAliasManager) GetChecksums(ctx context.Context, _ string, algos AlgorithmList) (ChecksumList, error) {
	return m.manager.GetChecksums(ctx, m.resource, algos)
}

func (m *InMemoryAliasManager) AddChecksum(_ string, checksum Checksum) {
	m.manager.AddChecksum(m.resource, checksum)
}
