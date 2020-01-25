package model

import "hash"

type Checksum struct {
	Type   string
	Bytes  string
	Hasher func() hash.Hash
}

type ChecksumMap map[string]Checksum
