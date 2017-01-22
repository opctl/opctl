package updater

import (
	"github.com/equinox-io/equinox"
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
	Context("Check", func() {
		Context("error occurs", func() {
			// note: this test assumes network connectivity is disabled resulting in all http requests error'ing
			It("should return the error", func() {
				/* arrange  */
				objectUnderTest := newEquinoxClient()

				/* act */
				_, actualError := objectUnderTest.Check("dummyAppId", equinox.Options{})

				/* assert */
				Expect(actualError).ToNot(BeNil())
			})
		})
	})
})
