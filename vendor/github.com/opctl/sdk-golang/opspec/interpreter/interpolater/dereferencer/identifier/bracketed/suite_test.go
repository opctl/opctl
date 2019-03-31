package bracketed

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "opspec/interpreter/interpolater/dereferencer/identifier/bracketed")
}
