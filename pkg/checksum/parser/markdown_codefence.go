package parser

import (
	"regexp"

	"github.com/dpb587/gget/pkg/checksum"
)

var markdownCodefence = regexp.MustCompile("(?mU)```([^`]+)```")

func ImportMarkdownCodefence(m checksum.WriteableManager, content []byte) {
	contentSubmatches := markdownCodefence.FindAllSubmatch(content, -1)

	if len(contentSubmatches) == 0 {
		return
	}

	for _, contentSubmatch := range contentSubmatches {
		ImportLines(m, contentSubmatch[1])
	}
}
