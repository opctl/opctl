package inputs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("inputs.Interpret doesn't error", func() {
		It("should return expected result", func() {
			/* arrange */
			providedArgName := "argName"

			providedInputArgs := map[string]interface{}{
				providedArgName: "",
			}

			providedParams := map[string]*model.Param{
				providedArgName: &model.Param{
					String: &model.StringParam{},
				},
			}

			expectedInputs := map[string]*model.Value{
				providedArgName: &model.Value{
					String: new(string),
				},
			}

			/* act */
			actualResult, actualErr := Interpret(
				providedInputArgs,
				providedParams,
				"dummyOpPath",
				map[string]*model.Value{
					providedArgName: &model.Value{
						String: new(string),
					},
				},
				"dummyOpScratchDir",
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(actualResult).To(Equal(expectedInputs))
		})
	})
})
