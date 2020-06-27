package checksum

import (
	"context"
)

type InMemoryManager struct {
	resourceChecksums map[string]ChecksumList
}

var _ Manager = &InMemoryManager{}

func NewInMemoryManager() WriteableManager {
	return &InMemoryManager{
		resourceChecksums: map[string]ChecksumList{},
	}
}

func (m *InMemoryManager) GetChecksums(ctx context.Context, resource string, algos AlgorithmList) (ChecksumList, error) {
	res, found := m.resourceChecksums[resource]
	if !found {
		return nil, nil
	}

	if algos != nil {
		res = res.FilterAlgorithms(algos)
	}

	return res, nil
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
