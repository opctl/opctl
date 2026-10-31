package k8s

import (
	"context"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
)

type discardCloser struct{ io.Writer }

func (discardCloser) Close() error { return nil }

var _ = Describe("Go op bootstrap", func() {
	It("fails explicitly before calling Kubernetes when separate streams are needed", func() {
		runtime := _containerRuntime{}
		_, err := runtime.RunContainer(containerruntime.WithSeparateOutput(context.Background()), nil, "root", nil, discardCloser{io.Discard}, discardCloser{io.Discard})
		Expect(err).To(MatchError("Go op bootstrap is not supported yet with the Kubernetes container runtime"))
	})
})
