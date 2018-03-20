package param

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"io/ioutil"
)

var _ = Context("Validate", func() {
	Context("param.Dir not nil", func() {
		Context("value nil", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &model.Value{}
				providedParam := &model.Param{
					Dir: &model.DirParam{},
				}

				expectedErrors := []error{
					errors.New("dir required"),
				}

				objectUnderTest := NewValidator()

				/* act */
				actualErrors := objectUnderTest.Validate(
					providedValue,
					providedParam,
				)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("value not nil", func() {
			Context("value.Dir not nil", func() {
				It("should call fs.Stat w/ expected args", func() {

					/* arrange */
					providedValueDir := "dummyDir"
					providedValue := &model.Value{
						Dir: &providedValueDir,
					}
					providedParam := &model.Param{
						Dir: &model.DirParam{},
					}

					fakeOS := new(ios.Fake)
					// error to trigger immediate return
					fakeOS.StatReturns(nil, errors.New("dummyError"))

					objectUnderTest := _validator{
						os: fakeOS,
					}

					/* act */
					objectUnderTest.Validate(
						providedValue,
						providedParam,
					)

					/* assert */
					Expect(fakeOS.StatArgsForCall(0)).To(Equal(*providedValue.Dir))

				})
				Context("fs.Stat errors", func() {
					It("should return expected errors", func() {

						/* arrange */
						providedValueDir := "dummyDir"
						providedValue := &model.Value{
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

						objectUnderTest := _validator{
							os: fakeOS,
						}

						/* act */
						actualErrors := objectUnderTest.Validate(
							providedValue,
							providedParam,
						)

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

							providedValue := &model.Value{
								Dir: &tmpDirPath,
							}
							providedParam := &model.Param{
								Dir: &model.DirParam{},
							}

							expectedErrors := []error{}

							objectUnderTest := NewValidator()

							/* act */
							actualErrors := objectUnderTest.Validate(
								providedValue,
								providedParam,
							)

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

							providedValue := &model.Value{
								Dir: &tmpFilePath,
							}
							providedParam := &model.Param{
								Dir: &model.DirParam{},
							}

							expectedErrors := []error{
								fmt.Errorf("%v not a dir", tmpFilePath),
							}

							objectUnderTest := NewValidator()

							/* act */
							actualErrors := objectUnderTest.Validate(
								providedValue,
								providedParam,
							)

							/* assert */
							Expect(actualErrors).To(Equal(expectedErrors))

						})
					})
				})
			})
		})
	})
})
