package input

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

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

			expectedError := fmt.Errorf("unable to bind to '%v'; '%v' not a defined input", providedName, providedName)

			/* act */
			_, actualError := Interpret(
				providedName,
				"dummyValue",
				nil,
				map[string]*model.Value{},
				"dummyScratchDir",
			)

			/* assert */
			Expect(actualError).To(Equal(expectedError))
		})
	})
	Context("Implicit arg", func() {
		Context("Ref not in scope", func() {
			It("should return expected error", func() {
				/* arrange */
				providedName := "dummyName"

				expectedError := fmt.Errorf("unable to bind to '%v' via implicit ref; '%v' not in scope", providedName, providedName)

				/* act */
				_, actualError := Interpret(
					providedName,
					"",
					&model.Param{},
					map[string]*model.Value{},
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
	Context("Arg is string", func() {
		Context("Input is array", func() {
			It("should return expected results", func() {
				name := "name"
				providedScope := map[string]*model.Value{
					name: {Array: new([]interface{})},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)

				expectedResult, err := array.Interpret(providedScope, providedExpression)
				if nil != err {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.Param{Array: &model.ArrayParam{}},
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
				providedScope := map[string]*model.Value{
					name: {Boolean: new(bool)},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)

				expectedResult, err := boolean.Interpret(providedScope, providedExpression)
				if nil != err {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.Param{Boolean: &model.BooleanParam{}},
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
				tmpDir := os.TempDir()

				providedScope := map[string]*model.Value{
					name: {Link: &tmpDir},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)
				providedScratchDirPath := "dummyScratchDir"

				expectedResult, err := dir.Interpret(providedScope, providedExpression, providedScratchDirPath, true)
				if nil != err {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.Param{Dir: &model.DirParam{}},
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
				tmpFile, err := ioutil.TempFile("", "")
				if nil != err {
					panic(err)
				}

				tmpFilePath := tmpFile.Name()

				name := "name"

				providedScope := map[string]*model.Value{
					name: {Link: &tmpFilePath},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)
				providedScratchDirPath := "dummyScratchDir"

				expectedResult, err := file.Interpret(providedScope, providedExpression, providedScratchDirPath, true)
				if nil != err {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.Param{File: &model.FileParam{}},
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
				providedScope := map[string]*model.Value{
					name: &model.Value{Number: new(float64)},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)

				expectedResult, err := number.Interpret(providedScope, providedExpression)
				if nil != err {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.Param{Number: &model.NumberParam{}},
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
				providedScope := map[string]*model.Value{
					name: &model.Value{Object: new(map[string]interface{})},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)

				expectedResult, err := object.Interpret(providedScope, providedExpression)
				if nil != err {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.Param{Object: &model.ObjectParam{}},
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
				providedScope := map[string]*model.Value{
					name: &model.Value{String: new(string)},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)

				expectedResult, err := str.Interpret(providedScope, providedExpression)
				if nil != err {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.Param{String: &model.StringParam{}},
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
						&model.Param{Array: &model.ArrayParam{}},
						map[string]*model.Value{},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualError).To(Equal(errors.New("unable to bind 'name' to '$(name)'; error was: 'unable to interpret $(name) to array; error was unable to interpret 'name' as reference; 'name' not in scope'")))
				})
			})
			It("should return expected result", func() {
				name := "name"
				providedScope := map[string]*model.Value{
					name: &model.Value{Socket: new(string)},
				}
				providedExpression := fmt.Sprintf("$(%s)", name)

				expectedResult, err := reference.Interpret(providedExpression, providedScope, nil)
				if nil != err {
					panic(err)
				}

				/* act */
				actualResult, actualError := Interpret(
					name,
					providedExpression,
					&model.Param{Socket: &model.SocketParam{}},
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
