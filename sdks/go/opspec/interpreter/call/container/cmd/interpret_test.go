package cmd

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("containerCallSpecCmd not empty", func() {
		Context("Contains reference", func() {
			It("should return expected dcg.Cmd", func() {
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
