package pkg

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"os"
)

var _ = Context("pkg", func() {
	wd, err := os.Getwd()
	if nil != err {
		panic(err)
	}
	Context("ListContents", func() {
		It("should call ioutil.ReadDir w/ expected args", func() {
			/* arrange */
			providedPkgPath := "dummyPkgPath"

			fakeIOUtil := new(iioutil.Fake)

			// error to trigger immediate return
			fakeIOUtil.ReadDirReturns(nil, errors.New("dummyError"))

			objectUnderTest := _Pkg{
				ioUtil: fakeIOUtil,
			}

			/* act */
			objectUnderTest.ListContents(providedPkgPath)

			/* assert */
			Expect(fakeIOUtil.ReadDirArgsForCall(0)).To(Equal(providedPkgPath))
		})
		Context("ioutil.ReadDir errors", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("dummyError")

				fakeIOUtil := new(iioutil.Fake)
				fakeIOUtil.ReadDirReturns(nil, expectedError)

				objectUnderTest := _Pkg{
					ioUtil: fakeIOUtil,
				}

				/* act */
				_, actualError := objectUnderTest.ListContents("")

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})
		Context("ioutil.ReadDir doesn't error", func() {
			It("should return expected contentList", func() {
				/* arrange */
				rootPkgPath := fmt.Sprintf("%v/testdata/listContents", wd)
				expectedContents := []*model.PkgContent{
					{
						Path: fmt.Sprintf("%v/dir1/file2.txt", rootPkgPath),
						Size: 34,
					},
					{
						Path: fmt.Sprintf("%v/file1.txt", rootPkgPath),
						Size: 18,
					},
				}

				objectUnderTest := _Pkg{
					ioUtil: iioutil.New(),
				}

				/* act */
				actualContents, err := objectUnderTest.ListContents(rootPkgPath)
				if nil != err {
					panic(err)
				}

				/* assert */
				Expect(actualContents).To(Equal(expectedContents))

			})
		})

	})
})
