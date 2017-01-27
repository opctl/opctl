package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/eventbus"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/validate"
)

var _ = Describe("opCaller", func() {
	Context("newOpCaller", func() {
		It("should return opCaller", func() {
			/* arrange/act/assert */
			Expect(newOpCaller(
				new(bundle.FakeBundle),
				new(eventbus.FakeEventBus),
				newNodeRepo(),
				new(fakeCaller),
				new(uniquestring.FakeUniqueStringFactory),
				validate.New(),
			)).Should(Not(BeNil()))
		})
	})
})
