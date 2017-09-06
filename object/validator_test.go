package object

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"strings"
)

var _ = Context("Validate", func() {
	objectUnderTest := newValidator()
	Context("value nil", func() {
		It("should return expected errors", func() {

			/* arrange */
			expectedErrors := []error{
				errors.New("object required"),
			}

			objectUnderTest := newValidator()

			/* act */
			actualErrors := objectUnderTest.Validate(
				nil,
				&model.ObjectConstraints{},
			)

			/* assert */
			Expect(actualErrors).To(Equal(expectedErrors))

		})
	})
	Context("value not nil", func() {
		Context("AdditionalProperties constraint", func() {
			Context("value props don't match properties, or patternProperties", func() {
				Context("value props meet AdditionalProperties constraint", func() {

					It("returns no errors", func() {

						/* arrange */
						providedValue := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
						}
						providedConstraints := &model.ObjectConstraints{
							AdditionalProperties: &model.JSONSchema{
								MinLength: 2,
							},
						}

						expectedErrors := []error{}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("some value props meet AdditionalProperties constraint", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValueProp1Value := "dummyProp1Value"
						providedValue := map[string]interface{}{
							"dummyProp1Name": providedValueProp1Value,
							"dummyProp2Name": "dummyProp2Value",
						}

						providedConstraints := &model.ObjectConstraints{
							AdditionalProperties: &model.JSONSchema{
								Pattern: providedValueProp1Value,
							},
						}

						expectedErrors := []error{
							fmt.Errorf("Does not match pattern '%v'", providedValueProp1Value),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
				Context("no value props meet AdditionalProperties constraint", func() {

					It("returns expected errors", func() {

						/* arrange */
						providedValue := map[string]interface{}{
							"dummyProp1Name": "dummyProp1Value",
							"dummyProp2Name": "dummyProp2Value",
						}

						pattern := "dummyPattern"
						providedConstraints := &model.ObjectConstraints{
							AdditionalProperties: &model.JSONSchema{
								Pattern: pattern,
							},
						}

						expectedErrors := []error{
							fmt.Errorf("Does not match pattern '%v'", pattern),
							fmt.Errorf("Does not match pattern '%v'", pattern),
						}

						/* act */
						actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

						/* assert */
						Expect(actualErrors).To(Equal(expectedErrors))

					})
				})
			})
		})
		Context("AllOf constraint", func() {
			Context("value meets all AllOf constraints", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						AllOf: []*model.ObjectConstraints{
							{
								MinProperties: 0,
							},
							{
								MaxProperties: 1,
							},
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value doesn't meet all AllOf constraints", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
						"dummyProp2Name": "dummyProp2Value",
					}

					providedConstraints := &model.ObjectConstraints{
						AllOf: []*model.ObjectConstraints{
							{
								MinProperties: 1,
							},
							{
								MaxProperties: 1,
							},
						},
					}

					expectedErrors := []error{
						errors.New(`Must have at most 1 properties`),
						errors.New("Must validate all the schemas (allOf)"),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("AnyOf constraint", func() {
			Context("value meets an AnyOf constraint", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						AnyOf: []*model.ObjectConstraints{
							{
								MinProperties: 1,
							},
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value doesn't meet an AnyOf constraint", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{}

					providedConstraints := &model.ObjectConstraints{
						AnyOf: []*model.ObjectConstraints{
							{
								MinProperties: 1,
							},
						},
					}

					expectedErrors := []error{
						errors.New("Must validate at least one schema (anyOf)"),
						errors.New("Must have at least 1 properties"),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("Enum constraint", func() {
			Context("value in enum", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						Enum: []map[string]interface{}{
							providedValue,
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value not in enum", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						Enum: []map[string]interface{}{
							{
								"dummyName": "dummyValue",
							},
						},
					}

					expectedError, err := json.Marshal(providedConstraints.Enum[0])
					if nil != err {
						panic(err.Error())
					}

					expectedErrors := []error{
						fmt.Errorf(
							`must be one of the following: %v`,
							string(expectedError),
						),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("MaxProperties constraint", func() {
			Context("value prop count == MaxProperties", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						MaxProperties: 1,
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value prop count > MaxProperties", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
						"dummyProp2Name": "dummyProp2Value",
					}

					providedConstraints := &model.ObjectConstraints{
						MaxProperties: 1,
					}

					expectedErrors := []error{
						fmt.Errorf(
							"Must have at most %v properties",
							providedConstraints.MaxProperties,
						),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value prop count < MaxProperties", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						MaxProperties: 2,
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("MinProperties constraint", func() {
			Context("value prop count == MinProperties", func() {

				It("should return no errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						MinProperties: 1,
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value prop count < MinProperties", func() {

				It("should return expected errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						MinProperties: 2,
					}

					expectedErrors := []error{
						fmt.Errorf(
							"Must have at least %v properties",
							providedConstraints.MinProperties,
						),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value prop count > MinProperties", func() {

				It("should return no errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
						"dummyProp2Name": "dummyProp2Value",
					}

					providedConstraints := &model.ObjectConstraints{
						MinProperties: 1,
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("Not constraint", func() {
			Context("value matches", func() {

				It("should return expected errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						Not: &model.ObjectConstraints{
							MinProperties: 1,
						},
					}

					expectedErrors := []error{
						errors.New("Must not validate the schema (not)"),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value doesn't match", func() {

				It("should return no errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						Not: &model.ObjectConstraints{
							MinProperties: 2,
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("OneOf constraint", func() {
			Context("value meets a single OneOf constraint", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						OneOf: []*model.ObjectConstraints{
							{
								MinProperties: 1,
							},
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value meets no OneOf constraints", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						OneOf: []*model.ObjectConstraints{
							{
								MinProperties: 2,
							},
						},
					}

					expectedErrors := []error{
						errors.New("Must validate one and only one schema (oneOf)"),
						errors.New("Must have at least 2 properties"),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value meets multiple OneOf constraints", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
						"dummyProp2Name": "dummyProp2Value",
					}

					providedConstraints := &model.ObjectConstraints{
						OneOf: []*model.ObjectConstraints{
							{
								MinProperties: 2,
							},
							{
								MaxProperties: 2,
							},
						},
					}

					expectedErrors := []error{
						errors.New("Must validate one and only one schema (oneOf)"),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
		Context("PatternProperties constraint", func() {
			Context("value props meet all Properties constraints", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValueProp1Name := "dummyProp1Name"
					providedValue := map[string]interface{}{
						providedValueProp1Name: "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						PatternProperties: map[string]*model.JSONSchema{
							providedValueProp1Name: {MinLength: 2},
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value props meet some Properties constraints", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValueProp1Name := "dummyProp1Name"
					providedValueProp2Name := "dummyProp2Name"
					providedValue := map[string]interface{}{
						providedValueProp1Name: "dummyProp1Value",
						providedValueProp2Name: "dummyProp2Value",
					}

					maxLength := 1
					providedConstraints := &model.ObjectConstraints{
						PatternProperties: map[string]*model.JSONSchema{
							providedValueProp1Name: {MinLength: 2},
							providedValueProp2Name: {MaxLength: maxLength},
						},
					}

					patterns, err := json.Marshal([]string{providedValueProp1Name, providedValueProp2Name})
					if err != nil {
						panic(err.Error())
					}

					expectedErrors := []error{
						fmt.Errorf("String length must be less than or equal to %v", maxLength),
						fmt.Errorf("Property \"%v\" does not match pattern %v", providedValueProp2Name, string(patterns)),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					// fix pattern order (maps in go not stable)
					OutOfOrderPatterns, err := json.Marshal([]string{providedValueProp2Name, providedValueProp1Name})
					if err != nil {
						panic(err.Error())
					}
					for errIndex, err := range actualErrors {
						actualErrors[errIndex] = errors.New(
							strings.Replace(err.Error(), string(OutOfOrderPatterns), string(patterns), -1),
						)
					}
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value props meet no Properties constraints", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValueProp1Name := "dummyProp1Name"
					providedValueProp2Name := "dummyProp2Name"
					providedValue := map[string]interface{}{
						providedValueProp1Name: "dummyProp1Value",
						providedValueProp2Name: "dummyProp2Value",
					}

					minLength := 100
					maxLength := 1
					providedConstraints := &model.ObjectConstraints{
						PatternProperties: map[string]*model.JSONSchema{
							providedValueProp1Name: {MinLength: minLength},
							providedValueProp2Name: {MaxLength: maxLength},
						},
					}

					patterns, err := json.Marshal([]string{providedValueProp1Name, providedValueProp2Name})
					if err != nil {
						panic(err.Error())
					}

					expectedErrors := []error{
						fmt.Errorf("String length must be greater than or equal to %v", minLength),
						fmt.Errorf("Property \"%v\" does not match pattern %v", providedValueProp1Name, string(patterns)),
						fmt.Errorf("String length must be less than or equal to %v", maxLength),
						fmt.Errorf("Property \"%v\" does not match pattern %v", providedValueProp2Name, string(patterns)),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					// fix pattern order (maps in go not stable)
					OutOfOrderPatterns, err := json.Marshal([]string{providedValueProp2Name, providedValueProp1Name})
					if err != nil {
						panic(err.Error())
					}
					for errIndex, err := range actualErrors {
						actualErrors[errIndex] = errors.New(
							strings.Replace(err.Error(), string(OutOfOrderPatterns), string(patterns), -1),
						)
					}
					Expect(actualErrors).To(ConsistOf(expectedErrors))

				})
			})
		})
		Context("Properties constraint", func() {
			Context("value props meet all Properties constraints", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValueProp1Name := "dummyProp1Name"
					providedValue := map[string]interface{}{
						providedValueProp1Name: "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						Properties: map[string]*model.JSONSchema{
							providedValueProp1Name: {MinLength: 2},
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value props meet some Properties constraints", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValueProp1Name := "dummyProp1Name"
					providedValueProp2Name := "dummyProp2Name"
					providedValue := map[string]interface{}{
						providedValueProp1Name: "dummyProp1Value",
						providedValueProp2Name: "dummyProp2Value",
					}

					maxLength := 1
					providedConstraints := &model.ObjectConstraints{
						Properties: map[string]*model.JSONSchema{
							providedValueProp1Name: {MinLength: 2},
							providedValueProp2Name: {MaxLength: maxLength},
						},
					}

					expectedErrors := []error{
						fmt.Errorf("String length must be less than or equal to %v", maxLength),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value props meet no Properties constraints", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValueProp1Name := "dummyProp1Name"
					providedValueProp2Name := "dummyProp2Name"
					providedValue := map[string]interface{}{
						providedValueProp1Name: "dummyProp1Value",
						providedValueProp2Name: "dummyProp2Value",
					}

					minLength := 100
					maxLength := 1
					providedConstraints := &model.ObjectConstraints{
						Properties: map[string]*model.JSONSchema{
							providedValueProp1Name: {MinLength: minLength},
							providedValueProp2Name: {MaxLength: maxLength},
						},
					}

					expectedErrors := []error{
						fmt.Errorf("String length must be greater than or equal to %v", minLength),
						fmt.Errorf("String length must be less than or equal to %v", maxLength),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(ConsistOf(expectedErrors))

				})
			})
		})
		Context("Required constraint", func() {
			Context("value contains all required props", func() {

				It("returns no errors", func() {

					/* arrange */
					providedValueProp1Name := "dummyProp1Name"
					providedValue := map[string]interface{}{
						providedValueProp1Name: "dummyProp1Value",
					}

					providedConstraints := &model.ObjectConstraints{
						Required: []string{
							providedValueProp1Name,
						},
					}

					expectedErrors := []error{}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value contains no required props", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValue := map[string]interface{}{
						"dummyProp1Name": "dummyProp1Value",
					}

					missingPropName := "missingPropName"
					providedConstraints := &model.ObjectConstraints{
						Required: []string{
							missingPropName,
						},
					}

					expectedErrors := []error{
						fmt.Errorf(
							"%v is required",
							missingPropName,
						),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
			Context("value contains some required props", func() {

				It("returns expected errors", func() {

					/* arrange */
					providedValueProp1Name := "dummyProp1Name"
					providedValue := map[string]interface{}{
						providedValueProp1Name: "dummyProp1Value",
						"dummyProp2Name":       "dummyProp2Value",
					}

					missingPropName := "missingPropName"
					providedConstraints := &model.ObjectConstraints{
						Required: []string{
							providedValueProp1Name,
							missingPropName,
						},
					}

					expectedErrors := []error{
						fmt.Errorf(
							"%v is required",
							missingPropName,
						),
					}

					/* act */
					actualErrors := objectUnderTest.Validate(providedValue, providedConstraints)

					/* assert */
					Expect(actualErrors).To(Equal(expectedErrors))

				})
			})
		})
	})
})
