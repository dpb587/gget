package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
)

func GuessChecksum(cs string) (Checksum, error) {
	var hasher func() hash.Hash
	var name string

	switch len(cs) {
	case 32:
		name = "md5"
		hasher = md5.New
	case 40:
		name = "sha1"
		hasher = sha1.New
	case 64:
		name = "sha256"
		hasher = sha256.New
	case 128:
		name = "sha512"
		hasher = sha512.New
	default:
		return Checksum{}, fmt.Errorf("unrecognized checksum: %s", cs)
	}

	return Checksum{
		Type:   name,
		Bytes:  cs,
		Hasher: hasher,
	}, nil
}
