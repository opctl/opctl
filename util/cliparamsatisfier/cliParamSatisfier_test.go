package cliparamsatisfier

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/clioutput"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/validate"
)

var _ = Context("parameterSatisfier", func() {
	Context("Satisfy", func() {
		It("should call inputSourcer.Source w/ expected args for each input", func() {
			/* arrange */
			providedInputSourcer := new(FakeInputSourcer)
			providedInputs := map[string]*model.Param{
				"input1": {String: &model.StringParam{}},
				"input2": {String: &model.StringParam{}},
			}

			expectedInputNames := map[string]struct{}{
				"input1": struct{}{},
				"input2": struct{}{},
			}

			objectUnderTest := New(
				new(cliexiter.Fake),
				new(clioutput.Fake),
				new(validate.Fake),
			)

			/* act */
			objectUnderTest.Satisfy(providedInputSourcer, providedInputs)

			/* assert */
			actualInputNames := map[string]struct{}{}
			for callIndex := 0; callIndex < providedInputSourcer.SourceCallCount(); callIndex++ {
				actualInputName := providedInputSourcer.SourceArgsForCall(callIndex)
				actualInputNames[actualInputName] = struct{}{}
			}

			Expect(actualInputNames).To(Equal(expectedInputNames))
		})
	})
})
