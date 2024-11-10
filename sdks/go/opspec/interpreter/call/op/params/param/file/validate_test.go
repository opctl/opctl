package file

import (
	"errors"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Validate", func() {
	Context("value.File nil", func() {
		It("should return expected errors", func() {

			/* arrange */
			providedValue := &ipld.Node{}

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
	Context("value.File isn't empty", func() {
		Context("fs.Stat errors", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValueFile := "dummyFile"
				providedValue := &ipld.Node{
					File: &providedValueFile,
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
					tmpFile, err := os.CreateTemp("", "")
					if err != nil {
						panic(err)
					}

					tmpFilePath := tmpFile.Name()

					providedValue := &ipld.Node{
						File: &tmpFilePath,
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
					tmpDirPath, err := os.MkdirTemp("", "")
					if err != nil {
						panic(err)
					}

					providedValue := &ipld.Node{
						File: &tmpDirPath,
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
