package core

import (
	"context"

	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	loopFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/fakes"
	iterationFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/loop/iteration/fakes"
	parallelloopFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/parallelloop/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uniquestringFakes "github.com/opctl/opctl/sdks/go/internal/uniquestring/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("parallelLoopCaller", func() {
	Context("newParallelLoopCaller", func() {
		It("should return parallelLoopCaller", func() {
			/* arrange/act/assert */
			Expect(newParallelLoopCaller(
				new(FakeCaller),
				new(FakePubSub),
			)).To(Not(BeNil()))
		})
	})

	Context("Call", func() {
		Context("initial dcgParallelLoop.Range empty", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeParallelLoopInterpreter := new(parallelloopFakes.FakeInterpreter)
				fakeParallelLoopInterpreter.InterpretReturns(
					&model.DCGParallelLoopCall{
						Range: &model.Value{
							Array: new([]interface{}),
						},
					},
					nil,
				)

				fakeCaller := new(FakeCaller)

				objectUnderTest := _parallelLoopCaller{
					caller:                  fakeCaller,
					loopDeScoper:            new(loopFakes.FakeDeScoper),
					parallelLoopInterpreter: fakeParallelLoopInterpreter,
					iterationScoper:         new(iterationFakes.FakeScoper),
					pubSub:                  new(FakePubSub),
					uniqueStringFactory:     new(uniquestringFakes.FakeUniqueStringFactory),
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					"id",
					map[string]*model.Value{},
					model.SCGParallelLoopCall{},
					new(modelFakes.FakeDataHandle),
					nil,
					"rootOpID",
				)

				/* assert */
				Expect(fakeCaller.CallCallCount()).To(Equal(0))
			})
		})
		It("should call caller.Call w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedScope := map[string]*model.Value{}
			index := "index"
			providedSCGParallelLoopCall := model.SCGParallelLoopCall{
				Vars: &model.SCGLoopVars{
					Index: &index,
				},
				Run: model.SCG{
					Container: new(model.SCGContainerCall),
				},
			}
			providedOpHandle := new(modelFakes.FakeDataHandle)
			providedParentCallIDValue := "providedParentCallID"
			providedParentCallID := &providedParentCallIDValue
			providedRootOpID := "providedRootOpID"

			loopRangeValue := []interface{}{"value1", "value2"}
			loopRange := &model.Value{
				Array: &loopRangeValue,
			}

			fakeParallelLoopInterpreter := new(parallelloopFakes.FakeInterpreter)
			fakeParallelLoopInterpreter.InterpretReturnsOnCall(
				0,
				&model.DCGParallelLoopCall{
					Range: loopRange,
					Vars: &model.DCGLoopVars{
						Index: &index,
					},
				},
				nil,
			)

			fakeParallelLoopInterpreter.InterpretReturnsOnCall(
				1,
				&model.DCGParallelLoopCall{
					Range: loopRange,
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

			callID := "callID"
			expectedErrorMessage := "expectedErrorMessage"

			fakeCaller := new(FakeCaller)
			eventChannel := make(chan model.Event, 100)
			fakeCaller.CallStub = func(context.Context, string, map[string]*model.Value, *model.SCG, model.DataHandle, *string, string) {
				eventChannel <- model.Event{
					CallEnded: &model.CallEndedEvent{
						CallID: callID,
						Error: &model.CallEndedEventError{
							Message: expectedErrorMessage,
						},
					},
				}
			}

			fakePubSub := new(FakePubSub)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
			fakeUniqueStringFactory.ConstructReturns(callID, nil)

			objectUnderTest := _parallelLoopCaller{
				caller:                  fakeCaller,
				loopDeScoper:            new(loopFakes.FakeDeScoper),
				parallelLoopInterpreter: fakeParallelLoopInterpreter,
				iterationScoper:         fakeIterationScoper,
				pubSub:                  fakePubSub,
				uniqueStringFactory:     fakeUniqueStringFactory,
			}

			/* act */
			objectUnderTest.Call(
				providedCtx,
				"id",
				providedScope,
				providedSCGParallelLoopCall,
				providedOpHandle,
				providedParentCallID,
				providedRootOpID,
			)

			/* assert */
			_,
				actualCallID,
				actualScope,
				actualSCG,
				actualOpHandle,
				actualParentCallID,
				actualRootOpID := fakeCaller.CallArgsForCall(0)

			Expect(actualCallID).To(Equal(callID))
			Expect(actualScope).To(Equal(expectedScope))
			Expect(actualSCG).To(Equal(&providedSCGParallelLoopCall.Run))
			Expect(actualOpHandle).To(Equal(providedOpHandle))
			Expect(actualParentCallID).To(Equal(providedParentCallID))
			Expect(actualRootOpID).To(Equal(providedRootOpID))
		})
	})
})
