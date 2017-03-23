package pkg

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("validator", func() {
	wd, err := os.Getwd()
	if nil != err {
		panic(err)
	}
	objectUnderTest := newValidator()
	Context("invalid__yml", func() {
		It("should return expected errs", func() {
			/* arrange */
			expectedErrs := []error{errors.New(
				"Error validating pkg. Details: yaml: did not find expected alphabetic or numeric character",
			)}

			/* act */
			actualErrs := objectUnderTest.Validate(fmt.Sprintf("%v/testdata/validate/invalid__yml", wd))

			/* assert */
			Expect(actualErrs).To(Equal(expectedErrs))
		})
	})
	Context("invalid_inputs_type", func() {
		It("should return expected errs", func() {

			/* arrange */
			expectedErrs := []error{errors.New(
				"inputs: Invalid type. Expected: object, given: array",
			)}

			/* act */
			actualErrs := objectUnderTest.Validate(fmt.Sprintf("%v/testdata/validate/invalid_inputs_type", wd))

			/* assert */
			Expect(actualErrs).To(Equal(expectedErrs))
		})
	})
	Context("invalid_outputs_type", func() {
		It("should return expected errs", func() {

			/* arrange */
			expectedErrs := []error{errors.New(
				"outputs: Invalid type. Expected: object, given: array",
			)}

			/* act */
			actualErrs := objectUnderTest.Validate(fmt.Sprintf("%v/testdata/validate/invalid_outputs_type", wd))

			/* assert */
			Expect(actualErrs).To(Equal(expectedErrs))
		})
	})
	Context("invalid_run_type", func() {
		It("should return expected errs", func() {

			/* arrange */
			expectedErrs := []error{
				errors.New("run: Must validate one and only one schema (oneOf)"),
				errors.New("run: Invalid type. Expected: object, given: array"),
			}

			/* act */
			actualErrs := objectUnderTest.Validate(fmt.Sprintf("%v/testdata/validate/invalid_run_type", wd))

			/* assert */
			Expect(actualErrs).To(Equal(expectedErrs))
		})
	})
	Context("valid__all", func() {
		It("should return no errors", func() {
			/* act */
			actualErrs := objectUnderTest.Validate(fmt.Sprintf("%v/testdata/validate/valid__all", wd))

			/* assert */
			Expect(actualErrs).To(BeNil())
		})
	})

})
