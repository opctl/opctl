package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

func TestOutputs(t *testing.T) { RegisterFailHandler(Fail); RunSpecs(t, "scenario output assertions") }

var _ = Describe("boxed completion assertions", func() {
	completion := func(path string, outputs map[string]*model.Value) model.Event {
		return model.Event{CallEnded: &model.CallEnded{Call: model.Call{ID: "root", RootID: "root", Op: &model.OpCall{BaseCall: model.BaseCall{OpPath: path}}}, Outcome: model.OpOutcomeSucceeded, Outputs: outputs}}
	}
	It("ignores other roots, nested calls, and bootstrap completions and compares replayed structured boxes", func() {
		events := make(chan model.Event, 4)
		events <- model.Event{CallEnded: &model.CallEnded{Call: model.Call{ID: "bootstrap", RootID: "bootstrap", Container: &model.ContainerCall{}}}}
		nested := completion("/op", nil)
		parent := "parent"
		nested.CallEnded.Call.ParentID = &parent
		events <- nested
		events <- completion("/other", nil)
		var replayed model.Event
		Expect(json.Unmarshal([]byte(`{"callEnded":{"call":{"id":"root","rootId":"root","op":{"opPath":"/op"}},"outcome":"SUCCEEDED","outputs":{"result":{"object":{"left":{"string":"a"},"items":{"array":[{"number":7}]}}}}}}`), &replayed)).To(Succeed())
		events <- replayed
		Expect(assertOutputs(context.Background(), events, "/op", []byte(`{"result":{"object":{"items":{"array":[{"number":7.0}]},"left":{"string":"a"}}}}`))).To(Succeed())
	})
	It("treats an authored empty object as an assertion and rejects extra output keys", func() {
		events := make(chan model.Event, 1)
		events <- completion("/op", nil)
		Expect(assertOutputs(context.Background(), events, "/op", []byte(`{}`))).To(Succeed())
		text := "value"
		events <- completion("/op", map[string]*model.Value{"extra": {String: &text}})
		Expect(assertOutputs(context.Background(), events, "/op", []byte(`{}`))).To(MatchError(ContainSubstring("do not match")))
	})
	It("times out even if the event stream remains open and rejects malformed expectations", func() {
		events := make(chan model.Event)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		Expect(assertOutputs(ctx, events, "/op", []byte(`{}`))).To(MatchError(ContainSubstring("waiting")))
		Expect(assertOutputs(context.Background(), events, "/op", []byte(`null`))).To(HaveOccurred())
	})
})
