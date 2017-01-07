package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

var _ = Describe("setCollectionDescription", func() {

	fakeWorkDirPathGetter := new(fakeWorkDirPathGetter)
	workDirPath := ""
	fakeWorkDirPathGetter.GetReturns(workDirPath)

	Context("Execute", func() {
		It("should invoke bundle.SetCollectionDescription with expected args", func() {
			/* arrange */
			fakeBundle := new(bundle.FakeBundle)

			expectedReq := model.SetCollectionDescriptionReq{
				PathToCollection: path.Join(workDirPath, ".opspec"),
				Description:      "dummyOpDescription",
			}

			objectUnderTest := _core{
				bundle:            fakeBundle,
				workDirPathGetter: fakeWorkDirPathGetter,
			}

			/* act */
			objectUnderTest.SetCollectionDescription(expectedReq.Description)

			/* assert */

			Expect(fakeBundle.SetCollectionDescriptionArgsForCall(0)).Should(Equal(expectedReq))
		})
		It("should return errors from bundle.SetCollectionDescription", func() {
			/* arrange */
			fakeBundle := new(bundle.FakeBundle)
			expectedError := errors.New("dummyError")
			fakeBundle.SetCollectionDescriptionReturns(expectedError)

			fakeExiter := new(fakeExiter)

			objectUnderTest := _core{
				bundle:            fakeBundle,
				exiter:            fakeExiter,
				workDirPathGetter: fakeWorkDirPathGetter,
			}

			/* act */
			objectUnderTest.SetCollectionDescription("")

			/* assert */
			Expect(fakeExiter.ExitArgsForCall(0)).
				Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))
		})
	})
})
