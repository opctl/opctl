package dircopier

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/filecopier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Context("dirCopier", func() {
	Context("New", func() {
		It("should return DirCopier", func() {
			/* arrange/act/assert */
			Expect(New()).Should(Not(BeNil()))
		})
	})
	Context("OS", func() {

		It("should call fs.Stat w/ expected args", func() {
			/* arrange */
			providedSrcPath := "dummySrcPath"

			fakeOS := new(ios.Fake)
			// trigger exit
			fakeOS.StatReturns(nil, errors.New("dummyError"))

			objectUnderTest := dirCopier{
				os: fakeOS,
			}

			/* act */
			objectUnderTest.OS(providedSrcPath, "dummyDstPath")

			/* assert */
			Expect(fakeOS.StatArgsForCall(0)).To(Equal(providedSrcPath))
		})
		Context("fs.Stat errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, expectedError)

				objectUnderTest := dirCopier{
					os: fakeOS,
				}

				/* act */
				actualError := objectUnderTest.OS("dummySrcPath", "dummyDstPath")

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("fs.Stat doesn't error", func() {
			Context("src isn't dir", func() {
				It("should return expected error", func() {
					/* arrange */
					fakeOS := new(ios.Fake)
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
					fakeOS.StatReturns(srcFileInfo, nil)

					providedSrcPath := "dummySrcPath"
					expectedError := fmt.Errorf("%v is not a dir", providedSrcPath)

					objectUnderTest := dirCopier{
						os: fakeOS,
					}

					/* act */
					actualError := objectUnderTest.OS(providedSrcPath, "dummyDstPath")

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("src is dir", func() {
				It("should call os.MkdirAll w/ expected args", func() {
					/* arrange */
					providedDstPath := "dummyDstPath"

					fakeOS := new(ios.Fake)
					// create a real srcDir; no good way to stub os.FileInfo
					srcDirInfo, err := os.Stat(os.TempDir())
					if nil != err {
						panic(err)
					}
					fakeOS.StatReturns(srcDirInfo, nil)

					// trigger exit
					fakeOS.MkdirAllReturns(errors.New("dummyError"))

					objectUnderTest := dirCopier{
						os: fakeOS,
					}

					/* act */
					objectUnderTest.OS("dummySrcPath", providedDstPath)

					/* assert */
					actualDstDirPath, actualDstDirMode := fakeOS.MkdirAllArgsForCall(0)
					Expect(actualDstDirPath).To(Equal(providedDstPath))
					Expect(actualDstDirMode).To(Equal(srcDirInfo.Mode()))
				})
				Context("os.MkdirAll errors", func() {
					It("should return expected error", func() {
						/* arrange */
						fakeOS := new(ios.Fake)
						// create a real srcDir; no good way to stub os.FileInfo
						srcDirInfo, err := os.Stat(os.TempDir())
						if nil != err {
							panic(err)
						}
						fakeOS.StatReturns(srcDirInfo, nil)

						expectedError := errors.New("dummyError")

						// trigger exit
						fakeOS.MkdirAllReturns(expectedError)

						objectUnderTest := dirCopier{
							os: fakeOS,
						}

						/* act */
						actualError := objectUnderTest.OS("dummySrcPath", "dummyDstPath")

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
							expectedSrcFilePath := filepath.Join(srcDirPath, fileName)
							expectedDstFilePath := filepath.Join(dstDirPath, fileName)
							file, err := os.Create(expectedSrcFilePath)
							defer file.Close()
							if nil != err {
								panic(err)
							}

							fakeFileCopier := new(filecopier.Fake)

							objectUnderTest := dirCopier{
								os:         ios.New(),
								ioutil:     iioutil.New(),
								fileCopier: fakeFileCopier,
							}

							/* act */
							objectUnderTest.OS(srcDirPath, dstDirPath)

							/* assert */
							actualSrcFilePath, actualDstFilePath := fakeFileCopier.OSArgsForCall(0)
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
								srcChildDirPath := filepath.Join(srcDirPath, childDirName)
								dstChildDirPath := filepath.Join(dstDirPath, childDirName)
								err = os.Mkdir(srcChildDirPath, 0600)
								if nil != err {
									panic(err)
								}

								// create a real srcChildDirFile
								childDirFileFileName := "file.dummy"
								expectedSrcChildDirFilePath := filepath.Join(srcChildDirPath, childDirFileFileName)
								expectedDstChildDirFilePath := filepath.Join(dstChildDirPath, childDirFileFileName)
								childDirFile, err := os.Create(expectedSrcChildDirFilePath)
								defer childDirFile.Close()
								if nil != err {
									panic(err)
								}

								fakeFileCopier := new(filecopier.Fake)

								objectUnderTest := dirCopier{
									os:         ios.New(),
									ioutil:     iioutil.New(),
									fileCopier: fakeFileCopier,
								}

								/* act */
								objectUnderTest.OS(srcDirPath, dstDirPath)

								/* assert */
								actualSrcFilePath, actualDstFilePath := fakeFileCopier.OSArgsForCall(0)
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
