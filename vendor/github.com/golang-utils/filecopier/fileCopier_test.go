package filecopier

import (
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	Context("OS", func() {
		It("should call fs.Open w/ expected args", func() {
			/* arrange */
			providedSrcPath := "dummySrcPath"

			fakeOS := new(ios.Fake)
			// trigger exit
			fakeOS.OpenReturns(nil, errors.New("dummyError"))

			objectUnderTest := fileCopier{
				os: fakeOS,
			}

			/* act */
			objectUnderTest.OS(providedSrcPath, "dummyDstPath")

			/* assert */
			Expect(fakeOS.OpenArgsForCall(0)).To(Equal(providedSrcPath))
		})
		Context("fs.Open errors", func() {
			It("returns expected error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")

				fakeOS := new(ios.Fake)
				fakeOS.OpenReturns(nil, expectedError)

				objectUnderTest := fileCopier{
					os: fakeOS,
				}

				/* act */
				actualError := objectUnderTest.OS("dummySrcPath", "dummyDstPath")

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("fs.Open doesn't error", func() {
			It("should call fs.Stat w/ expected args", func() {
				/* arrange */
				providedSrcPath := "dummySrcPath"
				providedDstPath := "dummyDstPath"

				fakeOS := new(ios.Fake)
				// trigger exit
				fakeOS.StatReturns(nil, errors.New("dummyError"))

				objectUnderTest := fileCopier{
					os: fakeOS,
				}

				/* act */
				objectUnderTest.OS(providedSrcPath, providedDstPath)

				/* assert */
				Expect(fakeOS.StatArgsForCall(0)).To(Equal(providedSrcPath))
			})
			Context("fs.Stat errors", func() {
				It("returns expected error", func() {
					/* arrange */
					expectedError := errors.New("dummyError")

					fakeOS := new(ios.Fake)
					fakeOS.StatReturns(nil, expectedError)

					objectUnderTest := fileCopier{
						os: fakeOS,
					}

					/* act */
					actualError := objectUnderTest.OS("dummySrcPath", "dummyDstPath")

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("fs.Stat doesn't error", func() {
				It("should call fs.Create w/ expected args", func() {
					/* arrange */
					providedDstPath := "dummyDstPath"

					fakeOS := new(ios.Fake)
					// trigger exit
					fakeOS.CreateReturns(nil, errors.New("dummyError"))

					objectUnderTest := fileCopier{
						os: fakeOS,
					}

					/* act */
					objectUnderTest.OS("dummySrcPath", providedDstPath)

					/* assert */
					Expect(fakeOS.CreateArgsForCall(0)).To(Equal(providedDstPath))
				})
				Context("fs.Create errors", func() {
					It("returns expected error", func() {
						/* arrange */
						expectedError := errors.New("dummyError")

						fakeOS := new(ios.Fake)
						fakeOS.CreateReturns(nil, expectedError)

						objectUnderTest := fileCopier{
							os: fakeOS,
						}

						/* act */
						actualError := objectUnderTest.OS("dummySrcPath", "dummyDstPath")

						/* assert */
						Expect(actualError).To(Equal(expectedError))
					})
				})
				Context("fs.Create doesn't error", func() {
					It("should call fs.Chmod w/ expected args", func() {
						/* arrange */
						providedDstPath := "dummyDstPath"

						fakeOS := new(ios.Fake)

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

						fakeOS.StatReturns(srcFileInfo, nil)

						// trigger exit
						fakeOS.ChmodReturns(errors.New("dummyError"))

						objectUnderTest := fileCopier{
							os: fakeOS,
						}

						/* act */
						objectUnderTest.OS("dummySrcPath", providedDstPath)

						/* assert */
						actualDstPath, actualMode := fakeOS.ChmodArgsForCall(0)

						Expect(actualDstPath).To(Equal(providedDstPath))
						Expect(actualMode).To(Equal(srcFileInfo.Mode()))
					})
					Context("os.Chmod errors", func() {
						It("should return expected error", func() {
							/* arrange */
							fakeOS := new(ios.Fake)

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
							fakeOS.StatReturns(srcFileInfo, nil)

							objectUnderTest := fileCopier{
								os: fakeOS,
							}

							expectedError := errors.New("dummyError")
							fakeOS.ChmodReturns(expectedError)

							/* act */
							actualError := objectUnderTest.OS("dummySrcPath", "dummyDstPath")

							/* assert */
							Expect(actualError).To(Equal(expectedError))
						})
					})
					Context("os.Chmod doesn't error", func() {
						It("doesn't error", func() {

							/* arrange */
							fakeOS := new(ios.Fake)

							// create a real srcFile; no good way to stub os.FileInfo
							srcFile, err := ioutil.TempFile("", "fileCopier_test")
							defer srcFile.Close()
							if nil != err {
								panic(err)
							}
							fakeOS.OpenReturns(srcFile, nil)

							srcFileInfo, err := os.Stat(srcFile.Name())
							if nil != err {
								panic(err)
							}
							fakeOS.StatReturns(srcFileInfo, nil)

							// create a real dstFile; no good way to stub os.FileInfo
							dstFile, err := ioutil.TempFile("", "fileCopier_test")
							defer srcFile.Close()
							if nil != err {
								panic(err)
							}
							fakeOS.CreateReturns(dstFile, nil)

							objectUnderTest := fileCopier{
								os: fakeOS,
							}

							/* act */
							actualError := objectUnderTest.OS("dummySrcPath", "dummyDstPath")

							/* assert */
							Expect(actualError).To(BeNil())
						})
					})
				})
			})
		})
	})
})
