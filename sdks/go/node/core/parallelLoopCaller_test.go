package core

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uniquestringFakes "github.com/opctl/opctl/sdks/go/internal/uniquestring/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
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
		Context("initial callParallelLoop.Range empty", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeCaller := new(FakeCaller)

				objectUnderTest := _parallelLoopCaller{
					caller:              fakeCaller,
					pubSub:              new(FakePubSub),
					uniqueStringFactory: new(uniquestringFakes.FakeUniqueStringFactory),
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					"id",
					map[string]*model.Value{},
					model.ParallelLoopCallSpec{
						Range: []interface{}{},
					},
					"dummyOpPath",
					nil,
					"rootCallID",
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
			providedParallelLoopCallSpec := model.ParallelLoopCallSpec{
				Range: []interface{}{
					"one",
				},
				Vars: &model.LoopVarsSpec{
					Index: &index,
				},
				Run: model.CallSpec{
					Container: &model.ContainerCallSpec{
						Image: &model.ContainerCallImageSpec{
							Ref: "ref",
						},
					},
				},
			}
			providedOpPath := "providedOpPath"
			providedParentCallIDValue := "providedParentCallID"
			providedParentCallID := &providedParentCallIDValue
			providedRootCallID := "providedRootCallID"

			expectedScope := map[string]*model.Value{
				index: &model.Value{Number: new(float64)},
			}

			callID := "callID"

			fakeCaller := new(FakeCaller)

			eventChannel := make(chan model.Event, 100)
			fakeCaller.CallStub = func(
				context.Context,
				string,
				map[string]*model.Value,
				*model.CallSpec,
				string,
				*string,
				string,
			) (
				map[string]*model.Value,
				error,
			) {
				eventChannel <- model.Event{
					CallEnded: &model.CallEnded{
						Call: model.Call{
							ID: callID,
						},
						Error: &model.CallEndedError{
							Message: "message",
						},
					},
				}

				return nil, nil
			}

			fakePubSub := new(FakePubSub)
			fakePubSub.SubscribeReturns(eventChannel, nil)

			fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
			fakeUniqueStringFactory.ConstructReturns(callID, nil)

			objectUnderTest := _parallelLoopCaller{
				caller:              fakeCaller,
				pubSub:              fakePubSub,
				uniqueStringFactory: fakeUniqueStringFactory,
			}

			/* act */
			objectUnderTest.Call(
				providedCtx,
				"id",
				providedScope,
				providedParallelLoopCallSpec,
				providedOpPath,
				providedParentCallID,
				providedRootCallID,
			)

			/* assert */
			_,
				actualCallID,
				actualScope,
				actualCallSpec,
				actualOpPath,
				actualParentCallID,
				actualRootCallID := fakeCaller.CallArgsForCall(0)

			Expect(actualCallID).To(Equal(callID))
			Expect(actualScope).To(Equal(expectedScope))
			Expect(actualCallSpec).To(Equal(&providedParallelLoopCallSpec.Run))
			Expect(actualOpPath).To(Equal(providedOpPath))
			Expect(actualParentCallID).To(Equal(providedParentCallID))
			Expect(actualRootCallID).To(Equal(providedRootCallID))
		})
	})
})
