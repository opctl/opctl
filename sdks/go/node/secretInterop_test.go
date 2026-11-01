package node

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/outputs"
)

var _ = Describe("secret op interop", func() {
	It("carries renamed secrets through YAML, serial, parallel, and both loops with private snapshots", func() {
		dir, err := os.MkdirTemp("", "secret-interop-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		Expect(os.WriteFile(filepath.Join(dir, "op.yml"), []byte("name: passthrough\ninputs:\n  renamed:\n    string: {}\noutputs:\n  renamed:\n    string: {}\nrun:\n  serial:\n    - container:\n        image: {ref: alpine}\n        envVars:\n          TOKEN: $(renamed)\n"), 0600)).To(Succeed())
		db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()
		bus := pubsub.New(db)
		runtime := new(fakes.FakeContainerCaller)
		clr := newCaller(runtime, dir, bus, &_stateStore{}, &goBootstrapCache{})
		token := "roundtrip-private-value"
		items := []interface{}{token, token}
		iteration := "$(item)"
		op := &model.CallSpec{Op: &model.OpCallSpec{Ref: dir, PullCreds: &model.CredsSpec{}, Inputs: map[string]interface{}{"renamed": "$(token)"}, Outputs: map[string]string{"renamed": "$(alias)"}}}
		consume := func(ref string) *model.CallSpec {
			return &model.CallSpec{Container: &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}, EnvVars: map[string]interface{}{"TOKEN": ref}}}
		}
		parallel := []*model.CallSpec{consume("$(alias)"), consume("$(alias)")}
		children := []*model.CallSpec{op, {Parallel: &parallel}, {SerialLoop: &model.SerialLoopCallSpec{Range: "$(items)", Vars: &model.LoopVarsSpec{Value: &iteration}, Run: *consume("$(item)")}}, {ParallelLoop: &model.ParallelLoopCallSpec{Range: "$(items)", Vars: &model.LoopVarsSpec{Value: &iteration}, Run: *consume("$(item)")}}}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ctx = withPrivateGoInputs(ctx, map[string]*model.ParamSpec{"token": {String: &model.StringParamSpec{IsSecret: true}}, "items": {Array: &model.ArrayParamSpec{IsSecret: true}}})
		result, err := clr.Call(ctx, "root", map[string]*model.Value{"token": {String: &token}, "items": {Array: &items}}, &model.CallSpec{Serial: &children}, dir, nil, "root")
		Expect(err).NotTo(HaveOccurred())
		Expect(result["alias"].Sensitive).To(BeTrue())
		Expect(*result["alias"].String).To(Equal(token))
		Expect(runtime.CallCallCount()).To(Equal(7))
		for i := 0; i < runtime.CallCallCount(); i++ {
			_, call, _, _, _ := runtime.CallArgsForCall(i)
			Expect(call.EnvVars["TOKEN"]).To(Equal(token))
		}
		replay, cancelReplay := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancelReplay()
		events, err := bus.Subscribe(replay, model.EventFilter{Roots: []string{"root"}})
		Expect(err).NotTo(HaveOccurred())
		for {
			select {
			case <-replay.Done():
				return
			case event := <-events:
				data, err := json.Marshal(event)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).NotTo(ContainSubstring(token))
			}
		}
	})
	It("reports sensitive coercion and constraint failures without their bytes", func() {
		// These public collectors are also used at both op boundaries.
		token := "sensitive-invalid-number"
		for _, param := range []*model.ParamSpec{{Number: &model.NumberParamSpec{}}, {String: &model.StringParamSpec{Constraints: model.Constraints{"enum": []interface{}{"allowed"}}}}} {
			err := validateSensitiveOutput(token, param)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).NotTo(ContainSubstring(token))
		}
	})
})

func validateSensitiveOutput(token string, param *model.ParamSpec) error {
	_, err := outputs.Interpret(map[string]*model.Value{"out": {String: &token, Sensitive: true}}, map[string]*model.ParamSpec{"out": param}, nil, "", "")
	return err
}

var _ = Describe("nested Go privacy activation", func() {
	It("keeps a secret Go result private through its enclosing YAML root", func() {
		dir, err := os.MkdirTemp("", "nested-go-secret-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		Expect(os.Mkdir(filepath.Join(dir, "go"), 0700)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(dir, "go", "op.go"), []byte("package main"), 0600)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(dir, "op.yml"), []byte("name: root\noutputs:\n  alias:\n    string: {}\nrun:\n  serial:\n    - op:\n        ref: ./go\n        outputs:\n          result: $(alias)\n    - container:\n        image: {ref: alpine}\n        envVars:\n          TOKEN: $(alias)\n"), 0600)).To(Succeed())
		db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		bus := pubsub.New(db)
		runtime := new(fakes.FakeContainerCaller)
		clr := newCaller(runtime, dir, bus, newStateStore(ctx, db, bus), &goBootstrapCache{}).(*_caller)
		secret := "nested-go-private-result"
		bootstrap := newGoBootstrap()
		bootstrap.DescribeReturns(&model.OpSpec{Outputs: map[string]*model.ParamSpec{"result": {String: &model.StringParamSpec{Default: secret, IsSecret: true}}}}, nil)
		bootstrap.PlanReturns(&goop.Plan{OpSpec: &model.OpSpec{Run: &model.CallSpec{Container: &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}}}}}, nil)
		clr.opCaller = _opCaller{caller: clr, bootstrapper: bootstrap, dataDirPath: dir}
		returned, err := clr.Call(ctx, "root", nil, &model.CallSpec{Op: &model.OpCallSpec{Ref: dir, PullCreds: &model.CredsSpec{}, Outputs: map[string]string{"alias": "$(alias)"}}}, dir, nil, "root")
		Expect(err).NotTo(HaveOccurred())
		Expect(returned["alias"].Sensitive).To(BeTrue())
		Expect(*returned["alias"].String).To(Equal(secret))
		_, last, _, _, _ := runtime.CallArgsForCall(runtime.CallCallCount() - 1)
		Expect(last.EnvVars["TOKEN"]).To(Equal(secret))
		replay, cancelReplay := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancelReplay()
		events, err := bus.Subscribe(replay, model.EventFilter{})
		Expect(err).NotTo(HaveOccurred())
		for {
			select {
			case <-replay.Done():
				return
			case event := <-events:
				data, err := json.Marshal(event)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).NotTo(ContainSubstring(secret))
			}
		}
	})
})
