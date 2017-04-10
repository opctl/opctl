package validate

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/virtual-go/fs"
	"io/ioutil"
)

var _ = Describe("Param", func() {
	Context("invoked w/ non-nil param.File", func() {
		Context("value.File is empty", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &model.Data{}
				providedParam := &model.Param{
					File: &model.FileParam{},
				}

				expectedErrors := []error{
					errors.New("File required"),
				}

				objectUnderTest := New()

				/* act */
				actualErrors := objectUnderTest.Param(providedValue, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("value nil", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedParam := &model.Param{
					File: &model.FileParam{},
				}

				expectedErrors := []error{
					errors.New("File required"),
				}

				objectUnderTest := New()

				/* act */
				actualErrors := objectUnderTest.Param(nil, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("value.File isn't empty", func() {
			It("should call fs.Stat w/ expected args", func() {

				/* arrange */
				providedValue := &model.Data{
					File: "dummyValue",
				}
				providedParam := &model.Param{
					File: &model.FileParam{},
				}

				fakeFS := new(fs.Fake)
				// error to trigger immediate return
				fakeFS.StatReturns(nil, errors.New("dummyError"))

				objectUnderTest := validate{
					fs: fakeFS,
				}

				/* act */
				objectUnderTest.Param(providedValue, providedParam)

				/* assert */
				Expect(fakeFS.StatArgsForCall(0)).To(Equal(providedValue.File))

			})
			Context("fs.Stat errors", func() {
				It("should return expected errors", func() {

					/* arrange */
					providedValue := &model.Data{
						File: "dummyValue",
					}
					providedParam := &model.Param{
						File: &model.FileParam{},
					}

					expectedErrors := []error{
						errors.New("dummyError"),
					}

					fakeFs := new(fs.Fake)
					fakeFs.StatReturns(nil, expectedErrors[0])

					objectUnderTest := validate{
						fs: fakeFs,
					}

					/* act */
					actualErrors := objectUnderTest.Param(providedValue, providedParam)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})

			})
			Context("fs.Stat doesn't error", func() {
				Context("FileInfo.IsDir returns false", func() {
					It("should return no errors", func() {

						/* arrange */
						// no good way to fake fileinfo
						tmpFile, err := ioutil.TempFile("", "")
						if nil != err {
							panic(err)
						}

						tmpFilePath := tmpFile.Name()

						providedValue := &model.Data{
							File: tmpFilePath,
						}
						providedParam := &model.Param{
							File: &model.FileParam{},
						}

						expectedErrors := []error{}

						objectUnderTest := New()

						/* act */
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("FileInfo.IsDir returns true", func() {
					It("should return expected errors", func() {

						/* arrange */
						// no good way to fake fileinfo
						tmpDirPath, err := ioutil.TempDir("", "")
						if nil != err {
							panic(err)
						}

						providedValue := &model.Data{
							File: tmpDirPath,
						}
						providedParam := &model.Param{
							File: &model.FileParam{},
						}

						expectedErrors := []error{
							fmt.Errorf("%v not a file", tmpDirPath),
						}

						objectUnderTest := New()

						/* act */
						actualErrors := objectUnderTest.Param(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
		})
	})

})
