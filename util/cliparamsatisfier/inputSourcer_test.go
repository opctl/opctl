package cliparamsatisfier

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("inputSourcer", func() {
	Context("NewInputSourcer()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewInputSourcer([]InputSrc{})).Should(Not(BeNil()))
		})
	})
	Context("Source()", func() {
		It("should call 1st source w/ expected args", func() {
			/* arrange */
			providedInputName := "dummyInputName"

			fakeSource := new(FakeInputSrc)

			sources := []InputSrc{fakeSource}

			objectUnderTest := NewInputSourcer(sources)

			/* act */
			objectUnderTest.Source(providedInputName)

			/* assert */
			Expect(fakeSource.ReadArgsForCall(0)).To(Equal(providedInputName))
		})
		Context("source returns nil", func() {
			Context("2nd src doesn't exist", func() {
				It("should return nil", func() {
					/* arrange */
					fakeSource := new(FakeInputSrc)

					sources := []InputSrc{fakeSource}

					objectUnderTest := NewInputSourcer(sources)

					/* act */
					actualValue := objectUnderTest.Source("")

					/* assert */
					Expect(actualValue).To(BeNil())
				})
			})
			Context("2nd source exists", func() {
				It("should call 2nd src w/ expected args", func() {
					/* arrange */
					providedInputName := "dummyInputName"

					fakeSource2 := new(FakeInputSrc)

					sources := []InputSrc{new(FakeInputSrc), fakeSource2}

					objectUnderTest := NewInputSourcer(sources)

					/* act */
					objectUnderTest.Source(providedInputName)

					/* assert */
					Expect(fakeSource2.ReadArgsForCall(0)).To(Equal(providedInputName))
				})
			})
		})
		Context("1st source doesn't return nil", func() {
			It("should return value", func() {
				/* arrange */
				providedInputName := "dummyInputName"
				expectedValue := "dummyValue"

				fakeSource := new(FakeInputSrc)
				fakeSource.ReadReturns(&expectedValue)

				sources := []InputSrc{fakeSource}

				objectUnderTest := NewInputSourcer(sources)

				/* act */
				actualValue := objectUnderTest.Source(providedInputName)

				/* assert */
				Expect(actualValue).To(Equal(&expectedValue))
			})
		})
	})
})
