package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

var _ = Describe("createCollection", func() {

	fakeWorkDirPathGetter := new(fakeWorkDirPathGetter)
	workDirPath := ""
	fakeWorkDirPathGetter.GetReturns(workDirPath)

	Context("Execute", func() {
		It("should invoke bundle.CreateCollection with expected args", func() {
			/* arrange */
			fakeBundle := new(bundle.FakeBundle)

			expectedCollectionName := "dummyCollectionName"

			expectedReq := model.CreateCollectionReq{
				Path:        path.Join(workDirPath, expectedCollectionName),
				Name:        expectedCollectionName,
				Description: "dummyCollectionDescription",
			}

			objectUnderTest := _core{
				bundle:            fakeBundle,
				workDirPathGetter: fakeWorkDirPathGetter,
			}

			/* act */
			objectUnderTest.CreateCollection(expectedReq.Description, expectedReq.Name)

			/* assert */
			Expect(fakeBundle.CreateCollectionArgsForCall(0)).Should(Equal(expectedReq))
		})
		It("should return error from bundle.CreateCollection", func() {
			/* arrange */
			fakeBundle := new(bundle.FakeBundle)
			expectedError := errors.New("dummyError")
			fakeBundle.CreateCollectionReturns(expectedError)

			fakeExiter := new(fakeExiter)

			objectUnderTest := _core{
				bundle:            fakeBundle,
				exiter:            fakeExiter,
				workDirPathGetter: fakeWorkDirPathGetter,
			}

			/* act */
			objectUnderTest.CreateCollection("", "")

			/* assert */
			Expect(fakeExiter.ExitArgsForCall(0)).
				Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))
		})
	})
})
