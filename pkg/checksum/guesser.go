package checksum

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"

	"github.com/dpb587/gget/pkg/model"
)

func GuessChecksum(cs string) (model.Checksum, error) {
	var hasher func() hash.Hash
	var name string

	switch len(cs) {
	case 40:
		name = "sha1"
		hasher = sha1.New
	case 64:
		name = "sha256"
		hasher = sha256.New
	default:
		return model.Checksum{}, fmt.Errorf("unrecognized checksum: %s", cs)
	}

	return model.Checksum{
		Type:   name,
		Bytes:  cs,
		Hasher: hasher,
	}, nil
}
