package checksum

import "hash"

type Checksum struct {
	Type   string
	Bytes  string
	Hasher func() hash.Hash
}
