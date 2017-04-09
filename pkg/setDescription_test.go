package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetDescription", func() {

	It("should call manifestUnmarshaller w/ expected args", func() {
		/* arrange */
		providedReq := SetDescriptionReq{
			Path:        "dummyPath",
			Description: "dummyDescription",
		}

		fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)
		// return error to trigger immediate return
		fakeManifestUnmarshaller.UnmarshalReturns(nil, errors.New("dummyError"))

		objectUnderTest := pkg{
			manifestUnmarshaller: fakeManifestUnmarshaller,
		}

		/* act */
		objectUnderTest.SetDescription(providedReq)

		/* assert */
		Expect(fakeManifestUnmarshaller.UnmarshalArgsForCall(0)).To(Equal(providedReq.Path))
	})

})
