package coerce

import (
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("ToDir", func() {
	Context("Value is nil", func() {
		It("should return expected result", func() {
			/* arrange */

			/* act */
			actualValue, actualErr := ToDir(nil, "scratchDir")

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce null to dir"))
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
			actualValue, actualErr := ToDir(providedValue, "scratchDir")

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce array to dir: incompatible types"))
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
			actualValue, actualErr := ToDir(providedValue, "scratchDir")

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce boolean to dir: incompatible types"))
		})
	})
	Context("Value.Dir isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := &model.Value{
				Dir: new(string),
			}

			/* act */
			actualValue, actualErr := ToDir(providedValue, "scratchDir")

			/* assert */
			Expect(actualValue).To(Equal(providedValue))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Value.File isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := &model.Value{
				File: new(string),
			}

			/* act */
			actualValue, actualErr := ToDir(providedValue, "scratchDir")

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce file to dir: incompatible types"))
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
			actualValue, actualErr := ToDir(providedValue, "scratchDir")

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce number to dir: incompatible types"))
		})
	})
	Context("Value.Object isn't nil", func() {
		Context("Value.Object has descendant file", func() {
			It("should return expected result", func() {
				/* arrange */
				providedChildDirPath := "/childDir"
				providedObject := map[string]interface{}{
					providedChildDirPath: map[string]interface{}{
						"/file": map[string]interface{}{
							"data": "blah",
						},
					},
				}

				providedScratchDir, err := ioutil.TempDir("", "")
				if nil != err {
					panic(err)
				}

				providedValue := &model.Value{
					Object: &providedObject,
				}

				/* act */
				actualValue, actualErr := ToDir(providedValue, providedScratchDir)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(strings.HasPrefix(*actualValue.Dir, providedScratchDir)).To(BeTrue())

			})
		})
		Context("Value.Object has descendant dir", func() {
			It("should return expected result", func() {
				/* arrange */
				providedChildDirPath := "/childDir"
				providedObject := map[string]interface{}{
					providedChildDirPath: map[string]interface{}{},
				}

				providedScratchDir, err := ioutil.TempDir("", "")
				if nil != err {
					panic(err)
				}

				providedValue := &model.Value{
					Object: &providedObject,
				}

				/* act */
				actualValue, actualErr := ToDir(providedValue, providedScratchDir)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(strings.HasPrefix(*actualValue.Dir, providedScratchDir)).To(BeTrue())

			})
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
			actualValue, actualErr := ToDir(providedValue, "scratchDir")

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce socket to dir: incompatible types"))
		})
	})
	Context("Value.String isn't nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := &model.Value{
				String: new(string),
			}

			/* act */
			actualValue, actualErr := ToDir(providedValue, "scratchDir")

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce string to dir: incompatible types"))
		})
	})
	Context("Value.Array,Dir,File,Number,Dir,Socket,String nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValue := &model.Value{}

			/* act */
			actualValue, actualErr := ToDir(providedValue, "scratchDir")

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualErr).To(MatchError("unable to coerce '&{Array:<nil> Boolean:<nil> Dir:<nil> File:<nil> Number:<nil> Object:<nil> Socket:<nil> String:<nil>}' to dir"))
		})
	})
})
