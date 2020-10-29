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
				fakeSerialLoopInterpreter.InterpretReturns(&model.SerialLoopCall{Until: &until}, nil)

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
					model.SerialLoopCallSpec{},
					"dummyOpPath",
					nil,
					"rootCallID",
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
					&model.SerialLoopCall{
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
					model.SerialLoopCallSpec{},
					"dummyOpPath",
					nil,
					"rootCallID",
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
				providedSerialLoopCallSpec := model.SerialLoopCallSpec{
					Run: model.CallSpec{
						Container: new(model.ContainerCallSpec),
					},
					Vars: &model.LoopVarsSpec{
						Index: &index,
					},
				}
				providedOpPath := "providedOpPath"
				providedParentCallIDValue := "providedParentCallID"
				providedParentCallID := &providedParentCallIDValue
				providedRootCallID := "providedRootCallID"

				fakeSerialLoopInterpreter := new(serialloopFakes.FakeInterpreter)
				until := false
				fakeSerialLoopInterpreter.InterpretReturns(
					&model.SerialLoopCall{
						Until: &until,
						Vars: &model.LoopVars{
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
						CallEnded: &model.CallEnded{
							CallID: callID,
							Error: &model.CallEndedError{
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
					providedSerialLoopCallSpec,
					providedOpPath,
					providedParentCallID,
					providedRootCallID,
				)

				/* assert */
				actualCtx,
					actualCallID,
					actualScope,
					actualCallSpec,
					actualOpPath,
					actualParentCallID,
					actualRootCallID := fakeCaller.CallArgsForCall(0)

				Expect(actualCtx).To(Not(BeNil()))
				Expect(actualCallID).To(Equal(callID))
				Expect(actualScope).To(Equal(expectedScope))
				Expect(actualCallSpec).To(Equal(&providedSerialLoopCallSpec.Run))
				Expect(actualOpPath).To(Equal(providedOpPath))
				Expect(actualParentCallID).To(Equal(providedParentCallID))
				Expect(actualRootCallID).To(Equal(providedRootCallID))
			})
		})
	})
})
