package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdk/go/opspec"
	"github.com/opctl/opctl/util/cliexiter"
	"path/filepath"
)

var _ = Context("core", func() {
	Context("OpCreate", func() {
		It("should call opCreator.Create w/ expected args", func() {
			/* arrange */
			fakeOpCreator := new(op.FakeCreator)

			providedPath := "dummyPath"
			providedPkgName := "dummyPkgName"
			providedPkgDescription := "dummyPkgDescription"

			expectedPath := filepath.Join(providedPath, providedPkgName)
			expectedPkgName := providedPkgName
			expectedPkgDescription := providedPkgDescription

			objectUnderTest := _core{
				opCreator: fakeOpCreator,
			}

			/* act */
			objectUnderTest.OpCreate(providedPath, providedPkgDescription, providedPkgName)

			/* assert */
			actualPath,
				actualPkgName,
				actualPkgDescription := fakeOpCreator.CreateArgsForCall(0)

			Expect(actualPath).To(Equal(expectedPath))
			Expect(actualPkgName).To(Equal(expectedPkgName))
			Expect(actualPkgDescription).To(Equal(expectedPkgDescription))
		})
		Context("opCreator.Create errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeOpCreator := new(op.FakeCreator)
				expectedError := errors.New("dummyError")
				fakeOpCreator.CreateReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					opCreator: fakeOpCreator,
					cliExiter: fakeCliExiter,
				}

				/* act */
				objectUnderTest.OpCreate("", "", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

			})
		})
	})
})
