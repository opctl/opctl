package core

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
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
					new(FakeCallKiller),
					new(FakePubSub),
				),
			).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		closedEventChan := make(chan model.Event, 1000)
		close(closedEventChan)

		Context("Null SCG", func() {
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
					new(modelFakes.FakeDataHandle),
					nil,
					"dummyRootOpID",
				)
			})
		})

		It("should call callInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedCallID := "dummyCallID"
			providedScope := map[string]*model.Value{}
			providedSCG := &model.SCG{}
			providedOpHandle := new(modelFakes.FakeDataHandle)
			providedParentIDValue := "providedParentID"
			providedParentID := &providedParentIDValue
			providedRootOpID := "dummyRootOpID"

			fakeCallInterpreter := new(callFakes.FakeInterpreter)
			fakeCallInterpreter.InterpretReturns(
				&model.DCG{},
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
				providedSCG,
				providedOpHandle,
				providedParentID,
				providedRootOpID,
			)

			/* assert */
			actualScope,
				actualSCG,
				actualID,
				actualOpHandle,
				actualParentID,
				actualRootOpID := fakeCallInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualSCG).To(Equal(providedSCG))
			Expect(actualID).To(Equal(providedCallID))
			Expect(actualOpHandle).To(Equal(providedOpHandle))
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
					&model.DCG{
						If: &falseBoolean,
					},
					nil,
				)

				expectedEvent := model.Event{
					CallEnded: &model.CallEndedEvent{
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
					&model.SCG{},
					new(modelFakes.FakeDataHandle),
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

		Context("Container SCG", func() {
			It("should call containerCaller.Call w/ expected args", func() {
				/* arrange */
				fakeContainerCaller := new(FakeContainerCaller)

				providedScope := map[string]*model.Value{}
				providedSCG := &model.SCG{
					Container: &model.SCGContainerCall{},
				}

				expectedDCG := &model.DCG{
					Container: &model.DCGContainerCall{},
				}
				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(expectedDCG, nil)

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
					providedSCG,
					new(modelFakes.FakeDataHandle),
					nil,
					"dummyRootOpID",
				)

				/* assert */
				_,
					actualDCGContainerCall,
					actualScope,
					actualSCG := fakeContainerCaller.CallArgsForCall(0)

				Expect(actualDCGContainerCall).To(Equal(expectedDCG.Container))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualSCG).To(Equal(providedSCG.Container))
			})
		})

		Context("Op SCG", func() {
			It("should call opCaller.Call w/ expected args", func() {
				/* arrange */
				fakeOpCaller := new(FakeOpCaller)

				expectedDCG := &model.DCG{
					Op: &model.DCGOpCall{},
				}
				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(
					expectedDCG,
					nil,
				)

				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedSCG := &model.SCG{
					Op: &model.SCGOpCall{
						Ref: "dummyOpRef",
					},
				}
				providedOpHandle := new(modelFakes.FakeDataHandle)
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
					providedSCG,
					providedOpHandle,
					&providedParentID,
					providedRootOpID,
				)

				/* assert */
				_,
					actualDCGOpCall,
					actualScope,
					actualParentID,
					actualSCG := fakeOpCaller.CallArgsForCall(0)

				Expect(actualDCGOpCall).To(Equal(expectedDCG.Op))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualParentID).To(Equal(actualParentID))
				Expect(actualSCG).To(Equal(providedSCG.Op))
			})
		})

		Context("Parallel SCG", func() {
			It("should call parallelCaller.Call w/ expected args", func() {
				/* arrange */
				fakeParallelCaller := new(FakeParallelCaller)

				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedSCG := &model.SCG{
					Parallel: []*model.SCG{
						{Container: &model.SCGContainerCall{}},
					},
				}
				providedOpHandle := new(modelFakes.FakeDataHandle)
				providedRootOpID := "dummyRootOpID"

				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(
					&model.DCG{},
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
					providedSCG,
					providedOpHandle,
					nil,
					providedRootOpID,
				)

				/* assert */
				_,
					actualCallID,
					actualScope,
					actualRootOpID,
					actualOpHandle,
					actualSCG := fakeParallelCaller.CallArgsForCall(0)

				Expect(actualCallID).To(Equal(providedCallID))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualSCG).To(Equal(providedSCG.Parallel))
			})
		})

		Context("ParallelLoop SCG", func() {
			It("should call parallelLoopCaller.Call w/ expected args", func() {
				/* arrange */
				fakeParallelLoopCaller := new(FakeParallelLoopCaller)

				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedSCG := &model.SCG{
					ParallelLoop: &model.SCGParallelLoopCall{},
				}
				providedOpHandle := new(modelFakes.FakeDataHandle)
				providedRootOpID := "dummyRootOpID"
				providedParentID := "providedParentID"

				expectedDCG := &model.DCG{
					ParallelLoop: &model.DCGParallelLoopCall{},
				}
				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(expectedDCG, nil)

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
					providedSCG,
					providedOpHandle,
					&providedParentID,
					providedRootOpID,
				)

				/* assert */
				_,
					actualID,
					actualScope,
					actualSCGParallelLoopCall,
					actualOpHandle,
					actualParentID,
					actualRootOpID := fakeParallelLoopCaller.CallArgsForCall(0)

				Expect(actualID).To(Equal(providedCallID))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualSCGParallelLoopCall).To(Equal(*providedSCG.ParallelLoop))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(*actualParentID).To(Equal(providedParentID))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
			})
		})

		Context("Serial SCG", func() {

			It("should call serialCaller.Call w/ expected args", func() {
				/* arrange */
				fakeSerialCaller := new(FakeSerialCaller)

				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedSCG := &model.SCG{
					Serial: []*model.SCG{
						{Container: &model.SCGContainerCall{}},
					},
				}
				providedOpHandle := new(modelFakes.FakeDataHandle)
				providedRootOpID := "dummyRootOpID"

				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(
					&model.DCG{},
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
					providedSCG,
					providedOpHandle,
					nil,
					providedRootOpID,
				)

				/* assert */
				_,
					actualCallID,
					actualScope,
					actualRootOpID,
					actualOpHandle,
					actualSCG := fakeSerialCaller.CallArgsForCall(0)

				Expect(actualCallID).To(Equal(providedCallID))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualSCG).To(Equal(providedSCG.Serial))
			})
		})

		Context("SerialLoop SCG", func() {
			It("should call serialLoopCaller.Call w/ expected args", func() {
				/* arrange */
				fakeSerialLoopCaller := new(FakeSerialLoopCaller)

				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedSCG := &model.SCG{
					SerialLoop: &model.SCGSerialLoopCall{},
				}
				providedOpHandle := new(modelFakes.FakeDataHandle)
				providedRootOpID := "dummyRootOpID"
				providedParentID := "providedParentID"

				expectedDCG := &model.DCG{
					SerialLoop: &model.DCGSerialLoopCall{},
				}
				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(expectedDCG, nil)

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
					providedSCG,
					providedOpHandle,
					&providedParentID,
					providedRootOpID,
				)

				/* assert */
				_,
					actualID,
					actualScope,
					actualSCGSerialLoopCall,
					actualOpHandle,
					actualParentID,
					actualRootOpID := fakeSerialLoopCaller.CallArgsForCall(0)

				Expect(actualID).To(Equal(providedCallID))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualSCGSerialLoopCall).To(Equal(*providedSCG.SerialLoop))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(*actualParentID).To(Equal(providedParentID))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
			})
		})

		Context("No SCG", func() {
			It("should error", func() {
				/* arrange */
				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedSCG := &model.SCG{}
				providedOpHandle := new(modelFakes.FakeDataHandle)
				providedRootOpID := "dummyRootOpID"
				expectedError := fmt.Errorf("Invalid call graph %+v\n", providedSCG)

				fakeCallInterpreter := new(callFakes.FakeInterpreter)
				fakeCallInterpreter.InterpretReturns(
					&model.DCG{},
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
					providedSCG,
					providedOpHandle,
					nil,
					providedRootOpID,
				)

				/* assert */
				actualEvent := fakePubSub.PublishArgsForCall(0)

				Expect(actualEvent.CallEnded.Error.Message).To(Equal(expectedError.Error()))
			})
		})
	})
})
