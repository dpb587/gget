package checksum

import (
	"context"

	"github.com/pkg/errors"
)

type MultiManager struct {
	managers []Manager
}

var _ Manager = MultiManager{}

func NewMultiManager(managers ...Manager) Manager {
	return MultiManager{
		managers: managers,
	}
}

func (m MultiManager) GetChecksums(ctx context.Context, resource string, algos AlgorithmList) (ChecksumList, error) {
	var res ChecksumList

	for managerIdx, manager := range m.managers {
		checksums, err := manager.GetChecksums(ctx, resource, algos)
		if err != nil {
			return nil, errors.Wrapf(err, "manager %d", managerIdx)
		}

		res = append(res, checksums...)
	}

	return res.FilterAlgorithms(algos), nil
}
