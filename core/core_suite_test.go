package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "testing"
)

func TestCore(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "engine/core")
}
