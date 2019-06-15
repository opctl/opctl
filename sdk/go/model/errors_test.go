package model

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("ErrDataProviderAuthentication", func() {
	Context("Error", func() {
		It("should return expected result", func() {
			/* arrange */
			expectedResult := "Data pull failed due to invalid/lack of authentication"
			objectUnderTest := ErrDataProviderAuthentication{}

			/* act */
			actualResult := objectUnderTest.Error()

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))

		})
	})
})

var _ = Context("ErrDataProviderAuthorization", func() {
	Context("Error", func() {
		It("should return expected result", func() {
			/* arrange */
			expectedResult := "Data pull failed due to insufficient/lack of authorization"
			objectUnderTest := ErrDataProviderAuthorization{}

			/* act */
			actualResult := objectUnderTest.Error()

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))

		})
	})
})
var _ = Context("ErrDataRefResolution", func() {
	Context("Error", func() {
		It("should return expected result", func() {
			/* arrange */
			expectedResult := "Provider failed to resolve the requested data"
			objectUnderTest := ErrDataRefResolution{}

			/* act */
			actualResult := objectUnderTest.Error()

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))

		})
	})
})
