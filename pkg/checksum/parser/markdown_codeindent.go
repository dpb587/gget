package parser

import (
	"regexp"

	"github.com/dpb587/gget/pkg/checksum"
)

var codeindent = regexp.MustCompile("    ([a-f0-9]{40,64})\\s+([^\\s]+)")

func ImportMarkdownCodeindent(m checksum.ManagerSetter, content string) {
	contentSubmatch := codeindent.FindAllStringSubmatch(content, -1)

	if len(contentSubmatch) == 0 {
		return
	}

	for _, match := range contentSubmatch {
		checksum, err := checksum.GuessChecksum(match[1])
		if err != nil {
			continue
		}

		m.SetChecksum(match[2], checksum)
	}
}
