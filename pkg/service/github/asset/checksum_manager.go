package asset

import (
	"github.com/dpb587/gget/pkg/checksum"
	"github.com/dpb587/gget/pkg/model"
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

func (cm *ChecksumManager) GetAssetChecksum(asset string) (model.Checksum, bool, error) {
	err := cm.requireKnownChecksums()
	if err != nil {
		return model.Checksum{}, false, err
	}

	cs, found := cm.known[asset]
	if !found {
		return model.Checksum{}, false, nil
	}

	return cs, true, nil
}

func (cm *ChecksumManager) requireKnownChecksums() error {
	if cm.known != nil {
		return nil
	}

	cm.known = checksum.ParseReleaseNotes(cm.release.GetBody())

	return nil
}
