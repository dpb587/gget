package checksum_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dpb587/gget/pkg/checksum"
)

var _ = Describe("Guesser", func() {
	It("guesses md5", func() {
		cs, err := GuessChecksum("c299ca339c89f3d6b4580adf05b4a779")
		Expect(err).ToNot(HaveOccurred())
		Expect(cs.Type).To(Equal("md5"))
		Expect(cs.Bytes).To(Equal("c299ca339c89f3d6b4580adf05b4a779"))

		h := cs.Hasher()
		h.Write([]byte("hashable"))
		Expect(fmt.Sprintf("%x", h.Sum(nil))).To(Equal("dd16fe9f76604a7400d5e1fcf88afaca"))
	})

	It("guesses sha1", func() {
		cs, err := GuessChecksum("de7daa47ddfc899daba52efd07e0ec9b8f5a944f")
		Expect(err).ToNot(HaveOccurred())
		Expect(cs.Type).To(Equal("sha1"))
		Expect(cs.Bytes).To(Equal("de7daa47ddfc899daba52efd07e0ec9b8f5a944f"))

		h := cs.Hasher()
		h.Write([]byte("hashable"))
		Expect(fmt.Sprintf("%x", h.Sum(nil))).To(Equal("705d0123108e62b9a94842986e4f12c7ef0a9239"))
	})

	It("guesses sha256", func() {
		cs, err := GuessChecksum("bc894542e78dace00fc0357d4c591cc1e2193877636ad5b9da3c5dfc9b790320")
		Expect(err).ToNot(HaveOccurred())
		Expect(cs.Type).To(Equal("sha256"))
		Expect(cs.Bytes).To(Equal("bc894542e78dace00fc0357d4c591cc1e2193877636ad5b9da3c5dfc9b790320"))

		h := cs.Hasher()
		h.Write([]byte("hashable"))
		Expect(fmt.Sprintf("%x", h.Sum(nil))).To(Equal("3b9cd0cd920e355805a6a243c62628dce2bb62fc4c2e0269a824f8589d905d50"))
	})

	It("guesses sha512", func() {
		cs, err := GuessChecksum("f2e5b29f4e3167d7944a07cd3d1905ac3c98aeb6deb6e532bc77b1142f4cff288bb3785de7f398a3465e44bf15d243608b538b2736e0f233bf86470cc80532e3")
		Expect(err).ToNot(HaveOccurred())
		Expect(cs.Type).To(Equal("sha512"))
		Expect(cs.Bytes).To(Equal("f2e5b29f4e3167d7944a07cd3d1905ac3c98aeb6deb6e532bc77b1142f4cff288bb3785de7f398a3465e44bf15d243608b538b2736e0f233bf86470cc80532e3"))

		h := cs.Hasher()
		h.Write([]byte("hashable"))
		Expect(fmt.Sprintf("%x", h.Sum(nil))).To(Equal("768dd75ae44b3d5537d047ef454c15833326602e568c1dbc31e5198c0e9a76380b8e392df6625adcb1f1411ad520c15f514008f4306196059dbe726a9e64c4da"))
	})

	It("errors for unknown", func() {
		_, err := GuessChecksum("oy")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unrecognized checksum"))
	})
})
