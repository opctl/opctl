package core

import (
	"context"

	"github.com/opctl/sdk-golang/opspec/interpreter/call/loop"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/loop/iteration"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/serialloop"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

var _ = Context("serialLoopCaller", func() {
	Context("newSerialLoopCaller", func() {
		It("should return serialLoopCaller", func() {
			/* arrange/act/assert */
			Expect(newSerialLoopCaller(
				new(fakeCaller),
				new(pubsub.Fake),
			)).To(Not(BeNil()))
		})
	})

	Context("Call", func() {
		Context("initial dcgSerialLoop.Until true", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeSerialLoopInterpreter := new(serialloop.FakeInterpreter)
				until := true
				fakeSerialLoopInterpreter.InterpretReturns(&model.DCGSerialLoop{Until: &until}, nil)

				fakeCaller := new(fakeCaller)

				objectUnderTest := _serialLoopCaller{
					caller:                fakeCaller,
					loopDeScoper:          new(loop.FakeDeScoper),
					serialLoopInterpreter: fakeSerialLoopInterpreter,
					iterationScoper:       new(iteration.FakeScoper),
					pubSub:                new(pubsub.Fake),
					uniqueStringFactory:   new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					"id",
					map[string]*model.Value{},
					model.SCGSerialLoop{},
					new(data.FakeHandle),
					nil,
					"rootOpID",
				)

				/* assert */
				Expect(fakeCaller.CallCallCount()).To(Equal(0))
			})
		})
		Context("initial dcgSerialLoop.On empty", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeSerialLoopInterpreter := new(serialloop.FakeInterpreter)
				fakeSerialLoopInterpreter.InterpretReturns(
					&model.DCGSerialLoop{
						Range: &model.Value{
							Array: new([]interface{}),
						},
					},
					nil,
				)

				fakeCaller := new(fakeCaller)

				objectUnderTest := _serialLoopCaller{
					caller:                fakeCaller,
					loopDeScoper:          new(loop.FakeDeScoper),
					serialLoopInterpreter: fakeSerialLoopInterpreter,
					iterationScoper:       new(iteration.FakeScoper),
					pubSub:                new(pubsub.Fake),
					uniqueStringFactory:   new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					"id",
					map[string]*model.Value{},
					model.SCGSerialLoop{},
					new(data.FakeHandle),
					nil,
					"rootOpID",
				)

				/* assert */
				Expect(fakeCaller.CallCallCount()).To(Equal(0))
			})
		})
		Context("initial dcgSerialLoop.Until false", func() {
			It("should call caller.Call w/ expected args", func() {
				/* arrange */
				providedCtx := context.Background()
				providedScope := map[string]*model.Value{}
				index := "index"
				providedSCGSerialLoop := model.SCGSerialLoop{
					Run: model.SCG{
						Container: new(model.SCGContainerCall),
					},
					Vars: &model.SCGLoopVars{
						Index: &index,
					},
				}
				providedOpHandle := new(data.FakeHandle)
				providedParentCallIDValue := "providedParentCallID"
				providedParentCallID := &providedParentCallIDValue
				providedRootOpID := "providedRootOpID"

				fakeSerialLoopInterpreter := new(serialloop.FakeInterpreter)
				until := false
				fakeSerialLoopInterpreter.InterpretReturns(
					&model.DCGSerialLoop{
						Until: &until,
						Vars: &model.DCGLoopVars{
							Index: &index,
						},
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

				objectUnderTest := _serialLoopCaller{
					caller:                fakeCaller,
					loopDeScoper:          new(loop.FakeDeScoper),
					serialLoopInterpreter: fakeSerialLoopInterpreter,
					iterationScoper:       fakeIterationScoper,
					pubSub:                fakePubSub,
					uniqueStringFactory:   fakeUniqueStringFactory,
				}

				/* act */
				objectUnderTest.Call(
					providedCtx,
					"id",
					providedScope,
					providedSCGSerialLoop,
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
				Expect(actualSCG).To(Equal(&providedSCGSerialLoop.Run))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualParentCallID).To(Equal(providedParentCallID))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
			})
		})
	})
})
