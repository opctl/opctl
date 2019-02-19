package value

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
)

var _ = Context("Constructor", func() {
	Context("NewConstructor", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewConstructor()).Should(Not(BeNil()))
		})
	})
	Context("Construct", func() {
		Context("data is bool", func() {
			It("should return expected result", func() {
				/* arrange */
				providedData := true
				objectUnderTest := _constructor{}

				/* act */
				actualValue, actualErr := objectUnderTest.Construct(providedData)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Boolean: &providedData}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("data is float64", func() {
			It("should return expected result", func() {
				/* arrange */
				providedData := 2.2
				objectUnderTest := _constructor{}

				/* act */
				actualValue, actualErr := objectUnderTest.Construct(providedData)

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
				objectUnderTest := _constructor{}

				/* act */
				actualValue, actualErr := objectUnderTest.Construct(providedData)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Number: &providedDataAsFloat64}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("data is string", func() {
			It("should return expected result", func() {
				/* arrange */
				providedData := "dummyData"
				objectUnderTest := _constructor{}

				/* act */
				actualValue, actualErr := objectUnderTest.Construct(providedData)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{String: &providedData}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("data is map[string]interface{}", func() {
			It("should return expected result", func() {
				/* arrange */
				providedData := map[string]interface{}{"hello": 2}
				objectUnderTest := _constructor{}

				/* act */
				actualValue, actualErr := objectUnderTest.Construct(providedData)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Object: providedData}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("data is []interface{}", func() {
			It("should return expected result", func() {
				/* arrange */
				providedData := []interface{}{"hello"}
				objectUnderTest := _constructor{}

				/* act */
				actualValue, actualErr := objectUnderTest.Construct(providedData)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Array: providedData}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("data is unexpected type", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _constructor{}
				expectedErr := fmt.Errorf("unable to construct value; '%+v' unexpected type", nil)

				/* act */
				_, actualErr := objectUnderTest.Construct(nil)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
	})
})
