package pkg

import (
	"errors"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
	"gopkg.in/yaml.v2"
	"os"
)

var _ = Describe("SetDescription", func() {

	It("should call manifestUnmarshaller w/ expected args", func() {
		/* arrange */
		providedPkgPath := "dummyPkgPath"
		providedPkgDescription := "dummyPkgDescription"

		fakeManifest := new(manifest.Fake)
		// return error to trigger immediate return
		fakeManifest.UnmarshalReturns(nil, errors.New("dummyError"))

		objectUnderTest := _Pkg{
			manifest: fakeManifest,
		}

		/* act */
		objectUnderTest.SetDescription(providedPkgPath, providedPkgDescription)

		/* assert */
		Expect(fakeManifest.UnmarshalArgsForCall(0)).To(Equal(providedPkgPath))
	})
	Context("manifestUnmarshaller.Unmarshal errors", func() {
		It("should return error", func() {
			/* arrange */
			expectedError := errors.New("dummyError")

			fakeManifestUnmarshaller := new(manifest.Fake)
			// return error to trigger immediate return
			fakeManifestUnmarshaller.UnmarshalReturns(nil, errors.New("dummyError"))

			objectUnderTest := _Pkg{
				manifest: fakeManifestUnmarshaller,
			}

			/* act */
			actualError := objectUnderTest.SetDescription("", "")

			/* assert */
			Expect(actualError).To(Equal(expectedError))
		})
	})
	Context("manifestUnmarshaller.Unmarshal doesn't error", func() {
		It("should call ioutil.WriteFile w/ expected args", func() {
			/* arrange */
			providedManifestPath := "dummyManifestPath"
			providedDescription := "dummyDescription"

			expectedManifest := &model.PkgManifest{
				Name:        "dummyPkgName",
				Description: providedDescription,
			}

			fakeManifestUnmarshaller := new(manifest.Fake)
			fakeManifestUnmarshaller.UnmarshalReturns(expectedManifest, nil)

			expectedFilename := providedManifestPath
			expectedData, err := yaml.Marshal(expectedManifest)
			if nil != err {
				panic(err)
			}
			expectedPerm := os.FileMode(0777)

			fakeIOUtil := new(iioutil.Fake)

			objectUnderTest := _Pkg{
				manifest: fakeManifestUnmarshaller,
				ioUtil:   fakeIOUtil,
			}

			/* act */
			objectUnderTest.SetDescription(providedManifestPath, providedDescription)

			/* assert */
			actualFilename, actualData, actualPerm := fakeIOUtil.WriteFileArgsForCall(0)
			Expect(actualFilename).To(Equal(expectedFilename))
			Expect(actualData).To(Equal(expectedData))
			Expect(actualPerm).To(Equal(expectedPerm))
		})
	})

})
