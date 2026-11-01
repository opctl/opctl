package node

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	runtimeFakes "github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

var _ = Describe("DescribeOp", func() {
	It("reads YAML directly and bootstraps a Go interface as a root", func() {
		root, err := os.MkdirTemp("", "describe-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(root)
		bootstrap := new(fakes.FakeBootstrapper)
		spec := &model.OpSpec{Name: "go", Inputs: map[string]*model.ParamSpec{"env": {String: &model.StringParamSpec{Default: "dev"}}}}
		bootstrap.DescribeReturns(spec, nil)
		core := core{bootstrapper: bootstrap}
		req := model.DescribeOpReq{Op: model.StartOpReqOp{Ref: root, PullCreds: &model.Creds{}}}
		Expect(os.WriteFile(filepath.Join(root, "op.yml"), []byte("name: yaml"), 0600)).To(Succeed())
		actual, err := core.DescribeOp(context.Background(), req)
		Expect(err).NotTo(HaveOccurred())
		Expect(actual.Name).To(Equal("yaml"))
		Expect(bootstrap.DescribeCallCount()).To(Equal(0))
		Expect(os.Remove(filepath.Join(root, "op.yml"))).To(Succeed())
		Expect(os.WriteFile(filepath.Join(root, "op.go"), []byte("package main"), 0600)).To(Succeed())
		actual, err = core.DescribeOp(context.Background(), req)
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).To(Equal(spec))
		_, call, rootID := bootstrap.DescribeArgsForCall(0)
		Expect(call.OpPath).To(Equal(root))
		Expect(call.OpID).To(BeEmpty())
		Expect(rootID).To(BeEmpty())
		Expect(os.WriteFile(filepath.Join(root, "op.yml"), []byte("name: both"), 0600)).To(Succeed())
		_, err = core.DescribeOp(context.Background(), req)
		var declaration *OpDeclarationError
		Expect(errors.As(err, &declaration)).To(BeTrue())
	})
})

var _ = Describe("DescribeOp error classification", func() {
	for _, tc := range []struct {
		name                                    string
		data                                    string
		module                                  *string
		opSuffix                                string
		exit                                    int64
		failure                                 error
		cancelBefore, cancelDuring, declaration bool
	}{
		{name: "classifies missing go.mod as a declaration error", module: stringPtr(""), declaration: true},
		{name: "classifies missing go directive as a declaration error", module: stringPtr("module hello\n"), declaration: true},
		{name: "classifies malformed go.mod as a declaration error", module: stringPtr("not a module file"), declaration: true},
		{name: "classifies interpolated paths as a declaration error", opSuffix: "$(path)", declaration: true},
		{name: "classifies forbidden mounts as a declaration error", module: stringPtr("module hello\ngo 1.21\nreplace example.com/lib => /usr\n"), declaration: true},
		{name: "classifies bootstrap program exits as declaration errors", exit: 1, declaration: true},
		{name: "classifies bootstrap JSON parse failures as declaration errors", data: "not JSON", declaration: true},
		{name: "classifies bootstrap contract failures as declaration errors", data: `{"contractVersion":1,"opspecSchema":2}`, declaration: true},
		{name: "preserves runtime failures", failure: errors.New("Docker unavailable")},
		{name: "preserves runtime failures even with a nonzero exit", failure: errors.New("wait failed"), exit: 1},
		{name: "preserves context cancellation errors", failure: context.Canceled},
		{name: "preserves context deadline errors", failure: context.DeadlineExceeded},
		{name: "reports cancellation before bootstrap", cancelBefore: true, data: `{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"hello"}}`},
		{name: "reports cancellation when the runtime returns no error", cancelDuring: true, data: `{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"hello"}}`},
	} {
		tc := tc
		It(tc.name, func() {
			root, err := os.MkdirTemp("", "describe-error-")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(root)
			if tc.opSuffix != "" {
				root = filepath.Join(root, tc.opSuffix)
				Expect(os.MkdirAll(root, 0700)).To(Succeed())
			}
			Expect(os.WriteFile(filepath.Join(root, "op.go"), []byte("package main"), 0600)).To(Succeed())
			module := "module hello\ngo 1.21\n"
			if tc.module != nil {
				module = *tc.module
			}
			if module != "" {
				Expect(os.WriteFile(filepath.Join(root, "go.mod"), []byte(module), 0600)).To(Succeed())
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			runtime := new(runtimeFakes.FakeContainerRuntime)
			runtime.RunContainerStub = func(_ context.Context, _ *model.ContainerCall, _ string, _ pubsub.EventPublisher, stdout, stderr io.WriteCloser) (*int64, error) {
				defer stdout.Close()
				defer stderr.Close()
				_, err := io.WriteString(stdout, tc.data)
				Expect(err).NotTo(HaveOccurred())
				if tc.cancelDuring {
					cancel()
				}
				return &tc.exit, tc.failure
			}
			containerCaller := newContainerCaller(runtime, new(fakes.FakePubSub), nil)
			caller := new(fakes.FakeCaller)
			caller.CallStub = func(ctx context.Context, _ string, scope map[string]*model.Value, spec *model.CallSpec, _ string, _ *string, rootID string) (map[string]*model.Value, error) {
				return containerCaller.Call(ctx, &model.ContainerCall{Image: &model.ContainerCallImage{}}, scope, spec.Container, rootID)
			}
			core := core{bootstrapper: newBootstrapper(caller, filepath.Join(root, "data"))}
			if tc.cancelBefore {
				cancel()
			}
			_, err = core.DescribeOp(ctx, model.DescribeOpReq{Op: model.StartOpReqOp{Ref: root, PullCreds: &model.Creds{}}})
			Expect(err).To(HaveOccurred())
			var declaration *OpDeclarationError
			Expect(errors.As(err, &declaration)).To(Equal(tc.declaration))
			if tc.failure != nil {
				Expect(errors.Is(err, tc.failure)).To(BeTrue())
			}
			if tc.cancelBefore || tc.cancelDuring {
				Expect(errors.Is(err, context.Canceled)).To(BeTrue())
			}
		})
	}
})

func stringPtr(value string) *string { return &value }
