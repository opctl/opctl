package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("Get", func() {
	It("should call getter.Get w/ expected inputs", func() {
		/* arrange */
		providedBasePath := "dummyBasePath"
		providedPkgRef := "dummyPkgRef"

		expectedBasePath := providedBasePath
		expectedPkgRef := providedPkgRef

		fakeGetter := new(fakeGetter)

		objectUnderTest := &pkg{
			getter: fakeGetter,
		}

		/* act */
		objectUnderTest.Get(providedBasePath, providedPkgRef)

		/* assert */
		actualBasePath, actualPkgRef := fakeGetter.GetArgsForCall(0)
		Expect(actualBasePath).To(Equal(expectedBasePath))
		Expect(actualPkgRef).To(Equal(expectedPkgRef))

	})

	It("should return result of getter.Get", func() {

		/* arrange */
		expectedPkgManifest := &model.PkgManifest{Name: "dummyName"}
		expectedError := errors.New("UnmarshalError")

		fakeGetter := new(fakeGetter)
		fakeGetter.GetReturns(expectedPkgManifest, expectedError)

		objectUnderTest := &pkg{
			getter:    fakeGetter,
			validator: new(fakeValidator),
		}

		/* act */
		actualPkgManifest, actualError := objectUnderTest.Get("", "")

		/* assert */
		Expect(actualPkgManifest).To(Equal(expectedPkgManifest))
		Expect(actualError).To(Equal(expectedError))

	})
})
