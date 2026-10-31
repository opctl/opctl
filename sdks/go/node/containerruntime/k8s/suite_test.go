package k8s

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func Test(t *testing.T) { RegisterFailHandler(Fail); RunSpecs(t, "node/containerruntime/k8s") }
