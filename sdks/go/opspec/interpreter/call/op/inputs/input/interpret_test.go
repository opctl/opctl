package input

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/array"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/boolean"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/dir"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/file"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/number"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/object"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

var _ = Context("Interpret", func() {
	Context("param nil", func() {
		It("should return expected error", func() {
			/* arrange */
			providedName := "dummyName"

			expectedError := fmt.Sprintf("unable to bind to '%v': '%v' not a defined input", providedName, providedName)

			/* act */
			_, actualError := Interpret(
				providedName,
				"dummyValue",
				nil,
				map[string]*ipld.Node{},
				"dummyScratchDir",
			)

			/* assert */
			Expect(actualError).To(MatchError(expectedError))
		})
	})
	Context("Implicit arg", func() {
		Context("Ref not in scope", func() {
			It("should return expected error", func() {
				/* arrange */
				providedName := "dummyName"

				expectedError := fmt.Sprintf("unable to bind to '%v' via implicit ref: '%v' not in scope", providedName, providedName)

				/* act */
				_, actualError := Interpret(
					providedName,
					"",
					&model.ParamSpec{},
					map[string]*ipld.Node{},
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualError).To(MatchError(expectedError))
			})
		})
	})
	Context("Arg is string", func() {
		Context("Input is array", func() {
			It("should return expected results", func() {
				name := "name"
				providedScope := map[string]*ipld.Node{
					name: {Array: new([]interface{})},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)

				expectedResult, err := array.Interpret(providedScope, providedExpression)
				if err != nil {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.ParamSpec{Array: &model.ArrayParamSpec{}},
					providedScope,
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualError).To(BeNil())
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
		Context("Input is boolean", func() {
			It("should return expected results", func() {
				name := "name"
				providedScope := map[string]*ipld.Node{
					name: {Boolean: new(bool)},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)

				expectedResult, err := boolean.Interpret(providedScope, providedExpression)
				if err != nil {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.ParamSpec{Boolean: &model.BooleanParamSpec{}},
					providedScope,
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualError).To(BeNil())
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
		Context("Input is dir", func() {
			It("should return expected results", func() {
				name := "name"
				providedScope := map[string]*ipld.Node{
					name: {Dir: new(string)},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)
				providedScratchDirPath := "dummyScratchDir"

				expectedResult, err := dir.Interpret(providedScope, providedExpression, providedScratchDirPath, true)
				if err != nil {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.ParamSpec{Dir: &model.DirParamSpec{}},
					providedScope,
					providedScratchDirPath,
				)

				/* assert */
				Expect(actualError).To(BeNil())
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
		Context("Input is file", func() {
			It("should return expected results", func() {
				name := "name"
				providedScope := map[string]*ipld.Node{
					name: {File: new(string)},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)
				providedScratchDirPath := "dummyScratchDir"

				expectedResult, err := file.Interpret(providedScope, providedExpression, providedScratchDirPath, true)
				if err != nil {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.ParamSpec{File: &model.FileParamSpec{}},
					providedScope,
					providedScratchDirPath,
				)

				/* assert */
				Expect(actualError).To(BeNil())
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
		Context("Input is number", func() {
			It("should return expected results", func() {
				name := "name"
				providedScope := map[string]*ipld.Node{
					name: {Number: new(float64)},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)

				expectedResult, err := number.Interpret(providedScope, providedExpression)
				if err != nil {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.ParamSpec{Number: &model.NumberParamSpec{}},
					providedScope,
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualError).To(BeNil())
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
		Context("Input is object", func() {
			It("should return expected result", func() {
				name := "name"
				providedScope := map[string]*ipld.Node{
					name: {Object: new(map[string]interface{})},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)

				expectedResult, err := object.Interpret(providedScope, providedExpression)
				if err != nil {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.ParamSpec{Object: &model.ObjectParamSpec{}},
					providedScope,
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualError).To(BeNil())
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
		Context("Input is string", func() {
			It("should return expected result", func() {
				name := "name"
				providedScope := map[string]*ipld.Node{
					name: {String: new(string)},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)

				expectedResult, err := str.Interpret(providedScope, providedExpression)
				if err != nil {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.ParamSpec{String: &model.StringParamSpec{}},
					providedScope,
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualError).To(BeNil())
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
		Context("Input is socket", func() {
			Context("reference.Interpret errs", func() {
				It("should return expected error", func() {
					name := "name"

					/* act */
					_, actualError := Interpret(
						name,
						fmt.Sprintf("$(%s)", name),
						&model.ParamSpec{Array: &model.ArrayParamSpec{}},
						map[string]*ipld.Node{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualError).To(MatchError("unable to bind 'name' to '$(name)': unable to interpret $(name) to array: unable to interpret 'name' as reference: 'name' not in scope"))
				})
			})
			It("should return expected result", func() {
				name := "name"
				providedScope := map[string]*ipld.Node{
					name: {Socket: new(string)},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)

				expectedResult, err := reference.Interpret(providedExpression, providedScope, nil)
				if err != nil {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.ParamSpec{Socket: &model.SocketParamSpec{}},
					providedScope,
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualError).To(BeNil())
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
	})
})
