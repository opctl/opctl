package coerce

import (
	"errors"
	"fmt"
	"os"

	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("ToArray", func() {
	Context("Value is nil", func() {
		It("should return expected result", func() {
			/* arrange */

			/* act */
			actualValue, actualErr := ToArray(nil)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(Equal(errors.New("unable to coerce null to array")))
		})
	})
	Context("Value.Array isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			array := &[]interface{}{"dummyItem"}
			providedValue := &model.Value{
				Array: array,
			}

			/* act */
			actualValue, actualErr := ToArray(providedValue)

			/* assert */
			Expect(actualValue).To(Equal(providedValue))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.Boolean isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedBoolean := true
			providedValue := &model.Value{
				Boolean: &providedBoolean,
			}

			/* act */
			actualValue, actualErr := ToArray(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce boolean to array; incompatible types")))
		})
	})
	Context("Value.Link isn't nil", func() {
		Context("is unresolvable", func() {
			It("should return expected result", func() {
				/* arrange */
				nonExistentPath := "nonExistent"

				/* act */
				actualValue, actualErr := ToArray(
					&model.Value{Link: &nonExistentPath},
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(errors.New("unable to coerce link to array; error was stat nonExistent: no such file or directory")))
			})
		})
		Context("is dir", func() {
			It("should return expected result", func() {
				/* arrange */
				tmpDir := os.TempDir()
				providedValue := &model.Value{
					Link: &tmpDir,
				}

				/* act */
				actualValue, actualErr := ToArray(providedValue)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce dir to array; incompatible types")))
			})
		})
		Context("ioutil.ReadFile doesn't err", func() {
			Context("json.Unmarshal errs", func() {
				It("should return expected result", func() {
					/* arrange */

					tmpFile, err := ioutil.TempFile("", "")
					if nil != err {
						panic(err)
					}

					filePath := tmpFile.Name()

					/* act */
					actualValue, actualErr := ToArray(
						&model.Value{Link: &filePath},
					)

					/* assert */
					Expect(actualValue).To(BeNil())
					Expect(actualErr).To(Equal(errors.New("unable to coerce file to array; error was unexpected end of JSON input")))
				})
			})
			Context("json.Unmarshal doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					arrayItem := "arrayItem"
					arrayJSON := fmt.Sprintf(`["%s"]`, arrayItem)

					tmpFile, err := ioutil.TempFile("", "")
					if nil != err {
						panic(err)
					}
					filePath := tmpFile.Name()

					err = ioutil.WriteFile(filePath, []byte(arrayJSON), 0777)
					if nil != err {
						panic(err)
					}

					array := []interface{}{arrayItem}
					expectedValue := model.Value{Array: &array}

					/* act */
					actualValue, actualErr := ToArray(
						&model.Value{Link: &filePath},
					)

					/* assert */
					Expect(actualErr).To(BeNil())
					Expect(*actualValue).To(Equal(expectedValue))
				})
			})
		})
	})
	Context("Value.Number isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedNumber := 2.2
			providedValue := &model.Value{
				Number: &providedNumber,
			}

			/* act */
			actualValue, actualErr := ToArray(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce number to array; incompatible types")))
		})
	})
	Context("Value.Socket isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedSocket := "dummySocket"
			providedValue := &model.Value{
				Socket: &providedSocket,
			}

			/* act */
			actualValue, actualErr := ToArray(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce socket to array; incompatible types")))
		})
	})
	Context("Value.String isn't nil", func() {
		Context("json.Unmarshal errs", func() {
			It("should return expected result", func() {
				/* arrange */

				/* act */
				actualValue, actualErr := ToArray(
					&model.Value{String: new(string)},
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(errors.New("unable to coerce string to array; error was unexpected end of JSON input")))
			})
		})
		Context("json.Unmarshal doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */

				arrayItem := "arrayItem"
				arrayJSON := fmt.Sprintf(`["%s"]`, arrayItem)
				array := &[]interface{}{arrayItem}
				expectedValue := model.Value{Array: array}

				/* act */
				actualValue, actualErr := ToArray(
					&model.Value{String: &arrayJSON},
				)

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Link,Number,Array,String nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := &model.Value{}

			/* act */
			actualValue, actualErr := ToArray(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce '%+v' to array", providedValue)))
		})
	})
})
