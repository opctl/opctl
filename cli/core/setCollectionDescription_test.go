package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path/filepath"
)

var _ = Describe("setCollectionDescription", func() {
	Context("Execute", func() {
		Context("vos.Getwd errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeVos := new(vos.FakeVos)
				expectedError := errors.New("dummyError")
				fakeVos.GetwdReturns("", expectedError)

				fakeExiter := new(fakeExiter)

				objectUnderTest := _core{
					bundle: new(bundle.FakeBundle),
					exiter: fakeExiter,
					vos:    fakeVos,
				}

				/* act */
				objectUnderTest.SetCollectionDescription("")

				/* assert */
				Expect(fakeExiter.ExitArgsForCall(0)).
					Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("vos.Getwd doesn't error", func() {
			It("should call bundle.SetCollectionDescription with expected args", func() {
				/* arrange */
				fakeBundle := new(bundle.FakeBundle)
				wdReturnedFromVos := "dummyWorkDir"

				fakeVos := new(vos.FakeVos)
				fakeVos.GetwdReturns(wdReturnedFromVos, nil)

				expectedReq := model.SetCollectionDescriptionReq{
					PathToCollection: filepath.Join(wdReturnedFromVos, ".opspec"),
					Description:      "dummyOpDescription",
				}

				objectUnderTest := _core{
					bundle: fakeBundle,
					vos:    fakeVos,
				}

				/* act */
				objectUnderTest.SetCollectionDescription(expectedReq.Description)

				/* assert */

				Expect(fakeBundle.SetCollectionDescriptionArgsForCall(0)).Should(Equal(expectedReq))
			})
			Context("bundle.SetCollectionDescription errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeBundle := new(bundle.FakeBundle)
					expectedError := errors.New("dummyError")
					fakeBundle.SetCollectionDescriptionReturns(expectedError)

					fakeExiter := new(fakeExiter)

					objectUnderTest := _core{
						bundle: fakeBundle,
						exiter: fakeExiter,
						vos:    new(vos.FakeVos),
					}

					/* act */
					objectUnderTest.SetCollectionDescription("")

					/* assert */
					Expect(fakeExiter.ExitArgsForCall(0)).
						Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))
				})
			})
		})
	})
})
