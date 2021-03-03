package parallelloop

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("loopable.Interpret errs", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			_, actualError := Interpret(
				model.ParallelLoopCallSpec{
					Range: "range",
				},
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualError).To(Equal(errors.New("unable to coerce string to object; error was invalid character 'r' looking for beginning of value")))
		})
	})
	It("should return expected result", func() {
		/* arrange */
		identifier := "identifier"
		providedScgLoop := model.ParallelLoopCallSpec{
			Range: fmt.Sprintf("$(%s)", identifier),
		}
		providedScope := map[string]*model.Value{
			identifier: {Array: new([]interface{})},
		}

		/* act */
		actualResult, _ := Interpret(
			providedScgLoop,
			providedScope,
		)

		/* assert */
		Expect(*actualResult).To(Equal(
			model.ParallelLoopCall{
				Range: providedScope[identifier],
			},
		))
	})
})
