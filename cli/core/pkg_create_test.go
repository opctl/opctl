package core

import (
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
	"path"
)

var _ = Context("core", func() {
	Context("PkgCreate", func() {
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
				objectUnderTest.PkgCreate("", "", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("ios.Getwd doesn't error", func() {
			It("should call pkg.Create w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)

				providedPath := "dummyPath"
				providedPkgName := "dummyPkgName"
				providedPkgDescription := "dummyPkgDescription"
				wdReturnedFromIOS := "dummyWorkDir"

				expectedPath := path.Join(wdReturnedFromIOS, providedPath, providedPkgName)
				expectedPkgName := providedPkgName
				expectedPkgDescription := providedPkgDescription

				fakeIOS := new(ios.Fake)
				fakeIOS.GetwdReturns(wdReturnedFromIOS, nil)

				objectUnderTest := _core{
					pkg: fakePkg,
					os:  fakeIOS,
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
						os:        new(ios.Fake),
					}

					/* act */
					objectUnderTest.PkgCreate("", "", "")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

				})
			})
		})
	})
})
