package node

import (
	"context"
	"encoding/json"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
)

var _ = Describe("private Go input event bindings", func() {
	It("publishes secret bindings as handles while keeping runtime values", func() {
		call := &model.Call{Container: &model.ContainerCall{EnvVars: map[string]string{"TOKEN": "private-token", "PUBLIC": "dev"}, Files: map[string]string{"/config": "private-path"}}}
		spec := &model.CallSpec{Container: &model.ContainerCallSpec{EnvVars: map[string]interface{}{"TOKEN": "$(token)", "PUBLIC": "dev"}, Files: map[string]interface{}{"/config": "$(config)"}}}
		ctx := withPrivateGoInputs(context.Background(), map[string]*model.ParamSpec{"token": {String: &model.StringParamSpec{IsSecret: true}}, "config": {File: &model.FileParamSpec{IsSecret: true}}})
		snapshot := callForEvent(ctx, call, spec)
		encoded, err := json.Marshal(snapshot)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(encoded)).NotTo(ContainSubstring("private-"))
		Expect(snapshot.Container.EnvVars["TOKEN"]).To(Equal("$(token)"))
		Expect(snapshot.Container.EnvVars["PUBLIC"]).To(Equal("dev"))
		Expect(call.Container.EnvVars["TOKEN"]).To(Equal("private-token"))
		Expect(call.Container.Files["/config"]).To(Equal("private-path"))
		Expect(callForEvent(context.Background(), call, spec)).To(Equal(*call))
	})
	It("keeps secret values out of both start and completion events", func() {
		dir, err := os.MkdirTemp("", "private-go-events-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		publisher := new(fakes.FakePubSub)
		closed := make(chan model.Event)
		close(closed)
		publisher.SubscribeReturns(closed, nil)
		runtime := new(fakes.FakeContainerCaller)
		caller := _caller{pubSub: publisher, containerCaller: runtime, dataDirPath: dir}
		secret := "private-token"
		runtime.CallReturns(map[string]*model.Value{"token": {String: &secret}}, nil)
		ctx := withPrivateGoInputs(context.Background(), map[string]*model.ParamSpec{"token": {String: &model.StringParamSpec{IsSecret: true}}})
		returned, err := caller.Call(ctx, "child", map[string]*model.Value{"token": {String: &secret}}, &model.CallSpec{Container: &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}, EnvVars: map[string]interface{}{"TOKEN": "$(token)"}}}, dir, nil, "root")
		Expect(err).NotTo(HaveOccurred())
		Expect(publisher.PublishCallCount()).To(Equal(2))
		Expect(runtime.CallCallCount()).To(Equal(1))
		Expect(*returned["token"].String).To(Equal(secret))
		_, actual, _, _, _ := runtime.CallArgsForCall(0)
		Expect(actual.EnvVars["TOKEN"]).To(Equal(secret))
		for i := 0; i < publisher.PublishCallCount(); i++ {
			encoded, err := json.Marshal(publisher.PublishArgsForCall(i))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(encoded)).NotTo(ContainSubstring(secret))
		}
	})

})
