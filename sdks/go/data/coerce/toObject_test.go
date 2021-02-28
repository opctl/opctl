package coerce

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("ToObject", func() {
	Context("Value is nil", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			actualValue, actualErr := ToObject(nil)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.Array isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := &model.Value{
				Array: new([]interface{}),
			}

			/* act */
			actualValue, actualErr := ToObject(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce array to object; incompatible types")))
		})
	})
	Context("Value.Link isn't nil", func() {
		Context("is unresolvable", func() {
			It("should return expected result", func() {
				/* arrange */
				nonExistentPath := "nonExistent"

				/* act */
				actualValue, actualErr := ToObject(
					&model.Value{Link: &nonExistentPath},
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(errors.New("unable to coerce link to object; error was stat nonExistent: no such file or directory")))
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
				actualValue, actualErr := ToObject(providedValue)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce dir '%v' to object; incompatible types", tmpDir)))
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
					actualValue, actualErr := ToObject(
						&model.Value{Link: &filePath},
					)

					/* assert */
					Expect(actualValue).To(BeNil())
					Expect(actualErr).To(Equal(errors.New("unable to coerce file to object; error was unexpected end of JSON input")))
				})
			})
			Context("json.Unmarshal doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					mapEntryValue := "mapEntryValue"
					mapEntryKey := "mapKey"
					mapValueJSON := fmt.Sprintf(`{"%s": "%s"}`, mapEntryKey, mapEntryValue)

					tmpFile, err := ioutil.TempFile("", "")
					if nil != err {
						panic(err)
					}
					filePath := tmpFile.Name()

					err = ioutil.WriteFile(filePath, []byte(mapValueJSON), 0777)
					if nil != err {
						panic(err)
					}

					mapValue := map[string]interface{}{
						mapEntryKey: mapEntryValue,
					}
					expectedValue := model.Value{Object: &mapValue}

					/* act */
					actualValue, actualErr := ToObject(
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
			actualValue, actualErr := ToObject(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce number '%v' to object; incompatible types", providedNumber)))
		})
	})
	Context("Value.Object isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := &model.Value{
				Object: new(map[string]interface{}),
			}

			/* act */
			actualValue, actualErr := ToObject(providedValue)

			/* assert */
			Expect(actualValue).To(Equal(providedValue))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.String isn't nil", func() {
		Context("json.Unmarshal errs", func() {
			It("should return expected result", func() {
				/* arrange */

				/* act */
				actualValue, actualErr := ToObject(
					&model.Value{String: new(string)},
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(errors.New("unable to coerce string to object; error was unexpected end of JSON input")))
			})
		})
		Context("json.Unmarshal doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				mapEntryValue := "mapEntryValue"
				mapEntryKey := "mapKey"
				mapValueJSON := fmt.Sprintf(`{"%s": "%s"}`, mapEntryKey, mapEntryValue)

				mapValue := map[string]interface{}{
					mapEntryKey: mapEntryValue,
				}
				expectedValue := model.Value{Object: &mapValue}

				/* act */
				actualValue, actualErr := ToObject(
					&model.Value{String: &mapValueJSON},
				)

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Array,Link,Number,Object,String nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := &model.Value{}

			/* act */
			actualValue, actualErr := ToObject(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce '%+v' to object", providedValue)))
		})
	})
})
