package vruntime

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"runtime"
)

var _ = Context("vruntime", func() {
	Context("New", func() {
		It("should return Vruntime", func() {
			/* arrange/act/assert */
			Expect(New()).To(Not(BeNil()))
		})
	})
	Context("GOOS", func() {
		It("should return runtime.GOOS", func() {
			/* arrange */
			expectedResult := runtime.GOOS

			objectUnderTest := New()

			/* act */
			actualResult := objectUnderTest.GOOS()

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))
		})
	})
})
