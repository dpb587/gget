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

func (m MultiManager) GetChecksum(ctx context.Context, resource string) (Checksum, error) {
	var checksums []Checksum

	for managerIdx, manager := range m.managers {
		checksum, err := manager.GetChecksum(ctx, resource)
		if err != nil {
			return nil, errors.Wrapf(err, "manager %d", managerIdx)
		} else if checksum != nil {
			checksums = append(checksums, checksum)
		}
	}

	return StrongestChecksum(checksums), nil
}
