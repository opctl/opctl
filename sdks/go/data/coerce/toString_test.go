package coerce

import (
	"io/ioutil"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("ToString", func() {
	Context("Value is nil", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			actualValue, actualErr := ToString(nil)

			/* assert */
			Expect(*actualValue).To(Equal(model.Value{String: new(string)}))
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
					&model.Value{Array: new([]interface{})},
				)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{String: &jsonArray}))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Boolean isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedBoolean := true
			providedValue := &model.Value{
				Boolean: &providedBoolean,
			}

			booleanString := strconv.FormatBool(providedBoolean)
			expectedValue := model.Value{String: &booleanString}

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
			providedValue := &model.Value{
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
		Context("ioutil.ReadFile errs", func() {
			It("should return expected result", func() {
				/* arrange */
				/* act */
				actualValue, actualErr := ToString(
					&model.Value{File: new(string)},
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(MatchError("unable to coerce file to string: open : no such file or directory"))
			})
		})
		Context("ioutil.ReadFile doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				tmpFile, err := ioutil.TempFile("", "")
				if nil != err {
					panic(err)
				}
				filePath := tmpFile.Name()

				expectedString := "expectedString"

				err = ioutil.WriteFile(filePath, []byte(expectedString), 0777)
				if nil != err {
					panic(err)
				}

				expectedValue := model.Value{String: &expectedString}

				/* act */
				actualValue, actualErr := ToString(
					&model.Value{File: &filePath},
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
			providedValue := &model.Value{
				Number: &providedNumber,
			}

			numberString := strconv.FormatFloat(providedNumber, 'f', -1, 64)
			expectedValue := model.Value{String: &numberString}

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
				expectedValue := model.Value{String: &marshaledString}

				/* act */
				actualValue, actualErr := ToString(
					&model.Value{Object: new(map[string]interface{})},
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
			providedValue := model.Value{
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
			providedValue := &model.Value{}

			/* act */
			actualValue, actualErr := ToString(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce '&{Array:<nil> Boolean:<nil> Dir:<nil> File:<nil> Number:<nil> Object:<nil> Socket:<nil> String:<nil>}' to string"))
		})
	})
})
