package dir

import (
	"fmt"
	"os"

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
				scratchDir, err := os.MkdirTemp("", "")
				if err != nil {
					panic(err)
				}

				/* act */
				_, actualErr := Interpret(
					map[string]*ipld.Node{
						identifier: {
							Socket: new(string),
						},
					},
					fmt.Sprintf("$(%s)", identifier),
					scratchDir,
					true,
				)

				/* assert */
				Expect(actualErr).To(MatchError("unable to interpret $(identifier) to dir: unable to coerce socket to dir: incompatible types"))

			})
		})
		Context("reference.Interpret doesn't error", func() {
			Context("value.Dir nil", func() {
				It("should return expected result", func() {
					/* arrange */
					identifier := "identifier"
					providedScope := map[string]*ipld.Node{
						identifier: {Dir: nil},
					}
					providedExpression := fmt.Sprintf("$(%s)", identifier)
					scratchDir, err := os.MkdirTemp("", "")
					if err != nil {
						panic(err)
					}

					/* act */
					_, actualErr := Interpret(
						providedScope,
						providedExpression,
						scratchDir,
						true,
					)

					/* assert */
					Expect(actualErr).To(MatchError("unable to interpret $(identifier) to dir: unable to coerce '&{Array:<nil> Boolean:<nil> Dir:<nil> File:<nil> Number:<nil> Object:<nil> Socket:<nil> String:<nil>}' to dir"))
				})
			})
			Context("value.Dir not nil", func() {
				It("should return expected result", func() {
					/* arrange */
					identifier := "identifier"
					providedScope := map[string]*ipld.Node{
						identifier: {Dir: new(string)},
					}
					scratchDir, err := os.MkdirTemp("", "")
					if err != nil {
						panic(err)
					}

					/* act */
					actualResult, actualErr := Interpret(
						providedScope,
						fmt.Sprintf("$(%s)", identifier),
						scratchDir,
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
