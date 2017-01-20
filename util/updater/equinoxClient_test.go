package updater

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("equinoxClient", func() {
	Context("New", func() {
		It("should return equinoxClient", func() {
			/* arrange/act/assert */
			Expect(newEquinoxClient()).Should(Not(BeNil()))
		})
	})
})
