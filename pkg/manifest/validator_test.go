package manifest

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Context("Validator", func() {
	wd, err := os.Getwd()
	if nil != err {
		panic(err)
	}
	objectUnderTest := New()
	Context("invalid__yml", func() {
		It("should return expected errs", func() {
			/* arrange */
			expectedErrs := []error{errors.New(
				"yaml: did not find expected alphabetic or numeric character",
			)}

			manifestFile, err := os.Open(fmt.Sprintf("%v/testdata/validate/invalid__yml/op.yml", wd))
			if nil != err {
				Fail(err.Error())
			}

			manifestBytes, err := ioutil.ReadAll(manifestFile)
			if nil != err {
				Fail(err.Error())
			}

			/* act */
			actualErrs := objectUnderTest.Validate(manifestBytes)

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

			manifestFile, err := os.Open(fmt.Sprintf("%v/testdata/validate/invalid_inputs_type/op.yml", wd))
			if nil != err {
				Fail(err.Error())
			}

			manifestBytes, err := ioutil.ReadAll(manifestFile)
			if nil != err {
				Fail(err.Error())
			}

			/* act */
			actualErrs := objectUnderTest.Validate(manifestBytes)

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

			manifestFile, err := os.Open(fmt.Sprintf("%v/testdata/validate/invalid_outputs_type/op.yml", wd))
			if nil != err {
				Fail(err.Error())
			}

			manifestBytes, err := ioutil.ReadAll(manifestFile)
			if nil != err {
				Fail(err.Error())
			}

			/* act */
			actualErrs := objectUnderTest.Validate(manifestBytes)

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

			manifestFile, err := os.Open(fmt.Sprintf("%v/testdata/validate/invalid_run_type/op.yml", wd))
			if nil != err {
				Fail(err.Error())
			}

			manifestBytes, err := ioutil.ReadAll(manifestFile)
			if nil != err {
				Fail(err.Error())
			}

			/* act */
			actualErrs := objectUnderTest.Validate(manifestBytes)

			/* assert */
			Expect(actualErrs).To(Equal(expectedErrs))
		})
	})
	Context("valid__all", func() {
		It("should return no errors", func() {

			/* arrange */
			manifestFile, err := os.Open(fmt.Sprintf("%v/testdata/validate/valid__all/op.yml", wd))
			if nil != err {
				Fail(err.Error())
			}

			manifestBytes, err := ioutil.ReadAll(manifestFile)
			if nil != err {
				Fail(err.Error())
			}

			/* act */
			actualErrs := objectUnderTest.Validate(manifestBytes)

			/* assert */
			Expect(actualErrs).To(BeNil())
		})
	})

})
