package pkg

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ErrAuthenticationFailed", func() {
	Describe("Error", func() {
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
