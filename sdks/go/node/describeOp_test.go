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
