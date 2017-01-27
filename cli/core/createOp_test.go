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

var _ = Describe("createOp", func() {
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
				objectUnderTest.CreateOp("", "", "")

				/* assert */
				Expect(fakeExiter.ExitArgsForCall(0)).
					Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("vos.Getwd doesn't error", func() {
			It("should call bundle.CreateOp w/ expected args", func() {
				/* arrange */
				fakeBundle := new(bundle.FakeBundle)

				providedCollection := "dummyCollection"
				providedName := "dummyName"
				wdReturnedFromVos := "dummyWorkDir"

				fakeVos := new(vos.FakeVos)
				fakeVos.GetwdReturns(wdReturnedFromVos, nil)

				expectedReq := model.CreateOpReq{
					Path:        filepath.Join(wdReturnedFromVos, providedCollection, providedName),
					Name:        providedName,
					Description: "dummyOpDescription",
				}

				objectUnderTest := _core{
					bundle: fakeBundle,
					vos:    fakeVos,
				}

				/* act */
				objectUnderTest.CreateOp(providedCollection, expectedReq.Description, expectedReq.Name)

				/* assert */

				Expect(fakeBundle.CreateOpArgsForCall(0)).Should(Equal(expectedReq))
			})
			Context("bundle.CreateOp errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeBundle := new(bundle.FakeBundle)
					expectedError := errors.New("dummyError")
					fakeBundle.CreateOpReturns(expectedError)

					fakeExiter := new(fakeExiter)

					objectUnderTest := _core{
						bundle: fakeBundle,
						exiter: fakeExiter,
						vos:    new(vos.FakeVos),
					}

					/* act */
					objectUnderTest.CreateOp("", "", "")

					/* assert */
					Expect(fakeExiter.ExitArgsForCall(0)).
						Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))

				})
			})
		})
	})
})
