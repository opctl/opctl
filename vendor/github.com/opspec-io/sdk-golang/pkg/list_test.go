package pkg

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
	"os"
)

var _ = Context("List", func() {
	wd, err := os.Getwd()
	if nil != err {
		panic(err)
	}
	It("should call ioutil.ReadDir w/ expected args", func() {
		/* arrange */
		providedDirPath := "dummyDirPath"

		fakeIOUtil := new(iioutil.Fake)

		// error to trigger immediate return
		fakeIOUtil.ReadDirReturns(nil, errors.New("dummyError"))

		objectUnderTest := _Pkg{
			ioUtil:   fakeIOUtil,
			manifest: new(manifest.Fake),
		}

		/* act */
		objectUnderTest.List(providedDirPath)

		/* assert */
		Expect(fakeIOUtil.ReadDirArgsForCall(0)).To(Equal(providedDirPath))
	})
	Context("ioutil.ReadDir errors", func() {
		It("should be returned", func() {

			/* arrange */
			expectedError := errors.New("dummyError")

			fakeIOUtil := new(iioutil.Fake)
			fakeIOUtil.ReadDirReturns(nil, expectedError)

			objectUnderTest := _Pkg{
				ioUtil:   fakeIOUtil,
				manifest: new(manifest.Fake),
			}

			/* act */
			_, actualError := objectUnderTest.List("")

			/* assert */
			Expect(actualError).To(Equal(expectedError))

		})
	})
	Context("ioutil.ReadDir doesn't error", func() {
		It("should call manifest.Unmarshal for each childDir", func() {
			/* arrange */
			rootPkgPath := fmt.Sprintf("%v/testdata/list", wd)

			fakeManifest := new(manifest.Fake)

			objectUnderTest := _Pkg{
				ioUtil:   iioutil.New(),
				os:       ios.New(),
				manifest: fakeManifest,
			}

			/* act */
			_, err := objectUnderTest.List(rootPkgPath)
			if nil != err {
				panic(err)
			}

			/* assert */
			Expect(fakeManifest.UnmarshalCallCount()).To(Equal(2))

		})
	})

})
