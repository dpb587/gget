package parser

import (
	"encoding/hex"
	"regexp"

	"github.com/dpb587/gget/pkg/checksum"
)

var codeindent = regexp.MustCompile("    ([a-f0-9]{40,64})\\s+([^\\s]+)")

func ImportMarkdownCodeindent(m checksum.WriteableManager, content []byte) {
	contentSubmatch := codeindent.FindAllStringSubmatch(string(content), -1)

	if len(contentSubmatch) == 0 {
		return
	}

	for _, match := range contentSubmatch {
		hashBytes, err := hex.DecodeString(match[1])
		if err != nil {
			// TODO log?
			continue
		}

		checksum, err := checksum.GuessChecksum(hashBytes)
		if err != nil {
			// TODO log?
			continue
		}

		m.AddChecksum(match[2], checksum)
	}
}
