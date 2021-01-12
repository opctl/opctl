package port

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("value.Interpret errs", func() {
		It("should return expected err", func() {
			/* act */
			_, actualErr := Interpret(
				map[string]*model.Value{},
				"$()",
			)

			/* assert */
			Expect(actualErr).To(Equal(errors.New("unable to interpret $() to port; error was unable to interpret '' as reference; '' not in scope")))

		})
	})
	Context("value.Interpret doesn't err", func() {
		It("should return expected result", func() {
			/* arrange */
			identifier := "identifier"

			number := uint16(8000)
			expectedValue := model.Value{Port: &number}

			/* act */
			actualPort, actualErr := Interpret(
				map[string]*model.Value{
					identifier: &expectedValue,
				},
				fmt.Sprintf("$(%s)", identifier),
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(*actualPort).To(Equal(expectedValue))
		})
	})
})
