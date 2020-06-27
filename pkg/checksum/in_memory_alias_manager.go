package checksum

import "context"

// InMemoryAliasManager enforces a specific file name is used for any added checksums.
//
// Namely useful for name-based, *.sha256-type files are used and the contents may have been generated using a different
// file name, but there is high certainty about the subject. The caller should always be using the expected name.
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

func (m *InMemoryAliasManager) GetChecksums(ctx context.Context, resource string, algos AlgorithmList) (ChecksumList, error) {
	if m.resource != resource {
		return nil, nil
	}

	return m.manager.GetChecksums(ctx, m.resource, algos)
}

func (m *InMemoryAliasManager) AddChecksum(_ string, checksum Checksum) {
	m.manager.AddChecksum(m.resource, checksum)
}
