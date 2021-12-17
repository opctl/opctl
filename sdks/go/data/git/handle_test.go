package git

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("handle", func() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	Context("GetContent", func() {
		It("should not err", func() {
			/* arrange */

			providedOpPath := wd
			providedContentPath := "testdata/file1.txt"

			objectUnderTest := handle{
				path: providedOpPath,
			}

			/* act */
			_, actualErr := objectUnderTest.GetContent(nil, providedContentPath)

			/* assert */
			Expect(actualErr).To(BeNil())
		})
	})
	Context("ListDescendants", func() {
		Context("os.ReadDir errors", func() {
			It("should be returned", func() {

				/* arrange */
				providedPath := "doesnt-exist"

				objectUnderTest := handle{
					path: providedPath,
				}

				/* act */
				_, actualError := objectUnderTest.ListDescendants(nil)

				/* assert */
				Expect(actualError.Error()).To(Equal("open doesnt-exist: no such file or directory"))

			})
		})
		Context("os.ReadDir doesn't error", func() {
			It("should return expected contentList", func() {
				/* arrange */
				rootOpPath := filepath.Join(wd, "../testdata/listDescendants")

				dirStat, err := os.Stat(filepath.Join(rootOpPath, "/dir1"))
				if err != nil {
					panic(err)
				}

				fileStat, err := os.Stat(filepath.Join(rootOpPath, "/dir1/file2.txt"))
				if err != nil {
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
					path: rootOpPath,
				}

				/* act */
				actualContents, err := objectUnderTest.ListDescendants(nil)
				if err != nil {
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
