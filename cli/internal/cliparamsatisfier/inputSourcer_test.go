package cliparamsatisfier

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/opctl/opctl/cli/internal/cliparamsatisfier/internal/fakes"
)

var _ = Describe("inputSourcer", func() {
	Context("NewInputSourcer()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewInputSourcer()).To(Not(BeNil()))
		})
	})
	Context("Source()", func() {
		It("should call 1st source w/ expected args", func() {
			/* arrange */
			providedInputName := "dummyInputName"

			fakeSource := new(FakeInputSrc)

			objectUnderTest := NewInputSourcer(fakeSource)

			/* act */
			objectUnderTest.Source(providedInputName)

			/* assert */
			Expect(fakeSource.ReadStringArgsForCall(0)).To(Equal(providedInputName))
		})
		Context("source fails", func() {
			Context("2nd src doesn't exist", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeSource := new(FakeInputSrc)

					objectUnderTest := NewInputSourcer(fakeSource)

					/* act */
					actualValue, actualOk := objectUnderTest.Source("")

					/* assert */
					Expect(actualValue).To(BeNil())
					Expect(actualOk).To(BeFalse())
				})
			})
			Context("2nd source exists", func() {
				It("should call 2nd src w/ expected args", func() {
					/* arrange */
					providedInputName := "dummyInputName"

					fakeSource2 := new(FakeInputSrc)

					objectUnderTest := NewInputSourcer(new(FakeInputSrc), fakeSource2)

					/* act */
					objectUnderTest.Source(providedInputName)

					/* assert */
					Expect(fakeSource2.ReadStringArgsForCall(0)).To(Equal(providedInputName))
				})
			})
		})
		Context("1st source succeeds", func() {
			It("should return expected result", func() {
				/* arrange */
				providedInputName := "dummyInputName"
				expectedValue := "dummyValue"

				fakeSource := new(FakeInputSrc)
				fakeSource.ReadStringReturns(&expectedValue, true)

				objectUnderTest := NewInputSourcer(fakeSource)

				/* act */
				actualValue, actualOk := objectUnderTest.Source(providedInputName)

				/* assert */
				Expect(actualValue).To(Equal(&expectedValue))
				Expect(actualOk).To(BeTrue())
			})
		})
	})
})
