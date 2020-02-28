package parser

import (
	"regexp"
	"strings"

	"github.com/dpb587/gget/pkg/checksum"
)

var markdownCodefence = regexp.MustCompile("(?mU)```([^`]+)```")

func ImportMarkdownCodefence(m checksum.ManagerSetter, content string) {
	contentSubmatches := markdownCodefence.FindAllStringSubmatch(content, -1)

	if len(contentSubmatches) == 0 {
		return
	}

	for _, contentSubmatch := range contentSubmatches {
		checksums := strings.Split(strings.TrimSpace(contentSubmatch[1]), "\n")

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
}
