package coerce

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ipld/go-ipld-prime"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("ToNumber", func() {
	Context("Value is nil", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			actualValue, actualErr := ToNumber(nil)

			/* assert */
			Expect(*actualValue).To(Equal(ipld.Node{Number: new(float64)}))
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
			providedValue := ipld.Node{
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
		Context("os.ReadFile errs", func() {
			It("should return expected result", func() {
				/* arrange */
				/* act */
				actualValue, actualErr := ToNumber(
					ipld.Node{File: new(string)},
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(MatchError("unable to coerce file to number: open : no such file or directory"))
			})
		})
		Context("os.ReadFile doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				tmpFile, err := os.CreateTemp("", "")
				if err != nil {
					panic(err)
				}

				number := 2.0
				filePath := tmpFile.Name()
				err = os.WriteFile(filePath, []byte(fmt.Sprintf("%v", number)), 0777)
				if err != nil {
					panic(err)
				}

				/* act */
				actualValue, actualErr := ToNumber(
					ipld.Node{File: &filePath},
				)

				/* assert */
				Expect(*actualValue).To(Equal(ipld.Node{Number: &number}))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Number isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedNumber := float64(2.2)
			providedValue := ipld.Node{
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
			providedValue := ipld.Node{
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
			providedValue := ipld.Node{
				String: &providedString,
			}

			parsedNumber, err := strconv.ParseFloat(providedString, 64)
			if err != nil {
				panic(err.Error)
			}

			/* act */
			actualValue, actualErr := ToNumber(providedValue)

			/* assert */
			Expect(*actualValue).To(Equal(ipld.Node{Number: &parsedNumber}))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.Array,Value.Dir,File,Number,Object,Number nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := ipld.Node{}

			/* act */
			actualNumber, actualErr := ToNumber(providedValue)

			/* assert */
			Expect(actualNumber).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce '&{Array:<nil> Boolean:<nil> Dir:<nil> File:<nil> Number:<nil> Object:<nil> Socket:<nil> String:<nil>}' to number"))
		})
	})
})
