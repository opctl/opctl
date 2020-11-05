package opfile

import (
	"errors"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Unmarshal", func() {
	Context("Validate returns errors", func() {
		It("should return the expected error", func() {
			/* arrange */
			/* act */
			_, actualError := Unmarshal([]byte("&"))

			/* assert */
			Expect(actualError).To(Equal(errors.New("\n-\n  Error(s):\n    - error converting YAML to JSON: yaml: did not find expected alphabetic or numeric character\n-")))
		})
	})
	Context("Validator.Validate doesn't return errors", func() {

		XIt("should return expected opFile", func() {

			/* arrange */
			paramDefault := "dummyDefault"
			dummyParams := map[string]*model.Param{
				"dummyName": {
					String: &model.StringParam{
						Constraints: map[string]interface{}{
							"MinLength": 0,
							"MaxLength": 1000,
							"Pattern":   "dummyPattern",
							"Format":    "dummyFormat",
							"Enum":      []interface{}{"dummyEnumItem1"},
						},
						Default:     &paramDefault,
						Description: "dummyDescription",
						IsSecret:    true,
					},
				},
			}

			expectedOpFile := &model.OpSpec{
				Description: "dummyDescription",
				Inputs:      dummyParams,
				Name:        "dummyName",
				Outputs:     dummyParams,
				Run: &model.CallSpec{
					Op: &model.OpCallSpec{
						Ref: "dummyOpRef",
					},
				},
				Version: "dummyVersion",
			}
			providedBytes, err := yaml.Marshal(expectedOpFile)
			if nil != err {
				panic(err.Error())
			}

			/* act */
			actualOpFile, _ := Unmarshal(providedBytes)

			/* assert */
			Expect(*actualOpFile).To(Equal(*expectedOpFile))

		})
	})
})
