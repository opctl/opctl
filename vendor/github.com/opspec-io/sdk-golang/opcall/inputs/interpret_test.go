package inputs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"path/filepath"
)

var _ = Context("inputs", func() {
	Context("Interpret", func() {
		Context("input param w/out arg", func() {
			Context("param is array", func() {
				Context("default exists", func() {
					It("should set input to default", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := []interface{}{}

						providedInputParams := map[string]*model.Param{
							providedInputName: {Array: &model.ArrayParam{Default: providedInputDefault}},
						}

						objectUnderTest := _Inputs{
							data: new(data.Fake),
						}

						expectedInputs := map[string]*model.Value{
							providedInputName: {Array: providedInputDefault},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]interface{}{},
							providedInputParams,
							new(pkg.FakeHandle),
							"dummyPkgRef",
							map[string]*model.Value{},
							"dummyOpScratchDir",
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
					})
					It("should call data.Validate w/ expected args", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := []interface{}{}

						providedInputParams := map[string]*model.Param{
							providedInputName: {Array: &model.ArrayParam{Default: providedInputDefault}},
						}

						expectedParam := providedInputParams[providedInputName]
						expectedInput := &model.Value{Array: providedInputDefault}

						fakeData := new(data.Fake)

						objectUnderTest := _Inputs{
							data: fakeData,
						}

						/* act */
						objectUnderTest.Interpret(
							map[string]interface{}{},
							providedInputParams,
							new(pkg.FakeHandle),
							"dummyPkgRef",
							map[string]*model.Value{},
							"dummyOpScratchDir",
						)

						/* assert */
						actualInput, actualParam := fakeData.ValidateArgsForCall(0)

						Expect(actualInput).To(Equal(expectedInput))
						Expect(actualParam).To(Equal(expectedParam))
					})
				})
			})
			Context("param is dir", func() {
				Context("default exists", func() {
					It("should set input to default", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "/pkgDirDefault"

						providedInputParams := map[string]*model.Param{
							providedInputName: {Dir: &model.DirParam{Default: &providedInputDefault}},
						}
						providedPkgRef := "dummyPkgRef"

						objectUnderTest := _Inputs{
							data: new(data.Fake),
						}

						expectedInputValue := filepath.Join(providedPkgRef, providedInputDefault)
						expectedInputs := map[string]*model.Value{
							providedInputName: {Dir: &expectedInputValue},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]interface{}{},
							providedInputParams,
							new(pkg.FakeHandle),
							providedPkgRef,
							map[string]*model.Value{},
							"dummyOpScratchDir",
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
					})
					It("should call data.Validate w/ expected args", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "/pkgDirDefault"
						providedPkgRef := "dummyPkgRef"

						providedInputParams := map[string]*model.Param{
							providedInputName: {Dir: &model.DirParam{Default: &providedInputDefault}},
						}

						expectedParam := providedInputParams[providedInputName]

						expectedInputValue := filepath.Join(providedPkgRef, providedInputDefault)
						expectedInput := &model.Value{Dir: &expectedInputValue}

						fakeData := new(data.Fake)

						objectUnderTest := _Inputs{
							data: fakeData,
						}

						/* act */
						objectUnderTest.Interpret(
							map[string]interface{}{},
							providedInputParams,
							new(pkg.FakeHandle),
							providedPkgRef,
							map[string]*model.Value{},
							"dummyOpScratchDir",
						)

						/* assert */
						actualInput, actualParam := fakeData.ValidateArgsForCall(0)

						Expect(actualInput).To(Equal(expectedInput))
						Expect(actualParam).To(Equal(expectedParam))
					})
				})
			})
			Context("param is file", func() {
				Context("default exists", func() {
					It("should set input to default", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "/pkgFileDefault"
						providedPkgRef := "dummyPkgRef"

						providedInputParams := map[string]*model.Param{
							providedInputName: {File: &model.FileParam{Default: &providedInputDefault}},
						}

						objectUnderTest := _Inputs{
							data: new(data.Fake),
						}

						expectedInputValue := filepath.Join(providedPkgRef, providedInputDefault)
						expectedInputs := map[string]*model.Value{
							providedInputName: {File: &expectedInputValue},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]interface{}{},
							providedInputParams,
							new(pkg.FakeHandle),
							providedPkgRef,
							map[string]*model.Value{},
							"dummyOpScratchDir",
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
					})
					It("should call data.Validate w/ expected args", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "/pkgFileDefault"
						providedPkgRef := "dummyPkgRef"

						providedInputParams := map[string]*model.Param{
							providedInputName: {File: &model.FileParam{Default: &providedInputDefault}},
						}

						expectedParam := providedInputParams[providedInputName]

						expectedInputValue := filepath.Join(providedPkgRef, providedInputDefault)
						expectedInput := &model.Value{File: &expectedInputValue}

						fakeData := new(data.Fake)

						objectUnderTest := _Inputs{
							data: fakeData,
						}

						/* act */
						objectUnderTest.Interpret(
							map[string]interface{}{},
							providedInputParams,
							new(pkg.FakeHandle),
							providedPkgRef,
							map[string]*model.Value{},
							"dummyOpScratchDir",
						)

						/* assert */
						actualInput, actualParam := fakeData.ValidateArgsForCall(0)

						Expect(actualInput).To(Equal(expectedInput))
						Expect(actualParam).To(Equal(expectedParam))
					})
				})
			})
			Context("param is number", func() {
				Context("default exists", func() {
					It("should set input to default", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := 2.2

						providedInputParams := map[string]*model.Param{
							providedInputName: {Number: &model.NumberParam{Default: &providedInputDefault}},
						}

						objectUnderTest := _Inputs{
							data: new(data.Fake),
						}

						expectedInputs := map[string]*model.Value{
							providedInputName: {Number: &providedInputDefault},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]interface{}{},
							providedInputParams,
							new(pkg.FakeHandle),
							"dummyPkgRef",
							map[string]*model.Value{},
							"dummyOpScratchDir",
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
					})
					It("should call data.Validate w/ expected args", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := 2.2

						providedInputParams := map[string]*model.Param{
							providedInputName: {Number: &model.NumberParam{Default: &providedInputDefault}},
						}

						expectedParam := providedInputParams[providedInputName]
						expectedInput := &model.Value{Number: &providedInputDefault}

						fakeData := new(data.Fake)

						objectUnderTest := _Inputs{
							data: fakeData,
						}

						/* act */
						objectUnderTest.Interpret(
							map[string]interface{}{},
							providedInputParams,
							new(pkg.FakeHandle),
							"dummyPkgRef",
							map[string]*model.Value{},
							"dummyOpScratchDir",
						)

						/* assert */
						actualInput, actualParam := fakeData.ValidateArgsForCall(0)

						Expect(actualInput).To(Equal(expectedInput))
						Expect(actualParam).To(Equal(expectedParam))
					})
				})
			})
			Context("param is object", func() {
				Context("default exists", func() {
					It("should set input to default", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := map[string]interface{}{}

						providedInputParams := map[string]*model.Param{
							providedInputName: {Object: &model.ObjectParam{Default: providedInputDefault}},
						}

						objectUnderTest := _Inputs{
							data: new(data.Fake),
						}

						expectedInputs := map[string]*model.Value{
							providedInputName: {Object: providedInputDefault},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]interface{}{},
							providedInputParams,
							new(pkg.FakeHandle),
							"dummyPkgRef",
							map[string]*model.Value{},
							"dummyOpScratchDir",
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
					})
					It("should call data.Validate w/ expected args", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := map[string]interface{}{}

						providedInputParams := map[string]*model.Param{
							providedInputName: {Object: &model.ObjectParam{Default: providedInputDefault}},
						}

						expectedParam := providedInputParams[providedInputName]
						expectedInput := &model.Value{Object: providedInputDefault}

						fakeData := new(data.Fake)

						objectUnderTest := _Inputs{
							data: fakeData,
						}

						/* act */
						objectUnderTest.Interpret(
							map[string]interface{}{},
							providedInputParams,
							new(pkg.FakeHandle),
							"dummyPkgRef",
							map[string]*model.Value{},
							"dummyOpScratchDir",
						)

						/* assert */
						actualInput, actualParam := fakeData.ValidateArgsForCall(0)

						Expect(actualInput).To(Equal(expectedInput))
						Expect(actualParam).To(Equal(expectedParam))
					})
				})
			})
			Context("param is string", func() {
				Context("default exists", func() {
					It("should set input to default", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "inputDefault"

						providedInputParams := map[string]*model.Param{
							providedInputName: {String: &model.StringParam{Default: &providedInputDefault}},
						}

						objectUnderTest := _Inputs{
							data: new(data.Fake),
						}

						expectedInputs := map[string]*model.Value{
							providedInputName: {String: &providedInputDefault},
						}

						/* act */
						actualInputs, _ := objectUnderTest.Interpret(
							map[string]interface{}{},
							providedInputParams,
							new(pkg.FakeHandle),
							"dummyPkgRef",
							map[string]*model.Value{},
							"dummyOpScratchDir",
						)

						/* assert */
						Expect(actualInputs).To(Equal(expectedInputs))
					})
					It("should call data.Validate w/ expected args", func() {
						/* arrange */
						providedInputName := "inputName"
						providedInputDefault := "inputDefault"

						providedInputParams := map[string]*model.Param{
							providedInputName: {String: &model.StringParam{Default: &providedInputDefault}},
						}

						expectedParam := providedInputParams[providedInputName]
						expectedInput := &model.Value{String: &providedInputDefault}

						fakeData := new(data.Fake)

						objectUnderTest := _Inputs{
							data: fakeData,
						}

						/* act */
						objectUnderTest.Interpret(
							map[string]interface{}{},
							providedInputParams,
							new(pkg.FakeHandle),
							"dummyPkgRef",
							map[string]*model.Value{},
							"dummyOpScratchDir",
						)

						/* assert */
						actualInput, actualParam := fakeData.ValidateArgsForCall(0)

						Expect(actualInput).To(Equal(expectedInput))
						Expect(actualParam).To(Equal(expectedParam))
					})
				})
			})
		})
		Context("input arg", func() {
			It("should call argInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedArgName := "argName"
				providedArgValue := "argValue"

				providedInputArgs := map[string]interface{}{
					providedArgName: providedArgValue,
				}

				providedInputParams := map[string]*model.Param{
					providedArgName: {String: &model.StringParam{}},
				}

				expectedParam := providedInputParams[providedArgName]

				providedParentPkgHandle := new(pkg.FakeHandle)

				providedScope := map[string]*model.Value{
					"scopeRef1Name": {},
				}

				providedOpScratchDir := "dummyOpScratchDir"

				fakeArgInterpreter := new(fakeArgInterpreter)

				objectUnderTest := _Inputs{
					argInterpreter: fakeArgInterpreter,
					data:           new(data.Fake),
				}

				/* act */
				objectUnderTest.Interpret(
					providedInputArgs,
					providedInputParams,
					providedParentPkgHandle,
					"dummyPkgRef",
					providedScope,
					providedOpScratchDir,
				)

				/* assert */
				actualArgName,
					actualArgValue,
					actualParam,
					actualParentPkgHandle,
					actualScope,
					actualOpScratchDir := fakeArgInterpreter.InterpretArgsForCall(0)

				Expect(actualArgName).To(Equal(providedArgName))
				Expect(actualArgValue).To(Equal(providedArgValue))
				Expect(actualParam).To(Equal(expectedParam))
				Expect(actualParentPkgHandle).To(Equal(providedParentPkgHandle))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualOpScratchDir).To(Equal(providedOpScratchDir))

			})
			Context("argInterpreter.Interpret doesn't error", func() {
				It("should call data.Validate w/ expected args", func() {
					/* arrange */
					providedArgName := "argName"

					providedInputArgs := map[string]interface{}{
						providedArgName: "",
					}

					providedInputParams := map[string]*model.Param{
						providedArgName: nil,
					}

					expectedParam := providedInputParams[providedArgName]

					expectedInput := &model.Value{String: new(string)}
					fakeArgInterpreter := new(fakeArgInterpreter)
					fakeArgInterpreter.InterpretReturns(expectedInput, nil)

					fakeData := new(data.Fake)

					objectUnderTest := _Inputs{
						argInterpreter: fakeArgInterpreter,
						data:           fakeData,
					}

					/* act */
					objectUnderTest.Interpret(
						providedInputArgs,
						providedInputParams,
						new(pkg.FakeHandle),
						"dummyPkgRef",
						map[string]*model.Value{},
						"dummyOpScratchDir",
					)

					/* assert */
					actualInput, actualParam := fakeData.ValidateArgsForCall(0)

					Expect(actualInput).To(Equal(expectedInput))
					Expect(actualParam).To(Equal(expectedParam))
				})
			})
		})
	})
})
