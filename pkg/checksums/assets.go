package checksums

import (
	"crypto/sha1"
	"crypto/sha256"
	"hash"
)

type ReleaseAsset struct {
	SHA  string
	Name string
}

func (ra ReleaseAsset) NewHash() hash.Hash {
	switch len(ra.SHA) {
	case 40:
		return sha1.New()
	case 64:
		return sha256.New()
	}

	panic("TODO unknown checksum")
}

type ReleaseAssets []ReleaseAsset

func (ras ReleaseAssets) GetByName(name string) (ReleaseAsset, bool) {
	for _, ra := range ras {
		if ra.Name != name {
			continue
		}

		return ra, true
	}

	return ReleaseAsset{}, false
}
