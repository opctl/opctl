package params

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Validate", func() {
	It("should return expected result", func() {
		/* arrange */

		providedValues := map[string]*model.Value{
			"expectedName1": &model.Value{
				String: new(string),
			},
		}

		providedParams := map[string]*model.Param{
			"expectedName1": &model.Param{
				String: &model.StringParam{
					Constraints: map[string]interface{}{
						"minLength": 10,
					},
				},
			},
		}

		/* act */
		actualErr := Validate(
			providedValues,
			providedParams,
		)

		/* assert */
		Expect(actualErr).To(Equal(errors.New("\n-\n  validation error(s):\n\n    - expectedName1: String length must be greater than or equal to 10\n\n-")))
	})
})
