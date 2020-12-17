package clicolorer

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("cliColorer", func() {
	Context("New", func() {
		It("should return CliColorer", func() {
			/* arrange/act/assert */
			Expect(New()).To(Not(BeNil()))
		})
	})
	Context("Disable", func() {
		Context("Attention", func() {
			It("should cause calls to Attention() to not color string", func() {
				/* arrange */
				objectUnderTest := New()
				providedString := "providedString"
				expectedString := providedString

				/* act */
				objectUnderTest.DisableColor()

				/* assert */
				Expect(objectUnderTest.Attention(providedString)).To(Equal(expectedString))
			})
		})
		Context("Error", func() {
			It("should cause calls to Error() to not color string", func() {
				/* arrange */
				objectUnderTest := New()
				providedString := "providedString"

				/* act */
				objectUnderTest.DisableColor()

				/* assert */
				Expect(objectUnderTest.Error(providedString)).To(Equal(providedString))
			})
		})
		Context("Info", func() {
			It("should cause calls to Info() to not color string", func() {
				/* arrange */
				objectUnderTest := New()
				providedString := "providedString"

				/* act */
				objectUnderTest.DisableColor()

				/* assert */
				Expect(objectUnderTest.Info(providedString)).To(Equal(providedString))
			})
		})
		Context("Success", func() {
			It("should cause calls to Success() to not color string", func() {
				/* arrange */
				objectUnderTest := New()
				providedString := "providedString"

				/* act */
				objectUnderTest.DisableColor()

				/* assert */
				Expect(objectUnderTest.Success(providedString)).To(Equal(providedString))
			})
		})
	})
	Context("Attention", func() {
		It("should return expected string", func() {
			/* arrange */
			objectUnderTest := New()
			providedString := "providedString"
			expectedString := fmt.Sprintf("\x1b[93;1m%s\x1b[0m", fmt.Sprint(providedString))

			/* act */
			actualString := objectUnderTest.Attention(providedString)

			/* assert */
			Expect(fmt.Sprintf("%+q", actualString)).To(Equal(fmt.Sprintf("%+q", expectedString)))
		})
	})
	Context("Error", func() {
		It("should return expected string", func() {
			/* arrange */
			objectUnderTest := New()
			providedString := "providedString"
			expectedString := fmt.Sprintf("\x1b[91;1m%s\x1b[0m", fmt.Sprint(providedString))

			/* act */
			actualString := objectUnderTest.Error(providedString)

			/* assert */
			Expect(fmt.Sprintf("%+q", actualString)).To(Equal(fmt.Sprintf("%+q", expectedString)))
		})
	})
	Context("Info", func() {
		It("should return expected string", func() {
			/* arrange */
			objectUnderTest := New()
			providedString := "providedString"
			expectedString := fmt.Sprintf("\x1b[96;1m%s\x1b[0m", fmt.Sprint(providedString))

			/* act */
			actualString := objectUnderTest.Info(providedString)

			/* assert */
			Expect(fmt.Sprintf("%+q", actualString)).To(Equal(fmt.Sprintf("%+q", expectedString)))
		})
	})
	Context("Success", func() {
		It("should return expected string", func() {
			/* arrange */
			objectUnderTest := New()
			providedString := "providedString"
			expectedString := fmt.Sprintf("\x1b[92;1m%s\x1b[0m", providedString)

			/* act */
			actualString := objectUnderTest.Success(providedString)

			/* assert */
			Expect(fmt.Sprintf("%+q", actualString)).To(Equal(fmt.Sprintf("%+q", expectedString)))
		})
	})
})
