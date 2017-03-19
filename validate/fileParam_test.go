package validate

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("Param", func() {
	objectUnderTest := New()
	Context("invoked w/ non-nil param.File", func() {
		Context("& non-empty value.File", func() {
			It("should return no errors", func() {

				/* arrange */
				providedValue := &model.Data{
					File: "dummyValue",
				}
				providedParam := &model.Param{
					File: &model.FileParam{},
				}

				expectedErrors := []error{}

				/* act */
				actualErrors := objectUnderTest.Param(providedValue, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("& empty value.File", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &model.Data{}
				providedParam := &model.Param{
					File: &model.FileParam{},
				}

				expectedErrors := []error{
					errors.New("File required"),
				}

				/* act */
				actualErrors := objectUnderTest.Param(providedValue, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("& nil value", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedParam := &model.Param{
					File: &model.FileParam{},
				}

				expectedErrors := []error{
					errors.New("File required"),
				}

				/* act */
				actualErrors := objectUnderTest.Param(nil, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
	})

})
