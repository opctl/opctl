package core

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	callFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("caller", func() {
	Context("newCaller", func() {
		It("should return caller", func() {
			/* arrange/act/assert */
			Expect(
				newCaller(
					new(callFakes.FakeInterpreter),
					new(FakeContainerCaller),
					"dummyDataDir",
					new(FakeCallStore),
					new(FakePubSub),
				),
			).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		closedEventChan := make(chan model.Event, 1000)
		close(closedEventChan)

		Context("Null CallSpec", func() {
			It("should not throw", func() {
				/* arrange */
				fakeContainerCaller := new(FakeContainerCaller)

				/* act */
				objectUnderTest := _caller{
					callInterpreter: new(callFakes.FakeInterpreter),
					callStore:       new(FakeCallStore),
					containerCaller: fakeContainerCaller,
					pubSub:          new(FakePubSub),
				}

				/* assert */
				objectUnderTest.Call(
					context.Background(),
					"dummyCallID",
					map[string]*model.Value{},
					nil,
					"dummyOpPath",
					nil,
					"dummyRootOpID",
				)
			})
		})

		It("should call callInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedCallID := "dummyCallID"
			providedScope := map[string]*model.Value{}
			providedCallSpec := &model.CallSpec{}
			providedOpPath := "providedOpPath"
			providedParentIDValue := "providedParentID"
			providedParentID := &providedParentIDValue
			providedRootOpID := "dummyRootOpID"

			fakeCallInterpreter := new(callFakes.FakeInterpreter)
			fakeCallInterpreter.InterpretReturns(
				&model.Call{},
				nil,
			)

			fakePubSub := new(FakePubSub)
			// ensure eventChan closed so call exits
			fakePubSub.SubscribeReturns(closedEventChan, nil)

			objectUnderTest := _caller{
				callInterpreter: fakeCallInterpreter,
				callStore:       new(FakeCallStore),
				containerCaller: new(FakeContainerCaller),
				pubSub:          fakePubSub,
			}

			/* act */
			objectUnderTest.Call(
				context.Background(),
				providedCallID,
				providedScope,
				providedCallSpec,
				providedOpPath,
				providedParentID,
				providedRootOpID,
			)

			/* assert */
			actualScope,
				actualCallSpec,
				actualID,
				actualOpPath,
				actualParentID,
				actualRootOpID := fakeCallInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualCallSpec).To(Equal(providedCallSpec))
			Expect(actualID).To(Equal(providedCallID))
			Expect(actualOpPath).To(Equal(providedOpPath))
			Expect(actualParentID).To(Equal(providedParentID))
			Expect(actualRootOpID).To(Equal(providedRootOpID))
		})
		Context("callInterpreter.Interpret result.If falsy", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedCallID := "dummyCallID"
				providedRootOpID := "dummyRootOpID"

				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				falseBoolean := false
				fakeCallInterpreter.InterpretReturns(
					&model.Call{
						If: &falseBoolean,
					},
					nil,
				)

				expectedEvent := model.Event{
					CallEnded: &model.CallEnded{
						CallID:     providedCallID,
						RootCallID: providedRootOpID,
					},
				}

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
					callInterpreter: fakeCallInterpreter,
					callStore:       new(FakeCallStore),
					containerCaller: new(FakeContainerCaller),
					pubSub:          fakePubSub,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedCallID,
					map[string]*model.Value{},
					&model.CallSpec{},
					"dummyOpPath",
					nil,
					providedRootOpID,
				)

				/* assert */
				actualEvent := fakePubSub.PublishArgsForCall(0)

				// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
				Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
				// set temporal fields to expected vals since they're already asserted
				actualEvent.Timestamp = expectedEvent.Timestamp

				Expect(actualEvent).To(Equal(expectedEvent))
			})
		})

		Context("Container CallSpec", func() {
			It("should call containerCaller.Call w/ expected args", func() {
				/* arrange */
				fakeContainerCaller := new(FakeContainerCaller)

				providedScope := map[string]*model.Value{}
				providedCallSpec := &model.CallSpec{
					Container: &model.ContainerCallSpec{},
				}

				expectedCall := &model.Call{
					Container: &model.ContainerCall{},
				}
				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(expectedCall, nil)

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
					callInterpreter: fakeCallInterpreter,
					callStore:       new(FakeCallStore),
					containerCaller: fakeContainerCaller,
					pubSub:          fakePubSub,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					"dummyCallID",
					providedScope,
					providedCallSpec,
					"dummyOpPath",
					nil,
					"dummyRootOpID",
				)

				/* assert */
				_,
					actualContainerCall,
					actualScope,
					actualCallSpec := fakeContainerCaller.CallArgsForCall(0)

				Expect(actualContainerCall).To(Equal(expectedCall.Container))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualCallSpec).To(Equal(providedCallSpec.Container))
			})
		})

		Context("Op CallSpec", func() {
			It("should call opCaller.Call w/ expected args", func() {
				/* arrange */
				fakeOpCaller := new(FakeOpCaller)

				expectedCall := &model.Call{
					Op: &model.OpCall{},
				}
				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(
					expectedCall,
					nil,
				)

				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedCallSpec := &model.CallSpec{
					Op: &model.OpCallSpec{
						Ref: "dummyOpRef",
					},
				}
				providedOpPath := "providedOpPath"
				providedParentID := "providedParentID"
				providedRootOpID := "dummyRootOpID"

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
					callInterpreter: fakeCallInterpreter,
					callStore:       new(FakeCallStore),
					opCaller:        fakeOpCaller,
					pubSub:          fakePubSub,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedCallID,
					providedScope,
					providedCallSpec,
					providedOpPath,
					&providedParentID,
					providedRootOpID,
				)

				/* assert */
				_,
					actualOpCall,
					actualScope,
					actualParentID,
					actualCallSpec := fakeOpCaller.CallArgsForCall(0)

				Expect(actualOpCall).To(Equal(expectedCall.Op))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualParentID).To(Equal(actualParentID))
				Expect(actualCallSpec).To(Equal(providedCallSpec.Op))
			})
		})

		Context("Parallel CallSpec", func() {
			It("should call parallelCaller.Call w/ expected args", func() {
				/* arrange */
				fakeParallelCaller := new(FakeParallelCaller)

				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedCallSpec := &model.CallSpec{
					Parallel: &[]*model.CallSpec{
						{Container: &model.ContainerCallSpec{}},
					},
				}
				providedOpPath := "providedOpPath"
				providedRootOpID := "dummyRootOpID"

				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(
					&model.Call{},
					nil,
				)

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
					callInterpreter: fakeCallInterpreter,
					callStore:       new(FakeCallStore),
					parallelCaller:  fakeParallelCaller,
					pubSub:          fakePubSub,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedCallID,
					providedScope,
					providedCallSpec,
					providedOpPath,
					nil,
					providedRootOpID,
				)

				/* assert */
				_,
					actualCallID,
					actualScope,
					actualRootOpID,
					actualOpPath,
					actualCallSpec := fakeParallelCaller.CallArgsForCall(0)

				Expect(actualCallID).To(Equal(providedCallID))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
				Expect(actualOpPath).To(Equal(providedOpPath))
				Expect(actualCallSpec).To(Equal(*providedCallSpec.Parallel))
			})
		})

		Context("ParallelLoop CallSpec", func() {
			It("should call parallelLoopCaller.Call w/ expected args", func() {
				/* arrange */
				fakeParallelLoopCaller := new(FakeParallelLoopCaller)

				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedCallSpec := &model.CallSpec{
					ParallelLoop: &model.ParallelLoopCallSpec{},
				}
				providedOpPath := "providedOpPath"
				providedRootOpID := "dummyRootOpID"
				providedParentID := "providedParentID"

				expectedCall := &model.Call{
					ParallelLoop: &model.ParallelLoopCall{},
				}
				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(expectedCall, nil)

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
					callInterpreter:    fakeCallInterpreter,
					callStore:          new(FakeCallStore),
					parallelLoopCaller: fakeParallelLoopCaller,
					pubSub:             fakePubSub,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedCallID,
					providedScope,
					providedCallSpec,
					providedOpPath,
					&providedParentID,
					providedRootOpID,
				)

				/* assert */
				_,
					actualID,
					actualScope,
					actualParallelLoopCallSpec,
					actualOpPath,
					actualParentID,
					actualRootOpID := fakeParallelLoopCaller.CallArgsForCall(0)

				Expect(actualID).To(Equal(providedCallID))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualParallelLoopCallSpec).To(Equal(*providedCallSpec.ParallelLoop))
				Expect(actualOpPath).To(Equal(providedOpPath))
				Expect(*actualParentID).To(Equal(providedParentID))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
			})
		})

		Context("Serial CallSpec", func() {

			It("should call serialCaller.Call w/ expected args", func() {
				/* arrange */
				fakeSerialCaller := new(FakeSerialCaller)

				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedCallSpec := &model.CallSpec{
					Serial: &[]*model.CallSpec{
						{Container: &model.ContainerCallSpec{}},
					},
				}
				providedOpPath := "providedOpPath"
				providedRootOpID := "dummyRootOpID"

				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(
					&model.Call{},
					nil,
				)

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
					callInterpreter: fakeCallInterpreter,
					callStore:       new(FakeCallStore),
					containerCaller: new(FakeContainerCaller),
					pubSub:          fakePubSub,
					serialCaller:    fakeSerialCaller,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedCallID,
					providedScope,
					providedCallSpec,
					providedOpPath,
					nil,
					providedRootOpID,
				)

				/* assert */
				_,
					actualCallID,
					actualScope,
					actualRootOpID,
					actualOpPath,
					actualCallSpec := fakeSerialCaller.CallArgsForCall(0)

				Expect(actualCallID).To(Equal(providedCallID))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
				Expect(actualOpPath).To(Equal(providedOpPath))
				Expect(actualCallSpec).To(Equal(*providedCallSpec.Serial))
			})
		})

		Context("SerialLoop CallSpec", func() {
			It("should call serialLoopCaller.Call w/ expected args", func() {
				/* arrange */
				fakeSerialLoopCaller := new(FakeSerialLoopCaller)

				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedCallSpec := &model.CallSpec{
					SerialLoop: &model.SerialLoopCallSpec{},
				}
				providedOpPath := "providedOpPath"
				providedRootOpID := "dummyRootOpID"
				providedParentID := "providedParentID"

				expectedCall := &model.Call{
					SerialLoop: &model.SerialLoopCall{},
				}
				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(expectedCall, nil)

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
					callInterpreter:  fakeCallInterpreter,
					callStore:        new(FakeCallStore),
					serialLoopCaller: fakeSerialLoopCaller,
					pubSub:           fakePubSub,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedCallID,
					providedScope,
					providedCallSpec,
					providedOpPath,
					&providedParentID,
					providedRootOpID,
				)

				/* assert */
				_,
					actualID,
					actualScope,
					actualSerialLoopCallSpec,
					actualOpPath,
					actualParentID,
					actualRootOpID := fakeSerialLoopCaller.CallArgsForCall(0)

				Expect(actualID).To(Equal(providedCallID))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualSerialLoopCallSpec).To(Equal(*providedCallSpec.SerialLoop))
				Expect(actualOpPath).To(Equal(providedOpPath))
				Expect(*actualParentID).To(Equal(providedParentID))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
			})
		})

		Context("No CallSpec", func() {
			It("should error", func() {
				/* arrange */
				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedCallSpec := &model.CallSpec{}
				expectedError := fmt.Errorf("Invalid call graph %+v\n", providedCallSpec)

				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(
					&model.Call{},
					nil,
				)

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
					callInterpreter: fakeCallInterpreter,
					callStore:       new(FakeCallStore),
					pubSub:          fakePubSub,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedCallID,
					providedScope,
					providedCallSpec,
					"providedOpPath",
					nil,
					"dummyRootOpID",
				)

				/* assert */
				actualEvent := fakePubSub.PublishArgsForCall(0)

				Expect(actualEvent.CallEnded.Error.Message).To(Equal(expectedError.Error()))
			})
		})
	})
})
