package github

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/dpb587/gget/pkg/checksum/parser"
	"github.com/dpb587/gget/pkg/service/github/asset"
	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type ReleaseChecksumManager struct {
	client            *github.Client
	releaseOwner      string
	releaseRepository string
	release           *github.RepositoryRelease

	known            *checksum.InMemoryManager
	attemptedDynamic map[string]struct{}
}

var _ checksum.Manager = &ReleaseChecksumManager{}

func NewReleaseChecksumManager(client *github.Client, releaseOwner, releaseRepository string, release *github.RepositoryRelease) *ReleaseChecksumManager {
	return &ReleaseChecksumManager{
		client:            client,
		releaseOwner:      releaseOwner,
		releaseRepository: releaseRepository,
		release:           release,
		attemptedDynamic:  map[string]struct{}{},
	}
}

func (cm *ReleaseChecksumManager) GetChecksum(ctx context.Context, resource string) (checksum.Checksum, bool, error) {
	err := cm.requireOptimistic(ctx)
	if err != nil {
		return checksum.Checksum{}, false, err
	}

	res, found, _ := cm.known.GetChecksum(ctx, resource) // never errs
	if found {
		return res, true, nil
	}

	if _, found := cm.attemptedDynamic[resource]; !found {
		cm.attemptedDynamic[resource] = struct{}{}

		for _, releaseAsset := range cm.release.Assets {
			switch releaseAsset.GetName() {
			case fmt.Sprintf("%s.sha1", resource), fmt.Sprintf("%s.sha256", resource), fmt.Sprintf("%s.sha512", resource):
				// good
			default:
				continue
			}

			err := cm.loadReleaseAsset(ctx, releaseAsset)
			if err != nil {
				// TODO log errors.Wrap(err, "loading checksum asset")
				continue
			}
		}

		return cm.GetChecksum(ctx, resource)
	}

	return checksum.Checksum{}, false, nil
}

func (cm *ReleaseChecksumManager) requireOptimistic(ctx context.Context) error {
	if cm.known != nil {
		return nil
	}

	cm.known = checksum.NewInMemoryManager()

	parser.ImportMarkdown(cm.known, cm.release.GetBody())

	for _, releaseAsset := range cm.release.Assets {
		if !strings.HasSuffix(releaseAsset.GetName(), "checksums.txt") {
			continue
		}

		err := cm.loadReleaseAsset(ctx, releaseAsset)
		if err != nil {
			// TODO log errors.Wrap(err, "loading checksum asset")
			continue
		}
	}

	return nil
}

func (cm *ReleaseChecksumManager) loadReleaseAsset(ctx context.Context, releaseAsset github.ReleaseAsset) error {
	resource := asset.NewResource(cm.client, cm.releaseOwner, cm.releaseRepository, releaseAsset, cm.known)
	fh, err := resource.Open(ctx)
	if err != nil {
		return errors.Wrap(err, "opening resource")
	}

	defer fh.Close()

	buf, err := ioutil.ReadAll(fh)
	if err != nil {
		return errors.Wrap(err, "reading")
	}

	parser.ImportLines(cm.known, string(buf))

	return nil
}
