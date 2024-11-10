package coerce

import (
	"os"

	"github.com/ipld/go-ipld-prime"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("ToBoolean", func() {
	Context("Value is nil", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			actualValue, actualErr := ToBoolean(nil)

			/* assert */
			Expect(*actualValue).To(Equal(ipld.Node{Boolean: new(bool)}))
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
					ipld.Node{
						Array: new([]interface{}),
					},
				)

				/* assert */
				Expect(*actualValue).To(Equal(ipld.Node{Boolean: &expectedBoolean}))
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
					ipld.Node{
						Array: array,
					},
				)

				/* assert */
				Expect(*actualValue).To(Equal(ipld.Node{Boolean: &expectedBoolean}))
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
				ipld.Node{Dir: new(string)},
			)

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce dir to boolean: incompatible types"))
		})
	})
	Context("Value.File isn't nil", func() {
		Context("os.ReadFile errs", func() {
			It("should return expected result", func() {
				/* arrange */
				/* act */
				actualValue, actualErr := ToBoolean(
					ipld.Node{File: new(string)},
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(MatchError("unable to coerce file to boolean: open : no such file or directory"))
			})
		})
		Context("os.ReadFile doesn't err", func() {
			Context("File content truthy", func() {
				It("should return expected result", func() {
					/* arrange */
					tmpFile, err := os.CreateTemp("", "")
					if err != nil {
						panic(err)
					}

					filePath := tmpFile.Name()
					err = os.WriteFile(filePath, []byte("true"), 0777)
					if err != nil {
						panic(err)
					}

					expectedBoolean := true

					/* act */
					actualValue, actualErr := ToBoolean(
						ipld.Node{File: &filePath},
					)

					/* assert */
					Expect(*actualValue).To(Equal(ipld.Node{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
			Context("File content falsy", func() {
				It("should return expected result", func() {
					/* arrange */
					tmpFile, err := os.CreateTemp("", "")
					if err != nil {
						panic(err)
					}

					filePath := tmpFile.Name()
					err = os.WriteFile(filePath, []byte("false"), 0777)
					if err != nil {
						panic(err)
					}

					expectedBoolean := false

					/* act */
					actualValue, actualErr := ToBoolean(
						ipld.Node{File: &filePath},
					)

					/* assert */
					Expect(*actualValue).To(Equal(ipld.Node{Boolean: &expectedBoolean}))
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
				providedValue := ipld.Node{
					Number: &providedNumber,
				}

				expectedBoolean := false

				/* act */
				actualValue, actualErr := ToBoolean(providedValue)

				/* assert */
				Expect(*actualValue).To(Equal(ipld.Node{Boolean: &expectedBoolean}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Number != 0", func() {
			It("should return expected result", func() {
				/* arrange */
				providedNumber := 1.0
				providedValue := ipld.Node{
					Number: &providedNumber,
				}

				expectedBoolean := true

				/* act */
				actualValue, actualErr := ToBoolean(providedValue)

				/* assert */
				Expect(*actualValue).To(Equal(ipld.Node{Boolean: &expectedBoolean}))
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
					ipld.Node{
						Object: new(map[string]interface{}),
					},
				)

				/* assert */
				Expect(*actualValue).To(Equal(ipld.Node{Boolean: &expectedBoolean}))
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
					ipld.Node{
						Object: object,
					},
				)

				/* assert */
				Expect(*actualValue).To(Equal(ipld.Node{Boolean: &expectedBoolean}))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Value.Array,Boolean,Dir,File,Number,Object,Boolean nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := ipld.Node{}

			/* act */
			actualBoolean, actualErr := ToBoolean(providedValue)

			/* assert */
			Expect(actualBoolean).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce '&{Array:<nil> Boolean:<nil> Dir:<nil> File:<nil> Number:<nil> Object:<nil> Socket:<nil> String:<nil>}' to boolean"))
		})
	})
})
