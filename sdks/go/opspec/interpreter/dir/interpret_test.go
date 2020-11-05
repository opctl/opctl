package dir

import (
	"os"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("expression is scope ref", func() {
		Context("reference.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				identifier := "identifier"
				/* act */
				_, actualErr := Interpret(
					map[string]*model.Value{
						identifier: &model.Value{
							Socket: new(string),
						},
					},
					fmt.Sprintf("$(%s)", identifier),
					os.TempDir(),
					true,
				)

				/* assert */
				Expect(actualErr.Error()).To(Equal("unable to interpret $(identifier) to dir"))

			})
		})
		Context("reference.Interpret doesn't error", func() {
			Context("value.Dir nil", func() {
				It("should return expected result", func() {
					/* arrange */
					identifier := "identifier"
					providedScope := map[string]*model.Value{
						identifier: &model.Value{Dir: nil},
					}
					providedExpression := fmt.Sprintf("$(%s)", identifier)

					/* act */
					_, actualErr := Interpret(
						providedScope,
						providedExpression,
						os.TempDir(),
						true,
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("unable to interpret %+v to dir", providedExpression)))
				})
			})
			Context("value.Dir not nil", func() {
				It("should return expected result", func() {
					/* arrange */
					identifier := "identifier"
					providedScope := map[string]*model.Value{
						identifier: &model.Value{Dir: new(string)},
					}

					/* act */
					actualResult, actualErr := Interpret(
						providedScope,
						fmt.Sprintf("$(%s)", identifier),
						os.TempDir(),
						true,
					)

					/* assert */
					Expect(actualErr).To(BeNil())
					Expect(actualResult).To(Equal(providedScope[identifier]))
				})
			})
		})
	})
})
