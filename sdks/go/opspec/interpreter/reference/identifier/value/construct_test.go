package value

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Construct", func() {
	Context("data is bool", func() {
		It("should return expected result", func() {
			/* arrange */
			providedData := true

			/* act */
			actualValue, actualErr := Construct(providedData)

			/* assert */
			Expect(*actualValue).To(Equal(model.Value{Boolean: &providedData}))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("data is float64", func() {
		It("should return expected result", func() {
			/* arrange */
			providedData := 2.2

			/* act */
			actualValue, actualErr := Construct(providedData)

			/* assert */
			Expect(*actualValue).To(Equal(model.Value{Number: &providedData}))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("data is int", func() {
		It("should return expected result", func() {
			/* arrange */
			providedData := 22
			providedDataAsFloat64 := float64(providedData)

			/* act */
			actualValue, actualErr := Construct(providedData)

			/* assert */
			Expect(*actualValue).To(Equal(model.Value{Number: &providedDataAsFloat64}))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("data is string", func() {
		It("should return expected result", func() {
			/* arrange */
			providedData := "dummyData"

			/* act */
			actualValue, actualErr := Construct(providedData)

			/* assert */
			Expect(*actualValue).To(Equal(model.Value{String: &providedData}))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("data is map[string]interface{}", func() {
		It("should return expected result", func() {
			/* arrange */
			providedData := map[string]interface{}{"hello": 2}

			/* act */
			actualValue, actualErr := Construct(providedData)

			/* assert */
			Expect(*actualValue).To(Equal(model.Value{Object: &providedData}))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("data is []interface{}", func() {
		It("should return expected result", func() {
			/* arrange */
			providedData := []interface{}{"hello"}

			/* act */
			actualValue, actualErr := Construct(providedData)

			/* assert */
			Expect(*actualValue).To(Equal(model.Value{Array: &providedData}))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("data is unexpected type", func() {
		It("should return expected result", func() {
			/* act */
			_, actualErr := Construct(nil)

			/* assert */
			Expect(actualErr).To(MatchError("unable to construct value: '<nil>' unexpected type"))
		})
	})
})
