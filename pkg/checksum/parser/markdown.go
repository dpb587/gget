package parser

import "github.com/dpb587/gget/pkg/checksum"

func ImportMarkdown(m checksum.WriteableManager, content []byte) {
	ImportMarkdownCodefence(m, content)
	ImportMarkdownCodeindent(m, content)
}
