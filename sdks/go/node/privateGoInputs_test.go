package node

import (
	"context"
	"encoding/json"
	"os"
	"strings"

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
	It("keeps shorthand env, dir, file, and bare socket bindings out of events", func() {
		dir, err := os.MkdirTemp("", "shorthand-private-path-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		file := dir + "/file"
		Expect(os.WriteFile(file, nil, 0600)).To(Succeed())
		secret := "shorthand-private-token"
		socket := dir + "/socket"
		publisher := new(fakes.FakePubSub)
		closed := make(chan model.Event)
		close(closed)
		publisher.SubscribeReturns(closed, nil)
		runtime := new(fakes.FakeContainerCaller)
		caller := _caller{pubSub: publisher, containerCaller: runtime, dataDirPath: dir}
		ctx := withPrivateGoInputs(context.Background(), map[string]*model.ParamSpec{
			"renamed":  {String: &model.StringParamSpec{IsSecret: true}},
			"artifact": {Dir: &model.DirParamSpec{IsSecret: true}},
			"config":   {File: &model.FileParamSpec{IsSecret: true}},
			"socket":   {Socket: &model.SocketParamSpec{IsSecret: true}},
		})
		scope := map[string]*model.Value{"renamed": {String: &secret}, "artifact": {Dir: &dir}, "config": {File: &file}, "socket": {Socket: &socket}}
		spec := &model.CallSpec{Container: &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}, EnvVars: map[string]interface{}{"renamed": nil}, Dirs: map[string]interface{}{"artifact": nil}, Files: map[string]interface{}{"config": nil}, Sockets: map[string]string{"/socket": "socket"}}}
		_, err = caller.Call(ctx, "child", scope, spec, dir, nil, "root")
		Expect(err).NotTo(HaveOccurred())
		_, actual, _, _, _ := runtime.CallArgsForCall(0)
		Expect(actual.EnvVars["renamed"]).To(Equal(secret))
		Expect(actual.Dirs["artifact"]).To(Equal(dir))
		Expect(actual.Files["config"]).To(Equal(file))
		Expect(actual.Sockets["/socket"]).To(Equal(socket))
		for i := 0; i < publisher.PublishCallCount(); i++ {
			encoded, err := json.Marshal(publisher.PublishArgsForCall(i))
			Expect(err).NotTo(HaveOccurred())
			Expect(strings.Contains(string(encoded), secret)).To(BeFalse())
			event := publisher.PublishArgsForCall(i)
			if event.CallStarted != nil {
				container := event.CallStarted.Call.Container
				Expect(container.Dirs["artifact"]).To(Equal("$(artifact)"))
				Expect(container.Files["config"]).To(Equal("$(config)"))
				Expect(container.Sockets["/socket"]).To(Equal("$(socket)"))
			}
		}
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
		runtime.CallReturns(map[string]*model.Value{"token": {String: &secret, Sensitive: true}}, nil)
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
			Expect(strings.Contains(string(encoded), secret)).To(BeFalse())
		}
	})

})

var _ = Describe("sensitive interpretation failures", func() {
	It("keeps sensitive command and credential conversion errors private while preserving YAML diagnostics", func() {
		dir, err := os.MkdirTemp("", "private-error-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		secret := "/private-interpretation-path"
		specs := []*model.CallSpec{
			{Container: &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}, Cmd: []interface{}{"$(source)"}}},
			{Op: &model.OpCallSpec{Ref: "./callee", PullCreds: &model.CredsSpec{Username: "user", Password: "$(source)"}}},
		}
		for _, private := range []bool{true, false} {
			for _, spec := range specs {
				bus := new(fakes.FakePubSub)
				runtime := new(fakes.FakeContainerCaller)
				clr := _caller{pubSub: bus, containerCaller: runtime, dataDirPath: dir, stateStore: &_stateStore{}}
				ctx := context.Background()
				if private {
					ctx = withPrivateGoInputs(ctx, map[string]*model.ParamSpec{"source": {Dir: &model.DirParamSpec{IsSecret: true}}})
				}
				_, err := clr.Call(ctx, "child", map[string]*model.Value{"source": {Dir: &secret}}, spec, "/op", nil, "root")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), secret)).To(Equal(!private))
				if private {
					Expect(err.Error()).To(ContainSubstring(`sensitive bindings ["source"]`))
				}
				Expect(runtime.CallCallCount()).To(BeZero())
				Expect(bus.PublishCallCount()).To(Equal(1))
				encoded, marshalErr := json.Marshal(bus.PublishArgsForCall(0))
				Expect(marshalErr).NotTo(HaveOccurred())
				Expect(strings.Contains(string(encoded), secret)).To(Equal(!private))
			}
		}
	})
})

var _ = Describe("private interpreted call metadata", func() {
	It("keeps resolved names, image credentials, and dynamic environments out of snapshots", func() {
		token := "metadata-private-token"
		expression := "$(token)"
		ctx := withPrivateGoInputs(context.Background(), map[string]*model.ParamSpec{"token": {String: &model.StringParamSpec{IsSecret: true}}})
		call := &model.Call{Name: &token, Container: &model.ContainerCall{Name: &token, Image: &model.ContainerCallImage{Ref: &token, PullCreds: &model.Creds{Username: "user", Password: token}}, EnvVars: map[string]string{"TOKEN": token}}}
		spec := &model.CallSpec{Name: &expression, Container: &model.ContainerCallSpec{Name: &expression, Image: &model.ContainerCallImageSpec{Ref: expression}, EnvVars: expression}}
		encoded, err := json.Marshal(callForEvent(ctx, call, spec))
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Contains(string(encoded), token)).To(BeFalse())
		Expect(*call.Name).To(Equal(token))
		Expect(call.Container.Image.PullCreds.Password).To(Equal(token))
		spec.Container.EnvVars = map[string]interface{}{expression: "public"}
		call.Container.EnvVars = map[string]string{token: "public"}
		encoded, err = json.Marshal(callForEvent(ctx, call, spec))
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Contains(string(encoded), token)).To(BeFalse())
	})
})
