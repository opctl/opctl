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
			stringParamName: {Array: &arrayValue},
		}

		providedParams := map[string]*model.Param{
			stringParamName: {String: &model.StringParam{}},
		}

		providedOpCallOutputs := map[string]string{}

		arrayValueAsString, err := coerce.ToString(providedArgs[stringParamName])
		if err != nil {
			panic(err)
		}

		expectedOutputs := map[string]*model.Value{
			stringParamName: arrayValueAsString,
		}

		/* act */
		actualOutputs, actualErr := Interpret(
			providedArgs,
			providedParams,
			providedOpCallOutputs,
			"opPath",
			"opScratchDir",
		)

		/* assert */
		Expect(actualOutputs).To(Equal(expectedOutputs))
		Expect(actualErr).To(BeNil())
	})
	Describe("ensures expected outputs match actual outputs", func() {
		It("detects inverted naming", func() {
			/* act */
			actualOutputs, actualErr := Interpret(
				map[string]*model.Value{},
				map[string]*model.Param{
					"bar": {String: &model.StringParam{}},
				},
				map[string]string{
					"foo": "$(bar)",
				},
				"opPath",
				"opScratchDir",
			)

			/* assert */
			Expect(actualOutputs).To(BeNil())
			Expect(actualErr).To(MatchError("unknown output 'foo', did you mean to use `bar: $(foo)`?"))
		})
		It("indicates what was expected when one output exists", func() {
			/* act */
			actualOutputs, actualErr := Interpret(
				map[string]*model.Value{},
				map[string]*model.Param{
					"bar": {String: &model.StringParam{}},
				},
				map[string]string{
					"foo": "",
				},
				"opPath",
				"opScratchDir",
			)

			/* assert */
			Expect(actualOutputs).To(BeNil())
			Expect(actualErr).To(MatchError("unknown output 'foo', expected 'bar'"))
		})
		It("indicates what was expected when one remaining output exists", func() {
			/* act */
			actualOutputs, actualErr := Interpret(
				map[string]*model.Value{},
				map[string]*model.Param{
					"x": {String: &model.StringParam{}},
					"y": {String: &model.StringParam{}},
					"z": {String: &model.StringParam{}},
				},
				map[string]string{
					"x": "",
					"y": "",
					"a": "",
				},
				"opPath",
				"opScratchDir",
			)

			/* assert */
			Expect(actualOutputs).To(BeNil())
			Expect(actualErr).To(MatchError("unknown output 'a', expected 'z'"))
		})
		It("indicates what was expected when multiple outputs exist", func() {
			/* act */
			actualOutputs, actualErr := Interpret(
				map[string]*model.Value{},
				map[string]*model.Param{
					"x": {String: &model.StringParam{}},
					"y": {String: &model.StringParam{}},
					"z": {String: &model.StringParam{}},
				},
				map[string]string{
					"a": "",
				},
				"opPath",
				"opScratchDir",
			)

			/* assert */
			Expect(actualOutputs).To(BeNil())
			Expect(actualErr).To(MatchError("unknown output 'a', expected one of [x, y, z]"))
		})
		It("indicates what was expected when multiple remaining outputs exist", func() {
			/* act */
			actualOutputs, actualErr := Interpret(
				map[string]*model.Value{},
				map[string]*model.Param{
					"x": {String: &model.StringParam{}},
					"y": {String: &model.StringParam{}},
					"z": {String: &model.StringParam{}},
				},
				map[string]string{
					"x": "",
					"a": "",
				},
				"opPath",
				"opScratchDir",
			)

			/* assert */
			Expect(actualOutputs).To(BeNil())
			Expect(actualErr).To(MatchError("unknown output 'a', expected one of [y, z]"))
		})
		It("provides information when it doesn't know what was expected", func() {
			/* act */
			actualOutputs, actualErr := Interpret(
				map[string]*model.Value{},
				map[string]*model.Param{},
				map[string]string{
					"a": "",
				},
				"opPath",
				"opScratchDir",
			)

			/* assert */
			Expect(actualOutputs).To(BeNil())
			Expect(actualErr).To(MatchError("unknown output 'a'"))
		})
	})
})
