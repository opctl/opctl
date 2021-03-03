package fs

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("_fs", func() {
	Context("TryResolve", func() {
		Context("dataRef is absolute path", func() {
			Context("doesnt exist", func() {
				It("should return err", func() {
					/* arrange */
					objectUnderTest := _fs{}

					/* act */
					actualHandle, actualError := objectUnderTest.TryResolve(
						context.Background(),
						"/doesnt-exist",
					)

					/* assert */
					Expect(actualHandle).To(BeNil())
					Expect(actualError).To(BeNil())
				})
			})
			Context("exists", func() {
				It("should return expected result", func() {
					/* arrange */
					file, err := ioutil.TempFile("", "")
					if nil != err {
						panic(err)
					}

					expectedHandle := newHandle(file.Name())

					objectUnderTest := _fs{}

					/* act */
					actualHandle, actualError := objectUnderTest.TryResolve(
						context.Background(),
						file.Name(),
					)

					/* assert */
					Expect(actualHandle).To(Equal(expectedHandle))
					Expect(actualError).To(BeNil())
				})
			})
		})
		Context("dataRef isn't absolute path", func() {
			Context("doesnt exist", func() {
				It("should return err", func() {
					/* arrange */
					objectUnderTest := _fs{}

					/* act */
					actualHandle, actualError := objectUnderTest.TryResolve(
						context.Background(),
						"doesnt-exist",
					)

					/* assert */
					Expect(actualHandle).To(BeNil())
					Expect(actualError).To(BeNil())
				})
			})
			Context("exists", func() {
				It("should return expected result", func() {
					/* arrange */
					basePath, err := os.Getwd()
					if nil != err {
						panic(err)
					}

					providedDataRef := "testdata/file1.txt"

					expectedHandle := newHandle(filepath.Join(
						basePath,
						providedDataRef,
					))

					objectUnderTest := _fs{
						basePaths: []string{basePath},
					}

					/* act */
					actualHandle, actualError := objectUnderTest.TryResolve(
						context.Background(),
						providedDataRef,
					)

					/* assert */
					Expect(actualHandle).To(Equal(expectedHandle))
					Expect(actualError).To(BeNil())
				})
			})
		})
	})
})
