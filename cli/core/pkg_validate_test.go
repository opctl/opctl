package core

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/vos"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
	"path"
)

var _ = Context("pkgValidate", func() {
	Context("Execute", func() {
		Context("vos.Getwd errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeVOS := new(vos.Fake)
				expectedError := errors.New("dummyError")
				fakeVOS.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:       new(pkg.Fake),
					cliExiter: fakeCliExiter,
					os:        fakeVOS,
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
			wdReturnedFromVOS := "dummyWorkDir"

			fakeVOS := new(vos.Fake)
			fakeVOS.GetwdReturns(wdReturnedFromVOS, nil)

			expectedPkgRef := path.Join(wdReturnedFromVOS, ".opspec", providedPkgRef)

			objectUnderTest := _core{
				pkg:       fakePkg,
				cliExiter: new(cliexiter.Fake),
				os:        fakeVOS,
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
					os:        new(vos.Fake),
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
					os:        new(vos.Fake),
				}

				/* act */
				objectUnderTest.PkgValidate(providedPkgRef)

				/* assert */

				Expect(fakeCliExiter.ExitArgsForCall(0)).To(Equal(expectedExitReq))
			})
		})
	})
})
