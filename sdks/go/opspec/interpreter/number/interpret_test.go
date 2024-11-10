package number

import (
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
				map[string]*ipld.Node{},
				"$()",
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret $() to number: unable to interpret '' as reference: '' not in scope"))

		})
	})
	Context("value.Interpret doesn't err", func() {
		It("should return expected result", func() {
			/* arrange */
			identifier := "identifier"

			number := 2.2
			expectedValue := ipld.Node{Number: &number}

			/* act */
			actualNumber, actualErr := Interpret(
				map[string]*ipld.Node{
					identifier: &expectedValue,
				},
				fmt.Sprintf("$(%s)", identifier),
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(*actualNumber).To(Equal(expectedValue))
		})
	})
})
