package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/virtual-go/vos"
	"path"
)

var _ = Context("pkgSetDescription", func() {
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
				objectUnderTest.PkgSetDescription("", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		It("should call pkg.SetDescription w/ expected args", func() {
			/* arrange */
			fakePkg := new(pkg.Fake)

			providedPkgRef := "dummyPkgRef"
			wdReturnedFromVos := "dummyWorkDir"

			fakeVos := new(vos.Fake)
			fakeVos.GetwdReturns(wdReturnedFromVos, nil)

			expectedReq := pkg.SetDescriptionReq{
				Path:        path.Join(wdReturnedFromVos, ".opspec", providedPkgRef),
				Description: "dummyPkgDescription",
			}

			objectUnderTest := _core{
				pkg: fakePkg,
				vos: fakeVos,
			}

			/* act */
			objectUnderTest.PkgSetDescription(expectedReq.Description, providedPkgRef)

			/* assert */

			Expect(fakePkg.SetDescriptionArgsForCall(0)).Should(Equal(expectedReq))
		})
		Context("pkg.SetDescription errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)
				expectedError := errors.New("dummyError")
				fakePkg.SetDescriptionReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:       fakePkg,
					cliExiter: fakeCliExiter,
					vos:       new(vos.Fake),
				}

				/* act */
				objectUnderTest.PkgSetDescription("", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
