package taggosdk

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTagGoSDK(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tag-go-sdk")
}
