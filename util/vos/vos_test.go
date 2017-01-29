package vos

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("vos", func() {
	Context("New", func() {
		It("should return Vos", func() {
			/* arrange/act/assert */
			Expect(New()).Should(Not(BeNil()))
		})
	})
	Context("Getenv proceeding Setenv", func() {
		It("should return value set by Setenv", func() {
			/* arrange */
			providedName := "dummyName"
			providedValue := "dummyValue"
			expectedValue := providedValue
			objectUnderTest := New()

			objectUnderTest.Setenv(providedName, providedValue)

			/* act */
			actualValue := objectUnderTest.Getenv(providedName)

			/* assert */
			Expect(actualValue).To(Equal(expectedValue))
		})
	})
})
