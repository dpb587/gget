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
	var algorithm Algorithm

	switch len(expected) {
	case 16:
		algorithm = MD5
		hasher = md5.New
	case 20:
		algorithm = SHA1
		hasher = sha1.New
	case 32:
		algorithm = SHA256
		hasher = sha256.New
	case 48:
		algorithm = SHA384
		hasher = sha512.New384
	case 64:
		algorithm = SHA512
		hasher = sha512.New
	default:
		return nil, fmt.Errorf("unrecognized hash: %s", expected)
	}

	return NewHashChecksum(algorithm, expected, hasher), nil
}
