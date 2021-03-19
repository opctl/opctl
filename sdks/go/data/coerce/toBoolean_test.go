package coerce

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("ToBoolean", func() {
	Context("Value is nil", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			actualValue, actualErr := ToBoolean(nil)

			/* assert */
			Expect(*actualValue).To(Equal(model.Value{Boolean: new(bool)}))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.Array isn't nil", func() {
		Context("Array empty", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedBoolean := false

				/* act */
				actualValue, actualErr := ToBoolean(
					&model.Value{
						Array: new([]interface{}),
					},
				)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Array not empty", func() {
			It("should return expected result", func() {
				/* arrange */
				array := &[]interface{}{
					"",
				}

				expectedBoolean := true

				/* act */
				actualValue, actualErr := ToBoolean(
					&model.Value{
						Array: array,
					},
				)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
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

			/* act */
			actualValue, actualErr := ToBoolean(providedValue)

			/* assert */
			Expect(actualValue).To(Equal(providedValue))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.Dir isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */

			/* act */
			actualValue, actualErr := ToBoolean(
				&model.Value{Dir: new(string)},
			)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce dir to boolean: incompatible types"))
		})
	})
	Context("Value.File isn't nil", func() {
		Context("ioutil.ReadFile errs", func() {
			It("should return expected result", func() {
				/* arrange */
				/* act */
				actualValue, actualErr := ToBoolean(
					&model.Value{File: new(string)},
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(MatchError("unable to coerce file to boolean: open : no such file or directory"))
			})
		})
		Context("ioutil.ReadFile doesn't err", func() {
			Context("File content truthy", func() {
				It("should return expected result", func() {
					/* arrange */
					tmpFile, err := ioutil.TempFile("", "")
					if err != nil {
						panic(err)
					}

					filePath := tmpFile.Name()
					err = ioutil.WriteFile(filePath, []byte("true"), 0777)
					if err != nil {
						panic(err)
					}

					expectedBoolean := true

					/* act */
					actualValue, actualErr := ToBoolean(
						&model.Value{File: &filePath},
					)

					/* assert */
					Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
			Context("File content falsy", func() {
				It("should return expected result", func() {
					/* arrange */
					tmpFile, err := ioutil.TempFile("", "")
					if err != nil {
						panic(err)
					}

					filePath := tmpFile.Name()
					err = ioutil.WriteFile(filePath, []byte("false"), 0777)
					if err != nil {
						panic(err)
					}

					expectedBoolean := false

					/* act */
					actualValue, actualErr := ToBoolean(
						&model.Value{File: &filePath},
					)

					/* assert */
					Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
	Context("Value.Number isn't nil", func() {
		Context("Number == 0", func() {
			It("should return expected result", func() {
				/* arrange */
				providedNumber := 0.0
				providedValue := &model.Value{
					Number: &providedNumber,
				}

				expectedBoolean := false

				/* act */
				actualValue, actualErr := ToBoolean(providedValue)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Number != 0", func() {
			It("should return expected result", func() {
				/* arrange */
				providedNumber := 1.0
				providedValue := &model.Value{
					Number: &providedNumber,
				}

				expectedBoolean := true

				/* act */
				actualValue, actualErr := ToBoolean(providedValue)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Object isn't nil", func() {
		Context("Object has no properties", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedBoolean := false

				/* act */
				actualValue, actualErr := ToBoolean(
					&model.Value{
						Object: new(map[string]interface{}),
					},
				)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Object has properties", func() {
			It("should return expected result", func() {
				/* arrange */
				object := &map[string]interface{}{
					"dummyProp": nil,
				}

				expectedBoolean := true

				/* act */
				actualValue, actualErr := ToBoolean(
					&model.Value{
						Object: object,
					},
				)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Array,Boolean,Dir,File,Number,Object,Boolean nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := &model.Value{}

			/* act */
			actualBoolean, actualErr := ToBoolean(providedValue)

			/* assert */
			Expect(actualBoolean).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce '&{Array:<nil> Boolean:<nil> Dir:<nil> File:<nil> Number:<nil> Object:<nil> Socket:<nil> String:<nil>}' to boolean"))
		})
	})
})
