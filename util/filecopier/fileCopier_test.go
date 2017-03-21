package filecopier

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/vfs"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

var _ = Context("fileCopier", func() {
	Context("New", func() {
		It("should return FileCopier", func() {
			/* arrange/act/assert */
			Expect(New()).Should(Not(BeNil()))
		})
	})
	Context("Fs", func() {
		It("should call fs.Open w/ expected args", func() {
			/* arrange */
			providedSrcPath := "dummySrcPath"

			fakeFs := new(vfs.Fake)
			// trigger exit
			fakeFs.OpenReturns(nil, errors.New("dummyError"))

			objectUnderTest := fileCopier{
				fs: fakeFs,
			}

			/* act */
			objectUnderTest.Fs(providedSrcPath, "dummyDstPath")

			/* assert */
			Expect(fakeFs.OpenArgsForCall(0)).To(Equal(providedSrcPath))
		})
		Context("fs.Open errors", func() {
			It("returns expected error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")

				fakeFs := new(vfs.Fake)
				fakeFs.OpenReturns(nil, expectedError)

				objectUnderTest := fileCopier{
					fs: fakeFs,
				}

				/* act */
				actualError := objectUnderTest.Fs("dummySrcPath", "dummyDstPath")

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("fs.Open doesn't error", func() {
			It("should call fs.Stat w/ expected args", func() {
				/* arrange */
				providedSrcPath := "dummySrcPath"
				providedDstPath := "dummyDstPath"

				fakeFs := new(vfs.Fake)
				// trigger exit
				fakeFs.StatReturns(nil, errors.New("dummyError"))

				objectUnderTest := fileCopier{
					fs: fakeFs,
				}

				/* act */
				objectUnderTest.Fs(providedSrcPath, providedDstPath)

				/* assert */
				Expect(fakeFs.StatArgsForCall(0)).To(Equal(providedSrcPath))
			})
			Context("fs.Stat errors", func() {
				It("returns expected error", func() {
					/* arrange */
					expectedError := errors.New("dummyError")

					fakeFs := new(vfs.Fake)
					fakeFs.StatReturns(nil, expectedError)

					objectUnderTest := fileCopier{
						fs: fakeFs,
					}

					/* act */
					actualError := objectUnderTest.Fs("dummySrcPath", "dummyDstPath")

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("fs.Stat doesn't error", func() {
				It("should call fs.Create w/ expected args", func() {
					/* arrange */
					providedDstPath := "dummyDstPath"

					fakeFs := new(vfs.Fake)
					// trigger exit
					fakeFs.CreateReturns(nil, errors.New("dummyError"))

					objectUnderTest := fileCopier{
						fs: fakeFs,
					}

					/* act */
					objectUnderTest.Fs("dummySrcPath", providedDstPath)

					/* assert */
					Expect(fakeFs.CreateArgsForCall(0)).To(Equal(providedDstPath))
				})
				Context("fs.Create errors", func() {
					It("returns expected error", func() {
						/* arrange */
						expectedError := errors.New("dummyError")

						fakeFs := new(vfs.Fake)
						fakeFs.CreateReturns(nil, expectedError)

						objectUnderTest := fileCopier{
							fs: fakeFs,
						}

						/* act */
						actualError := objectUnderTest.Fs("dummySrcPath", "dummyDstPath")

						/* assert */
						Expect(actualError).To(Equal(expectedError))
					})
				})
				Context("fs.Create doesn't error", func() {
					It("should call fs.Chmod w/ expected args", func() {
						/* arrange */
						providedDstPath := "dummyDstPath"

						fakeFs := new(vfs.Fake)

						// create a real srcFile; no good way to stub os.FileInfo
						srcFile, err := ioutil.TempFile("", "fileCopier_test")
						defer srcFile.Close()
						if nil != err {
							panic(err)
						}
						srcFileInfo, err := os.Stat(srcFile.Name())
						if nil != err {
							panic(err)
						}

						fakeFs.StatReturns(srcFileInfo, nil)

						// trigger exit
						fakeFs.ChmodReturns(errors.New("dummyError"))

						objectUnderTest := fileCopier{
							fs: fakeFs,
						}

						/* act */
						objectUnderTest.Fs("dummySrcPath", providedDstPath)

						/* assert */
						actualDstPath, actualMode := fakeFs.ChmodArgsForCall(0)

						Expect(actualDstPath).To(Equal(providedDstPath))
						Expect(actualMode).To(Equal(srcFileInfo.Mode()))
					})
					Context("os.Chmod errors", func() {
						It("should return expected error", func() {
							/* arrange */
							fakeFs := new(vfs.Fake)

							// create a real srcFile; no good way to stub os.FileInfo
							srcFile, err := ioutil.TempFile("", "fileCopier_test")
							defer srcFile.Close()
							if nil != err {
								panic(err)
							}
							srcFileInfo, err := os.Stat(srcFile.Name())
							if nil != err {
								panic(err)
							}
							fakeFs.StatReturns(srcFileInfo, nil)

							objectUnderTest := fileCopier{
								fs: fakeFs,
							}

							expectedError := errors.New("dummyError")
							fakeFs.ChmodReturns(expectedError)

							/* act */
							actualError := objectUnderTest.Fs("dummySrcPath", "dummyDstPath")

							/* assert */
							Expect(actualError).To(Equal(expectedError))
						})
					})
					Context("os.Chmod doesn't error", func() {
						It("doesn't error", func() {

							/* arrange */
							fakeFs := new(vfs.Fake)

							// create a real srcFile; no good way to stub os.FileInfo
							srcFile, err := ioutil.TempFile("", "fileCopier_test")
							defer srcFile.Close()
							if nil != err {
								panic(err)
							}
							fakeFs.OpenReturns(srcFile, nil)

							srcFileInfo, err := os.Stat(srcFile.Name())
							if nil != err {
								panic(err)
							}
							fakeFs.StatReturns(srcFileInfo, nil)

							// create a real dstFile; no good way to stub os.FileInfo
							dstFile, err := ioutil.TempFile("", "fileCopier_test")
							defer srcFile.Close()
							if nil != err {
								panic(err)
							}
							fakeFs.CreateReturns(dstFile, nil)

							objectUnderTest := fileCopier{
								fs: fakeFs,
							}

							/* act */
							actualError := objectUnderTest.Fs("dummySrcPath", "dummyDstPath")

							/* assert */
							Expect(actualError).To(BeNil())
						})
					})
				})
			})
		})
	})
})
