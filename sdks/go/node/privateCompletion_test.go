package node

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

var _ = Describe("private Go graph completions", func() {
	It("retains runtime values through every group and loop without publishing them", func() {
		dir, err := os.MkdirTemp("", "private-completion-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()
		bus := pubsub.New(db)
		runtime := new(fakes.FakeContainerCaller)
		token := "completion-private-token"
		runtime.CallReturns(map[string]*model.Value{"token": {String: &token, Sensitive: true}}, nil)
		leaf := &model.CallSpec{Container: &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}, EnvVars: map[string]interface{}{"TOKEN": "$(seed)"}}}
		list := []*model.CallSpec{leaf, leaf}
		cases := []*model.CallSpec{
			{Serial: &list}, {Parallel: &list},
			{SerialLoop: &model.SerialLoopCallSpec{Range: []interface{}{1, 2}, Run: *leaf}},
			{ParallelLoop: &model.ParallelLoopCallSpec{Range: []interface{}{1, 2}, Run: *leaf}},
		}
		for _, spec := range cases {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			ctx = withPrivateGoInputs(ctx, map[string]*model.ParamSpec{"token": {String: &model.StringParamSpec{IsSecret: true}}, "seed": {String: &model.StringParamSpec{IsSecret: true}}})
			caller := newCaller(runtime, dir, bus, &_stateStore{})
			results, err := caller.Call(ctx, "root", map[string]*model.Value{"seed": {String: &token}}, spec, dir, nil, "root")
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveKey("token"))
			Expect(*results["token"].String).To(Equal(token))
			cancel()
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		stream, err := bus.Subscribe(ctx, model.EventFilter{})
		Expect(err).NotTo(HaveOccurred())
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-stream:
				data, err := json.Marshal(event)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).NotTo(ContainSubstring(token))
			}
		}
	})
})
