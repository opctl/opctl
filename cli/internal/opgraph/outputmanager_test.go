package opgraph

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewOutputManager(t *testing.T) {
	g := NewGomegaWithT(t)

	objectUnderTest := NewOutputManager()

	g.Expect(objectUnderTest).NotTo(BeNil())
	_, err := objectUnderTest.getWidth()
	g.Expect(err).NotTo(BeNil(), "in tests, terminal width isn't available")

	g.Expect(objectUnderTest.Print("")).NotTo(BeNil())
	objectUnderTest.Clear()
}

func TestOutputManagerShortLines(t *testing.T) {
	g := NewGomegaWithT(t)

	// arrange
	var buff bytes.Buffer
	objectUnderTest := OutputManager{
		getWidth: func() (int, error) { return 80, nil },
		out:      &buff,
	}

	// act
	err := objectUnderTest.Print(`testing
the drinks should just be like a glass of water
inspired by space jam and who knows`)

	// assert
	g.Expect(err).To(BeNil())
	expected := `┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄
testing
the drinks should just be like a glass of water
inspired by space jam and who knows`
	g.Expect(buff.String()).To(Equal(expected))
	g.Expect(objectUnderTest.lastHeight).To(Equal(4))
}

func TestOutputManagerLongLines(t *testing.T) {
	g := NewGomegaWithT(t)

	// arrange
	var buff bytes.Buffer
	objectUnderTest := OutputManager{
		getWidth: func() (int, error) { return 80, nil },
		out:      &buff,
	}
	longLine := strings.Repeat("-", 80)
	longerLine := strings.Repeat("-", 81)

	// act
	err := objectUnderTest.Print(fmt.Sprintf(`testing
%s
%s
testing3`, longLine, longerLine))

	// assert
	g.Expect(err).To(BeNil())
	expected := fmt.Sprintf(`┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄
testing
--------------------------------------------------------------------------------
------------------------------------------------------------------------------…%s
testing3`, "\033[0m")
	g.Expect(buff.String()).To(Equal(expected))
	g.Expect(objectUnderTest.lastHeight).To(Equal(5))
}

func TestOutputManagerClearing(t *testing.T) {
	g := NewGomegaWithT(t)

	// arrange
	var buff bytes.Buffer
	objectUnderTest := OutputManager{
		getWidth: func() (int, error) { return 80, nil },
		out:      &buff,
	}

	// act
	objectUnderTest.Print(`testing
the drinks should just be like a glass of water
inspired by space jam and who knows`)
	ioutil.ReadAll(&buff)
	objectUnderTest.Clear()

	// assert
	expected := "\x1b[80D\x1b[K\x1b[1A\x1b[K\x1b[1A\x1b[K\x1b[1A\x1b[K"
	g.Expect(buff.String()).To(Equal(expected))
	g.Expect(objectUnderTest.lastHeight).To(Equal(4))
}
