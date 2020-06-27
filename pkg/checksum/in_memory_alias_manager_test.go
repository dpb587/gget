package checksum_test

import (
	"encoding/hex"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dpb587/gget/pkg/checksum"
)

var _ = Describe("InMemoryAliasManager", func() {
	It("enforces a name", func() {
		subject := NewInMemoryAliasManager("my-test-file.zip")

		hashHex, _ := hex.DecodeString("3b9cd0cd920e355805a6a243c62628dce2bb62fc4c2e0269a824f8589d905d50")
		cs, _ := GuessChecksum(hashHex)

		subject.AddChecksum("alternative-name", cs)

		checksums, err := subject.GetChecksums(nil, "my-test-file.zip", nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(checksums).To(HaveLen(1))
		Expect(checksums[0].Algorithm()).To(Equal(SHA256))

		checksums, err = subject.GetChecksums(nil, "alternative-name", nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(checksums).To(HaveLen(0))
	})
})
