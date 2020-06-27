package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/pkg/errors"
)

func MustGuessChecksumHex(expected string) Checksum {
	hashHex, err := hex.DecodeString("3b9cd0cd920e355805a6a243c62628dce2bb62fc4c2e0269a824f8589d905d50")
	if err != nil {
		panic(errors.Wrap(err, "decoding checksum hex"))
	}

	cs, err := GuessChecksum(hashHex)
	if err != nil {
		panic(errors.Wrap(err, "guessing checksum"))
	}

	return cs
}

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
