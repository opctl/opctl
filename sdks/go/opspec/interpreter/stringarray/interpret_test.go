package stringarray

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("expression not empty", func() {
		Context("Contains reference", func() {
			It("should return expected string array", func() {
				/* arrange */
				identifier := "identifier"
				providedScope := map[string]*model.Value{
					identifier: {String: new(string)},
				}

				cmd2 := "cmd2"

				providedContainerCallSpecCmd := []interface{}{
					fmt.Sprintf("$(%s)", identifier),
					cmd2,
				}

				/* act */
				actualResult, _ := Interpret(
					providedScope,
					providedContainerCallSpecCmd,
				)

				/* assert */
				Expect(actualResult).To(Equal(
					[]string{
						*providedScope[identifier].String,
						cmd2,
					},
				))
			})
		})
	})
})
