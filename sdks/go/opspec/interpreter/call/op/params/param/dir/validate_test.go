package dir

import (
	"errors"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Validate", func() {
	Context("value nil", func() {
		It("should return expected errors", func() {

			/* arrange */
			providedValue := &model.Value{}

			expectedErrors := []error{
				errors.New("dir required"),
			}

			/* act */
			actualErrors := Validate(
				providedValue,
			)

			/* assert */
			Expect(actualErrors).To(Equal(expectedErrors))

		})
	})
	Context("value not nil", func() {
		Context("value.Dir not nil", func() {
			Context("fs.Stat errors", func() {
				It("should return expected errors", func() {

					/* arrange */
					providedValueDir := "dummyDir"
					providedValue := &model.Value{
						Dir: &providedValueDir,
					}

					/* act */
					actualErrors := Validate(
						providedValue,
					)

					/* assert */
					Expect(actualErrors[0].Error()).To(Equal("stat dummyDir: no such file or directory"))

				})

			})
			Context("fs.Stat doesn't error", func() {
				Context("FileInfo.IsDir returns true", func() {
					It("should return no errors", func() {

						/* arrange */
						// no good way to fake fileinfo
						tmpDirPath, err := ioutil.TempDir("", "")
						if err != nil {
							panic(err)
						}

						providedValue := &model.Value{
							Dir: &tmpDirPath,
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
				Context("FileInfo.IsDir returns false", func() {
					It("should return expected errors", func() {

						/* arrange */
						// no good way to fake fileinfo
						tmpFile, err := ioutil.TempFile("", "")
						if err != nil {
							panic(err)
						}

						tmpFilePath := tmpFile.Name()

						providedValue := &model.Value{
							Dir: &tmpFilePath,
						}

						expectedErrors := []error{
							fmt.Errorf("%v not a dir", tmpFilePath),
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
})
