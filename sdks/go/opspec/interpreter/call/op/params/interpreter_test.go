package params

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	It("should return expected result", func() {
		/* arrange */
		defaultStrParamName := "defaultStrParamName"
		defaultStrParamValue := "defaultStrParamValue"
		numberParamName := "numberParamName"
		numberParamValue := 2.2

		providedParams := map[string]*model.Param{
			defaultStrParamName: &model.Param{
				String: &model.StringParam{
					Default: &defaultStrParamValue,
				},
			},
			numberParamName: &model.Param{
				Number: &model.NumberParam{},
			},
		}

		providedOpPath := "dummyOpPath"

		expectedOutputs := map[string]*model.Value{
			defaultStrParamName: &model.Value{
				String: &defaultStrParamValue,
			},
			numberParamName: &model.Value{
				Number: &numberParamValue,
			},
		}

		/* act */
		actualOutputs, actualErr := Interpret(
			map[string]*model.Value{
				defaultStrParamName: &model.Value{
					String: &defaultStrParamValue,
				},
				numberParamName: &model.Value{
					Number: &numberParamValue,
				},
			},
			providedParams,
			providedOpPath,
			"dummyOpScratchDir",
		)

		/* assert */

		Expect(actualErr).To(BeNil())
		Expect(actualOutputs).To(Equal(expectedOutputs))

	})
})
