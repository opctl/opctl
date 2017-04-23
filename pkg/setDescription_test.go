package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetDescription", func() {

	It("should call manifestUnmarshaller w/ expected args", func() {
		/* arrange */
		providedPkgPath := "dummyPkgPath"
		providedPkgDescription := "dummyPkgDescription"

		fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)
		// return error to trigger immediate return
		fakeManifestUnmarshaller.UnmarshalReturns(nil, errors.New("dummyError"))

		objectUnderTest := pkg{
			manifestUnmarshaller: fakeManifestUnmarshaller,
		}

		/* act */
		objectUnderTest.SetDescription(providedPkgPath, providedPkgDescription)

		/* assert */
		Expect(fakeManifestUnmarshaller.UnmarshalArgsForCall(0)).To(Equal(providedPkgPath))
	})
	Context("manifestUnmarshaller.Unmarshal errors", func() {
		It("should return error", func() {
			/* arrange */
			expectedError := errors.New("dummyError")

			fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)
			// return error to trigger immediate return
			fakeManifestUnmarshaller.UnmarshalReturns(nil, errors.New("dummyError"))

			objectUnderTest := pkg{
				manifestUnmarshaller: fakeManifestUnmarshaller,
			}

			/* act */
			actualError := objectUnderTest.SetDescription("", "")

			/* assert */
			Expect(actualError).To(Equal(expectedError))
		})
	})

})
