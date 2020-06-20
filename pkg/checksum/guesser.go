package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
)

func GuessChecksum(expected []byte) (Checksum, error) {
	var hasher func() hash.Hash
	var algorithm string

	switch len(expected) {
	case 16:
		algorithm = "md5"
		hasher = md5.New
	case 20:
		algorithm = "sha1"
		hasher = sha1.New
	case 32:
		algorithm = "sha256"
		hasher = sha256.New
	case 64:
		algorithm = "sha512"
		hasher = sha512.New
	default:
		return nil, fmt.Errorf("unrecognized hash: %s", expected)
	}

	return NewHashChecksum(algorithm, expected, hasher), nil
}
