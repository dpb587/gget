package github

import (
	"context"
	"io"
	"path/filepath"
	"strings"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/dpb587/gget/pkg/checksum/parser"
	"github.com/dpb587/gget/pkg/service/github/asset"
	"github.com/google/go-github/v29/github"
)

func NewReleaseChecksumManager(client *github.Client, releaseOwner, releaseRepository string, release *github.RepositoryRelease) checksum.Manager {
	literalManager := checksum.NewInMemoryManager()
	var deferredManagers []checksum.Manager

	// parse from release notes
	parser.ImportMarkdown(literalManager, []byte(release.GetBody()))

	// checksums from convention-based file names
	for _, releaseAsset := range release.Assets {
		algorithm, resource, useful := checkReleaseAssetChecksumBehavior(releaseAsset)
		if !useful {
			continue
		}

		opener := newReleaseAssetChecksumOpener(client, releaseOwner, releaseRepository, releaseAsset)

		var expectedAlgos checksum.AlgorithmList

		if algorithm != "" && algorithm != "unknown" {
			expectedAlgos = append(expectedAlgos, algorithm)
		}

		if resource != "" {
			literalManager.AddChecksum(
				resource,
				checksum.NewDeferredChecksum(
					parser.NewDeferredManager(checksum.NewInMemoryAliasManager(resource), expectedAlgos, opener),
					resource,
					algorithm,
				),
			)
		} else if algorithm != "" {
			deferredManagers = append(deferredManagers, parser.NewDeferredManager(checksum.NewInMemoryManager(), expectedAlgos, opener))
		}
	}

	return checksum.NewMultiManager(append([]checksum.Manager{literalManager}, deferredManagers...)...)
}

func checkReleaseAssetChecksumBehavior(releaseAsset github.ReleaseAsset) (checksum.Algorithm, string, bool) {
	name := releaseAsset.GetName()
	nameLower := strings.ToLower(name)
	ext := filepath.Ext(releaseAsset.GetName())
	extLower := strings.ToLower(strings.TrimPrefix(ext, "."))

	if extLower == "md5" || extLower == "sha1" || extLower == "sha256" || extLower == "sha384" || extLower == "sha512" {
		return checksum.Algorithm(extLower), strings.TrimSuffix(name, ext), true
	} else if nameLower == "md5sum" || nameLower == "md5sums" || nameLower == "md5sum.txt" || nameLower == "md5sums.txt" {
		return checksum.MD5, "", true
	} else if nameLower == "sha1sum" || nameLower == "sha1sums" || nameLower == "sha1sum.txt" || nameLower == "sha1sums.txt" {
		return checksum.SHA1, "", true
	} else if nameLower == "sha384sum" || nameLower == "sha384sums" || nameLower == "sha384sum.txt" || nameLower == "sha384sums.txt" {
		return checksum.SHA384, "", true
	} else if nameLower == "sha256sum" || nameLower == "sha256sums" || nameLower == "sha256sum.txt" || nameLower == "sha256sums.txt" {
		return checksum.SHA256, "", true
	} else if nameLower == "sha512sum" || nameLower == "sha512sums" || nameLower == "sha512sum.txt" || nameLower == "sha512sums.txt" {
		return checksum.SHA512, "", true
	} else if nameLower == "checksum" || nameLower == "checksums" || strings.HasSuffix(nameLower, "checksum.txt") || strings.HasSuffix(nameLower, "checksums.txt") {
		return checksum.Algorithm("unknown"), "", true
	}

	return "", "", false
}

func newReleaseAssetChecksumOpener(client *github.Client, releaseOwner, releaseRepository string, releaseAsset github.ReleaseAsset) func(context.Context) (io.ReadCloser, error) {
	return func(ctx context.Context) (io.ReadCloser, error) {
		resource := asset.NewResource(client, releaseOwner, releaseRepository, releaseAsset, nil) // TODO pass shared checksum manager

		return resource.Open(ctx)
	}
}
