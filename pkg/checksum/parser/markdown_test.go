package parser_test

import (
	"context"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dpb587/gget/pkg/checksum"
	. "github.com/dpb587/gget/pkg/checksum/parser"
)

var _ = Describe("Markdown", func() {
	var csm *checksum.InMemoryManager

	BeforeEach(func() {
		csm = checksum.NewInMemoryManager()
	})

	Context("code fences", func() {
		It("parses", func() {
			ImportMarkdownCodefence(csm, strings.Join([]string{
				"dear release note readers. here are your checksums",
				"```",
				"bc894542e78dace00fc0357d4c591cc1e2193877636ad5b9da3c5dfc9b790320  gget-0.1.1-darwin-amd64",
				"9b0731100e631ca92d5f6979f30e3e3cc275c84f466647462d7afa0819801348  gget-0.1.1-linux-amd64",
				"```",
			}, "\n"))

			cs, found, err := csm.GetChecksum(context.Background(), "gget-0.1.1-darwin-amd64")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(cs.Bytes).To(Equal("bc894542e78dace00fc0357d4c591cc1e2193877636ad5b9da3c5dfc9b790320"))

			cs, found, err = csm.GetChecksum(context.Background(), "gget-0.1.1-linux-amd64")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(cs.Bytes).To(Equal("9b0731100e631ca92d5f6979f30e3e3cc275c84f466647462d7afa0819801348"))
		})

		It("ignores unexpected data", func() {
			ImportMarkdownCodefence(csm, strings.Join([]string{
				"dear release note readers. here are your checksums",
				"",
				"```",
				"other",
				"```",
				"",
				"some other note",
				"",
				"```",
				"bc894542e78dace00fc0357d4c591cc1e2193877636ad5b9da3c5dfc9b790320  gget-0.1.1-darwin-amd64",
				"```",
			}, "\n"))

			cs, found, err := csm.GetChecksum(context.Background(), "gget-0.1.1-darwin-amd64")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(cs.Bytes).To(Equal("bc894542e78dace00fc0357d4c591cc1e2193877636ad5b9da3c5dfc9b790320"))
		})
	})

	Context("code indent", func() {
		It("parses", func() {
			ImportMarkdownCodeindent(csm, strings.Join([]string{
				"dear release note readers. here are your checksums",
				"",
				"    bc894542e78dace00fc0357d4c591cc1e2193877636ad5b9da3c5dfc9b790320  gget-0.1.1-darwin-amd64",
				"    9b0731100e631ca92d5f6979f30e3e3cc275c84f466647462d7afa0819801348  gget-0.1.1-linux-amd64",
			}, "\n"))

			cs, found, err := csm.GetChecksum(context.Background(), "gget-0.1.1-darwin-amd64")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(cs.Bytes).To(Equal("bc894542e78dace00fc0357d4c591cc1e2193877636ad5b9da3c5dfc9b790320"))

			cs, found, err = csm.GetChecksum(context.Background(), "gget-0.1.1-linux-amd64")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(cs.Bytes).To(Equal("9b0731100e631ca92d5f6979f30e3e3cc275c84f466647462d7afa0819801348"))
		})

		It("ignores unexpected data", func() {
			ImportMarkdownCodeindent(csm, strings.Join([]string{
				"dear release note readers. here are your checksums",
				"",
				"    other",
				"",
				"some other note",
				"",
				"    bc894542e78dace00fc0357d4c591cc1e2193877636ad5b9da3c5dfc9b790320  gget-0.1.1-darwin-amd64",
			}, "\n"))

			cs, found, err := csm.GetChecksum(context.Background(), "gget-0.1.1-darwin-amd64")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(cs.Bytes).To(Equal("bc894542e78dace00fc0357d4c591cc1e2193877636ad5b9da3c5dfc9b790320"))
		})
	})
})
