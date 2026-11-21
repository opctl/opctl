package model

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StartOpReq", func() {
	It("encodes interactive on the wire only when requested", func() {
		encoded, err := json.Marshal(StartOpReq{Interactive: true})
		Expect(err).NotTo(HaveOccurred())
		Expect(string(encoded)).To(ContainSubstring(`"interactive":true`))

		encoded, err = json.Marshal(StartOpReq{})
		Expect(err).NotTo(HaveOccurred())
		Expect(string(encoded)).NotTo(ContainSubstring("interactive"))
	})
})
