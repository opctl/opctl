package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

var _ = Describe("setOpDescription", func() {

	fakeWorkDirPathGetter := new(fakeWorkDirPathGetter)
	workDirPath := ""
	fakeWorkDirPathGetter.GetReturns(workDirPath)

	Context("Execute", func() {
		It("should invoke bundle.SetOpDescription with expected args", func() {
			/* arrange */
			fakeBundle := new(bundle.FakeBundle)

			providedCollection := "dummyCollection"
			providedName := "dummyOpName"

			expectedReq := model.SetOpDescriptionReq{
				PathToOp:    path.Join(workDirPath, providedCollection, providedName),
				Description: "dummyOpDescription",
			}

			objectUnderTest := _core{
				bundle:            fakeBundle,
				workDirPathGetter: fakeWorkDirPathGetter,
			}

			/* act */
			objectUnderTest.SetOpDescription(providedCollection, expectedReq.Description, providedName)

			/* assert */

			Expect(fakeBundle.SetOpDescriptionArgsForCall(0)).Should(Equal(expectedReq))
		})
		It("should return errors from bundle.SetOpDescription", func() {
			/* arrange */
			fakeBundle := new(bundle.FakeBundle)
			expectedError := errors.New("dummyError")
			fakeBundle.SetOpDescriptionReturns(expectedError)

			fakeExiter := new(fakeExiter)

			objectUnderTest := _core{
				bundle:            fakeBundle,
				exiter:            fakeExiter,
				workDirPathGetter: fakeWorkDirPathGetter,
			}

			/* act */
			objectUnderTest.SetOpDescription("", "", "")

			/* assert */
			Expect(fakeExiter.ExitArgsForCall(0)).
				Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))
		})
	})
})
