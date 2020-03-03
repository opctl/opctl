package op

import (
	"context"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	cliexiterFakes "github.com/opctl/opctl/cli/internal/cliexiter/fakes"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	. "github.com/opctl/opctl/sdks/go/model/fakes"
	. "github.com/opctl/opctl/sdks/go/opspec/fakes"
)

var _ = Context("Validater", func() {
	Context("Validate", func() {
		It("should call dataResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"

			fakeDataResolver := new(dataresolver.Fake)
			fakeDataResolver.ResolveReturns(nil)

			fakeOpValidator := new(FakeValidator)
			// error to trigger immediate return
			fakeOpValidator.ValidateReturns([]error{errors.New("dummyError")})

			objectUnderTest := _validater{
				cliExiter:    new(cliexiterFakes.FakeCliExiter),
				dataResolver: fakeDataResolver,
				opValidator:  fakeOpValidator,
			}

			/* act */
			objectUnderTest.Validate(
				context.Background(),
				providedPkgRef,
			)

			/* assert */
			actualPkgRef,
				actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)

			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualPullCreds).To(BeNil())
		})
		It("should call pkg.Validate w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()

			fakeOpValidator := new(FakeValidator)

			fakeDataResolver := new(dataresolver.Fake)
			fakeOpHandle := new(FakeDataHandle)
			fakeDataResolver.ResolveReturns(fakeOpHandle)

			objectUnderTest := _validater{
				cliExiter:    new(cliexiterFakes.FakeCliExiter),
				dataResolver: fakeDataResolver,
				opValidator:  fakeOpValidator,
			}

			/* act */
			objectUnderTest.Validate(
				providedCtx,
				"dummyPkgRef",
			)

			/* assert */
			actualCtx,
				actualOpHandle := fakeOpValidator.ValidateArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualOpHandle).To(Equal(fakeOpHandle))
		})
		Context("pkg.Validate returns errors", func() {
			It("should call cliExiter.Exit w/ expected args", func() {
				/* arrange */
				fakeOpValidator := new(FakeValidator)

				fakeDataResolver := new(dataresolver.Fake)

				fakeOpHandle := new(FakeDataHandle)
				fakeDataResolver.ResolveReturns(fakeOpHandle)

				errsReturnedFromValidate := []error{errors.New("dummyError")}
				fakeOpValidator.ValidateReturns(errsReturnedFromValidate)

				expectedExitReq := cliexiter.ExitReq{
					Message: fmt.Sprintf(`
-
  Error(s):
    - %v
-`, errsReturnedFromValidate[0]),
					Code: 1,
				}

				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

				objectUnderTest := _validater{
					cliExiter:    fakeCliExiter,
					dataResolver: fakeDataResolver,
					opValidator:  fakeOpValidator,
				}

				/* act */
				objectUnderTest.Validate(
					context.Background(),
					"dummyPkgRef",
				)

				/* assert */

				Expect(fakeCliExiter.ExitArgsForCall(0)).To(Equal(expectedExitReq))
			})
		})
		Context("pkg.Validate doesn't return errors", func() {
			It("should call cliExiter.Exit w/ expected args", func() {
				/* arrange */
				fakeOpValidator := new(FakeValidator)

				fakeOpHandle := new(FakeDataHandle)
				opRef := "dummyPkgRef"
				fakeOpHandle.RefReturns(opRef)

				fakeDataResolver := new(dataresolver.Fake)
				fakeDataResolver.ResolveReturns(fakeOpHandle)

				expectedExitReq := cliexiter.ExitReq{
					Message: fmt.Sprintf("%v is valid", opRef),
				}

				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

				objectUnderTest := _validater{
					cliExiter:    fakeCliExiter,
					dataResolver: fakeDataResolver,
					opValidator:  fakeOpValidator,
				}

				/* act */
				objectUnderTest.Validate(
					context.Background(),
					"dummyPkgRef",
				)

				/* assert */

				Expect(fakeCliExiter.ExitArgsForCall(0)).To(Equal(expectedExitReq))
			})
		})
	})
})
