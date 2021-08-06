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

		providedParams := map[string]*model.ParamSpec{
			defaultStrParamName: {
				String: &model.StringParamSpec{
					Default: defaultStrParamValue,
				},
			},
			numberParamName: {
				Number: &model.NumberParamSpec{},
			},
		}

		providedOpPath := "dummyOpPath"

		expectedOutputs := map[string]*model.Value{
			defaultStrParamName: {
				String: &defaultStrParamValue,
			},
			numberParamName: {
				Number: &numberParamValue,
			},
		}

		/* act */
		actualOutputs, actualErr := Interpret(
			map[string]*model.Value{
				defaultStrParamName: {
					String: &defaultStrParamValue,
				},
				numberParamName: {
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
	Describe("params.ApplyDefaults errors", func() {
		It("should return expected result", func() {
			/* arrange */
			providedParams := map[string]*model.ParamSpec{
				"param0": {
					String: &model.StringParamSpec{
						Default: "$(nonExistent)",
					},
				},
			}

			/* act */
			actualResult, actualErr := Interpret(
				map[string]*model.Value{},
				providedParams,
				"opPath",
				"opScratchDir",
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret $(nonExistent) to string: unable to interpret 'nonExistent' as reference: 'nonExistent' not in scope"))
			Expect(actualResult).To(BeNil())
		})
	})
})
