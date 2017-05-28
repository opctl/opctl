package handler

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSDK(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "node/api/handler")
}
