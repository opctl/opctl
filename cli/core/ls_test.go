package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"os"
	"path/filepath"
)

var _ = Context("listOpsInCollection", func() {
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
					writer:    os.Stdout,
				}

				/* act */
				objectUnderTest.ListOpsInCollection("")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("vos.Getwd doesn't error", func() {
			It("should call bundle.GetCollection w/ expected args", func() {
				/* arrange */
				fakeBundle := new(bundle.Fake)

				providedCollection := "dummyCollection"
				wdReturnedFromVos := "dummyWorkDir"

				fakeVos := new(vos.Fake)
				fakeVos.GetwdReturns(wdReturnedFromVos, nil)
				expectedPath := filepath.Join(wdReturnedFromVos, providedCollection)

				objectUnderTest := _core{
					bundle: fakeBundle,
					vos:    fakeVos,
					writer: os.Stdout,
				}

				/* act */
				objectUnderTest.ListOpsInCollection(providedCollection)

				/* assert */

				Expect(fakeBundle.GetCollectionArgsForCall(0)).Should(Equal(expectedPath))
			})
			Context("bundle.GetCollection errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeBundle := new(bundle.Fake)
					expectedError := errors.New("dummyError")
					fakeBundle.GetCollectionReturns(model.CollectionView{}, expectedError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						bundle:    fakeBundle,
						cliExiter: fakeCliExiter,
						vos:       new(vos.Fake),
						writer:    os.Stdout,
					}

					/* act */
					objectUnderTest.ListOpsInCollection("dummyCollection")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
				})
			})
		})
	})
})
