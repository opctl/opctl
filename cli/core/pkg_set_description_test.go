package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/managepackages"
	"github.com/opspec-io/sdk-golang/model"
	"path"
)

var _ = Context("pkgSetDescription", func() {
	Context("Execute", func() {
		Context("vos.Getwd errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeVos := new(vos.Fake)
				expectedError := errors.New("dummyError")
				fakeVos.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					managePackages: new(managepackages.Fake),
					cliExiter:      fakeCliExiter,
					vos:            fakeVos,
				}

				/* act */
				objectUnderTest.PkgSetDescription("", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		It("should call managepackages.SetPackageDescription w/ expected args", func() {
			/* arrange */
			fakeManagePackages := new(managepackages.Fake)

			providedPkgRef := "dummyPkgRef"
			wdReturnedFromVos := "dummyWorkDir"

			fakeVos := new(vos.Fake)
			fakeVos.GetwdReturns(wdReturnedFromVos, nil)

			expectedReq := model.SetPackageDescriptionReq{
				PathToOp:    path.Join(wdReturnedFromVos, ".opspec", providedPkgRef),
				Description: "dummyPkgDescription",
			}

			objectUnderTest := _core{
				managePackages: fakeManagePackages,
				vos:            fakeVos,
			}

			/* act */
			objectUnderTest.PkgSetDescription(expectedReq.Description, providedPkgRef)

			/* assert */

			Expect(fakeManagePackages.SetPackageDescriptionArgsForCall(0)).Should(Equal(expectedReq))
		})
		Context("managepackages.SetPackageDescription errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeManagePackages := new(managepackages.Fake)
				expectedError := errors.New("dummyError")
				fakeManagePackages.SetPackageDescriptionReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					managePackages: fakeManagePackages,
					cliExiter:      fakeCliExiter,
					vos:            new(vos.Fake),
				}

				/* act */
				objectUnderTest.PkgSetDescription("", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
