package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
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
					bundle:    new(bundle.Fake),
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
			It("should call bundle.CreateCollection w/ expected args", func() {
				/* arrange */
				fakeBundle := new(bundle.Fake)

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
					bundle: fakeBundle,
					vos:    fakeVos,
				}

				/* act */
				objectUnderTest.CreateCollection(expectedReq.Description, expectedReq.Name)

				/* assert */
				Expect(fakeBundle.CreateCollectionArgsForCall(0)).Should(Equal(expectedReq))
			})
			Context("bundle.CreateCollection errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeBundle := new(bundle.Fake)
					expectedError := errors.New("dummyError")
					fakeBundle.CreateCollectionReturns(expectedError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						bundle:    fakeBundle,
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
