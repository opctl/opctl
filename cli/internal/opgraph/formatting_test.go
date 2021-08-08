package opgraph

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestStripAnsi(t *testing.T) {
	g := NewGomegaWithT(t)

	str := "test"
	ansiStr := "\033[1Atest"
	withoutAnsi := stripAnsi(ansiStr)

	g.Expect(withoutAnsi).To(Equal(str), "stripped string is not equal to original")
}

func TestStripAnsi_noAnsi(t *testing.T) {
	g := NewGomegaWithT(t)

	str := "◉ ⠴ ./test"
	ansiStr := "◉ ⠴ [1m./test[0m"
	withoutAnsi := stripAnsi(ansiStr)

	g.Expect(withoutAnsi).To(Equal(str), "stripped string is not equal to original")
}

func TestStripAnsiToLength(t *testing.T) {
	g := NewGomegaWithT(t)

	ansiStr := "\033[1Atesting a string"
	stripped := stripAnsiToLength(ansiStr, 9)
	expected := "\033[1Atesting a"

	g.Expect(stripped).To(Equal(expected))
}

func TestStripAnsiToLength_escapeCodeInMid(t *testing.T) {
	g := NewGomegaWithT(t)

	ansiStr := "\033[1Atesting\033[0m a\033[1A string"
	stripped := stripAnsiToLength(ansiStr, 9)
	expected := "\033[1Atesting\033[0m a\033[1A"

	g.Expect(stripped).To(Equal(expected))
}

func TestStripAnsiToLength_noAnsi(t *testing.T) {
	g := NewGomegaWithT(t)

	ansiStr := "testing a string"
	stripped := stripAnsiToLength(ansiStr, 9)
	expected := "testing a"

	g.Expect(stripped).To(Equal(expected))
}
