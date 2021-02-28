package file

import (
	"errors"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Validate", func() {
	Context("value.Link nil", func() {
		It("should return expected errors", func() {

			/* arrange */
			providedValue := &model.Value{}

			expectedErrors := []error{
				errors.New("file required"),
			}

			/* act */
			actualErrors := Validate(
				providedValue,
			)

			/* assert */
			Expect(actualErrors).To(Equal(expectedErrors))

		})
	})
	Context("value.Link isn't empty", func() {
		Context("fs.Stat errors", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValueFile := "dummyFile"
				providedValue := &model.Value{
					Link: &providedValueFile,
				}

				/* act */
				actualErrors := Validate(
					providedValue,
				)

				/* assert */
				Expect(actualErrors[0].Error()).To(Equal("stat dummyFile: no such file or directory"))

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

					providedValue := &model.Value{
						Link: &tmpFilePath,
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := Validate(
						providedValue,
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

					providedValue := &model.Value{
						Link: &tmpDirPath,
					}

					expectedErrors := []error{
						fmt.Errorf("%v not a file", tmpDirPath),
					}

					/* act */
					actualErrors := Validate(
						providedValue,
					)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
	})
})
