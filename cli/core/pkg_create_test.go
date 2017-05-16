package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
	"path/filepath"
)

var _ = Context("core", func() {
	Context("PkgCreate", func() {
		It("should call pkg.Create w/ expected args", func() {
			/* arrange */
			fakePkg := new(pkg.Fake)

			providedPath := "dummyPath"
			providedPkgName := "dummyPkgName"
			providedPkgDescription := "dummyPkgDescription"

			expectedPath := filepath.Join(providedPath, providedPkgName)
			expectedPkgName := providedPkgName
			expectedPkgDescription := providedPkgDescription

			objectUnderTest := _core{
				pkg: fakePkg,
			}

			/* act */
			objectUnderTest.PkgCreate(providedPath, providedPkgDescription, providedPkgName)

			/* assert */
			actualPath,
				actualPkgName,
				actualPkgDescription := fakePkg.CreateArgsForCall(0)

			Expect(actualPath).To(Equal(expectedPath))
			Expect(actualPkgName).To(Equal(expectedPkgName))
			Expect(actualPkgDescription).To(Equal(expectedPkgDescription))
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
				}

				/* act */
				objectUnderTest.PkgCreate("", "", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

			})
		})
	})
})
