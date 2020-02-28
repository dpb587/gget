package parser

import "github.com/dpb587/gget/pkg/checksum"

func ParseMarkdown(content string) *checksum.InMemoryManager {
	m := checksum.NewInMemoryManager()

	ImportMarkdownCodefence(m, content)
	ImportMarkdownCodeindent(m, content)

	return m
}
