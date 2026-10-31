package node

import (
	"bytes"
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	runtimeFakes "github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
	"io"
	"strings"
)

var _ = Describe("container stdout capture", func() {
	It("captures bootstrap stdout privately while preserving stderr and pull progress", func() {
		var captured bytes.Buffer
		runtime := new(runtimeFakes.FakeContainerRuntime)
		events := new(fakes.FakePubSub)
		runtime.RunContainerStub = func(ctx context.Context, call *model.ContainerCall, root string, publisher pubsub.EventPublisher, stdout, stderr io.WriteCloser) (*int64, error) {
			defer stdout.Close()
			defer stderr.Close()
			Expect(containerruntime.SeparateOutputFrom(ctx)).To(BeTrue())
			publisher.Publish(model.Event{ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{Data: []byte("pull progress")}})
			_, err := io.WriteString(stdout, "private plan\n")
			Expect(err).NotTo(HaveOccurred())
			_, err = io.WriteString(stderr, "compile diagnostic\n")
			Expect(err).NotTo(HaveOccurred())
			return nil, nil
		}
		cc := newContainerCaller(runtime, events, nil)
		_, err := cc.Call(withStdoutCapture(context.Background(), &captured), &model.ContainerCall{Image: &model.ContainerCallImage{}}, nil, &model.ContainerCallSpec{}, "root")
		Expect(err).NotTo(HaveOccurred())
		Expect(captured.String()).To(Equal("private plan\n"))
		var diagnostics []byte
		for i := 0; i < events.PublishCallCount(); i++ {
			event := events.PublishArgsForCall(i)
			Expect(event.ContainerStdOutWrittenTo).To(BeNil())
			if event.ContainerStdErrWrittenTo != nil {
				diagnostics = append(diagnostics, event.ContainerStdErrWrittenTo.Data...)
			}
		}
		Expect(string(diagnostics)).To(ContainSubstring("compile diagnostic"))
		Expect(string(diagnostics)).To(ContainSubstring("pull progress"))
	})
	It("continues to publish stdout for ordinary container calls", func() {
		runtime := new(runtimeFakes.FakeContainerRuntime)
		events := new(fakes.FakePubSub)
		runtime.RunContainerStub = func(ctx context.Context, _ *model.ContainerCall, _ string, _ pubsub.EventPublisher, stdout, stderr io.WriteCloser) (*int64, error) {
			defer stdout.Close()
			defer stderr.Close()
			Expect(containerruntime.SeparateOutputFrom(ctx)).To(BeFalse())
			_, err := io.WriteString(stdout, "normal output\n")
			return nil, err
		}
		_, err := newContainerCaller(runtime, events, nil).Call(context.Background(), &model.ContainerCall{Image: &model.ContainerCallImage{}}, nil, &model.ContainerCallSpec{}, "root")
		Expect(err).NotTo(HaveOccurred())
		Expect(events.PublishCallCount()).To(Equal(1))
		Expect(string(events.PublishArgsForCall(0).ContainerStdOutWrittenTo.Data)).To(Equal("normal output\n"))
	})
	It("fails oversized capture even when a runtime ignores the writer error", func() {
		var captured bytes.Buffer
		runtime := new(runtimeFakes.FakeContainerRuntime)
		runtime.RunContainerStub = func(_ context.Context, _ *model.ContainerCall, _ string, _ pubsub.EventPublisher, stdout, stderr io.WriteCloser) (*int64, error) {
			defer stdout.Close()
			defer stderr.Close()
			_, err := io.WriteString(stdout, strings.Repeat("x", 32*1024*1024+1))
			Expect(err).To(HaveOccurred())
			return nil, nil
		}
		_, err := newContainerCaller(runtime, new(fakes.FakePubSub), nil).Call(withStdoutCapture(context.Background(), &captured), &model.ContainerCall{Image: &model.ContainerCallImage{}}, nil, &model.ContainerCallSpec{}, "root")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("32 MiB"))
		Expect(captured.Len()).To(Equal(0))
	})
})
