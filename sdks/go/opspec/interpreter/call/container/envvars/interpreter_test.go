package envvars

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	objectFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/object/fakes"
	strFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/str/fakes"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		It("should call objectInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			envVarName := "dummyEnvVar"
			providedScope := map[string]*model.Value{
				envVarName: nil,
			}
			providedSCGContainerCallEnvVars := "providedSCGContainerCallEnvVars"
			providedOpHandle := new(modelFakes.FakeDataHandle)

			fakeObjectInterpreter := new(objectFakes.FakeInterpreter)
			// err to trigger immediate return
			fakeObjectInterpreter.InterpretReturns(&model.Value{String: new(string)}, errors.New("dummyErr"))

			objectUnderTest := _interpreter{
				objectInterpreter: fakeObjectInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedSCGContainerCallEnvVars,
				providedOpHandle,
			)

			/* assert */
			actualScope,
				actualExpression,
				actualOpHandle := fakeObjectInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualExpression).To(Equal(providedSCGContainerCallEnvVars))
			Expect(actualOpHandle).To(Equal(actualOpHandle))
		})
		Context("objectInterpreter.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				providedSCGContainerCallEnvVars := "providedSCGContainerCallEnvVars"

				fakeObjectInterpreter := new(objectFakes.FakeInterpreter)

				interpretErr := errors.New("dummyError")
				fakeObjectInterpreter.InterpretReturns(nil, interpretErr)

				expectedErr := fmt.Errorf(
					"unable to interpret '%v' as envVars; error was %v",
					providedSCGContainerCallEnvVars,
					interpretErr,
				)

				objectUnderTest := _interpreter{
					objectInterpreter: fakeObjectInterpreter,
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCGContainerCallEnvVars,
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("objectInterpreter.Interpret doesn't err", func() {
			It("should call stringInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				envVarName := "dummyEnvVar"
				providedScope := map[string]*model.Value{
					envVarName: nil,
				}

				providedOpHandle := new(modelFakes.FakeDataHandle)

				fakeObjectInterpreter := new(objectFakes.FakeInterpreter)

				expectedEnvVarValue := "expectedEnvVarValue"
				interpretedValueObject := map[string]interface{}{
					envVarName: expectedEnvVarValue,
				}
				// err to trigger immediate return
				fakeObjectInterpreter.InterpretReturns(&model.Value{Object: &interpretedValueObject}, nil)

				fakeStrInterpreter := new(strFakes.FakeInterpreter)
				// err to trigger immediate return
				fakeStrInterpreter.InterpretReturns(nil, errors.New("dummyErr"))

				objectUnderTest := _interpreter{
					objectInterpreter: fakeObjectInterpreter,
					stringInterpreter: fakeStrInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedScope,
					map[string]interface{}{},
					providedOpHandle,
				)

				/* assert */
				actualScope,
					actualExpression,
					actualOpHandle := fakeStrInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualExpression).To(Equal(expectedEnvVarValue))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
			})
			Context("stringInterpreter.Interpret errs", func() {

				It("should return expected result", func() {
					/* arrange */
					envVarName := "dummyEnvVar"
					providedScope := map[string]*model.Value{
						envVarName: nil,
					}

					providedOpHandle := new(modelFakes.FakeDataHandle)

					fakeObjectInterpreter := new(objectFakes.FakeInterpreter)

					expectedEnvVarValue := "expectedEnvVarValue"
					interpretedValueObject := map[string]interface{}{
						envVarName: expectedEnvVarValue,
					}
					// err to trigger immediate return
					fakeObjectInterpreter.InterpretReturns(&model.Value{Object: &interpretedValueObject}, nil)

					err := errors.New("err")
					fakeStrInterpreter := new(strFakes.FakeInterpreter)
					// err to trigger immediate return
					fakeStrInterpreter.InterpretReturns(nil, err)

					expectedErr := fmt.Errorf(
						"unable to interpret %+v as value of env var '%v'; error was %v",
						expectedEnvVarValue,
						envVarName,
						err,
					)

					objectUnderTest := _interpreter{
						objectInterpreter: fakeObjectInterpreter,
						stringInterpreter: fakeStrInterpreter,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						providedScope,
						map[string]interface{}{},
						providedOpHandle,
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("stringInterpreter.Interpret doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					envVarName := "dummyEnvVar"
					providedScope := map[string]*model.Value{
						envVarName: nil,
					}

					providedOpHandle := new(modelFakes.FakeDataHandle)

					fakeObjectInterpreter := new(objectFakes.FakeInterpreter)

					interpretedValueObject := map[string]interface{}{
						envVarName: "envVarValue",
					}
					// err to trigger immediate return
					fakeObjectInterpreter.InterpretReturns(&model.Value{Object: &interpretedValueObject}, nil)

					expectedEnvVarValueString := "expectedEnvVarValueString"
					expectedEnvVarValue := &model.Value{String: &expectedEnvVarValueString}
					expectedResult := map[string]string{
						envVarName: expectedEnvVarValueString,
					}

					fakeStrInterpreter := new(strFakes.FakeInterpreter)
					// err to trigger immediate return
					fakeStrInterpreter.InterpretReturns(expectedEnvVarValue, nil)

					objectUnderTest := _interpreter{
						objectInterpreter: fakeObjectInterpreter,
						stringInterpreter: fakeStrInterpreter,
					}

					/* act */
					actualResult, _ := objectUnderTest.Interpret(
						providedScope,
						map[string]interface{}{},
						providedOpHandle,
					)

					/* assert */
					Expect(actualResult).To(Equal(expectedResult))
				})
			})
		})
	})
})
