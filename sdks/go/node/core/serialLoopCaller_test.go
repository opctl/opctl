package core

import (
	"context"
	"errors"

	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uniquestringFakes "github.com/opctl/opctl/sdks/go/internal/uniquestring/fakes"
	"github.com/opctl/opctl/sdks/go/model"
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
		Context("initial callSerialLoop.Until true", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeCaller := new(FakeCaller)

				objectUnderTest := _serialLoopCaller{
					caller:                fakeCaller,
					pubSub:                new(FakePubSub),
					uniqueStringFactory:   new(uniquestringFakes.FakeUniqueStringFactory),
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					"id",
					map[string]*model.Value{},
					model.SerialLoopCallSpec{
						Until: []*model.PredicateSpec{
							{
								Eq: &[]interface{}{
									true,
									true,
								},
							},
						},
					},
					"dummyOpPath",
					nil,
					"rootCallID",
				)

				/* assert */
				Expect(fakeCaller.CallCallCount()).To(Equal(0))
			})
		})
		Context("initial callSerialLoop.On empty", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeCaller := new(FakeCaller)

				objectUnderTest := _serialLoopCaller{
					caller:                fakeCaller,
					pubSub:                new(FakePubSub),
					uniqueStringFactory:   new(uniquestringFakes.FakeUniqueStringFactory),
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					"id",
					map[string]*model.Value{},
					model.SerialLoopCallSpec{
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
		Context("initial callSerialLoop.Until false", func() {
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

				expectedScope := map[string]*model.Value{
					index: &model.Value{Number: new(float64)},
				}

				fakeCaller := new(FakeCaller)
				fakeCaller.CallReturns(nil, errors.New(""))

				callID := "callID"

				fakePubSub := new(FakePubSub)
				eventChannel := make(chan model.Event, 100)
				fakePubSub.SubscribeStub = func(ctx context.Context, filter model.EventFilter) (<-chan model.Event, <-chan error) {
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

					return eventChannel, make(chan error)
				}

				fakeUniqueStringFactory := new(uniquestringFakes.FakeUniqueStringFactory)
				fakeUniqueStringFactory.ConstructReturns(callID, nil)

				objectUnderTest := _serialLoopCaller{
					caller:                fakeCaller,
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
