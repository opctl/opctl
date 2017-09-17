package lockfile

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/pscanary"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

var _ = Context("lockfile", func() {
	Context("New", func() {
		It("should return LockFile", func() {
			/* arrange/act/assert */
			Expect(New()).Should(Not(BeNil()))
		})
	})
	Context("Lock", func() {
		It("should call fs.MkdirAll w/ expected args", func() {
			/* arrange */
			providedFilePath := "/dummy/pid.lock"
			expectedPath := path.Dir(providedFilePath)
			expectedPerms := os.FileMode(0700)

			fakeOS := new(ios.Fake)

			objectUnderTest := lockfile{
				os:       fakeOS,
				psCanary: new(pscanary.Fake),
			}

			/* act */
			objectUnderTest.Lock(providedFilePath)

			/* assert */
			actualPath, actualPerms := fakeOS.MkdirAllArgsForCall(0)
			Expect(actualPath).To(Equal(expectedPath))
			Expect(actualPerms).To(Equal(expectedPerms))
		})
		Context("fs.mkdirAll errors", func() {
			It("should return expected err", func() {
				/* arrange */
				providedFilePath := "/dummy/pid.lock"
				expectedErr := errors.New("dummyError")

				fakeOS := new(ios.Fake)
				fakeOS.MkdirAllReturns(expectedErr)

				objectUnderTest := lockfile{
					os: fakeOS,
				}

				/* act */
				actualErr := objectUnderTest.Lock(providedFilePath)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("fs.mkdirAll doesn't error", func() {
			It("should call fs.Open w/ expected args", func() {
				/* arrange */
				providedFilePath := "/dummy/pid.lock"

				fakeOS := new(ios.Fake)

				objectUnderTest := lockfile{
					os:       fakeOS,
					psCanary: new(pscanary.Fake),
				}

				/* act */
				objectUnderTest.Lock(providedFilePath)

				/* assert */
				Expect(fakeOS.OpenArgsForCall(0)).To(Equal(providedFilePath))
			})
			Context("fs.Open errors", func() {
				It("should call fs.Create w/ expected args", func() {
					/* arrange */
					providedFilePath := "/dummy/pid.lock"

					fakeOS := new(ios.Fake)
					fakeOS.OpenReturns(nil, errors.New("dummyError"))

					objectUnderTest := lockfile{
						os:       fakeOS,
						psCanary: new(pscanary.Fake),
					}

					/* act */
					objectUnderTest.Lock(providedFilePath)

					/* assert */
					Expect(fakeOS.CreateArgsForCall(0)).To(Equal(providedFilePath))
				})

			})
			Context("fs.Open doesn't error", func() {
				Context("lockfile is invalid", func() {
					It("should create new lock file", func() {
						/* arrange */
						providedFilePath := "/dummy/pid.lock"

						fakeOS := new(ios.Fake)

						objectUnderTest := lockfile{
							os:       fakeOS,
							psCanary: new(pscanary.Fake),
						}

						/* act */
						objectUnderTest.Lock(providedFilePath)

						/* assert */
						Expect(fakeOS.CreateArgsForCall(0)).To(Equal(providedFilePath))
					})
				})
				Context("lockfile is valid", func() {
					Context("lock owner dead", func() {
						It("should create new lock file", func() {
							/* arrange */
							providedFilePath := "/dummy/pid.lock"

							fakeOS := new(ios.Fake)

							fakePsCanary := new(pscanary.Fake)
							fakePsCanary.IsAliveReturns(false)

							objectUnderTest := lockfile{
								os:       fakeOS,
								psCanary: fakePsCanary,
							}

							/* act */
							objectUnderTest.Lock(providedFilePath)

							/* assert */
							Expect(fakeOS.CreateArgsForCall(0)).To(Equal(providedFilePath))
						})
					})
					Context("lock owner alive", func() {
						It("should return expected error", func() {
							/* arrange */
							pIdFromFile := 1234
							expectedErr := fmt.Errorf("Unable to obtain lock; currently owned by PId: %v\n", pIdFromFile)

							fakeOS := new(ios.Fake)
							// create a real temp file; no good way to stub os.File
							tempFile, err := ioutil.TempFile("", "lockfile_test")
							defer tempFile.Close()
							if nil != err {
								panic(err)
							}
							_, err = tempFile.WriteString(strconv.Itoa(pIdFromFile))
							if nil != err {
								panic(err)
							}

							tempFile.Sync()
							tempFile.Seek(0, 0)

							fakeOS.OpenReturns(tempFile, nil)

							fakePsCanary := new(pscanary.Fake)
							fakePsCanary.IsAliveReturns(true)

							objectUnderTest := lockfile{
								os:       fakeOS,
								psCanary: fakePsCanary,
							}

							/* act */
							actualErr := objectUnderTest.Lock("/dummy/pid.lock")

							/* assert */
							Expect(actualErr).To(Equal(expectedErr))
						})
					})
				})
			})
		})
	})
})
