package uniquestring

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("uniqueStringFactory", func() {
	Context("Construct()", func() {
		It("should not return the same string in 100000 iterations", func() {

			/* arrange */
			objectUnderTest := NewUniqueStringFactory()
			stringsReturnedFromConstruct := map[string]bool{}

			/* act/assert */
			for i := 0; i < 100000; i++ {

				uniqueString := objectUnderTest.Construct()

				if _, ok := stringsReturnedFromConstruct[uniqueString]; ok {
					Fail("same string returned twice")
				} else {
					stringsReturnedFromConstruct[uniqueString] = true
				}

			}

		})
	})
})
