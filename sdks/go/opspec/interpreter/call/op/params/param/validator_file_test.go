package param

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/types"
	"io/ioutil"
)

var _ = Context("Validate", func() {
	Context("param.File not nil", func() {
		Context("value.File nil", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &types.Value{}
				providedParam := &types.Param{
					File: &types.FileParam{},
				}

				expectedErrors := []error{
					errors.New("file required"),
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
		Context("value.File isn't empty", func() {
			It("should call fs.Stat w/ expected args", func() {

				/* arrange */
				providedValueFile := "dummyFile"
				providedValue := &types.Value{
					File: &providedValueFile,
				}
				providedParam := &types.Param{
					File: &types.FileParam{},
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
				Expect(fakeOS.StatArgsForCall(0)).To(Equal(*providedValue.File))

			})
			Context("fs.Stat errors", func() {
				It("should return expected errors", func() {

					/* arrange */
					providedValueFile := "dummyFile"
					providedValue := &types.Value{
						File: &providedValueFile,
					}
					providedParam := &types.Param{
						File: &types.FileParam{},
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
				Context("FileInfo.IsDir returns false", func() {
					It("should return no errors", func() {

						/* arrange */
						// no good way to fake fileinfo
						tmpFile, err := ioutil.TempFile("", "")
						if nil != err {
							panic(err)
						}

						tmpFilePath := tmpFile.Name()

						providedValue := &types.Value{
							File: &tmpFilePath,
						}
						providedParam := &types.Param{
							File: &types.FileParam{},
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
				Context("FileInfo.IsDir returns true", func() {
					It("should return expected errors", func() {

						/* arrange */
						// no good way to fake fileinfo
						tmpDirPath, err := ioutil.TempDir("", "")
						if nil != err {
							panic(err)
						}

						providedValue := &types.Value{
							File: &tmpDirPath,
						}
						providedParam := &types.Param{
							File: &types.FileParam{},
						}

						expectedErrors := []error{
							fmt.Errorf("%v not a file", tmpDirPath),
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
