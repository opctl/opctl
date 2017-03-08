package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/managepackages"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path/filepath"
)

var _ = Context("createPackage", func() {
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
				objectUnderTest.CreatePackage("", "", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("vos.Getwd doesn't error", func() {
			It("should call managepackages.CreatePackage w/ expected args", func() {
				/* arrange */
				fakeManagePackages := new(managepackages.Fake)

				providedPath := "dummyPath"
				providedName := "dummyName"
				wdReturnedFromVos := "dummyWorkDir"

				fakeVos := new(vos.Fake)
				fakeVos.GetwdReturns(wdReturnedFromVos, nil)

				expectedReq := model.CreatePackageReq{
					Path:        filepath.Join(wdReturnedFromVos, providedPath, providedName),
					Name:        providedName,
					Description: "dummyPkgDescription",
				}

				objectUnderTest := _core{
					managePackages: fakeManagePackages,
					vos:            fakeVos,
				}

				/* act */
				objectUnderTest.CreatePackage(providedPath, expectedReq.Description, expectedReq.Name)

				/* assert */

				Expect(fakeManagePackages.CreatePackageArgsForCall(0)).Should(Equal(expectedReq))
			})
			Context("managepackages.CreatePackage errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeManagePackages := new(managepackages.Fake)
					expectedError := errors.New("dummyError")
					fakeManagePackages.CreatePackageReturns(expectedError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						managePackages: fakeManagePackages,
						cliExiter:      fakeCliExiter,
						vos:            new(vos.Fake),
					}

					/* act */
					objectUnderTest.CreatePackage("", "", "")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

				})
			})
		})
	})
})
