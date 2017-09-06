package file

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Context("Validate", func() {
	Context("value.File nil", func() {
		It("should return expected errors", func() {

			/* arrange */
			expectedErrors := []error{
				errors.New("file required"),
			}

			objectUnderTest := newValidator()

			/* act */
			actualErrors := objectUnderTest.Validate(
				nil,
			)

			/* assert */
			Expect(actualErrors).To(Equal(expectedErrors))

		})
	})
	Context("value isn't nil", func() {
		It("should call fs.Stat w/ expected args", func() {

			/* arrange */
			providedValue := "dummyFile"

			fakeOS := new(ios.Fake)
			// error to trigger immediate return
			fakeOS.StatReturns(nil, errors.New("dummyError"))

			objectUnderTest := _validator{
				os: fakeOS,
			}

			/* act */
			objectUnderTest.Validate(
				&providedValue,
			)

			/* assert */
			Expect(fakeOS.StatArgsForCall(0)).To(Equal(providedValue))

		})
		Context("fs.Stat errors", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := "dummyFile"

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
					&providedValue,
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

					providedValue := tmpFile.Name()

					expectedErrors := []error{}

					objectUnderTest := newValidator()

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("FileInfo.IsDir returns true", func() {
				It("should return expected errors", func() {

					/* arrange */
					// no good way to fake fileinfo
					providedValue, err := ioutil.TempDir("", "")
					if nil != err {
						panic(err)
					}

					expectedErrors := []error{
						fmt.Errorf("%v not a file", providedValue),
					}

					objectUnderTest := newValidator()

					/* act */
					actualErrors := objectUnderTest.Validate(
						&providedValue,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
	})
})
