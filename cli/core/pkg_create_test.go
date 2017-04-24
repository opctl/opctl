package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/virtual-go/vos"
	"path"
)

var _ = Context("core", func() {
	Context("PkgCreate", func() {
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
					vos:       fakeVOS,
				}

				/* act */
				objectUnderTest.PkgCreate("", "", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("vos.Getwd doesn't error", func() {
			It("should call pkg.Create w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)

				providedPath := "dummyPath"
				providedPkgName := "dummyPkgName"
				providedPkgDescription := "dummyPkgDescription"
				wdReturnedFromVOS := "dummyWorkDir"

				expectedPath := path.Join(wdReturnedFromVOS, providedPath, providedPkgName)
				expectedPkgName := providedPkgName
				expectedPkgDescription := providedPkgDescription

				fakeVOS := new(vos.Fake)
				fakeVOS.GetwdReturns(wdReturnedFromVOS, nil)

				objectUnderTest := _core{
					pkg: fakePkg,
					vos: fakeVOS,
				}

				/* act */
				objectUnderTest.PkgCreate(providedPath, providedPkgDescription, providedPkgName)

				/* assert */
				actualPath,
					actualPkgName,
					actualPkgDescription := fakePkg.CreateArgsForCall(0)

				Expect(actualPath).To(Equal(expectedPath))
				Expect(actualPkgName).To(Equal(expectedPkgName))
				Expect(actualPkgDescription).To(Equal(expectedPkgDescription))
			})
			Context("pkg.Create errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakePkg := new(pkg.Fake)
					expectedError := errors.New("dummyError")
					fakePkg.CreateReturns(expectedError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						pkg:       fakePkg,
						cliExiter: fakeCliExiter,
						vos:       new(vos.Fake),
					}

					/* act */
					objectUnderTest.PkgCreate("", "", "")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

				})
			})
		})
	})
})
