package core

import (
	"context"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/serialloop"
	"github.com/opctl/opctl/sdks/go/types"
	"github.com/opctl/opctl/sdks/go/util/pubsub"
	"github.com/opctl/opctl/sdks/go/util/uniquestring"
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
				fakeSerialLoopInterpreter.InterpretReturns(&types.DCGSerialLoopCall{Until: &until}, nil)

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
					map[string]*types.Value{},
					types.SCGSerialLoopCall{},
					new(data.FakeHandle),
					nil,
					"rootOpID",
				)

				/* assert */
				Expect(fakeCaller.CallCount()).To(Equal(0))
			})
		})
		Context("initial dcgSerialLoop.On empty", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeSerialLoopInterpreter := new(serialloop.FakeInterpreter)
				fakeSerialLoopInterpreter.InterpretReturns(
					&types.DCGSerialLoopCall{
						Range: &types.Value{
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
					map[string]*types.Value{},
					types.SCGSerialLoopCall{},
					new(data.FakeHandle),
					nil,
					"rootOpID",
				)

				/* assert */
				Expect(fakeCaller.CallCount()).To(Equal(0))
			})
		})
		Context("initial dcgSerialLoop.Until false", func() {
			It("should call caller.Call w/ expected args", func() {
				/* arrange */
				providedCtx := context.Background()
				providedScope := map[string]*types.Value{}
				index := "index"
				providedSCGSerialLoopCall := types.SCGSerialLoopCall{
					Run: types.SCG{
						Container: new(types.SCGContainerCall),
					},
					Vars: &types.SCGLoopVars{
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
					&types.DCGSerialLoopCall{
						Until: &until,
						Vars: &types.DCGLoopVars{
							Index: &index,
						},
					},
					nil,
				)

				fakeIterationScoper := new(iteration.FakeScoper)
				expectedScope := map[string]*types.Value{
					index: &types.Value{Number: new(float64)},
				}
				fakeIterationScoper.ScopeReturns(expectedScope, nil)

				fakeCaller := new(fakeCaller)

				callID := "callID"

				expectedErrorMessage := "expectedErrorMessage"
				fakePubSub := new(pubsub.Fake)
				eventChannel := make(chan types.Event, 100)
				fakePubSub.SubscribeStub = func(ctx context.Context, filter types.EventFilter) (<-chan types.Event, <-chan error) {
					eventChannel <- types.Event{
						CallEnded: &types.CallEndedEvent{
							CallID: callID,
							Error: &types.CallEndedEventError{
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
					providedSCGSerialLoopCall,
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
				Expect(actualSCG).To(Equal(&providedSCGSerialLoopCall.Run))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualParentCallID).To(Equal(providedParentCallID))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
			})
		})
	})
})
