package main

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "testing"
)

func TestEngine(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Engine Suite")
}
