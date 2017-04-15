package validate

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("Param", func() {
	objectUnderTest := New()
	Context("invoked w/ non-nil param.Socket", func() {
		Context("& non-empty value.Socket", func() {
			It("should return no errors", func() {

				/* arrange */
				providedValueSocket := "dummySocket"
				providedValue := &model.Data{
					Socket: &providedValueSocket,
				}
				providedParam := &model.Param{
					Socket: &model.SocketParam{},
				}

				expectedErrors := []error{}

				/* act */
				actualErrors := objectUnderTest.Param(providedValue, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
		Context("& empty value.Socket", func() {
			It("should return expected errors", func() {

				/* arrange */
				providedValue := &model.Data{}
				providedParam := &model.Param{
					Socket: &model.SocketParam{},
				}

				expectedErrors := []error{
					errors.New("Socket required"),
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
					Socket: &model.SocketParam{},
				}

				expectedErrors := []error{
					errors.New("Socket required"),
				}

				/* act */
				actualErrors := objectUnderTest.Param(nil, providedParam)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))

			})
		})
	})

})
