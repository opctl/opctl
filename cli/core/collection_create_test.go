package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/pkg"
	"path/filepath"
)

var _ = Context("createCollection", func() {
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
				objectUnderTest.CreateCollection("", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("vos.Getwd doesn't error", func() {
			It("should call pkg.CreateCollection w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)

				expectedCollectionName := "dummyCollectionName"
				wdReturnedFromVos := "dummyWorkDir"

				fakeVos := new(vos.Fake)
				fakeVos.GetwdReturns(wdReturnedFromVos, nil)

				expectedReq := model.CreateCollectionReq{
					Path:        filepath.Join(wdReturnedFromVos, expectedCollectionName),
					Name:        expectedCollectionName,
					Description: "dummyCollectionDescription",
				}

				objectUnderTest := _core{
					pkg: fakePkg,
					vos: fakeVos,
				}

				/* act */
				objectUnderTest.CreateCollection(expectedReq.Description, expectedReq.Name)

				/* assert */
				Expect(fakePkg.CreateCollectionArgsForCall(0)).Should(Equal(expectedReq))
			})
			Context("pkg.CreateCollection errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakePkg := new(pkg.Fake)
					expectedError := errors.New("dummyError")
					fakePkg.CreateCollectionReturns(expectedError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						pkg:       fakePkg,
						cliExiter: fakeCliExiter,
						vos:       new(vos.Fake),
					}

					/* act */
					objectUnderTest.CreateCollection("", "")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
				})
			})
		})
	})
})
