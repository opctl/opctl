package envvars

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/expression"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/pkg/errors"
)

var _ = Context("EnvVars", func() {
	Context("Interpret", func() {
		Context("implicitly bound", func() {
			Context("name not in scope", func() {
				It("should return expected error", func() {
					/* arrange */
					envVarName := "dummyEnvVarName"
					providedSCGContainerCallEnvVars := map[string]string{
						// implicitly bound
						envVarName: "",
					}

					expectedErr := fmt.Errorf("Unable to bind env var to '%v' via implicit ref; '%v' not in scope", envVarName, envVarName)

					objectUnderTest := _EnvVars{}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						providedSCGContainerCallEnvVars,
						new(pkg.FakeHandle),
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
			providedPkgHandle := new(pkg.FakeHandle)

			fakeExpression := new(expression.Fake)
			objectUnderTest := _EnvVars{
				expression: fakeExpression,
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				map[string]string{
					// implicitly bound to string
					envVarName: "",
				},
				providedPkgHandle,
			)

			/* assert */
			actualScope,
				actualExpression,
				actualPkgHandle := fakeExpression.EvalToStringArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualExpression).To(Equal(fmt.Sprintf("$(%v)", envVarName)))
			Expect(actualPkgHandle).To(Equal(actualPkgHandle))
		})
		Context("expression.EvalToString errs", func() {
			It("should return expected result", func() {
				/* arrange */
				envVarName := "dummyEnvVar"

				fakeExpression := new(expression.Fake)

				interpretErr := errors.New("dummyError")
				fakeExpression.EvalToStringReturns("", interpretErr)

				expectedErr := fmt.Errorf(
					"Unable to bind env var to '%v' via implicit ref; '%v' not in scope",
					envVarName,
					envVarName,
				)

				objectUnderTest := _EnvVars{
					expression: fakeExpression,
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{
						envVarName: nil,
					},
					map[string]string{
						// implicitly bound to string
						envVarName: "",
					},
					new(pkg.FakeHandle),
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
				fakeExpression.EvalToStringReturns(interpretedEnvVar, nil)

				expectedEnvVars := map[string]string{
					envVarName: interpretedEnvVar,
				}

				objectUnderTest := _EnvVars{
					expression: fakeExpression,
				}

				/* act */
				actualValue, _ := objectUnderTest.Interpret(
					map[string]*model.Value{
						envVarName: nil,
					},
					map[string]string{
						// implicitly bound to string
						envVarName: "",
					},
					new(pkg.FakeHandle),
				)

				/* assert */
				Expect(actualValue).To(Equal(expectedEnvVars))

			})
		})
	})
})
