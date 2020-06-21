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
	var csm checksum.WriteableManager
	var ctx context.Context

	BeforeEach(func() {
		csm = checksum.NewInMemoryManager()
	})

	Context("code fences", func() {
		It("parses", func() {
			ImportMarkdownCodefence(csm, []byte(strings.Join([]string{
				"dear release note readers. here are your checksums",
				"```",
				"9534cebfed045f466d446f45ff1d76e38aa94941ccdbbcd8a8b82e51657a579e  gget-0.1.1-darwin-amd64",
				"9d61c2edcdb8ed71d58d94970d7ef4aeacbe1ac4bce4aecb06e2f3d804caee4b  gget-0.1.1-linux-amd64",
				"```",
			}, "\n")))

			cs, err := csm.GetChecksum(ctx, "gget-0.1.1-darwin-amd64")
			Expect(err).ToNot(HaveOccurred())
			Expect(cs).ToNot(BeNil())

			csv, err := cs.NewVerifier(ctx)
			csv.Write([]byte("gget-0.1.1-darwin-amd64"))
			Expect(csv.Verify()).ToNot(HaveOccurred())

			cs, err = csm.GetChecksum(ctx, "gget-0.1.1-linux-amd64")
			Expect(err).ToNot(HaveOccurred())
			Expect(csv).ToNot(BeNil())

			csv, err = cs.NewVerifier(ctx)
			csv.Write([]byte("gget-0.1.1-linux-amd64"))
			Expect(csv.Verify()).ToNot(HaveOccurred())
		})

		It("ignores unexpected data", func() {
			ImportMarkdownCodefence(csm, []byte(strings.Join([]string{
				"dear release note readers. here are your checksums",
				"",
				"```",
				"other",
				"```",
				"",
				"some other note",
				"",
				"```",
				"9534cebfed045f466d446f45ff1d76e38aa94941ccdbbcd8a8b82e51657a579e  gget-0.1.1-darwin-amd64",
				"```",
			}, "\n")))

			cs, err := csm.GetChecksum(ctx, "gget-0.1.1-darwin-amd64")
			Expect(err).ToNot(HaveOccurred())
			Expect(cs).ToNot(BeNil())

			csv, err := cs.NewVerifier(ctx)
			csv.Write([]byte("gget-0.1.1-darwin-amd64"))
			Expect(csv.Verify()).ToNot(HaveOccurred())
		})
	})

	Context("code indent", func() {
		It("parses", func() {
			ImportMarkdownCodeindent(csm, []byte(strings.Join([]string{
				"dear release note readers. here are your checksums",
				"",
				"    9534cebfed045f466d446f45ff1d76e38aa94941ccdbbcd8a8b82e51657a579e  gget-0.1.1-darwin-amd64",
				"    9d61c2edcdb8ed71d58d94970d7ef4aeacbe1ac4bce4aecb06e2f3d804caee4b  gget-0.1.1-linux-amd64",
			}, "\n")))

			cs, err := csm.GetChecksum(ctx, "gget-0.1.1-darwin-amd64")
			Expect(err).ToNot(HaveOccurred())
			Expect(cs).ToNot(BeNil())

			csv, err := cs.NewVerifier(ctx)
			csv.Write([]byte("gget-0.1.1-darwin-amd64"))
			Expect(csv.Verify()).ToNot(HaveOccurred())

			cs, err = csm.GetChecksum(ctx, "gget-0.1.1-linux-amd64")
			Expect(err).ToNot(HaveOccurred())
			Expect(csv).ToNot(BeNil())

			csv, err = cs.NewVerifier(ctx)
			csv.Write([]byte("gget-0.1.1-linux-amd64"))
			Expect(csv.Verify()).ToNot(HaveOccurred())
		})

		It("ignores unexpected data", func() {
			ImportMarkdownCodeindent(csm, []byte(strings.Join([]string{
				"dear release note readers. here are your checksums",
				"",
				"    other",
				"",
				"some other note",
				"",
				"    9534cebfed045f466d446f45ff1d76e38aa94941ccdbbcd8a8b82e51657a579e  gget-0.1.1-darwin-amd64",
				"```",
			}, "\n")))

			cs, err := csm.GetChecksum(ctx, "gget-0.1.1-darwin-amd64")
			Expect(err).ToNot(HaveOccurred())
			Expect(cs).ToNot(BeNil())

			csv, err := cs.NewVerifier(ctx)
			csv.Write([]byte("gget-0.1.1-darwin-amd64"))
			Expect(csv.Verify()).ToNot(HaveOccurred())
		})
	})
})
