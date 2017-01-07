package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"os"
	"path"
)

var _ = Describe("listOpsInCollection", func() {

	fakeWorkDirPathGetter := new(fakeWorkDirPathGetter)
	workDirPath := ""
	fakeWorkDirPathGetter.GetReturns(workDirPath)

	Context("Execute", func() {
		It("should invoke bundle.GetCollection with expected args", func() {
			/* arrange */
			fakeBundle := new(bundle.FakeBundle)

			providedCollection := "dummyCollection"
			expectedPath := path.Join(workDirPath, providedCollection)

			objectUnderTest := _core{
				bundle:            fakeBundle,
				workDirPathGetter: fakeWorkDirPathGetter,
				writer:            os.Stdout,
			}

			/* act */
			objectUnderTest.ListOpsInCollection(providedCollection)

			/* assert */

			Expect(fakeBundle.GetCollectionArgsForCall(0)).Should(Equal(expectedPath))
		})
		It("should return errors from bundle.GetCollection", func() {
			/* arrange */
			fakeBundle := new(bundle.FakeBundle)
			expectedError := errors.New("dummyError")
			fakeBundle.GetCollectionReturns(model.CollectionView{}, expectedError)

			fakeExiter := new(fakeExiter)

			objectUnderTest := _core{
				bundle:            fakeBundle,
				exiter:            fakeExiter,
				workDirPathGetter: fakeWorkDirPathGetter,
				writer:            os.Stdout,
			}

			/* act */
			objectUnderTest.ListOpsInCollection("dummyCollection")

			/* assert */
			Expect(fakeExiter.ExitArgsForCall(0)).
				Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))
		})
	})
})
