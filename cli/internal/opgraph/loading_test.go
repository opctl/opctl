package opgraph

import (
	"testing"
	"unicode/utf8"

	. "github.com/onsi/gomega"
)

func TestDotLoadingSpinner(t *testing.T) {
	g := NewGomegaWithT(t)

	// arrange
	var objectUnderTest DotLoadingSpinner

	// act
	l := objectUnderTest.String()

	// assert
	_, size := utf8.DecodeRuneInString(l)
	g.Expect(size).To(Equal(3))
	g.Expect(len(l)).To(Equal(size))
}
