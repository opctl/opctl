package interpreter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestSdk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "opcall/input/interpreter")
}
