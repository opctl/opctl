package inputs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	It("should return expected result", func() {
		/* arrange */
		providedArgName := "argName"

		providedInputArgs := map[string]interface{}{
			providedArgName: "",
		}

		providedParams := map[string]*model.ParamSpec{
			providedArgName: {
				String: &model.StringParamSpec{},
			},
		}

		expectedInputs := map[string]*ipld.Node{
			providedArgName: {
				String: new(string),
			},
		}

		/* act */
		actualResult, actualErr := Interpret(
			providedInputArgs,
			providedParams,
			"dummyOpPath",
			map[string]*ipld.Node{
				providedArgName: {
					String: new(string),
				},
			},
			"dummyOpScratchDir",
		)

		/* assert */
		Expect(actualErr).To(BeNil())
		Expect(actualResult).To(Equal(expectedInputs))
	})
	Context("params.ApplyDefaults errors", func() {
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
				map[string]interface{}{},
				providedParams,
				"dummyOpPath",
				map[string]*ipld.Node{},
				"dummyOpScratchDir",
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret input defaults: unable to interpret $(nonExistent) to string: unable to interpret 'nonExistent' as reference: 'nonExistent' not in scope"))
			Expect(actualResult).To(BeNil())
		})
	})
})
