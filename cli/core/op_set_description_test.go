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

var _ = Context("setOpDescription", func() {
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
				objectUnderTest.SetOpDescription("", "", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		It("should call pkg.SetOpDescription w/ expected args", func() {
			/* arrange */
			fakePkg := new(pkg.Fake)

			providedCollection := "dummyCollection"
			providedName := "dummyOpName"
			wdReturnedFromVos := "dummyWorkDir"

			fakeVos := new(vos.Fake)
			fakeVos.GetwdReturns(wdReturnedFromVos, nil)

			expectedReq := model.SetOpDescriptionReq{
				PathToOp:    filepath.Join(wdReturnedFromVos, providedCollection, providedName),
				Description: "dummyOpDescription",
			}

			objectUnderTest := _core{
				pkg: fakePkg,
				vos: fakeVos,
			}

			/* act */
			objectUnderTest.SetOpDescription(providedCollection, expectedReq.Description, providedName)

			/* assert */

			Expect(fakePkg.SetOpDescriptionArgsForCall(0)).Should(Equal(expectedReq))
		})
		Context("pkg.SetOpDescription errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)
				expectedError := errors.New("dummyError")
				fakePkg.SetOpDescriptionReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:       fakePkg,
					cliExiter: fakeCliExiter,
					vos:       new(vos.Fake),
				}

				/* act */
				objectUnderTest.SetOpDescription("", "", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
