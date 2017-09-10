package envvars

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	stringPkg "github.com/opspec-io/sdk-golang/string"
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
		It("should call string.Interpret w/ expected args", func() {
			/* arrange */
			envVarName := "dummyEnvVar"
			providedScope := map[string]*model.Value{
				envVarName: nil,
			}
			providedPkgHandle := new(pkg.FakeHandle)

			fakeString := new(stringPkg.Fake)
			objectUnderTest := _EnvVars{
				string: fakeString,
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
				actualPkgHandle := fakeString.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualExpression).To(Equal(fmt.Sprintf("$(%v)", envVarName)))
			Expect(actualPkgHandle).To(Equal(actualPkgHandle))
		})
		Context("string.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				envVarName := "dummyEnvVar"

				fakeString := new(stringPkg.Fake)

				interpretErr := errors.New("dummyError")
				fakeString.InterpretReturns("", interpretErr)

				expectedErr := fmt.Errorf(
					"Unable to bind env var to '%v' via implicit ref; '%v' not in scope",
					envVarName,
					envVarName,
				)

				objectUnderTest := _EnvVars{
					string: fakeString,
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
		Context("string.Interpret doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				envVarName := "dummyEnvVar"

				fakeString := new(stringPkg.Fake)

				interpretedEnvVar := "dummyEnvVarValue"
				fakeString.InterpretReturns(interpretedEnvVar, nil)

				expectedEnvVars := map[string]string{
					envVarName: interpretedEnvVar,
				}

				objectUnderTest := _EnvVars{
					string: fakeString,
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
