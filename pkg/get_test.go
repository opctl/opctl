package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("Get", func() {
	It("should call manifestUnmarshaller.Unmarshal w/ expected inputs", func() {
		/* arrange */
		providedPkgRef := "dummyPkgRef"

		expectedPkgRef := providedPkgRef

		fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)

		objectUnderTest := &pkg{
			manifestUnmarshaller: fakeManifestUnmarshaller,
		}

		/* act */
		objectUnderTest.Get(providedPkgRef)

		/* assert */
		actualPkgRef := fakeManifestUnmarshaller.UnmarshalArgsForCall(0)
		Expect(actualPkgRef).To(Equal(expectedPkgRef))

	})

	It("should return result of manifestUnmarshaller.Unmarshal", func() {

		/* arrange */
		expectedPkgManifest := &model.PkgManifest{Name: "dummyName"}
		expectedError := errors.New("UnmarshalError")

		fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)
		fakeManifestUnmarshaller.UnmarshalReturns(expectedPkgManifest, expectedError)

		objectUnderTest := &pkg{
			manifestUnmarshaller: fakeManifestUnmarshaller,
		}

		/* act */
		actualPkgManifest, actualError := objectUnderTest.Get("")

		/* assert */
		Expect(actualPkgManifest).To(Equal(expectedPkgManifest))
		Expect(actualError).To(Equal(expectedError))

	})
})
