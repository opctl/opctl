package paramvalidator

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"io/ioutil"
)

var _ = Describe("Validate", func() {
	Context("invoked w/ non-nil param.Dir", func() {
		Context("value.Dir is empty", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &model.Data{}
				providedParam := &model.Param{
					Dir: &model.DirParam{},
				}

				expectedErrors := []error{
					errors.New("Dir required"),
				}

				objectUnderTest := New()

				/* act */
				actualErrors := objectUnderTest.Validate(providedValue, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("value nil", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedParam := &model.Param{
					Dir: &model.DirParam{},
				}

				expectedErrors := []error{
					errors.New("Dir required"),
				}

				objectUnderTest := New()

				/* act */
				actualErrors := objectUnderTest.Validate(nil, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("value.Dir isn't empty", func() {
			It("should call fs.Stat w/ expected args", func() {

				/* arrange */
				providedValueDir := "dummyDir"
				providedValue := &model.Data{
					Dir: &providedValueDir,
				}
				providedParam := &model.Param{
					Dir: &model.DirParam{},
				}

				fakeOS := new(ios.Fake)
				// error to trigger immediate return
				fakeOS.StatReturns(nil, errors.New("dummyError"))

				objectUnderTest := _ParamValidator{
					os: fakeOS,
				}

				/* act */
				objectUnderTest.Validate(providedValue, providedParam)

				/* assert */
				Expect(fakeOS.StatArgsForCall(0)).To(Equal(*providedValue.Dir))

			})
			Context("fs.Stat errors", func() {
				It("should return expected errors", func() {

					/* arrange */
					providedValueDir := "dummyDir"
					providedValue := &model.Data{
						Dir: &providedValueDir,
					}
					providedParam := &model.Param{
						Dir: &model.DirParam{},
					}

					expectedErrors := []error{
						errors.New("dummyError"),
					}

					fakeOS := new(ios.Fake)
					fakeOS.StatReturns(nil, expectedErrors[0])

					objectUnderTest := _ParamValidator{
						os: fakeOS,
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedParam)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})

			})
			Context("fs.Stat doesn't error", func() {
				Context("FileInfo.IsDir returns true", func() {
					It("should return no errors", func() {

						/* arrange */
						// no good way to fake fileinfo
						tmpDirPath, err := ioutil.TempDir("", "")
						if nil != err {
							panic(err)
						}

						providedValue := &model.Data{
							Dir: &tmpDirPath,
						}
						providedParam := &model.Param{
							Dir: &model.DirParam{},
						}

						expectedErrors := []error{}

						objectUnderTest := New()

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("FileInfo.IsDir returns false", func() {
					It("should return expected errors", func() {

						/* arrange */
						// no good way to fake fileinfo
						tmpFile, err := ioutil.TempFile("", "")
						if nil != err {
							panic(err)
						}

						tmpFilePath := tmpFile.Name()

						providedValue := &model.Data{
							Dir: &tmpFilePath,
						}
						providedParam := &model.Param{
							Dir: &model.DirParam{},
						}

						expectedErrors := []error{
							fmt.Errorf("%v not a dir", tmpFilePath),
						}

						objectUnderTest := New()

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
		})
	})

})
