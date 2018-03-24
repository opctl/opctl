package envvars

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	stringPkg "github.com/opspec-io/sdk-golang/op/interpreter/string"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		Context("implicitly bound", func() {
			Context("name not in scope", func() {
				It("should return expected error", func() {
					/* arrange */
					envVarName := "dummyEnvVarName"
					providedSCGContainerCallEnvVars := map[string]interface{}{
						// implicitly bound
						envVarName: nil,
					}

					expectedErr := fmt.Errorf("unable to bind env var to '%v' via implicit ref; '%v' not in scope", envVarName, envVarName)

					objectUnderTest := _interpreter{}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						providedSCGContainerCallEnvVars,
						new(data.FakeHandle),
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
		})
		It("should call stringInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			envVarName := "dummyEnvVar"
			providedScope := map[string]*model.Value{
				envVarName: nil,
			}
			providedOpHandle := new(data.FakeHandle)

			fakeStringInterpreter := new(stringPkg.FakeInterpreter)
			fakeStringInterpreter.InterpretReturns(&model.Value{String: new(string)}, nil)

			objectUnderTest := _interpreter{
				stringInterpreter: fakeStringInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				map[string]interface{}{
					// implicitly bound to string
					envVarName: nil,
				},
				providedOpHandle,
			)

			/* assert */
			actualScope,
				actualExpression,
				actualOpHandle := fakeStringInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualExpression).To(Equal(fmt.Sprintf("$(%v)", envVarName)))
			Expect(actualOpHandle).To(Equal(actualOpHandle))
		})
		Context("stringInterpreter.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				envVarName := "dummyEnvVar"

				fakeStringInterpreter := new(stringPkg.FakeInterpreter)

				interpretErr := errors.New("dummyError")
				fakeStringInterpreter.InterpretReturns(nil, interpretErr)

				expectedErr := fmt.Errorf(
					"unable to bind env var to '%v' via implicit ref; '%v' not in scope",
					envVarName,
					envVarName,
				)

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStringInterpreter,
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{
						envVarName: nil,
					},
					map[string]interface{}{
						// implicitly bound to string
						envVarName: "",
					},
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("stringInterpreter.Interpret doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				envVarName := "dummyEnvVar"

				fakeStringInterpreter := new(stringPkg.FakeInterpreter)

				interpretedEnvVar := "dummyEnvVarValue"
				fakeStringInterpreter.InterpretReturns(&model.Value{String: &interpretedEnvVar}, nil)

				expectedEnvVars := map[string]string{
					envVarName: interpretedEnvVar,
				}

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStringInterpreter,
				}

				/* act */
				actualValue, _ := objectUnderTest.Interpret(
					map[string]*model.Value{
						envVarName: nil,
					},
					map[string]interface{}{
						// implicitly bound to string
						envVarName: "",
					},
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualValue).To(Equal(expectedEnvVars))

			})
		})
	})
})
