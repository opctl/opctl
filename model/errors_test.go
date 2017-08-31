package model

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("ErrPkgPullAuthentication", func() {
	Context("Error", func() {
		It("should return expected result", func() {
			/* arrange */
			expectedResult := "Pkg pull failed due to invalid/lack of authentication"
			objectUnderTest := ErrPkgPullAuthentication{}

			/* act */
			actualResult := objectUnderTest.Error()

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))

		})
	})
})

var _ = Context("ErrPkgPullAuthorization", func() {
	Context("Error", func() {
		It("should return expected result", func() {
			/* arrange */
			expectedResult := "Pkg pull failed due to insufficient/lack of authorization"
			objectUnderTest := ErrPkgPullAuthorization{}

			/* act */
			actualResult := objectUnderTest.Error()

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))

		})
	})
})
var _ = Context("ErrPkgNotFound", func() {
	Context("Error", func() {
		It("should return expected result", func() {
			/* arrange */
			expectedResult := "Provider failed to resolve the requested pkg"
			objectUnderTest := ErrPkgNotFound{}

			/* act */
			actualResult := objectUnderTest.Error()

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))

		})
	})
})
