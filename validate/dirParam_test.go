package validate

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("Param", func() {
	objectUnderTest := New()
	Context("invoked w/ non-nil param.Dir", func() {
		Context("& non-empty value.Dir", func() {
			It("should return no errors", func() {

				/* arrange */
				providedValue := &model.Data{
					Dir: "dummyValue",
				}
				providedParam := &model.Param{
					Dir: &model.DirParam{},
				}

				expectedErrors := []error{}

				/* act */
				actualErrors := objectUnderTest.Param(providedValue, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("& empty value.Dir", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &model.Data{}
				providedParam := &model.Param{
					Dir: &model.DirParam{},
				}

				expectedErrors := []error{
					errors.New("Dir required"),
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
					Dir: &model.DirParam{},
				}

				expectedErrors := []error{
					errors.New("Dir required"),
				}

				/* act */
				actualErrors := objectUnderTest.Param(nil, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
	})

})
