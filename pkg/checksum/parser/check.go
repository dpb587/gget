package parser

import (
	"regexp"
	"strings"

	"github.com/dpb587/gget/pkg/checksum"
)

var lines = regexp.MustCompile(`^([a-f0-9]{40,64})\s+([^\s]+)$`)

func ImportLines(m checksum.ManagerSetter, content string) {
	checksums := strings.Split(strings.TrimSpace(content), "\n")

	for _, checksumLine := range checksums {
		checksumSplit := strings.Fields(strings.TrimSpace(checksumLine))
		if len(checksumSplit) != 2 {
			continue
		}

		if len(checksumSplit[0]) < 16 {
			continue
		}

		checksum, err := checksum.GuessChecksum(checksumSplit[0])
		if err != nil {
			continue
		}

		m.SetChecksum(checksumSplit[1], checksum)
	}
}
