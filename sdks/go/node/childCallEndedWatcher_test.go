package node

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

var _ = Context("serial callers", func() {
	for _, testCase := range []struct {
		name string
		call func(context.Context, caller, pubsub.PubSub, map[string]*model.Value) (map[string]*model.Value, error)
	}{
		{
			name: "serialCaller",
			call: func(ctx context.Context, caller caller, pubSub pubsub.PubSub, scope map[string]*model.Value) (map[string]*model.Value, error) {
				return newSerialCaller(caller, pubSub).Call(ctx, "callID", scope, "rootCallID", "opPath", []*model.CallSpec{{}, {}})
			},
		},
		{
			name: "serialLoopCaller",
			call: func(ctx context.Context, caller caller, pubSub pubsub.PubSub, scope map[string]*model.Value) (map[string]*model.Value, error) {
				return newSerialLoopCaller(caller, pubSub).Call(ctx, "callID", scope, model.SerialLoopCallSpec{Range: []interface{}{0, 1}}, "opPath", nil, "rootCallID")
			},
		},
	} {
		testCase := testCase
		Context(testCase.name, func() {
			It("drains events while children run and keeps only their first CallEnded", func() {
				pubSub := &serialCallerTestPubSub{events: make(chan model.Event)}
				fakeCaller := new(FakeCaller)
				initial, first, second := "initial", "first", "second"
				unchanged := true
				values := []*model.Value{{String: &initial}, {String: &first}, {String: &second}}
				fakeCaller.CallStub = func(ctx context.Context, id string, scope map[string]*model.Value, callSpec *model.CallSpec, opPath string, parentCallID *string, rootCallID string) (map[string]*model.Value, error) {
					index := fakeCaller.CallCallCount() - 1
					Expect(scope["value"]).To(Equal(values[index]))
					for i := 0; i < 100; i++ {
						pubSub.Publish(model.Event{ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{
							RootCallID: rootCallID,
							Data:       []byte("container log"),
						}})
					}
					pubSub.Publish(model.Event{CallEnded: &model.CallEnded{
						Call:  model.Call{ID: "otherChild", RootID: rootCallID},
						Error: &model.CallEndedError{Message: "unrelated error"},
					}})
					pubSub.Publish(model.Event{CallEnded: &model.CallEnded{
						Call:    model.Call{ID: id, RootID: rootCallID},
						Outputs: map[string]*model.Value{"value": values[index+1]},
						Outcome: model.OpOutcomeKilled,
					}})
					for i := 0; i < 3; i++ {
						pubSub.Publish(model.Event{CallEnded: &model.CallEnded{
							Call:  model.Call{ID: id, RootID: rootCallID},
							Error: &model.CallEndedError{Message: "duplicate error"},
						}})
					}
					return nil, nil
				}

				outputs, err := testCase.call(context.Background(), fakeCaller, pubSub, map[string]*model.Value{
					"value": values[0], "unchanged": {Boolean: &unchanged},
				})

				Expect(err).To(BeNil())
				Expect(fakeCaller.CallCallCount()).To(Equal(2))
				Expect(outputs).To(Equal(map[string]*model.Value{
					"value": values[2], "unchanged": {Boolean: &unchanged},
				}))
				Eventually(pubSub.events).Should(BeClosed())
			})

			It("stops on a direct child's CallEnded error", func() {
				pubSub := &serialCallerTestPubSub{events: make(chan model.Event)}
				fakeCaller := new(FakeCaller)
				fakeCaller.CallStub = func(ctx context.Context, id string, scope map[string]*model.Value, callSpec *model.CallSpec, opPath string, parentCallID *string, rootCallID string) (map[string]*model.Value, error) {
					pubSub.Publish(model.Event{CallEnded: &model.CallEnded{
						Call:  model.Call{ID: id, RootID: rootCallID},
						Error: &model.CallEndedError{Message: "child error"},
					}})
					return nil, nil
				}

				outputs, err := testCase.call(context.Background(), fakeCaller, pubSub, nil)

				Expect(err).To(MatchError("child error"))
				Expect(outputs).To(BeNil())
				Expect(fakeCaller.CallCallCount()).To(Equal(1))
				Eventually(pubSub.events).Should(BeClosed())
			})

			It("continues without an error when cancellation closes the subscription", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				pubSub := &serialCallerTestPubSub{events: make(chan model.Event)}
				fakeCaller := new(FakeCaller)
				fakeCaller.CallStub = func(ctx context.Context, id string, scope map[string]*model.Value, callSpec *model.CallSpec, opPath string, parentCallID *string, rootCallID string) (map[string]*model.Value, error) {
					cancel()
					return nil, nil
				}

				outputs, err := testCase.call(ctx, fakeCaller, pubSub, nil)

				Expect(err).To(BeNil())
				Expect(outputs).To(BeEmpty())
				Expect(fakeCaller.CallCallCount()).To(Equal(2))
				Eventually(pubSub.events).Should(BeClosed())
			})
		})
	}
})

type serialCallerTestPubSub struct {
	pubsub.PubSub
	events chan model.Event
}

func (ps *serialCallerTestPubSub) Subscribe(ctx context.Context, filter model.EventFilter) (<-chan model.Event, error) {
	go func() {
		<-ctx.Done()
		close(ps.events)
	}()
	return ps.events, nil
}

func (ps *serialCallerTestPubSub) Publish(event model.Event) {
	select {
	case ps.events <- event:
	case <-time.After(time.Second):
		Fail("subscription did not read an event while the child was running")
	}
}
