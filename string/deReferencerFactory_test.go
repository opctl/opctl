package string

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("deReferencerFactory", func() {
	Context("New", func() {
		It("should return expected result", func() {
			/* arrange */
			providedScope := map[string]*model.Value{"dummyName": {}}

			expectedResult := _deReferencer{
				coercer: newCoercer(),
				scope:   providedScope,
			}

			objectUnderTest := _deReferencerFactory{}

			/* act */
			actualResult := objectUnderTest.New(providedScope)

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))

		})
	})
})
