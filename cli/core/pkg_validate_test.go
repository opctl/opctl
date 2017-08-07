package core

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
)

var _ = Context("pkgValidate", func() {
	Context("Execute", func() {
		It("should call pkgResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"

			fakePkgResolver := new(fakePkgResolver)
			fakePkgResolver.ResolveReturns(nil)

			fakePkg := new(pkg.Fake)
			// error to trigger immediate return
			fakePkg.ValidateReturns([]error{errors.New("dummyError")})

			objectUnderTest := _core{
				pkgResolver: fakePkgResolver,
				pkg:         fakePkg,
				cliExiter:   new(cliexiter.Fake),
				os:          new(ios.Fake),
			}

			/* act */
			objectUnderTest.PkgValidate(providedPkgRef)

			/* assert */
			actualPkgRef, actualPullCreds := fakePkgResolver.ResolveArgsForCall(0)
			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualPullCreds).To(BeNil())
		})
		It("should call pkg.Validate w/ expected args", func() {
			/* arrange */
			fakePkg := new(pkg.Fake)

			fakePkgResolver := new(fakePkgResolver)

			fakePkgHandle := new(pkg.FakeHandle)
			fakePkgResolver.ResolveReturns(fakePkgHandle)

			objectUnderTest := _core{
				pkg:         fakePkg,
				pkgResolver: fakePkgResolver,
				cliExiter:   new(cliexiter.Fake),
				os:          new(ios.Fake),
			}

			/* act */
			objectUnderTest.PkgValidate("dummyPkgRef")

			/* assert */
			Expect(fakePkg.ValidateArgsForCall(0)).To(Equal(fakePkgHandle))
		})
		Context("pkg.Validate returns errors", func() {
			It("should call cliExiter.Exit w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)

				fakePkgResolver := new(fakePkgResolver)

				fakePkgHandle := new(pkg.FakeHandle)
				fakePkgResolver.ResolveReturns(fakePkgHandle)

				errsReturnedFromValidate := []error{errors.New("dummyError")}
				fakePkg.ValidateReturns(errsReturnedFromValidate)

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
					pkg:         fakePkg,
					pkgResolver: fakePkgResolver,
					cliExiter:   fakeCliExiter,
					os:          new(ios.Fake),
				}

				/* act */
				objectUnderTest.PkgValidate("dummyPkgRef")

				/* assert */

				Expect(fakeCliExiter.ExitArgsForCall(0)).To(Equal(expectedExitReq))
			})
		})
		Context("pkg.Validate doesn't return errors", func() {
			It("should call cliExiter.Exit w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)

				fakePkgHandle := new(pkg.FakeHandle)
				pkgRef := "dummyPkgRef"
				fakePkgHandle.RefReturns(pkgRef)

				fakePkgResolver := new(fakePkgResolver)
				fakePkgResolver.ResolveReturns(fakePkgHandle)

				expectedExitReq := cliexiter.ExitReq{
					Message: fmt.Sprintf("%v is valid", pkgRef),
				}

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:         fakePkg,
					pkgResolver: fakePkgResolver,
					cliExiter:   fakeCliExiter,
					os:          new(ios.Fake),
				}

				/* act */
				objectUnderTest.PkgValidate("dummyPkgRef")

				/* assert */

				Expect(fakeCliExiter.ExitArgsForCall(0)).To(Equal(expectedExitReq))
			})
		})
	})
})
