package coerce

import (
	"os"
	"strconv"

	"github.com/ipld/go-ipld-prime"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("ToString", func() {
	Context("Value is nil", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			actualValue, actualErr := ToString(nil)

			/* assert */
			Expect(*actualValue).To(Equal(ipld.Node{String: new(string)}))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.Array isn't nil", func() {
		Context("json.Marshal doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				jsonArray := "[]"
				/* act */
				actualValue, actualErr := ToString(
					ipld.Node{Array: new([]interface{})},
				)

				/* assert */
				Expect(*actualValue).To(Equal(ipld.Node{String: &jsonArray}))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Boolean isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedBoolean := true
			providedValue := ipld.Node{
				Boolean: &providedBoolean,
			}

			booleanString := strconv.FormatBool(providedBoolean)
			expectedValue := ipld.Node{String: &booleanString}

			/* act */
			actualValue, actualErr := ToString(providedValue)

			/* assert */
			Expect(*actualValue).To(Equal(expectedValue))
			Expect(actualErr).To(BeNil())
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
			actualValue, actualErr := ToString(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce dir 'dummyValue' to string: incompatible types"))
		})
	})
	Context("Value.File isn't nil", func() {
		Context("os.ReadFile errs", func() {
			It("should return expected result", func() {
				/* arrange */
				/* act */
				actualValue, actualErr := ToString(
					ipld.Node{File: new(string)},
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(MatchError("unable to coerce file to string: open : no such file or directory"))
			})
		})
		Context("os.ReadFile doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				tmpFile, err := os.CreateTemp("", "")
				if err != nil {
					panic(err)
				}
				filePath := tmpFile.Name()

				expectedString := "expectedString"

				err = os.WriteFile(filePath, []byte(expectedString), 0777)
				if err != nil {
					panic(err)
				}

				expectedValue := ipld.Node{String: &expectedString}

				/* act */
				actualValue, actualErr := ToString(
					ipld.Node{File: &filePath},
				)

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
				Expect(actualErr).To(BeNil())
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

			numberString := strconv.FormatFloat(providedNumber, 'f', -1, 64)
			expectedValue := ipld.Node{String: &numberString}

			/* act */
			actualValue, actualErr := ToString(providedValue)

			/* assert */
			Expect(*actualValue).To(Equal(expectedValue))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.Object isn't nil", func() {
		Context("json.Marshal doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */

				marshaledString := string("{}")
				expectedValue := ipld.Node{String: &marshaledString}

				/* act */
				actualValue, actualErr := ToString(
					ipld.Node{Object: new(map[string]interface{})},
				)

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.String isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedString := "dummyValue"
			providedValue := ipld.Node{
				String: &providedString,
			}

			/* act */
			actualValue, actualErr := ToString(&providedValue)

			/* assert */
			Expect(*actualValue).To(Equal(providedValue))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.Dir,File,Number,Object,String nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := ipld.Node{}

			/* act */
			actualValue, actualErr := ToString(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce '&{Array:<nil> Boolean:<nil> Dir:<nil> File:<nil> Number:<nil> Object:<nil> Socket:<nil> String:<nil>}' to string"))
		})
	})
})
