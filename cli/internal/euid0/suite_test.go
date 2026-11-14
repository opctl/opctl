package euid0

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEuid0(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "opctl/cli/internal/euid0")
}
