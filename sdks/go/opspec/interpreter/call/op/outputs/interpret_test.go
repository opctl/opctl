package outputs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	It("should return expected result", func() {
		/* arrange */
		arrayValue := []interface{}{"item"}
		stringParamName := "stringParamName"

		providedArgs := map[string]*model.Value{
			stringParamName: &model.Value{Array: &arrayValue},
		}

		providedParams := map[string]*model.Param{
			stringParamName: &model.Param{
				String: &model.StringParam{},
			},
		}

		arrayValueAsString, err := coerce.ToString(providedArgs[stringParamName])
		if nil != err {
			panic(err)
		}

		expectedOutputs := map[string]*model.Value{
			stringParamName: arrayValueAsString,
		}

		/* act */
		actualOutputs, actualErr := Interpret(
			providedArgs,
			providedParams,
			"opPath",
			"opScratchDir",
		)

		/* assert */
		Expect(actualOutputs).To(Equal(expectedOutputs))
		Expect(actualErr).To(BeNil())

	})
})
