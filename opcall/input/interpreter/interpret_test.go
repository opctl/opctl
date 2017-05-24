package interpreter

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall/input/validator"
	"path/filepath"
	"strconv"
)

var _ = Context("Interpreter", func() {
	Context("Interpret", func() {
		Context("name doesn't match a param", func() {
			It("should return expected error", func() {
				/* arrange */
				providedName := "dummyName"

				fakeValidator := new(validator.Fake)
				expectedError := fmt.Errorf("Unable to bind to '%v'. '%v' is not a defined input", providedName, providedName)

				objectUnderTest := _Interpreter{
					validator: fakeValidator,
				}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					providedName,
					"dummyValue",
					map[string]*model.Param{},
					map[string]*model.Data{},
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
					providedValue := ""
					providedParams := map[string]*model.Param{providedName: {}}

					fakeValidator := new(validator.Fake)
					expectedError := fmt.Errorf("Unable to bind to '%v' via implicit ref. '%v' is not in scope", providedName, providedName)

					objectUnderTest := _Interpreter{
						validator: fakeValidator,
					}

					/* act */
					_, actualError := objectUnderTest.Interpret(
						providedName,
						providedValue,
						providedParams,
						map[string]*model.Data{},
					)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("Ref in scope", func() {
				It("should call validate.Validate w/ expected args & return result", func() {
					/* arrange */
					providedName := "dummyName"
					providedValue := ""
					providedParams := map[string]*model.Param{providedName: {}}
					expectedValue := &model.Data{String: new(string)}
					providedScope := map[string]*model.Data{providedName: expectedValue}

					fakeValidator := new(validator.Fake)

					objectUnderTest := _Interpreter{
						validator: fakeValidator,
					}

					/* act */
					actualValue, actualError := objectUnderTest.Interpret(
						providedName,
						providedValue,
						providedParams,
						providedScope,
					)

					/* assert */
					Expect(actualValue).To(Equal(expectedValue))
					Expect(actualError).To(BeNil())
				})
			})
		})
		Context("Explicit arg", func() {
			Context("Ref not in scope", func() {
				It("should return expected error", func() {
					/* arrange */
					providedName := "dummyName"
					explicitRef := "dummyRef"
					providedValue := fmt.Sprintf("$(%v)", explicitRef)
					providedParams := map[string]*model.Param{providedName: {}}

					fakeValidator := new(validator.Fake)
					expectedError := fmt.Errorf("Unable to bind '%v' to '%v' via explicit ref. '%v' is not in scope", providedName, explicitRef, explicitRef)

					objectUnderTest := _Interpreter{
						validator: fakeValidator,
					}

					/* act */
					_, actualError := objectUnderTest.Interpret(
						providedName,
						providedValue,
						providedParams,
						map[string]*model.Data{},
					)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("Ref in scope", func() {
				It("should call validate.Validate w/ expected args & return result", func() {
					/* arrange */
					providedName := "dummyName"
					explicitRef := "dummyRef"
					providedValue := fmt.Sprintf("$(%v)", explicitRef)
					providedParams := map[string]*model.Param{providedName: {}}
					expectedValue := &model.Data{String: new(string)}
					providedScope := map[string]*model.Data{explicitRef: expectedValue}

					fakeValidator := new(validator.Fake)

					objectUnderTest := _Interpreter{
						validator: fakeValidator,
					}

					/* act */
					actualValue, actualError := objectUnderTest.Interpret(
						providedName,
						providedValue,
						providedParams,
						providedScope,
					)

					/* assert */
					Expect(actualValue).To(Equal(expectedValue))
					Expect(actualError).To(BeNil())
				})
			})
		})
		Context("Interpolated arg", func() {
			It("should call interpolater.Interpolate w/ expected args", func() {
				/* arrange */
				providedName := "dummyName"
				providedValue := "dummyValue"
				providedParams := map[string]*model.Param{providedName: {}}
				providedScope := map[string]*model.Data{"dummyScopeRef": {}}

				fakeInterpolater := new(interpolater.Fake)

				objectUnderTest := _Interpreter{
					interpolater: fakeInterpolater,
					validator:    new(validator.Fake),
				}

				/* act */
				objectUnderTest.Interpret(
					providedName,
					providedValue,
					providedParams,
					providedScope,
				)

				/* assert */
				actualTemplate, actualScope := fakeInterpolater.InterpolateArgsForCall(0)
				Expect(actualTemplate).To(Equal(providedValue))
				Expect(actualScope).To(Equal(providedScope))
			})
			Context("Input is string", func() {
				It("should call validate w/ expected args", func() {
					/* arrange */
					providedName := "dummyName"
					providedParams := map[string]*model.Param{providedName: {String: &model.StringParam{}}}
					expectedParam := providedParams[providedName]

					fakeInterpolater := new(interpolater.Fake)

					interpolatedValue := "dummyValue"
					fakeInterpolater.InterpolateReturns(interpolatedValue)

					fakeValidator := new(validator.Fake)

					objectUnderTest := _Interpreter{
						interpolater: fakeInterpolater,
						validator:    fakeValidator,
					}

					/* act */
					objectUnderTest.Interpret(
						providedName,
						"dummyValue",
						providedParams,
						map[string]*model.Data{},
					)

					/* assert */
					actualValue, actualParam := fakeValidator.ValidateArgsForCall(0)
					Expect(*actualValue).To(Equal(model.Data{String: &interpolatedValue}))
					Expect(actualParam).To(Equal(expectedParam))
				})
			})
			Context("Input is Dir", func() {
				It("should call validate w/ expected args", func() {
					/* arrange */
					providedName := "dummyName"
					providedParams := map[string]*model.Param{providedName: {Dir: &model.DirParam{}}}
					expectedParam := providedParams[providedName]

					fakeInterpolater := new(interpolater.Fake)

					interpolatedValue := "dummyValue"
					fakeInterpolater.InterpolateReturns(interpolatedValue)

					fakeValidator := new(validator.Fake)

					objectUnderTest := _Interpreter{
						interpolater: fakeInterpolater,
						validator:    fakeValidator,
					}

					/* act */
					objectUnderTest.Interpret(
						providedName,
						"dummyValue",
						providedParams,
						map[string]*model.Data{},
					)

					/* assert */
					actualValue, actualParam := fakeValidator.ValidateArgsForCall(0)
					Expect(*actualValue).To(Equal(model.Data{Dir: &interpolatedValue}))
					Expect(actualParam).To(Equal(expectedParam))
				})
				It("should root path", func() {
					/* arrange */
					providedName := "dummyName"
					providedParams := map[string]*model.Param{providedName: {Dir: &model.DirParam{}}}
					expectedParam := providedParams[providedName]

					fakeInterpolater := new(interpolater.Fake)

					expectedValue := fmt.Sprintf("%v%v", string(filepath.Separator), "dummyValue")
					interpolatedValue := fmt.Sprintf("..\\../%v../..\\", expectedValue)
					fakeInterpolater.InterpolateReturns(interpolatedValue)

					fakeValidator := new(validator.Fake)

					objectUnderTest := _Interpreter{
						interpolater: fakeInterpolater,
						validator:    fakeValidator,
					}

					/* act */
					objectUnderTest.Interpret(
						providedName,
						"dummyValue",
						providedParams,
						map[string]*model.Data{},
					)

					/* assert */
					actualValue, actualParam := fakeValidator.ValidateArgsForCall(0)
					Expect(*actualValue).To(Equal(model.Data{Dir: &expectedValue}))
					Expect(actualParam).To(Equal(expectedParam))
				})
			})
			Context("Input is Number", func() {
				It("should call validate w/ expected args", func() {
					/* arrange */
					providedName := "dummyName"
					providedParams := map[string]*model.Param{providedName: {Number: &model.NumberParam{}}}
					expectedParam := providedParams[providedName]

					fakeInterpolater := new(interpolater.Fake)

					interpolatedValue := "2.1"
					expectedValue, err := strconv.ParseFloat(interpolatedValue, 64)
					if nil != err {
						panic(err)
					}

					fakeInterpolater.InterpolateReturns(interpolatedValue)

					fakeValidator := new(validator.Fake)

					objectUnderTest := _Interpreter{
						interpolater: fakeInterpolater,
						validator:    fakeValidator,
					}

					/* act */
					objectUnderTest.Interpret(
						providedName,
						"dummyValue",
						providedParams,
						map[string]*model.Data{},
					)

					/* assert */
					actualValue, actualParam := fakeValidator.ValidateArgsForCall(0)
					Expect(*actualValue).To(Equal(model.Data{Number: &expectedValue}))
					Expect(actualParam).To(Equal(expectedParam))
				})
			})
			Context("Input is File", func() {
				It("should call validate w/ expected args", func() {
					/* arrange */
					providedName := "dummyName"
					providedParams := map[string]*model.Param{providedName: {File: &model.FileParam{}}}
					expectedParam := providedParams[providedName]

					fakeInterpolater := new(interpolater.Fake)

					interpolatedValue := "dummyValue"
					fakeInterpolater.InterpolateReturns(interpolatedValue)

					fakeValidator := new(validator.Fake)

					objectUnderTest := _Interpreter{
						interpolater: fakeInterpolater,
						validator:    fakeValidator,
					}

					/* act */
					objectUnderTest.Interpret(
						providedName,
						"dummyValue",
						providedParams,
						map[string]*model.Data{},
					)

					/* assert */
					actualValue, actualParam := fakeValidator.ValidateArgsForCall(0)
					Expect(*actualValue).To(Equal(model.Data{File: &interpolatedValue}))
					Expect(actualParam).To(Equal(expectedParam))
				})
				It("should root path", func() {
					/* arrange */
					providedName := "dummyName"
					providedParams := map[string]*model.Param{providedName: {File: &model.FileParam{}}}
					expectedParam := providedParams[providedName]

					fakeInterpolater := new(interpolater.Fake)

					expectedValue := fmt.Sprintf("%v%v", string(filepath.Separator), "dummyValue")
					interpolatedValue := fmt.Sprintf("..\\../%v../..\\", expectedValue)
					fakeInterpolater.InterpolateReturns(interpolatedValue)

					fakeValidator := new(validator.Fake)

					objectUnderTest := _Interpreter{
						interpolater: fakeInterpolater,
						validator:    fakeValidator,
					}

					/* act */
					objectUnderTest.Interpret(
						providedName,
						"dummyValue",
						providedParams,
						map[string]*model.Data{},
					)

					/* assert */
					actualValue, actualParam := fakeValidator.ValidateArgsForCall(0)
					Expect(*actualValue).To(Equal(model.Data{File: &expectedValue}))
					Expect(actualParam).To(Equal(expectedParam))
				})
			})
			Context("Input is Socket", func() {
				It("should return expected error", func() {
					/* arrange */
					providedName := "dummyName"
					providedParams := map[string]*model.Param{providedName: {Socket: &model.SocketParam{}}}

					fakeInterpolater := new(interpolater.Fake)

					interpolatedValue := "dummyValue"
					fakeInterpolater.InterpolateReturns(interpolatedValue)

					expectedError := fmt.Errorf("Unable to bind '%v' to '%v'; sockets must be passed by reference", providedName, interpolatedValue)

					fakeValidator := new(validator.Fake)

					objectUnderTest := _Interpreter{
						interpolater: fakeInterpolater,
						validator:    fakeValidator,
					}

					/* act */
					_, actualError := objectUnderTest.Interpret(
						providedName,
						"dummyValue",
						providedParams,
						map[string]*model.Data{},
					)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
		})
		Context("Validate errors", func() {
			It("should return error", func() {
				/* arrange */
				providedName := "dummyName"
				providedValue := "dummyValue"
				providedParams := map[string]*model.Param{providedName: {String: &model.StringParam{}}}

				fakeInterpolater := new(interpolater.Fake)

				interpolatedValue := "dummyValue"
				fakeInterpolater.InterpolateReturns(interpolatedValue)

				fakeValidator := new(validator.Fake)
				validationErrors := []error{errors.New("dummyError")}
				fakeValidator.ValidateReturns(validationErrors)

				expectedError := fmt.Errorf(`
-
  validation of the following input failed:

  Name: %v
  Value: %v
  Error(s):
    - %v

-`, providedName, providedValue, validationErrors[0].Error())

				objectUnderTest := _Interpreter{
					interpolater: fakeInterpolater,
					validator:    fakeValidator,
				}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					providedName,
					"dummyValue",
					providedParams,
					map[string]*model.Data{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
	})
})
