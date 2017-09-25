package data

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("Validate", func() {
	arrayUnderTest := newValidator()
	Context("param.Array not nil", func() {
		Context("value.Array nil", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &model.Value{}
				providedParam := &model.Param{
					Array: &model.ArrayParam{},
				}

				expectedErrors := []error{
					fmt.Errorf("unable to coerce '%+v' to array", providedValue),
				}

				arrayUnderTest := newValidator()

				/* act */
				actualErrors := arrayUnderTest.Validate(
					providedValue,
					providedParam,
				)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("value.Array not nil", func() {
			Context("Items constraint", func() {
				Context("value item meets constraint", func() {
					It("returns no errors", func() {

						/* arrange */
						providedValueArray := []interface{}{
							2,
						}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									Items: []interface{}{
										&model.TypeConstraints{
											Enum: []interface{}{2},
										},
									},
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value items meet constraints", func() {
					It("returns no errors", func() {

						/* arrange */
						providedValueArray := []interface{}{
							2,
							2,
						}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									Items: []interface{}{&model.TypeConstraints{
										Enum: []interface{}{2},
									},
									},
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value item doesn't meet constraint", func() {
					It("returns expected errors", func() {

						/* arrange */
						providedValueArray := []interface{}{
							2,
						}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									Items: []interface{}{
										&model.TypeConstraints{
											Enum: []interface{}{3},
										},
									},
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"0 must be one of the following: %v",
								3,
							),
						}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value items don't meet constraints", func() {
					It("returns expected errors", func() {

						/* arrange */
						providedValueArray := []interface{}{
							2,
							3,
						}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									Items: []interface{}{
										&model.TypeConstraints{
											Enum: []interface{}{2},
										},
										&model.TypeConstraints{
											Enum: []interface{}{2},
										},
									},
								},
							},
						}

						expectedErrors := []error{
							errors.New("1 must be one of the following: 2"),
						}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("AdditionalItems constraint", func() {
				Context("applicable value items meet constraint", func() {
					It("returns no errors", func() {

						/* arrange */
						providedValueArray := []interface{}{
							2,
							2,
						}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									Items: []interface{}{&model.TypeConstraints{
										Enum: []interface{}{2},
									},
									},
									AdditionalItems: &model.TypeConstraints{
										Enum: []interface{}{2},
									},
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("applicable value items don't meet constraint", func() {
					It("returns expected errors", func() {

						/* arrange */
						providedValueArray := []interface{}{
							2,
							3,
						}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									Items: []interface{}{&model.TypeConstraints{
										Enum: []interface{}{2},
									},
									},
									AdditionalItems: &model.TypeConstraints{
										Enum: []interface{}{2},
									},
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"1 must be one of the following: %v",
								providedParam.Array.Constraints.AdditionalItems.Enum[0],
							),
						}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("MaxItems constraint", func() {
				Context("value item count == MaxItems", func() {
					It("returns no errors", func() {

						/* arrange */
						providedValueArray := []interface{}{
							"dummyItem",
						}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									MaxItems: 1,
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value item count > MaxItems", func() {
					It("returns expected errors", func() {

						/* arrange */
						providedValueArray := []interface{}{
							"dummyItem1",
							"dummyItem2",
						}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									MaxItems: 1,
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"Array must have at most %v items",
								providedParam.Array.Constraints.MaxItems,
							),
						}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value item count < MaxItems", func() {
					It("returns no errors", func() {

						/* arrange */
						providedValueArray := []interface{}{
							"dummyItem",
						}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									MaxItems: 2,
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("MinItems constraint", func() {
				Context("value item count == MinItems", func() {
					It("should return no errors", func() {

						/* arrange */
						providedValueArray := []interface{}{"dummyItem"}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									MinItems: 1,
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value item count < MinItems", func() {

					It("should return expected errors", func() {

						/* arrange */
						providedValueArray := []interface{}{"dummyItem"}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									MinItems: 2,
								},
							},
						}

						expectedErrors := []error{
							fmt.Errorf(
								"Array must have at least %v items",
								providedParam.Array.Constraints.MinItems,
							),
						}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value item count > MinItems", func() {

					It("should return no errors", func() {

						/* arrange */
						providedValueArray := []interface{}{
							"dummyItem1",
							"dummyItem2",
						}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									MinItems: 1,
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
			Context("UniqueItems constraint", func() {
				Context("value contains no duplicates", func() {
					It("returns no errors", func() {

						/* arrange */
						providedValueArray := []interface{}{
							"dummyItem1",
							"dummyItem2",
						}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									UniqueItems: true,
								},
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("value contains a duplicate", func() {
					It("returns expected errors", func() {

						/* arrange */
						providedValueArray := []interface{}{
							"dummyItem1",
							"dummyItem1",
						}
						providedValue := &model.Value{
							Array: providedValueArray,
						}
						providedParam := &model.Param{
							Array: &model.ArrayParam{
								Constraints: &model.ArrayConstraints{
									UniqueItems: true,
								},
							},
						}

						expectedErrors := []error{
							errors.New("array items must be unique"),
						}

						/* act */
						actualErrors := arrayUnderTest.Validate(providedValue, providedParam)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
		})
	})
})
