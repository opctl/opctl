package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/virtual-go/vioutil"
)

var _ = Describe("List", func() {

	Context("when ioutil.ReadDir returns an error", func() {

		It("should be returned", func() {

			/* arrange */
			expectedError := errors.New("dummyError")

			fakeIOUtil := new(vioutil.Fake)
			fakeIOUtil.ReadDirReturns(nil, expectedError)

			objectUnderTest := pkg{
				ioUtil:               fakeIOUtil,
				manifestUnmarshaller: new(fakeManifestUnmarshaller),
			}

			/* act */
			_, actualError := objectUnderTest.List("")

			/* assert */
			Expect(actualError).To(Equal(expectedError))

		})

	})

})
