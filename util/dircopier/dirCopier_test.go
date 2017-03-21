package dircopier

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/filecopier"
	"github.com/opctl/opctl/util/vfs"
	osfs "github.com/opctl/opctl/util/vfs/os"
	"io/ioutil"
	"os"
	"path"
)

var _ = Context("dirCopier", func() {
	Context("New", func() {
		It("should return DirCopier", func() {
			/* arrange/act/assert */
			Expect(New()).Should(Not(BeNil()))
		})
	})
	Context("Fs", func() {

		It("should call fs.Stat w/ expected args", func() {
			/* arrange */
			providedSrcPath := "dummySrcPath"

			fakeFs := new(vfs.Fake)
			// trigger exit
			fakeFs.StatReturns(nil, errors.New("dummyError"))

			objectUnderTest := dirCopier{
				fs: fakeFs,
			}

			/* act */
			objectUnderTest.Fs(providedSrcPath, "dummyDstPath")

			/* assert */
			Expect(fakeFs.StatArgsForCall(0)).To(Equal(providedSrcPath))
		})
		Context("fs.Stat errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")

				fakeFs := new(vfs.Fake)
				fakeFs.StatReturns(nil, expectedError)

				objectUnderTest := dirCopier{
					fs: fakeFs,
				}

				/* act */
				actualError := objectUnderTest.Fs("dummySrcPath", "dummyDstPath")

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("fs.Stat doesn't error", func() {
			Context("src isn't dir", func() {
				It("should return expected error", func() {
					/* arrange */
					fakeFs := new(vfs.Fake)
					// create a real srcFile; no good way to stub os.FileInfo
					srcFile, err := ioutil.TempFile("", "dirCopier_test")
					defer srcFile.Close()
					if nil != err {
						panic(err)
					}

					srcFileInfo, err := os.Stat(srcFile.Name())
					if nil != err {
						panic(err)
					}
					fakeFs.StatReturns(srcFileInfo, nil)

					providedSrcPath := "dummySrcPath"
					expectedError := fmt.Errorf("%v is not a dir", providedSrcPath)

					objectUnderTest := dirCopier{
						fs: fakeFs,
					}

					/* act */
					actualError := objectUnderTest.Fs(providedSrcPath, "dummyDstPath")

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("src is dir", func() {
				It("should call os.MkdirAll w/ expected args", func() {
					/* arrange */
					providedDstPath := "dummyDstPath"

					fakeFs := new(vfs.Fake)
					// create a real srcDir; no good way to stub os.FileInfo
					srcDirInfo, err := os.Stat(os.TempDir())
					if nil != err {
						panic(err)
					}
					fakeFs.StatReturns(srcDirInfo, nil)

					// trigger exit
					fakeFs.MkdirAllReturns(errors.New("dummyError"))

					objectUnderTest := dirCopier{
						fs: fakeFs,
					}

					/* act */
					objectUnderTest.Fs("dummySrcPath", providedDstPath)

					/* assert */
					actualDstDirPath, actualDstDirMode := fakeFs.MkdirAllArgsForCall(0)
					Expect(actualDstDirPath).To(Equal(providedDstPath))
					Expect(actualDstDirMode).To(Equal(srcDirInfo.Mode()))
				})
				Context("os.MkdirAll errors", func() {
					It("should return expected error", func() {
						/* arrange */
						fakeFs := new(vfs.Fake)
						// create a real srcDir; no good way to stub os.FileInfo
						srcDirInfo, err := os.Stat(os.TempDir())
						if nil != err {
							panic(err)
						}
						fakeFs.StatReturns(srcDirInfo, nil)

						expectedError := errors.New("dummyError")

						// trigger exit
						fakeFs.MkdirAllReturns(expectedError)

						objectUnderTest := dirCopier{
							fs: fakeFs,
						}

						/* act */
						actualError := objectUnderTest.Fs("dummySrcPath", "dummyDstPath")

						/* assert */
						Expect(actualError).To(Equal(expectedError))
					})
				})
				Context("os.MkdirAll doesn't error", func() {
					Context("srcDir contains a file", func() {
						It("should call fileCopier.Fs w/ expected args", func() {
							/* arrange */
							// create a real srcDir & dstDir
							srcDirPath, err := ioutil.TempDir("", "dirCopier_test")
							if nil != err {
								panic(err)
							}
							dstDirPath, err := ioutil.TempDir("", "dirCopier_test")
							if nil != err {
								panic(err)
							}

							// create a real file
							fileName := "file.dummy"
							expectedSrcFilePath := path.Join(srcDirPath, fileName)
							expectedDstFilePath := path.Join(dstDirPath, fileName)
							file, err := os.Create(expectedSrcFilePath)
							defer file.Close()
							if nil != err {
								panic(err)
							}

							fakeFileCopier := new(filecopier.Fake)

							objectUnderTest := dirCopier{
								fs:         osfs.New(),
								fileCopier: fakeFileCopier,
							}

							/* act */
							objectUnderTest.Fs(srcDirPath, dstDirPath)

							/* assert */
							actualSrcFilePath, actualDstFilePath := fakeFileCopier.FsArgsForCall(0)
							Expect(actualSrcFilePath).To(Equal(expectedSrcFilePath))
							Expect(actualDstFilePath).To(Equal(expectedDstFilePath))
						})
					})
					Context("srcDir contains a dir", func() {
						Context("dir contains a file", func() {
							It("should call fileCopier.Fs w/ expected args", func() {
								/* arrange */
								// create a real srcDir & dstDir
								srcDirPath, err := ioutil.TempDir("", "dirCopier_test")
								if nil != err {
									panic(err)
								}
								dstDirPath, err := ioutil.TempDir("", "dirCopier_test")
								if nil != err {
									panic(err)
								}

								// create a real srcChildDir
								childDirName := "dummyDir"
								srcChildDirPath := path.Join(srcDirPath, childDirName)
								dstChildDirPath := path.Join(dstDirPath, childDirName)
								err = os.Mkdir(srcChildDirPath, 0600)
								if nil != err {
									panic(err)
								}

								// create a real srcChildDirFile
								childDirFileFileName := "file.dummy"
								expectedSrcChildDirFilePath := path.Join(srcChildDirPath, childDirFileFileName)
								expectedDstChildDirFilePath := path.Join(dstChildDirPath, childDirFileFileName)
								childDirFile, err := os.Create(expectedSrcChildDirFilePath)
								defer childDirFile.Close()
								if nil != err {
									panic(err)
								}

								fakeFileCopier := new(filecopier.Fake)

								objectUnderTest := dirCopier{
									fs:         osfs.New(),
									fileCopier: fakeFileCopier,
								}

								/* act */
								objectUnderTest.Fs(srcDirPath, dstDirPath)

								/* assert */
								actualSrcFilePath, actualDstFilePath := fakeFileCopier.FsArgsForCall(0)
								Expect(actualSrcFilePath).To(Equal(expectedSrcChildDirFilePath))
								Expect(actualDstFilePath).To(Equal(expectedDstChildDirFilePath))
							})
						})
					})
				})
			})
		})
	})
})
