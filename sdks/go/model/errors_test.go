package model

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("ErrDataProviderAuthentication", func() {
	Context("Error", func() {
		It("should return expected result", func() {
			/* arrange */
			expectedResult := "unauthenticated"
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
			expectedResult := "unauthorized"
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
			expectedResult := "not found"
			objectUnderTest := ErrDataRefResolution{}

			/* act */
			actualResult := objectUnderTest.Error()

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))

		})
	})
})
