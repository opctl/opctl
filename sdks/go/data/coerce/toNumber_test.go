package coerce

import (
	"fmt"
	"io/ioutil"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("ToNumber", func() {
	Context("Value is nil", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			actualValue, actualErr := ToNumber(nil)

			/* assert */
			Expect(*actualValue).To(Equal(model.Value{Number: new(float64)}))
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
			actualValue, actualErr := ToNumber(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce array to number: incompatible types"))
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
			actualValue, actualErr := ToNumber(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce dir to number: incompatible types"))
		})
	})
	Context("Value.File isn't nil", func() {
		Context("ioutil.ReadFile errs", func() {
			It("should return expected result", func() {
				/* arrange */
				/* act */
				actualValue, actualErr := ToNumber(
					&model.Value{File: new(string)},
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(MatchError("unable to coerce file to number: open : no such file or directory"))
			})
		})
		Context("ioutil.ReadFile doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				tmpFile, err := ioutil.TempFile("", "")
				if nil != err {
					panic(err)
				}

				number := 2.0
				filePath := tmpFile.Name()
				err = ioutil.WriteFile(filePath, []byte(fmt.Sprintf("%v", number)), 0777)
				if nil != err {
					panic(err)
				}

				/* act */
				actualValue, actualErr := ToNumber(
					&model.Value{File: &filePath},
				)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Number: &number}))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Number isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedNumber := float64(2.2)
			providedValue := &model.Value{
				Number: &providedNumber,
			}

			/* act */
			actualValue, actualErr := ToNumber(providedValue)

			/* assert */
			Expect(actualValue).To(Equal(providedValue))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.Object isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := &model.Value{
				Object: new(map[string]interface{}),
			}

			/* act */
			actualValue, actualErr := ToNumber(providedValue)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce object to number: incompatible types"))
		})
	})
	Context("Value.String isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedString := "2.2"
			providedValue := &model.Value{
				String: &providedString,
			}

			parsedNumber, err := strconv.ParseFloat(providedString, 64)
			if nil != err {
				panic(err.Error)
			}

			/* act */
			actualValue, actualErr := ToNumber(providedValue)

			/* assert */
			Expect(*actualValue).To(Equal(model.Value{Number: &parsedNumber}))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.Array,Value.Dir,File,Number,Object,Number nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := &model.Value{}

			/* act */
			actualNumber, actualErr := ToNumber(providedValue)

			/* assert */
			Expect(actualNumber).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce '&{Array:<nil> Boolean:<nil> Dir:<nil> File:<nil> Number:<nil> Object:<nil> Socket:<nil> String:<nil>}' to number"))
		})
	})
})
