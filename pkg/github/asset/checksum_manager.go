package asset

import (
	"fmt"

	"github.com/dpb587/ghet/pkg/checksum"
	"github.com/dpb587/ghet/pkg/model"
	"github.com/google/go-github/v29/github"
)

type ChecksumManager struct {
	release *github.RepositoryRelease
	known   model.ChecksumMap
}

func NewChecksumManager(release *github.RepositoryRelease) *ChecksumManager {
	return &ChecksumManager{
		release: release,
	}
}

func (cm *ChecksumManager) GetAssetChecksum(asset string) (model.Checksum, error) {
	err := cm.requireKnownChecksums()
	if err != nil {
		return model.Checksum{}, err
	}

	cs, found := cm.known[asset]
	if !found {
		return model.Checksum{}, fmt.Errorf("expected checksum: no checksum found")
	}

	return cs, nil
}

func (cm *ChecksumManager) requireKnownChecksums() error {
	if cm.known != nil {
		return nil
	}

	cm.known = checksum.ParseReleaseNotes(cm.release.GetBody())

	return nil
}
