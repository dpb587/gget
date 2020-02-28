package github

import (
	"context"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/dpb587/gget/pkg/checksum/parser"
	"github.com/google/go-github/v29/github"
)

type ReleaseChecksumManager struct {
	client  *github.Client
	release *github.RepositoryRelease

	known *checksum.InMemoryManager
}

var _ checksum.Manager = &ReleaseChecksumManager{}

func NewReleaseChecksumManager(client *github.Client, release *github.RepositoryRelease) *ReleaseChecksumManager {
	return &ReleaseChecksumManager{
		client:  client,
		release: release,
	}
}

func (cm *ReleaseChecksumManager) GetChecksum(ctx context.Context, resource string) (checksum.Checksum, bool, error) {
	err := cm.requireOptimistic()
	if err != nil {
		return checksum.Checksum{}, false, err
	}

	res, found, _ := cm.known.GetChecksum(ctx, resource) // never errs
	if found {
		return res, true, nil
	}

	// TODO on-demand
	return checksum.Checksum{}, false, nil
}

func (cm *ReleaseChecksumManager) requireOptimistic() error {
	if cm.known != nil {
		return nil
	}

	manager := parser.ParseMarkdown(cm.release.GetBody())
	if manager == nil {
		manager = checksum.NewInMemoryManager()
	}

	cm.known = manager

	return nil
}
