package bundle

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/util/fs"
	"os"
	"path"
	"time"
)

// dummy fileInfo
type fileInfo struct {
	name    string      // base name of the file
	size    int64       // length in bytes for regular files; system-dependent for others
	mode    os.FileMode // file mode bits
	modTime time.Time   // modification time
	isDir   bool        // abbreviation for Mode().IsDir()
	sys     interface{} // underlying data source (can return nil)
}

func (this fileInfo) Name() string {
	return this.name
}

func (this fileInfo) Size() int64 {
	return this.size
}
func (this fileInfo) Mode() os.FileMode {
	return this.mode
}
func (this fileInfo) ModTime() time.Time {
	return this.modTime
}
func (this fileInfo) IsDir() bool {
	return this.isDir
}
func (this fileInfo) Sys() interface{} {
	return &this.sys
}

var _ = Describe("_tryResolveDefaultCollection", func() {

	Context("Execute", func() {

		Context("when FileSystem.ListChildFileInfosOfDir returns an error", func() {

			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("ListChildFileInfosOfDirError")

				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.ListChildFileInfosOfDirReturns(nil, expectedError)

				objectUnderTest := &_bundle{
					fileSystem: fakeFileSystem,
				}

				/* act */
				_, actualError := objectUnderTest.TryResolveDefaultCollection(
					model.TryResolveDefaultCollectionReq{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})

		})

		Context("when default op collection doesn't exist", func() {

			It("should return an empty pathToDefaultCollection", func() {

				/* arrange */
				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.ListChildFileInfosOfDirReturns(nil, nil)

				objectUnderTest := &_bundle{
					fileSystem: fakeFileSystem,
				}

				/* act */
				pathToDefaultCollection, _ := objectUnderTest.TryResolveDefaultCollection(
					model.TryResolveDefaultCollectionReq{},
				)

				/* assert */
				Expect(pathToDefaultCollection).To(BeEmpty())

			})

			It("should return a nil err", func() {

				/* arrange */
				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.ListChildFileInfosOfDirReturns(nil, nil)

				objectUnderTest := &_bundle{
					fileSystem: fakeFileSystem,
				}

				/* act */
				_, err := objectUnderTest.TryResolveDefaultCollection(
					model.TryResolveDefaultCollectionReq{},
				)

				/* assert */
				Expect(err).To(BeNil())

			})

		})
		Context("when default op collection is in provided PathToDir", func() {

			It("should return expected pathToDefaultCollection", func() {

				/* arrange */
				providedPathToDir := "/dummy/path"

				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.ListChildFileInfosOfDirReturns(
					[]os.FileInfo{
						fileInfo{
							name:  NameOfDefaultOpCollectionDir,
							isDir: true,
						},
					},
					nil,
				)

				expectedPathToDefaultCollection := path.Join(providedPathToDir, NameOfDefaultOpCollectionDir)

				objectUnderTest := &_bundle{
					fileSystem: fakeFileSystem,
				}

				/* act */
				actualPathToDefaultCollection, _ := objectUnderTest.TryResolveDefaultCollection(
					model.TryResolveDefaultCollectionReq{
						PathToDir: providedPathToDir,
					},
				)

				/* assert */
				Expect(actualPathToDefaultCollection).To(Equal(expectedPathToDefaultCollection))

			})

			It("should return a nil err", func() {

				/* arrange */
				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.ListChildFileInfosOfDirReturns(
					[]os.FileInfo{
						fileInfo{
							name:  NameOfDefaultOpCollectionDir,
							isDir: true,
						},
					},
					nil,
				)

				objectUnderTest := &_bundle{
					fileSystem: fakeFileSystem,
				}

				/* act */
				_, err := objectUnderTest.TryResolveDefaultCollection(
					model.TryResolveDefaultCollectionReq{},
				)

				/* assert */
				Expect(err).To(BeNil())

			})

		})

		It("should call FileSystem.ListChildFileInfosOfDir with expected args", func() {

			/* arrange */
			expectedListChildFileInfosOfDirArg := "/dummy/path"

			fakeFileSystem := new(fs.Fake)
			fakeFileSystem.ListChildFileInfosOfDirReturns(nil, nil)

			objectUnderTest := &_bundle{
				fileSystem: fakeFileSystem,
			}

			/* act */
			objectUnderTest.TryResolveDefaultCollection(
				model.TryResolveDefaultCollectionReq{
					PathToDir: expectedListChildFileInfosOfDirArg,
				},
			)

			/* assert */
			actualListChildFileInfosOfDirArg := fakeFileSystem.ListChildFileInfosOfDirArgsForCall(0)
			Expect(actualListChildFileInfosOfDirArg).To(Equal(expectedListChildFileInfosOfDirArg))

		})

	})

})
