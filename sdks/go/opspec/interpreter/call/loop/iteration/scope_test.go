package iteration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Scope", func() {
	Context("callSpec.Vars != nil", func() {
		Context("callSpec.Range == nil", func() {
			Context("callSpec.Vars.Index != nil", func() {
				It("should return expected result", func() {
					/* arrange */
					indexValue := 2
					indexValueAsFloat64 := float64(indexValue)

					indexName := "indexName"

					expectedScope := map[string]*ipld.Node{
						indexName: &ipld.Node{Number: &indexValueAsFloat64},
					}

					/* act */
					actualScope, _ := Scope(
						indexValue,
						map[string]*ipld.Node{},
						nil,
						&model.LoopVarsSpec{
							Index: &indexName,
						},
					)

					/* assert */
					Expect(actualScope).To(Equal(expectedScope))
				})
			})
		})
		Context("callSpec.Range != nil", func() {
			Context("loopable errs", func() {

				It("should return expected result", func() {
					/* arrange */
					providedLoopRange := "providedLoopRange"

					providedScope := map[string]*ipld.Node{
						"name1": &ipld.Node{String: new(string)},
					}

					/* act */
					_, actualErr := Scope(
						0,
						providedScope,
						providedLoopRange,
						&model.LoopVarsSpec{
							Index: new(string),
						},
					)

					/* assert */
					Expect(actualErr).To(MatchError("unable to coerce string to object: invalid character 'p' looking for beginning of value"))
				})
			})
		})
	})
})
