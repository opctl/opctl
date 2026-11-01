package node

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	runtimefakes "github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

var _ = Describe("container stream bindings", func() {
	It("keeps public logs raw, captures separately, and suppresses secret stream events", func() {
		for _, secret := range []bool{false, true} {
			runtime := new(runtimefakes.FakeContainerRuntime)
			bus := new(fakes.FakePubSub)
			runtime.RunContainerStub = func(ctx context.Context, c *model.ContainerCall, _ string, p pubsub.EventPublisher, out, stderr io.WriteCloser) (*int64, error) {
				defer out.Close()
				defer stderr.Close()
				Expect(containerruntime.SeparateOutputFrom(ctx)).To(BeTrue())
				Expect(containerruntime.BootstrapFrom(ctx)).To(BeFalse())
				out.Write([]byte("version-private-marker\r"))
				out.Write([]byte("\n"))
				stderr.Write([]byte("diagnostic\n"))
				return nil, nil
			}
			cc := newContainerCaller(runtime, bus, &_stateStore{})
			ctx := context.Background()
			if secret {
				ctx = withOutputIntent(ctx, map[string]*model.ParamSpec{"version": {String: &model.StringParamSpec{IsSecret: true}}})
			}
			result, err := cc.Call(ctx, &model.ContainerCall{Image: &model.ContainerCallImage{}}, nil, &model.ContainerCallSpec{Stdout: "$(version)", Stderr: "$(diagnostic)"}, "root")
			Expect(err).NotTo(HaveOccurred())
			Expect(*result["version"].String).To(Equal("version-private-marker"))
			Expect(result["version"].Sensitive).To(Equal(secret))
			var raw []byte
			for i := 0; i < bus.PublishCallCount(); i++ {
				e := bus.PublishArgsForCall(i)
				if e.ContainerStdOutWrittenTo != nil {
					raw = append(raw, e.ContainerStdOutWrittenTo.Data...)
				}
				data, _ := json.Marshal(e)
				if secret {
					Expect(string(data)).NotTo(ContainSubstring("version-private-marker"))
				}
			}
			if secret {
				Expect(raw).To(BeEmpty())
			} else {
				Expect(string(raw)).To(Equal("version-private-marker\r\n"))
			}
		}
	})
	It("fails atomically on overflow even when the runtime ignores the writer error", func() {
		runtime := new(runtimefakes.FakeContainerRuntime)
		runtime.RunContainerStub = func(ctx context.Context, _ *model.ContainerCall, _ string, _ pubsub.EventPublisher, out, stderr io.WriteCloser) (*int64, error) {
			defer out.Close()
			defer stderr.Close()
			out.Write([]byte(strings.Repeat("x", containerResultLimit+1)))
			Expect(ctx.Err()).To(HaveOccurred())
			return nil, ctx.Err()
		}
		cc := newContainerCaller(runtime, new(fakes.FakePubSub), &_stateStore{})
		result, err := cc.Call(context.Background(), &model.ContainerCall{Image: &model.ContainerCallImage{}}, nil, &model.ContainerCallSpec{Stdout: "$(v)"}, "root")
		Expect(result).To(BeNil())
		Expect(err).To(MatchError(ContainSubstring("stdout exceeds 1 MiB")))
	})
	It("closes logging pipes after runtime start failure and rejects bootstrap result bindings", func() {
		runtime := new(runtimefakes.FakeContainerRuntime)
		runtime.RunContainerReturns(nil, errors.New("start failed"))
		cc := newContainerCaller(runtime, new(fakes.FakePubSub), &_stateStore{})
		_, err := cc.Call(context.Background(), &model.ContainerCall{Image: &model.ContainerCallImage{}}, nil, &model.ContainerCallSpec{Stderr: "$(v)"}, "root")
		Expect(err).To(MatchError("start failed"))
		_, err = cc.Call(withStdoutCapture(context.Background(), io.Discard), &model.ContainerCall{Image: &model.ContainerCallImage{}}, nil, &model.ContainerCallSpec{Stdout: "$(v)"}, "root")
		Expect(err).To(MatchError(ContainSubstring("bootstrap")))
	})
})
