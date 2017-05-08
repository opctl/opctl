package core

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
	"path"
)

var _ = Context("pkgValidate", func() {
	Context("Execute", func() {
		Context("ios.Getwd errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeIOS := new(ios.Fake)
				expectedError := errors.New("dummyError")
				fakeIOS.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:       new(pkg.Fake),
					cliExiter: fakeCliExiter,
					os:        fakeIOS,
				}

				/* act */
				objectUnderTest.PkgValidate("")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		It("should call pkg.Validate w/ expected args", func() {
			/* arrange */
			fakePkg := new(pkg.Fake)

			providedPkgRef := "dummyPkgRef"
			wdReturnedFromIOS := "dummyWorkDir"

			fakeIOS := new(ios.Fake)
			fakeIOS.GetwdReturns(wdReturnedFromIOS, nil)

			expectedPkgRef := path.Join(wdReturnedFromIOS, ".opspec", providedPkgRef)

			objectUnderTest := _core{
				pkg:       fakePkg,
				cliExiter: new(cliexiter.Fake),
				os:        fakeIOS,
			}

			/* act */
			objectUnderTest.PkgValidate(providedPkgRef)

			/* assert */

			Expect(fakePkg.ValidateArgsForCall(0)).To(Equal(expectedPkgRef))
		})
		Context("pkg.Validate returns errors", func() {
			It("should call cliExiter.Exit w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)
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
					pkg:       fakePkg,
					cliExiter: fakeCliExiter,
					os:        new(ios.Fake),
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
				providedPkgRef := "dummyPkgRef"

				fakePkg := new(pkg.Fake)
				errsReturnedFromValidate := []error{}
				fakePkg.ValidateReturns(errsReturnedFromValidate)

				expectedExitReq := cliexiter.ExitReq{
					Message: fmt.Sprintf("%v is valid", providedPkgRef),
				}

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:       fakePkg,
					cliExiter: fakeCliExiter,
					os:        new(ios.Fake),
				}

				/* act */
				objectUnderTest.PkgValidate(providedPkgRef)

				/* assert */

				Expect(fakeCliExiter.ExitArgsForCall(0)).To(Equal(expectedExitReq))
			})
		})
	})
})
