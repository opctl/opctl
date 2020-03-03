package op

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	cliexiterFakes "github.com/opctl/opctl/cli/internal/cliexiter/fakes"
	. "github.com/opctl/opctl/sdks/go/opspec/fakes"
	"path/filepath"
)

var _ = Context("Creater", func() {
	Context("Create", func() {
		It("should call opCreator.Create w/ expected args", func() {
			/* arrange */
			fakeOpCreator := new(FakeCreator)

			providedPath := "dummyPath"
			providedPkgName := "dummyPkgName"
			providedPkgDescription := "dummyPkgDescription"

			expectedPath := filepath.Join(providedPath, providedPkgName)
			expectedPkgName := providedPkgName
			expectedPkgDescription := providedPkgDescription

			objectUnderTest := _creater{
				opCreator: fakeOpCreator,
			}

			/* act */
			objectUnderTest.Create(providedPath, providedPkgDescription, providedPkgName)

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
				fakeOpCreator := new(FakeCreator)
				expectedError := errors.New("dummyError")
				fakeOpCreator.CreateReturns(expectedError)

				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

				objectUnderTest := _creater{
					opCreator: fakeOpCreator,
					cliExiter: fakeCliExiter,
				}

				/* act */
				objectUnderTest.Create("", "", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

			})
		})
	})
})
