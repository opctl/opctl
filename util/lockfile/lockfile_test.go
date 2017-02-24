package lockfile

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/pscanary"
	"github.com/opspec-io/opctl/util/vfs"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/pkg/errors"
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

			fakeFs := new(vfs.Fake)

			objectUnderTest := lockfile{
				fs:       fakeFs,
				psCanary: new(pscanary.Fake),
				os:       new(vos.Fake),
			}

			/* act */
			objectUnderTest.Lock(providedFilePath)

			/* assert */
			actualPath, actualPerms := fakeFs.MkdirAllArgsForCall(0)
			Expect(actualPath).To(Equal(expectedPath))
			Expect(actualPerms).To(Equal(expectedPerms))
		})
		Context("fs.mkdirAll errors", func() {
			It("should return expected err", func() {
				/* arrange */
				providedFilePath := "/dummy/pid.lock"
				expectedErr := errors.New("dummyError")

				fakeFs := new(vfs.Fake)
				fakeFs.MkdirAllReturns(expectedErr)

				objectUnderTest := lockfile{
					fs: fakeFs,
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

				fakeFs := new(vfs.Fake)

				objectUnderTest := lockfile{
					fs:       fakeFs,
					psCanary: new(pscanary.Fake),
					os:       new(vos.Fake),
				}

				/* act */
				objectUnderTest.Lock(providedFilePath)

				/* assert */
				Expect(fakeFs.OpenArgsForCall(0)).To(Equal(providedFilePath))
			})
			Context("fs.Open errors", func() {
				It("should call fs.Create w/ expected args", func() {
					/* arrange */
					providedFilePath := "/dummy/pid.lock"

					fakeFs := new(vfs.Fake)
					fakeFs.OpenReturns(nil, errors.New("dummyError"))

					objectUnderTest := lockfile{
						fs:       fakeFs,
						psCanary: new(pscanary.Fake),
						os:       new(vos.Fake),
					}

					/* act */
					objectUnderTest.Lock(providedFilePath)

					/* assert */
					Expect(fakeFs.CreateArgsForCall(0)).To(Equal(providedFilePath))
				})

			})
			Context("fs.Open doesn't error", func() {
				Context("lockfile is invalid", func() {
					It("should create new lock file", func() {
						/* arrange */
						providedFilePath := "/dummy/pid.lock"

						fakeFs := new(vfs.Fake)

						objectUnderTest := lockfile{
							fs:       fakeFs,
							psCanary: new(pscanary.Fake),
							os:       new(vos.Fake),
						}

						/* act */
						objectUnderTest.Lock(providedFilePath)

						/* assert */
						Expect(fakeFs.CreateArgsForCall(0)).To(Equal(providedFilePath))
					})
				})
				Context("lockfile is valid", func() {
					Context("lock owner dead", func() {
						It("should create new lock file", func() {
							/* arrange */
							providedFilePath := "/dummy/pid.lock"

							fakeFs := new(vfs.Fake)

							fakePsCanary := new(pscanary.Fake)
							fakePsCanary.IsAliveReturns(false)

							objectUnderTest := lockfile{
								fs:       fakeFs,
								psCanary: fakePsCanary,
								os:       new(vos.Fake),
							}

							/* act */
							objectUnderTest.Lock(providedFilePath)

							/* assert */
							Expect(fakeFs.CreateArgsForCall(0)).To(Equal(providedFilePath))
						})
					})
					Context("lock owner alive", func() {
						It("should return expected error", func() {
							/* arrange */
							pIdFromFile := 1234
							expectedErr := fmt.Errorf("Unable to obtain lock; currently owned by PId: %v\n", pIdFromFile)

							fakeFs := new(vfs.Fake)
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

							fakeFs.OpenReturns(tempFile, nil)

							fakePsCanary := new(pscanary.Fake)
							fakePsCanary.IsAliveReturns(true)

							objectUnderTest := lockfile{
								fs:       fakeFs,
								psCanary: fakePsCanary,
								os:       new(vos.Fake),
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
