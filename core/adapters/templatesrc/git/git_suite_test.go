package git

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "testing"
)

func TestGit(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "engine/core/adapters/templatesrc/git")
}
