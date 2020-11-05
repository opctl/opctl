package bracketed

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("CoerceToArrayOrObject", func() {
	Context("Is coercible to array", func() {
		It("should call coerce.ToArray w/ expected args", func() {
			/* arrange */

			providedData := model.Value{Array: new([]interface{})}

			/* act */
			actualResult, actualErr := CoerceToArrayOrObject(
				&providedData,
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(*actualResult).To(Equal(model.Value{Array: new([]interface{})}))
		})
	})
})
