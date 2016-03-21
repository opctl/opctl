package osfilesys

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "testing"
)

func TestOsfilesys(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Osfilesys Suite")
}
