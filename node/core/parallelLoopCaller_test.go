package core

import (
	"context"

	"github.com/opctl/sdk-golang/opspec/interpreter/call/loop"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/loop/iteration"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/parallelloop"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

var _ = Context("parallelLoopCaller", func() {
	Context("newParallelLoopCaller", func() {
		It("should return parallelLoopCaller", func() {
			/* arrange/act/assert */
			Expect(newParallelLoopCaller(
				new(fakeCaller),
				new(pubsub.Fake),
			)).To(Not(BeNil()))
		})
	})

	Context("Call", func() {
		Context("initial dcgParallelLoop.Range empty", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeParallelLoopInterpreter := new(parallelloop.FakeInterpreter)
				fakeParallelLoopInterpreter.InterpretReturns(
					&model.DCGParallelLoopCall{
						Range: &model.Value{
							Array: new([]interface{}),
						},
					},
					nil,
				)

				fakeCaller := new(fakeCaller)

				objectUnderTest := _parallelLoopCaller{
					caller:                  fakeCaller,
					loopDeScoper:            new(loop.FakeDeScoper),
					parallelLoopInterpreter: fakeParallelLoopInterpreter,
					iterationScoper:         new(iteration.FakeScoper),
					pubSub:                  new(pubsub.Fake),
					uniqueStringFactory:     new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					"id",
					map[string]*model.Value{},
					model.SCGParallelLoopCall{},
					new(data.FakeHandle),
					nil,
					"rootOpID",
				)

				/* assert */
				Expect(fakeCaller.CallCount()).To(Equal(0))
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
			providedOpHandle := new(data.FakeHandle)
			providedParentCallIDValue := "providedParentCallID"
			providedParentCallID := &providedParentCallIDValue
			providedRootOpID := "providedRootOpID"

			loopRangeValue := []interface{}{"value1", "value2"}
			loopRange := &model.Value{
				Array: &loopRangeValue,
			}

			fakeParallelLoopInterpreter := new(parallelloop.FakeInterpreter)
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

			fakeIterationScoper := new(iteration.FakeScoper)
			expectedScope := map[string]*model.Value{
				index: &model.Value{Number: new(float64)},
			}
			fakeIterationScoper.ScopeReturns(expectedScope, nil)

			callID := "callID"
			expectedErrorMessage := "expectedErrorMessage"

			fakeCaller := new(fakeCaller)
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

			fakePubSub := new(pubsub.Fake)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeUniqueStringFactory := new(uniquestring.Fake)
			fakeUniqueStringFactory.ConstructReturns(callID, nil)

			objectUnderTest := _parallelLoopCaller{
				caller:                  fakeCaller,
				loopDeScoper:            new(loop.FakeDeScoper),
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
