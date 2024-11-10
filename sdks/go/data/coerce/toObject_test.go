package coerce

import (
	"fmt"
	"os"

	"github.com/ipld/go-ipld-prime"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			providedValue := ipld.Node{
				Array: new([]interface{}),
			}

			/* act */
			actualValue, actualErr := ToObject(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce array to object: incompatible types"))
		})
	})
	Context("Value.Dir isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedDir := "dummyValue"
			providedValue := ipld.Node{
				Dir: &providedDir,
			}

			/* act */
			actualValue, actualErr := ToObject(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce dir 'dummyValue' to object: incompatible types"))
		})
	})
	Context("Value.File isn't nil", func() {
		Context("os.ReadFile errs", func() {
			It("should return expected result", func() {
				/* arrange */
				nonExistentPath := "nonExistent"

				/* act */
				actualValue, actualErr := ToObject(
					ipld.Node{File: &nonExistentPath},
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(MatchError("unable to coerce file to object: open nonExistent: no such file or directory"))
			})
		})
		Context("os.ReadFile doesn't err", func() {
			Context("json.Unmarshal errs", func() {
				It("should return expected result", func() {
					/* arrange */

					tmpFile, err := os.CreateTemp("", "")
					if err != nil {
						panic(err)
					}

					filePath := tmpFile.Name()

					/* act */
					actualValue, actualErr := ToObject(
						ipld.Node{File: &filePath},
					)

					/* assert */
					Expect(actualValue).To(BeNil())
					Expect(actualErr).To(MatchError("unable to coerce file to object: unexpected end of JSON input"))
				})
			})
			Context("json.Unmarshal doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					mapEntryValue := "mapEntryValue"
					mapEntryKey := "mapKey"
					mapValueJSON := fmt.Sprintf(`{"%s": "%s"}`, mapEntryKey, mapEntryValue)

					tmpFile, err := os.CreateTemp("", "")
					if err != nil {
						panic(err)
					}
					filePath := tmpFile.Name()

					err = os.WriteFile(filePath, []byte(mapValueJSON), 0777)
					if err != nil {
						panic(err)
					}

					mapValue := map[string]interface{}{
						mapEntryKey: mapEntryValue,
					}
					expectedValue := ipld.Node{Object: &mapValue}

					/* act */
					actualValue, actualErr := ToObject(
						ipld.Node{File: &filePath},
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
			providedValue := ipld.Node{
				Number: &providedNumber,
			}

			/* act */
			actualValue, actualErr := ToObject(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce number '2.2' to object: incompatible types"))
		})
	})
	Context("Value.Object isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := ipld.Node{
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
					ipld.Node{String: new(string)},
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(MatchError("unable to coerce string to object: unexpected end of JSON input"))
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
				expectedValue := ipld.Node{Object: &mapValue}

				/* act */
				actualValue, actualErr := ToObject(
					ipld.Node{String: &mapValueJSON},
				)

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Array,Value.Dir,File,Number,Object,String nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := ipld.Node{}

			/* act */
			actualValue, actualErr := ToObject(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce '&{Array:<nil> Boolean:<nil> Dir:<nil> File:<nil> Number:<nil> Object:<nil> Socket:<nil> String:<nil>}' to object"))
		})
	})
})
