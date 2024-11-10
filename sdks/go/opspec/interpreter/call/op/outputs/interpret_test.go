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

		providedArgs := map[string]*ipld.Node{
			stringParamName: {Array: &arrayValue},
		}

		providedParams := map[string]*model.ParamSpec{
			stringParamName: {String: &model.StringParamSpec{}},
		}

		providedOpCallOutputs := map[string]string{}

		arrayValueAsString, err := coerce.ToString(providedArgs[stringParamName])
		if err != nil {
			panic(err)
		}

		expectedOutputs := map[string]*ipld.Node{
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
				map[string]*ipld.Node{},
				providedParams,
				map[string]string{},
				"opPath",
				"opScratchDir",
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret $(nonExistent) to string: unable to interpret 'nonExistent' as reference: 'nonExistent' not in scope"))
			Expect(actualResult).To(BeNil())
		})
	})
	Describe("deprecated output binding syntax", func() {
		It("should return expected result", func() {
			/* arrange */
			arrayValue := []interface{}{"item"}
			stringParamName := "stringParamName"

			providedArgs := map[string]*ipld.Node{
				stringParamName: {Array: &arrayValue},
			}

			providedParams := map[string]*model.ParamSpec{
				stringParamName: {String: &model.StringParamSpec{}},
			}

			providedOpCallOutputs := map[string]string{
				"myVar": stringParamName,
			}

			arrayValueAsString, err := coerce.ToString(providedArgs[stringParamName])
			if err != nil {
				panic(err)
			}

			expectedOutputs := map[string]*ipld.Node{
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
	})
	Describe("ensures expected outputs match actual outputs", func() {
		It("indicates what was expected when one output exists", func() {
			/* act */
			actualOutputs, actualErr := Interpret(
				map[string]*ipld.Node{},
				map[string]*model.ParamSpec{
					"bar": {String: &model.StringParamSpec{}},
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
				map[string]*ipld.Node{},
				map[string]*model.ParamSpec{
					"x": {String: &model.StringParamSpec{}},
					"y": {String: &model.StringParamSpec{}},
					"z": {String: &model.StringParamSpec{}},
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
				map[string]*ipld.Node{},
				map[string]*model.ParamSpec{
					"x": {String: &model.StringParamSpec{}},
					"y": {String: &model.StringParamSpec{}},
					"z": {String: &model.StringParamSpec{}},
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
				map[string]*ipld.Node{},
				map[string]*model.ParamSpec{
					"x": {String: &model.StringParamSpec{}},
					"y": {String: &model.StringParamSpec{}},
					"z": {String: &model.StringParamSpec{}},
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
				map[string]*ipld.Node{},
				map[string]*model.ParamSpec{},
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
