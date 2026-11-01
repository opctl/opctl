package node

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
)

var _ = Describe("Go op interop", func() {
	It("guards secret YAML inputs before publishing a resolved call or running it", func() {
		dir, err := os.MkdirTemp("", "go-interop-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		Expect(os.WriteFile(filepath.Join(dir, "op.yml"), []byte("name: secret\ninputs:\n  renamed:\n    string: {}\nrun:\n  container:\n    image: {ref: alpine}\n"), 0600)).To(Succeed())
		bus := new(fakes.FakePubSub)
		stream := make(chan model.Event)
		close(stream)
		bus.SubscribeReturns(stream, nil)
		op := new(fakes.FakeOpCaller)
		clr := _caller{pubSub: bus, opCaller: op, stateStore: &_stateStore{}, dataDirPath: dir}
		token := "interop-private-token"
		ctx := withPrivateGoInputs(context.Background(), map[string]*model.ParamSpec{"token": {String: &model.StringParamSpec{IsSecret: true}}})
		_, err = clr.Call(ctx, "child", map[string]*model.Value{"token": {String: &token}}, &model.CallSpec{Op: &model.OpCallSpec{Ref: dir, PullCreds: &model.CredsSpec{}, Inputs: map[string]interface{}{"renamed": "$(token)"}}}, dir, nil, "root")
		Expect(err).To(MatchError(ContainSubstring("secret interop")))
		Expect(op.CallCallCount()).To(BeZero())
		for i := 0; i < bus.PublishCallCount(); i++ {
			data, err := json.Marshal(bus.PublishArgsForCall(i))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).NotTo(ContainSubstring(token))
		}
	})
	It("copies nested credential snapshots and keeps handle expressions", func() {
		child := &model.CallSpec{Op: &model.OpCallSpec{Ref: "repo#v1", PullCreds: &model.CredsSpec{Username: "literal-user", Password: "literal-password"}}}
		children := []*model.CallSpec{child}
		call := &model.Call{Serial: children}
		ctx := withPrivateGoInputs(context.Background(), nil)
		snapshot := callForEvent(ctx, call, &model.CallSpec{Serial: &children})
		encoded, err := json.Marshal(snapshot)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(encoded)).NotTo(ContainSubstring("literal-password"))
		Expect(child.Op.PullCreds.Password).To(Equal("literal-password"))
		child.Op.PullCreds.Password = "$(token)"
		snapshot = callForEvent(ctx, call, &model.CallSpec{Serial: &children})
		Expect(snapshot.Serial[0].Op.PullCreds.Password).To(Equal("$(token)"))
	})
})
