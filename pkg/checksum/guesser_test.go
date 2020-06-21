package checksum_test

import (
	"context"
	"encoding/hex"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dpb587/gget/pkg/checksum"
)

var _ = Describe("Guesser", func() {
	It("guesses md5", func() {
		hashHex, _ := hex.DecodeString("dd16fe9f76604a7400d5e1fcf88afaca")

		cs, err := GuessChecksum(hashHex)
		Expect(err).ToNot(HaveOccurred())
		Expect(cs.Algorithm()).To(Equal("md5"))

		h, err := cs.NewVerifier(context.Background())
		Expect(err).ToNot(HaveOccurred())
		h.Write([]byte("hashable"))
		Expect(h.Verify()).ToNot(HaveOccurred())
	})

	It("guesses sha1", func() {
		hashHex, _ := hex.DecodeString("705d0123108e62b9a94842986e4f12c7ef0a9239")

		cs, err := GuessChecksum(hashHex)
		Expect(err).ToNot(HaveOccurred())
		Expect(cs.Algorithm()).To(Equal("sha1"))

		h, err := cs.NewVerifier(context.Background())
		Expect(err).ToNot(HaveOccurred())
		h.Write([]byte("hashable"))
		Expect(h.Verify()).ToNot(HaveOccurred())
	})

	It("guesses sha256", func() {
		hashHex, _ := hex.DecodeString("3b9cd0cd920e355805a6a243c62628dce2bb62fc4c2e0269a824f8589d905d50")

		cs, err := GuessChecksum(hashHex)
		Expect(err).ToNot(HaveOccurred())
		Expect(cs.Algorithm()).To(Equal("sha256"))

		h, err := cs.NewVerifier(context.Background())
		Expect(err).ToNot(HaveOccurred())
		h.Write([]byte("hashable"))
		Expect(h.Verify()).ToNot(HaveOccurred())
	})

	It("guesses sha512", func() {
		hashHex, _ := hex.DecodeString("768dd75ae44b3d5537d047ef454c15833326602e568c1dbc31e5198c0e9a76380b8e392df6625adcb1f1411ad520c15f514008f4306196059dbe726a9e64c4da")

		cs, err := GuessChecksum(hashHex)
		Expect(err).ToNot(HaveOccurred())
		Expect(cs.Algorithm()).To(Equal("sha512"))

		h, err := cs.NewVerifier(context.Background())
		Expect(err).ToNot(HaveOccurred())
		h.Write([]byte("hashable"))
		Expect(h.Verify()).ToNot(HaveOccurred())
	})

	It("errors for unknown", func() {
		hashHex, _ := hex.DecodeString("dead")

		_, err := GuessChecksum(hashHex)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unrecognized hash"))
	})
})
