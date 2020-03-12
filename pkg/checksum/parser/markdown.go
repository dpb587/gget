package parser

import "github.com/dpb587/gget/pkg/checksum"

func ImportMarkdown(m checksum.ManagerSetter, content string) {
	ImportMarkdownCodefence(m, content)
	ImportMarkdownCodeindent(m, content)
}
