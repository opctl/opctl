package git

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("handle", func() {

	Context("GetContent", func() {

		It("should call os.Open w/ expected args", func() {
			/* arrange */
			providedOpPath := "dummyOpPath"
			providedContentPath := "dummyContentPath"

			fakeOS := new(ios.Fake)

			objectUnderTest := handle{
				os:   fakeOS,
				path: providedOpPath,
			}

			/* act */
			objectUnderTest.GetContent(nil, providedContentPath)

			/* assert */
			Expect(fakeOS.OpenArgsForCall(0)).To(Equal(filepath.Join(providedOpPath, providedContentPath)))
		})
	})

	wd, err := os.Getwd()
	if nil != err {
		panic(err)
	}
	Context("ListDescendants", func() {
		It("should call ioutil.ReadDir w/ expected args", func() {
			/* arrange */
			providedOpPath := "dummyOpPath"

			fakeIOUtil := new(iioutil.Fake)

			// error to trigger immediate return
			fakeIOUtil.ReadDirReturns(nil, errors.New("dummyError"))

			objectUnderTest := handle{
				ioUtil: fakeIOUtil,
				path:   providedOpPath,
			}

			/* act */
			objectUnderTest.ListDescendants(nil)

			/* assert */
			Expect(fakeIOUtil.ReadDirArgsForCall(0)).To(Equal(providedOpPath))
		})
		Context("ioutil.ReadDir errors", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("dummyError")

				fakeIOUtil := new(iioutil.Fake)
				fakeIOUtil.ReadDirReturns(nil, expectedError)

				objectUnderTest := handle{
					ioUtil: fakeIOUtil,
				}

				/* act */
				_, actualError := objectUnderTest.ListDescendants(nil)

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})
		Context("ioutil.ReadDir doesn't error", func() {
			It("should return expected contentList", func() {
				/* arrange */
				rootOpPath := filepath.Join(wd, "../testdata/listDescendants")

				dirStat, err := os.Stat(filepath.Join(rootOpPath, "/dir1"))
				if nil != err {
					panic(err)
				}

				fileStat, err := os.Stat(filepath.Join(rootOpPath, "/dir1/file2.txt"))
				if nil != err {
					panic(err)
				}

				expectedContents := []*model.DirEntry{
					{
						Mode: fileStat.Mode(),
						Path: "/dir1/file2.txt",
						Size: 34,
					},
					{
						Mode: dirStat.Mode(),
						Path: "/dir1",
						Size: dirStat.Size(),
					},
					{
						Mode: fileStat.Mode(),
						Path: "/file1.txt",
						Size: 18,
					},
				}

				objectUnderTest := handle{
					ioUtil: iioutil.New(),
					path:   rootOpPath,
				}

				/* act */
				actualContents, err := objectUnderTest.ListDescendants(nil)
				if nil != err {
					panic(err)
				}

				/* assert */
				Expect(actualContents).To(Equal(expectedContents))

			})
		})

	})
	Context("Ref", func() {
		It("should return expected ref", func() {
			/* arrange */
			dataRef := "dummyDataRef"

			objectUnderTest := handle{
				dataRef: dataRef,
			}

			/* act */
			actualRef := objectUnderTest.Ref()

			/* assert */
			Expect(actualRef).To(Equal(dataRef))
		})
	})
})
