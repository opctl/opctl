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
				fakeVos := new(vos.Fake)
				expectedError := errors.New("dummyError")
				fakeVos.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:       new(pkg.Fake),
					cliExiter: fakeCliExiter,
					vos:       fakeVos,
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
				wdReturnedFromVos := "dummyWorkDir"

				fakeVos := new(vos.Fake)
				fakeVos.GetwdReturns(wdReturnedFromVos, nil)

				expectedReq := pkg.CreateReq{
					Path:        filepath.Join(wdReturnedFromVos, providedPath, providedName),
					Name:        providedName,
					Description: "dummyPkgDescription",
				}

				objectUnderTest := _core{
					pkg: fakePkg,
					vos: fakeVos,
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
