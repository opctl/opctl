package cliparamsatisfier

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc/fakes"
)

var _ = Describe("Go run input sources", func() {
	It("never replaces an explicit invalid value with a default", func() {
		defaults, prompt := new(fakes.FakeInputSrc), new(fakes.FakeInputSrc)
		defaults.ReadStringReturns(nil, true)
		source := goInputSourcer{provided: NewInputSourcer(New(nil).NewSliceInputSrc([]string{"env=staging"}, "=")), defaults: defaults, prompt: prompt, seen: map[string]bool{}}
		value, ok := source.Source("env")
		Expect(ok).To(BeTrue())
		Expect(*value).To(Equal("staging"))
		_, ok = source.Source("env")
		Expect(ok).To(BeFalse())
		Expect(defaults.ReadStringCallCount()).To(Equal(0))
		Expect(prompt.ReadStringCallCount()).To(Equal(1))
	})
	It("uses defaults for missing values and prompts for required values", func() {
		defaults, prompt := new(fakes.FakeInputSrc), new(fakes.FakeInputSrc)
		defaults.ReadStringReturnsOnCall(0, nil, true)
		source := goInputSourcer{provided: NewInputSourcer(), defaults: defaults, prompt: prompt, seen: map[string]bool{}}
		value, ok := source.Source("env")
		Expect(ok).To(BeTrue())
		Expect(value).To(BeNil())
		Expect(prompt.ReadStringCallCount()).To(Equal(0))
		_, ok = source.Source("token")
		Expect(ok).To(BeFalse())
		Expect(prompt.ReadStringCallCount()).To(Equal(1))
	})
})
