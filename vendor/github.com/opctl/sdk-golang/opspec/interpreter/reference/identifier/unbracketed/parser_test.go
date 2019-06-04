package unbracketed

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("Parser", func() {
	Context("NewParser", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewParser()).Should(Not(BeNil()))
		})
	})
	Context("Parse", func() {
		Context("ref doesn't contain '.' or '['", func() {
			It("should return expected result", func() {
				/* arrange */
				noDotOrBracket := "noDotOrBracket"

				objectUnderTest := _parser{}

				/* act */
				actualRef, actualRefRemainder := objectUnderTest.Parse(
					noDotOrBracket,
				)

				/* assert */
				Expect(actualRef).To(Equal(noDotOrBracket))
				Expect(actualRefRemainder).To(BeEmpty())

			})
		})
		Context("ref contains '.' then '[' then '/'", func() {
			It("should return expected result", func() {
				/* arrange */
				dotThenBracket := "ta.da["

				expectedRef := "ta"
				expectedRefRemainder := ".da["

				objectUnderTest := _parser{}

				/* act */
				actualRef, actualRefRemainder := objectUnderTest.Parse(
					dotThenBracket,
				)

				/* assert */
				Expect(actualRef).To(Equal(expectedRef))
				Expect(actualRefRemainder).To(Equal(expectedRefRemainder))

			})
		})
		Context("ref contains '[' then '.' then '/'", func() {
			It("should return expected result", func() {
				/* arrange */
				dotThenBracket := "ta[da."

				expectedRef := "ta"
				expectedRefRemainder := "[da."

				objectUnderTest := _parser{}

				/* act */
				actualRef, actualRefRemainder := objectUnderTest.Parse(
					dotThenBracket,
				)

				/* assert */
				Expect(actualRef).To(Equal(expectedRef))
				Expect(actualRefRemainder).To(Equal(expectedRefRemainder))

			})
		})
		Context("ref contains '/' then '.' then '['", func() {
			It("should return expected result", func() {
				/* arrange */
				dotThenBracket := "ta/da.["

				expectedRef := "ta"
				expectedRefRemainder := "/da.["

				objectUnderTest := _parser{}

				/* act */
				actualRef, actualRefRemainder := objectUnderTest.Parse(
					dotThenBracket,
				)

				/* assert */
				Expect(actualRef).To(Equal(expectedRef))
				Expect(actualRefRemainder).To(Equal(expectedRefRemainder))

			})
		})
	})
})
