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

var _ = Describe("setOpDescription", func() {
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
				objectUnderTest.SetOpDescription("", "", "")

				/* assert */
				Expect(fakeExiter.ExitArgsForCall(0)).
					Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		It("should call bundle.SetOpDescription with expected args", func() {
			/* arrange */
			fakeBundle := new(bundle.FakeBundle)

			providedCollection := "dummyCollection"
			providedName := "dummyOpName"
			wdReturnedFromVos := "dummyWorkDir"

			fakeVos := new(vos.FakeVos)
			fakeVos.GetwdReturns(wdReturnedFromVos, nil)

			expectedReq := model.SetOpDescriptionReq{
				PathToOp:    filepath.Join(wdReturnedFromVos, providedCollection, providedName),
				Description: "dummyOpDescription",
			}

			objectUnderTest := _core{
				bundle: fakeBundle,
				vos:    fakeVos,
			}

			/* act */
			objectUnderTest.SetOpDescription(providedCollection, expectedReq.Description, providedName)

			/* assert */

			Expect(fakeBundle.SetOpDescriptionArgsForCall(0)).Should(Equal(expectedReq))
		})
		Context("bundle.SetOpDescription errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeBundle := new(bundle.FakeBundle)
				expectedError := errors.New("dummyError")
				fakeBundle.SetOpDescriptionReturns(expectedError)

				fakeExiter := new(fakeExiter)

				objectUnderTest := _core{
					bundle: fakeBundle,
					exiter: fakeExiter,
					vos:    new(vos.FakeVos),
				}

				/* act */
				objectUnderTest.SetOpDescription("", "", "")

				/* assert */
				Expect(fakeExiter.ExitArgsForCall(0)).
					Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
