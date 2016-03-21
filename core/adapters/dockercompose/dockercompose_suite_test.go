package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "testing"
)

func TestDockerCompose(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Dockercompose Suite")
}
