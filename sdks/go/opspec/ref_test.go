package opspec

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestRefToName(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(RefToName("$(foo)")).To(Equal("foo"))
}

func TestNameToRef(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(NameToRef("foo")).To(Equal("$(foo)"))
}
