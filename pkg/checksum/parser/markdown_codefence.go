package parser

import (
	"regexp"

	"github.com/dpb587/gget/pkg/checksum"
)

var markdownCodefence = regexp.MustCompile("(?mU)```([^`]+)```")

func ImportMarkdownCodefence(m checksum.ManagerSetter, content string) {
	contentSubmatches := markdownCodefence.FindAllStringSubmatch(content, -1)

	if len(contentSubmatches) == 0 {
		return
	}

	for _, contentSubmatch := range contentSubmatches {
		ImportLines(m, contentSubmatch[1])
	}
}
