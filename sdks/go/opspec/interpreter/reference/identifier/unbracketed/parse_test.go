package unbracketed

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("Parse", func() {
	Context("ref doesn't contain '.' or '['", func() {
		It("should return expected result", func() {
			/* arrange */
			noDotOrBracket := "noDotOrBracket"

			/* act */
			actualRef, actualRefRemainder := Parse(
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

			/* act */
			actualRef, actualRefRemainder := Parse(
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

			/* act */
			actualRef, actualRefRemainder := Parse(
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

			/* act */
			actualRef, actualRefRemainder := Parse(
				dotThenBracket,
			)

			/* assert */
			Expect(actualRef).To(Equal(expectedRef))
			Expect(actualRefRemainder).To(Equal(expectedRefRemainder))

		})
	})
})
