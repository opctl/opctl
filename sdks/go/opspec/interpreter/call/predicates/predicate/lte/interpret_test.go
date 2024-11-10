package lte

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("number.Interpret errs", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			_, actualError := Interpret(
				[]interface{}{nil},
				map[string]*ipld.Node{},
			)

			/* assert */
			Expect(actualError).To(MatchError("unable to interpret <nil> to number: unable to interpret <nil> as value: unsupported type"))
		})
	})
	Context("items are in ascending order", func() {
		It("should return true", func() {
			/* arrange */
			/* act */
			actualResult, _ := Interpret(
				[]interface{}{
					1,
					2,
				},
				map[string]*ipld.Node{},
			)

			/* assert */
			Expect(actualResult).To(BeTrue())
		})
	})
	Context("items are in ascending and equal order", func() {
		It("should return true", func() {
			/* arrange */
			/* act */
			actualResult, _ := Interpret(
				[]interface{}{
					1,
					1,
					2,
				},
				map[string]*ipld.Node{},
			)

			/* assert */
			Expect(actualResult).To(BeTrue())
		})
	})
	Context("items aren't in ascending or equal order", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			actualResult, _ := Interpret(
				[]interface{}{
					2,
					1,
				},
				map[string]*ipld.Node{},
			)

			/* assert */
			Expect(actualResult).To(BeFalse())
		})
	})
})
