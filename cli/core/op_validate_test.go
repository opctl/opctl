package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/op"
)

var _ = Context("opValidate", func() {
	Context("Execute", func() {
		It("should call dataResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"

			fakeDataResolver := new(fakeDataResolver)
			fakeDataResolver.ResolveReturns(nil)

			fakeOpValidator := new(op.FakeValidator)
			// error to trigger immediate return
			fakeOpValidator.ValidateReturns([]error{errors.New("dummyError")})

			objectUnderTest := _core{
				cliExiter:    new(cliexiter.Fake),
				dataResolver: fakeDataResolver,
				opValidator:  fakeOpValidator,
				os:           new(ios.Fake),
			}

			/* act */
			objectUnderTest.OpValidate(
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

			fakeOpValidator := new(op.FakeValidator)

			fakeDataResolver := new(fakeDataResolver)
			fakeOpHandle := new(data.FakeHandle)
			fakeDataResolver.ResolveReturns(fakeOpHandle)

			objectUnderTest := _core{
				cliExiter:    new(cliexiter.Fake),
				dataResolver: fakeDataResolver,
				opValidator:  fakeOpValidator,
				os:           new(ios.Fake),
			}

			/* act */
			objectUnderTest.OpValidate(
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
				fakeOpValidator := new(op.FakeValidator)

				fakeDataResolver := new(fakeDataResolver)

				fakeOpHandle := new(data.FakeHandle)
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

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					cliExiter:    fakeCliExiter,
					dataResolver: fakeDataResolver,
					opValidator:  fakeOpValidator,
					os:           new(ios.Fake),
				}

				/* act */
				objectUnderTest.OpValidate(
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
				fakeOpValidator := new(op.FakeValidator)

				fakeOpHandle := new(data.FakeHandle)
				opRef := "dummyPkgRef"
				fakeOpHandle.RefReturns(opRef)

				fakeDataResolver := new(fakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle)

				expectedExitReq := cliexiter.ExitReq{
					Message: fmt.Sprintf("%v is valid", opRef),
				}

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					cliExiter:    fakeCliExiter,
					dataResolver: fakeDataResolver,
					opValidator:  fakeOpValidator,
					os:           new(ios.Fake),
				}

				/* act */
				objectUnderTest.OpValidate(
					context.Background(),
					"dummyPkgRef",
				)

				/* assert */

				Expect(fakeCliExiter.ExitArgsForCall(0)).To(Equal(expectedExitReq))
			})
		})
	})
})
