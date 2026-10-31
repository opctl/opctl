package node

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
)

var _ = Describe("Go op calls", func() {
	var dir string
	BeforeEach(func() {
		var err error
		dir, err = os.MkdirTemp("", "go-op-")
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(filepath.Join(dir, "op.go"), []byte("package main"), 0600)).To(Succeed())
	})
	AfterEach(func() { Expect(os.RemoveAll(dir)).To(Succeed()) })
	It("runs the planned op graph as the op's child call", func() {
		caller := new(fakes.FakeCaller)
		bootstrap := new(fakes.FakeBootstrapper)
		plan := &model.OpSpec{Run: &model.CallSpec{Container: &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine:3"}}}}
		bootstrap.PlanReturns(plan, nil)
		opCall := &model.OpCall{BaseCall: model.BaseCall{OpPath: dir}, OpID: "op", ChildCallID: "child"}
		_, err := (_opCaller{caller: caller, bootstrapper: bootstrap}).Call(context.Background(), opCall, nil, "root", &model.OpCallSpec{})
		Expect(err).NotTo(HaveOccurred())
		Expect(bootstrap.PlanCallCount()).To(Equal(1))
		_, _, scope, spec, _, parent, root := caller.CallArgsForCall(0)
		Expect(spec).To(Equal(plan.Run))
		Expect(parent).To(Equal(&opCall.OpID))
		Expect(root).To(Equal("root"))
		Expect(*scope["./"].Dir).To(Equal(dir))
	})
	It("does not run a child when bootstrap fails", func() {
		caller := new(fakes.FakeCaller)
		bootstrap := new(fakes.FakeBootstrapper)
		failure := errors.New("bootstrap failed")
		bootstrap.PlanReturns(nil, failure)
		_, err := (_opCaller{caller: caller, bootstrapper: bootstrap}).Call(context.Background(), &model.OpCall{BaseCall: model.BaseCall{OpPath: dir}}, nil, "root", &model.OpCallSpec{})
		Expect(err).To(Equal(failure))
		Expect(caller.CallCallCount()).To(Equal(0))
	})
	It("rejects both declarations before bootstrap or a child call", func() {
		Expect(os.WriteFile(filepath.Join(dir, "op.yml"), []byte("name: both"), 0600)).To(Succeed())
		caller := new(fakes.FakeCaller)
		bootstrap := new(fakes.FakeBootstrapper)
		_, err := (_opCaller{caller: caller, bootstrapper: bootstrap}).Call(context.Background(), &model.OpCall{BaseCall: model.BaseCall{OpPath: dir}}, nil, "root", &model.OpCallSpec{})
		Expect(err).To(HaveOccurred())
		Expect(bootstrap.PlanCallCount()).To(Equal(0))
		Expect(caller.CallCallCount()).To(Equal(0))
	})
	It("starts a root Go op without implicit output bindings", func() {
		caller := new(fakes.FakeCaller)
		done := make(chan struct{})
		caller.CallStub = func(_ context.Context, _ string, _ map[string]*model.Value, _ *model.CallSpec, _ string, _ *string, _ string) (map[string]*model.Value, error) {
			defer close(done)
			return nil, nil
		}
		id, err := (core{caller: caller}).StartOp(context.Background(), model.StartOpReq{Op: model.StartOpReqOp{Ref: dir, PullCreds: &model.Creds{}}})
		Expect(err).NotTo(HaveOccurred())
		Expect(id).To(HaveLen(32))
		Eventually(done).Should(BeClosed())
		_, _, _, spec, _, _, _ := caller.CallArgsForCall(0)
		Expect(spec.Op.Outputs).To(BeEmpty())
	})
})
