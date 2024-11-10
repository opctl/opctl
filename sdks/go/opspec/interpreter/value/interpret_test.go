package value

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("expression is bool", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValueExpression := false

			/* act */
			actualValue, _ := Interpret(
				providedValueExpression,
				map[string]*ipld.Node{},
			)

			/* assert */
			Expect(actualValue).To(Equal(ipld.Node{Boolean: &providedValueExpression}))
		})
	})
	Context("expression is float64", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValueExpression := 2.2

			/* act */
			actualValue, _ := Interpret(
				providedValueExpression,
				map[string]*ipld.Node{},
			)

			/* assert */
			Expect(actualValue).To(Equal(ipld.Node{Number: &providedValueExpression}))
		})
	})
	Context("expression is int", func() {
		It("should return expected result", func() {
			/* arrange */
			providedValueExpression := 2
			expectedNumber := float64(providedValueExpression)

			/* act */
			actualValue, _ := Interpret(
				providedValueExpression,
				map[string]*ipld.Node{},
			)

			/* assert */
			Expect(actualValue).To(Equal(ipld.Node{Number: &expectedNumber}))
		})
	})
	Context("expression is map[string]interface{}", func() {
	})
	Context("expression is []interface{}", func() {
	})
	Context("expression is string", func() {
		Context("interpolater.Interpolate errs", func() {
			It("should return expected err", func() {
				/* arrange */
				/* act */
				_, actualErr := Interpret(
					"$()",
					map[string]*ipld.Node{},
				)

				/* assert */
				Expect(actualErr).To(MatchError("unable to interpret '' as reference: '' not in scope"))

			})
		})
	})
	It("should return expected result", func() {
		/* arrange */
		identifier := "identifier"
		stringValue := ipld.Node{String: new(string)}

		/* act */
		actualValue, actualErr := Interpret(
			"$(identifier)",
			map[string]*ipld.Node{
				identifier: &stringValue,
			},
		)

		/* assert */
		Expect(actualValue).To(Equal(stringValue))
		Expect(actualErr).To(BeNil())
	})
})
