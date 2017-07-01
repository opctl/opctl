package pkg

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("ErrAuthenticationFailed", func() {
	Context("Error", func() {
		It("should return expected result", func() {
			/* arrange */
			expectedResult := "Authentication failed while attempting to transport package"
			objectUnderTest := ErrAuthenticationFailed{}

			/* act */
			actualResult := objectUnderTest.Error()

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))

		})
	})
})
