package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/virtual-go/vioutil"
)

var _ = Describe("listPackagesInDir", func() {

	Context("ListPackagesInDir", func() {

		Context("when ioutil.ReadDir returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("dummyError")

				fakeIOUtil := new(vioutil.Fake)
				fakeIOUtil.ReadDirReturns(nil, expectedError)

				objectUnderTest := pkg{
					ioUtil:      fakeIOUtil,
					viewFactory: new(fakeViewFactory),
					yaml:        new(format.Fake),
				}

				/* act */
				_, actualError := objectUnderTest.ListPackagesInDir("/dummy/path")

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})
	})

})
