package envvars

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression"
	"github.com/pkg/errors"
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
		It("should call expression.EvalToString w/ expected args", func() {
			/* arrange */
			envVarName := "dummyEnvVar"
			providedScope := map[string]*model.Value{
				envVarName: nil,
			}
			providedOpHandle := new(data.FakeHandle)

			fakeExpression := new(expression.Fake)
			fakeExpression.EvalToStringReturns(&model.Value{String: new(string)}, nil)

			objectUnderTest := _interpreter{
				expression: fakeExpression,
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
				actualOpHandle := fakeExpression.EvalToStringArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualExpression).To(Equal(fmt.Sprintf("$(%v)", envVarName)))
			Expect(actualOpHandle).To(Equal(actualOpHandle))
		})
		Context("expression.EvalToString errs", func() {
			It("should return expected result", func() {
				/* arrange */
				envVarName := "dummyEnvVar"

				fakeExpression := new(expression.Fake)

				interpretErr := errors.New("dummyError")
				fakeExpression.EvalToStringReturns(nil, interpretErr)

				expectedErr := fmt.Errorf(
					"unable to bind env var to '%v' via implicit ref; '%v' not in scope",
					envVarName,
					envVarName,
				)

				objectUnderTest := _interpreter{
					expression: fakeExpression,
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
		Context("expression.EvalToString doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				envVarName := "dummyEnvVar"

				fakeExpression := new(expression.Fake)

				interpretedEnvVar := "dummyEnvVarValue"
				fakeExpression.EvalToStringReturns(&model.Value{String: &interpretedEnvVar}, nil)

				expectedEnvVars := map[string]string{
					envVarName: interpretedEnvVar,
				}

				objectUnderTest := _interpreter{
					expression: fakeExpression,
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
