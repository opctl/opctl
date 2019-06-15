package core

import (
	"context"

	"github.com/opctl/sdk-golang/opspec/interpreter/call/loop/iteration"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/loop"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

var _ = Context("looper", func() {
	Context("newLooper", func() {
		It("should return looper", func() {
			/* arrange/act/assert */
			Expect(newLooper(
				new(fakeCaller),
				new(pubsub.Fake),
			)).To(Not(BeNil()))
		})
	})

	Context("Loop", func() {
		Context("initial dcgLoop.Until true", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeLoopInterpreter := new(loop.FakeInterpreter)
				until := true
				fakeLoopInterpreter.InterpretReturns(&model.DCGLoop{Until: &until}, nil)

				fakeCaller := new(fakeCaller)

				objectUnderTest := _looper{
					caller:              fakeCaller,
					loopDeScoper:        new(loop.FakeDeScoper),
					loopInterpreter:     fakeLoopInterpreter,
					iterationScoper:     new(iteration.FakeScoper),
					pubSub:              new(pubsub.Fake),
					uniqueStringFactory: new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.Loop(
					context.Background(),
					"id",
					map[string]*model.Value{},
					&model.SCG{Loop: &model.SCGLoop{}},
					new(data.FakeHandle),
					nil,
					"rootOpID",
				)

				/* assert */
				Expect(fakeCaller.CallCallCount()).To(Equal(0))
			})
		})
		Context("initial dcgLoop.For.Each empty", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeLoopInterpreter := new(loop.FakeInterpreter)
				fakeLoopInterpreter.InterpretReturns(
					&model.DCGLoop{
						For: &model.DCGLoopFor{
							Each: &model.Value{
								Array: new([]interface{}),
							},
						},
					},
					nil,
				)

				fakeCaller := new(fakeCaller)

				objectUnderTest := _looper{
					caller:              fakeCaller,
					loopDeScoper:        new(loop.FakeDeScoper),
					loopInterpreter:     fakeLoopInterpreter,
					iterationScoper:     new(iteration.FakeScoper),
					pubSub:              new(pubsub.Fake),
					uniqueStringFactory: new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.Loop(
					context.Background(),
					"id",
					map[string]*model.Value{},
					&model.SCG{Loop: &model.SCGLoop{}},
					new(data.FakeHandle),
					nil,
					"rootOpID",
				)

				/* assert */
				Expect(fakeCaller.CallCallCount()).To(Equal(0))
			})
		})
		Context("initial dcgLoop.Until false", func() {
			It("should call caller.Call w/ expected args", func() {
				/* arrange */
				providedCtx := context.Background()
				providedScope := map[string]*model.Value{}
				index := "index"
				providedSCG := &model.SCG{
					Loop: &model.SCGLoop{
						Index: &index,
					},
				}
				providedOpHandle := new(data.FakeHandle)
				providedParentCallIDValue := "providedParentCallID"
				providedParentCallID := &providedParentCallIDValue
				providedRootOpID := "providedRootOpID"

				fakeLoopInterpreter := new(loop.FakeInterpreter)
				until := false
				fakeLoopInterpreter.InterpretReturns(
					&model.DCGLoop{
						Until: &until,
						Index: providedSCG.Loop.Index,
					},
					nil,
				)

				fakeIterationScoper := new(iteration.FakeScoper)
				expectedScope := map[string]*model.Value{
					index: &model.Value{Number: new(float64)},
				}
				fakeIterationScoper.ScopeReturns(expectedScope, nil)

				fakeCaller := new(fakeCaller)

				callID := "callID"

				expectedErrorMessage := "expectedErrorMessage"
				fakePubSub := new(pubsub.Fake)
				eventChannel := make(chan model.Event, 100)
				fakePubSub.SubscribeStub = func(ctx context.Context, filter model.EventFilter) (<-chan model.Event, <-chan error) {
					eventChannel <- model.Event{
						CallEnded: &model.CallEndedEvent{
							CallID: callID,
							Error: &model.CallEndedEventError{
								Message: expectedErrorMessage,
							},
						},
					}

					return eventChannel, make(chan error)
				}

				fakeUniqueStringFactory := new(uniquestring.Fake)
				fakeUniqueStringFactory.ConstructReturns(callID, nil)

				objectUnderTest := _looper{
					caller:              fakeCaller,
					loopDeScoper:        new(loop.FakeDeScoper),
					loopInterpreter:     fakeLoopInterpreter,
					iterationScoper:     fakeIterationScoper,
					pubSub:              fakePubSub,
					uniqueStringFactory: fakeUniqueStringFactory,
				}

				/* act */
				objectUnderTest.Loop(
					providedCtx,
					"id",
					providedScope,
					providedSCG,
					providedOpHandle,
					providedParentCallID,
					providedRootOpID,
				)

				/* assert */
				actualCtx,
					actualCallID,
					actualScope,
					actualSCG,
					actualOpHandle,
					actualParentCallID,
					actualRootOpID := fakeCaller.CallArgsForCall(0)

				Expect(actualCtx).To(Equal(providedCtx))
				Expect(actualCallID).To(Equal(callID))
				Expect(actualScope).To(Equal(expectedScope))
				Expect(actualSCG).To(Equal(providedSCG))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualParentCallID).To(Equal(providedParentCallID))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
			})
		})
	})
})
