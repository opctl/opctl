package core

import (
	"context"

	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	loopFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/fakes"
	iterationFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uniquestringFakes "github.com/opctl/opctl/sdks/go/internal/uniquestring/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	serialloopFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/serialloop/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("serialLoopCaller", func() {
	Context("newSerialLoopCaller", func() {
		It("should return serialLoopCaller", func() {
			/* arrange/act/assert */
			Expect(newSerialLoopCaller(
				new(FakeCaller),
				new(FakePubSub),
			)).To(Not(BeNil()))
		})
	})

	Context("Call", func() {
		Context("initial dcgSerialLoop.Until true", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeSerialLoopInterpreter := new(serialloopFakes.FakeInterpreter)
				until := true
				fakeSerialLoopInterpreter.InterpretReturns(&model.DCGSerialLoopCall{Until: &until}, nil)

				fakeCaller := new(FakeCaller)

				objectUnderTest := _serialLoopCaller{
					caller:                fakeCaller,
					loopDeScoper:          new(loopFakes.FakeDeScoper),
					serialLoopInterpreter: fakeSerialLoopInterpreter,
					iterationScoper:       new(iterationFakes.FakeScoper),
					pubSub:                new(FakePubSub),
					uniqueStringFactory:   new(uniquestringFakes.FakeUniqueStringFactory),
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					"id",
					map[string]*model.Value{},
					model.SCGSerialLoopCall{},
					new(modelFakes.FakeDataHandle),
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
				fakeSerialLoopInterpreter := new(serialloopFakes.FakeInterpreter)
				fakeSerialLoopInterpreter.InterpretReturns(
					&model.DCGSerialLoopCall{
						Range: &model.Value{
							Array: new([]interface{}),
						},
					},
					nil,
				)

				fakeCaller := new(FakeCaller)

				objectUnderTest := _serialLoopCaller{
					caller:                fakeCaller,
					loopDeScoper:          new(loopFakes.FakeDeScoper),
					serialLoopInterpreter: fakeSerialLoopInterpreter,
					iterationScoper:       new(iterationFakes.FakeScoper),
					pubSub:                new(FakePubSub),
					uniqueStringFactory:   new(uniquestringFakes.FakeUniqueStringFactory),
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					"id",
					map[string]*model.Value{},
					model.SCGSerialLoopCall{},
					new(modelFakes.FakeDataHandle),
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
				providedSCGSerialLoopCall := model.SCGSerialLoopCall{
					Run: model.SCG{
						Container: new(model.SCGContainerCall),
					},
					Vars: &model.SCGLoopVars{
						Index: &index,
					},
				}
				providedOpHandle := new(modelFakes.FakeDataHandle)
				providedParentCallIDValue := "providedParentCallID"
				providedParentCallID := &providedParentCallIDValue
				providedRootOpID := "providedRootOpID"

				fakeSerialLoopInterpreter := new(serialloopFakes.FakeInterpreter)
				until := false
				fakeSerialLoopInterpreter.InterpretReturns(
					&model.DCGSerialLoopCall{
						Until: &until,
						Vars: &model.DCGLoopVars{
							Index: &index,
						},
					},
					nil,
				)

				fakeIterationScoper := new(iterationFakes.FakeScoper)
				expectedScope := map[string]*model.Value{
					index: &model.Value{Number: new(float64)},
				}
				fakeIterationScoper.ScopeReturns(expectedScope, nil)

				fakeCaller := new(FakeCaller)

				callID := "callID"

				expectedErrorMessage := "expectedErrorMessage"
				fakePubSub := new(FakePubSub)
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

				fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
				fakeUniqueStringFactory.ConstructReturns(callID, nil)

				objectUnderTest := _serialLoopCaller{
					caller:                fakeCaller,
					loopDeScoper:          new(loopFakes.FakeDeScoper),
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
