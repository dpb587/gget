package parser

import (
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/dpb587/gget/pkg/checksum"
)

var lines = regexp.MustCompile(`^([a-f0-9]{40,64})\s+([^\s]+)$`)

func ImportLines(m checksum.WriteableManager, content []byte) {
	checksums := strings.Split(strings.TrimSpace(string(content)), "\n")

	for _, checksumLine := range checksums {
		checksumSplit := strings.Fields(strings.TrimSpace(checksumLine))
		if len(checksumSplit) != 2 {
			continue
		}

		if len(checksumSplit[0]) < 16 {
			continue
		}

		hashBytes, err := hex.DecodeString(checksumSplit[0])
		if err != nil {
			// TODO log?
			continue
		}

		checksum, err := checksum.GuessChecksum(hashBytes)
		if err != nil {
			// TODO log?
			continue
		}

		m.AddChecksum(checksumSplit[1], checksum)
	}
}
