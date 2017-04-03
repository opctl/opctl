package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/virtual-go/vos"
	"path/filepath"
)

var _ = Context("create", func() {
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
					vos:       fakeVOS,
				}

				/* act */
				objectUnderTest.Create("", "", "")

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
				providedName := "dummyName"
				wdReturnedFromVOS := "dummyWorkDir"

				fakeVOS := new(vos.Fake)
				fakeVOS.GetwdReturns(wdReturnedFromVOS, nil)

				expectedReq := pkg.CreateReq{
					Path:        filepath.Join(wdReturnedFromVOS, providedPath, providedName),
					Name:        providedName,
					Description: "dummyPkgDescription",
				}

				objectUnderTest := _core{
					pkg: fakePkg,
					vos: fakeVOS,
				}

				/* act */
				objectUnderTest.Create(providedPath, expectedReq.Description, expectedReq.Name)

				/* assert */

				Expect(fakePkg.CreateArgsForCall(0)).Should(Equal(expectedReq))
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
					objectUnderTest.Create("", "", "")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

				})
			})
		})
	})
})
