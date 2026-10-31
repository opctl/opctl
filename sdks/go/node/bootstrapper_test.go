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
	containerInterpreter "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container"
)

var _ = Describe("Go op bootstrapper", func() {
	var root, opDir, dataDir string
	var fixture []byte
	BeforeEach(func() {
		var err error
		root, err = os.MkdirTemp("", "bootstrap-")
		Expect(err).NotTo(HaveOccurred())
		opDir = filepath.Join(root, "hello")
		dataDir = filepath.Join(root, "data $(secret)")
		Expect(os.MkdirAll(opDir, 0700)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(opDir, "go.mod"), []byte("module hello\ngo 1.21\n"), 0600)).To(Succeed())
		fixture, err = os.ReadFile("../opspec/sdk/testdata/plan.hello.json")
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() { Expect(os.RemoveAll(root)).To(Succeed()) })
	It("runs a nested bootstrap call with a request file and private response capture", func() {
		caller := new(fakes.FakeCaller)
		opCall := &model.OpCall{OpID: "op", BaseCall: model.BaseCall{OpPath: opDir}}
		var requestPath string
		caller.CallStub = func(ctx context.Context, id string, scope map[string]*model.Value, spec *model.CallSpec, path string, parent *string, rootID string) (map[string]*model.Value, error) {
			Expect(id).To(HaveLen(32))
			Expect(parent).To(Equal(&opCall.OpID))
			Expect(rootID).To(Equal("root"))
			Expect(path).To(Equal(opDir))
			Expect(spec.Name).NotTo(BeNil())
			Expect(*spec.Name).To(Equal("bootstrap plan"))
			requestPath = *scope["request"].File
			request, err := os.ReadFile(requestPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(request).To(MatchJSON(`{"contractVersion":1,"phase":"plan","opspecSchema":{"min":1,"max":1},"inputs":{"env":"dev"},"knownOutputs":{}}`))
			interpreted, err := containerInterpreter.Interpret(scope, spec.Container, id, path, dataDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(interpreted.Files["/opctl/request.json"]).To(Equal(requestPath))
			Expect(interpreted.Dirs[opDir]).To(Equal(opDir))
			Expect(interpreted.Cmd).To(Equal([]string{"go", "run", ".", "plan"}))
			capture := stdoutCaptureFrom(ctx)
			Expect(capture).NotTo(BeNil())
			_, err = capture.Write(fixture)
			Expect(err).NotTo(HaveOccurred())
			return nil, nil
		}
		plan, err := newBootstrapper(caller, dataDir).Plan(context.Background(), opCall, "root", map[string]interface{}{"env": "dev"})
		Expect(err).NotTo(HaveOccurred())
		Expect(plan.Name).To(Equal("hello"))
		Expect(caller.CallCallCount()).To(Equal(1))
		Expect(requestPath).NotTo(BeAnExistingFile())
		Expect(filepath.Dir(requestPath)).NotTo(BeAnExistingFile())
	})
	It("cleans up the request on runtime and response failures", func() {
		for _, result := range []struct {
			data []byte
			err  error
		}{{nil, errors.New("compile failed")}, {[]byte("invalid JSON"), nil}, {[]byte(`{"contractVersion":1,"opspecSchema":2}`), nil}} {
			caller := new(fakes.FakeCaller)
			var request string
			caller.CallStub = func(ctx context.Context, _ string, scope map[string]*model.Value, _ *model.CallSpec, _ string, _ *string, _ string) (map[string]*model.Value, error) {
				request = *scope["request"].File
				_, err := stdoutCaptureFrom(ctx).Write(result.data)
				Expect(err).NotTo(HaveOccurred())
				return nil, result.err
			}
			_, err := newBootstrapper(caller, dataDir).Plan(context.Background(), &model.OpCall{BaseCall: model.BaseCall{OpPath: opDir}, OpID: "op"}, "root", nil)
			Expect(err).To(HaveOccurred())
			Expect(request).NotTo(BeAnExistingFile())
			if result.err != nil {
				Expect(errors.Is(err, result.err)).To(BeTrue())
			}
		}
	})
	It("fails before making a call when the module cannot be prepared", func() {
		Expect(os.Remove(filepath.Join(opDir, "go.mod"))).To(Succeed())
		caller := new(fakes.FakeCaller)
		_, err := newBootstrapper(caller, dataDir).Plan(context.Background(), &model.OpCall{BaseCall: model.BaseCall{OpPath: opDir}}, "root", nil)
		Expect(err).To(HaveOccurred())
		Expect(caller.CallCallCount()).To(Equal(0))
	})
	It("returns stderr diagnostics on bootstrap failure without exposing stdout", func() {
		caller := new(fakes.FakeCaller)
		failure := errors.New("nonzero container exit code: 1")
		caller.CallStub = func(ctx context.Context, _ string, _ map[string]*model.Value, _ *model.CallSpec, _ string, _ *string, _ string) (map[string]*model.Value, error) {
			capture := stdoutCaptureFrom(ctx)
			_, err := capture.Write([]byte("private plan output"))
			Expect(err).NotTo(HaveOccurred())
			_, err = capture.diagnostics.Write([]byte("input field X: unknown opctl tag key minimun\n"))
			Expect(err).NotTo(HaveOccurred())
			return nil, failure
		}
		_, err := newBootstrapper(caller, dataDir).Describe(context.Background(), &model.OpCall{BaseCall: model.BaseCall{OpPath: opDir}}, "")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unknown opctl tag key minimun"))
		Expect(err.Error()).NotTo(ContainSubstring("private plan output"))
		Expect(errors.Is(err, failure)).To(BeTrue())
	})
	It("nests describe under the op without input values", func() {
		caller := new(fakes.FakeCaller)
		caller.CallStub = func(ctx context.Context, id string, scope map[string]*model.Value, spec *model.CallSpec, path string, parent *string, rootID string) (map[string]*model.Value, error) {
			Expect(*parent).To(Equal("op"))
			Expect(rootID).To(Equal("root"))
			Expect(*spec.Name).To(Equal("bootstrap describe"))
			Expect(spec.Container.Cmd).To(Equal([]interface{}{"go", "run", ".", "describe"}))
			request, err := os.ReadFile(*scope["request"].File)
			Expect(err).NotTo(HaveOccurred())
			Expect(request).To(MatchJSON(`{"contractVersion":1,"phase":"describe","opspecSchema":{"min":1,"max":1},"inputs":{},"knownOutputs":{}}`))
			data, err := os.ReadFile("../opspec/sdk/testdata/describe.inputs.json")
			Expect(err).NotTo(HaveOccurred())
			_, err = stdoutCaptureFrom(ctx).Write(data)
			Expect(err).NotTo(HaveOccurred())
			return nil, nil
		}
		spec, err := newBootstrapper(caller, dataDir).Describe(context.Background(), &model.OpCall{BaseCall: model.BaseCall{OpPath: opDir}, OpID: "op"}, "root")
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Inputs["token"].String.IsSecret).To(BeTrue())
	})

})
