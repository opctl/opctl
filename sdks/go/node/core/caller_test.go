package core

import (
	"path/filepath"
	"os"
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("caller", func() {
	Context("newCaller", func() {
		It("should return caller", func() {
			/* arrange/act/assert */
			Expect(
				newCaller(
					new(FakeContainerCaller),
					"dummyDataDir",
					new(FakeStateStore),
					new(FakePubSub),
				),
			).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		closedEventChan := make(chan model.Event, 1000)
		close(closedEventChan)

		Context("Nil CallSpec", func() {
			It("should not throw", func() {
				/* arrange */
				fakeContainerCaller := new(FakeContainerCaller)

				/* act */
				objectUnderTest := _caller{
					containerCaller: fakeContainerCaller,
					dataDirPath: os.TempDir(),
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
					"dummyRootCallID",
				)
			})
		})

		Context("callInterpreter.Interpret result.If falsy", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedCallID := "dummyCallID"
				providedOpPath := "testdata/startOp"
				providedRootCallID := "dummyRootCallID"

				predicateSpec := []interface{}{
					true,
					true,
				}

				ifSpec := []*model.PredicateSpec{
					&model.PredicateSpec{
						Eq: &predicateSpec,
					},
				}

				expectedIf := true
				expectedEvent := model.Event{
					CallStarted: &model.CallStarted{
						Call: model.Call{
							ID:     providedCallID,
							If:     &expectedIf,
							RootID: providedRootCallID,
							Serial: []*model.CallSpec{},
						},
						Ref: providedOpPath,
					},
				}

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				fakeSerialCaller := new(FakeSerialCaller)

				objectUnderTest := _caller{
					containerCaller: new(FakeContainerCaller),
					dataDirPath: os.TempDir(),
					pubSub:          fakePubSub,
					serialCaller:    fakeSerialCaller,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedCallID,
					map[string]*model.Value{},
					&model.CallSpec{
						If:     &ifSpec,
						Serial: &[]*model.CallSpec{},
					},
					providedOpPath,
					nil,
					providedRootCallID,
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
				providedCallID := "providedCallID"
				providedOpPath := "providedOpPath"
				fakeContainerCaller := new(FakeContainerCaller)

				providedScope := map[string]*model.Value{}
				imageSpec := &model.ContainerCallImageSpec{
					Ref: "docker.io/library/ref",
				}
				providedCallSpec := &model.CallSpec{
					Container: &model.ContainerCallSpec{
						Image: imageSpec,
					},
				}
				providedRootCallID := "providedRootCallID"

				expectedCall := &model.Call{
					Container: &model.ContainerCall{
						BaseCall: model.BaseCall{
							OpPath: providedOpPath,
						},
						ContainerID: providedCallID,
						Cmd:         []string{},
						Dirs:        map[string]string{},
						Files:       map[string]string{},
						Image: &model.ContainerCallImage{
							Ref: &imageSpec.Ref,
						},
						Sockets: map[string]string{},
					},
				}

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
					containerCaller: fakeContainerCaller,
					dataDirPath: os.TempDir(),
					pubSub:          fakePubSub,
				}

				/* act */
				_, actualErr := objectUnderTest.Call(
					context.Background(),
					providedCallID,
					providedScope,
					providedCallSpec,
					providedOpPath,
					nil,
					providedRootCallID,
				)

				/* assert */
				Expect(actualErr).To(BeNil())
				_,
					actualContainerCall,
					actualScope,
					actualCallSpec,
					actualRootCallID := fakeContainerCaller.CallArgsForCall(0)

				Expect(*actualContainerCall).To(Equal(*expectedCall.Container))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualCallSpec).To(Equal(providedCallSpec.Container))
				Expect(actualRootCallID).To(Equal(providedRootCallID))
			})
		})

		Context("Op CallSpec", func() {
			It("should call opCaller.Call w/ expected args", func() {
				/* arrange */
				fakeOpCaller := new(FakeOpCaller)
				
				wd, err := os.Getwd()
				if nil != err {
					panic(err)
				}
				providedOpPath := filepath.Join(wd, "testdata/caller")

				providedCallID := "dummyCallID"
				providedScope := map[string]*model.Value{}
				providedCallSpec := &model.CallSpec{
					Op: &model.OpCallSpec{
						Ref: providedOpPath,
					},
				}
				providedParentID := "providedParentID"
				providedRootCallID := "dummyRootCallID"

				expectedCall := &model.Call{
					Op: &model.OpCall{
						BaseCall: model.BaseCall{
							OpPath: providedOpPath,
						},
						OpID:   providedCallID,
						Inputs: map[string]*model.Value{},
					},
				}

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
					dataDirPath: os.TempDir(),
					opCaller: fakeOpCaller,
					pubSub:   fakePubSub,
				}

				/* act */
				_, actualErr := objectUnderTest.Call(
					context.Background(),
					providedCallID,
					providedScope,
					providedCallSpec,
					providedOpPath,
					&providedParentID,
					providedRootCallID,
				)

				/* assert */
				Expect(actualErr).To(BeNil())
				_,
					actualOpCall,
					actualScope,
					actualParentID,
					actualRootCallID,
					actualCallSpec := fakeOpCaller.CallArgsForCall(0)

				// ignore ChildCallID since autogenerated unique
				actualOpCall.ChildCallID = expectedCall.Op.ChildCallID

				Expect(*actualOpCall).To(Equal(*expectedCall.Op))
				Expect(actualScope).To(Equal(providedScope))
				Expect(*actualParentID).To(Equal(providedParentID))
				Expect(actualRootCallID).To(Equal(providedRootCallID))
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
				providedRootCallID := "dummyRootCallID"

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
					parallelCaller: fakeParallelCaller,
					pubSub:         fakePubSub,
				}

				/* act */
				objectUnderTest.Call(
					context.Background(),
					providedCallID,
					providedScope,
					providedCallSpec,
					providedOpPath,
					nil,
					providedRootCallID,
				)

				/* assert */
				_,
					actualCallID,
					actualScope,
					actualRootCallID,
					actualOpPath,
					actualCallSpec := fakeParallelCaller.CallArgsForCall(0)

				Expect(actualCallID).To(Equal(providedCallID))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualRootCallID).To(Equal(providedRootCallID))
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
				providedRootCallID := "dummyRootCallID"
				providedParentID := "providedParentID"

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
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
					providedRootCallID,
				)

				/* assert */
				_,
					actualID,
					actualScope,
					actualParallelLoopCallSpec,
					actualOpPath,
					actualParentID,
					actualRootCallID := fakeParallelLoopCaller.CallArgsForCall(0)

				Expect(actualID).To(Equal(providedCallID))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualParallelLoopCallSpec).To(Equal(*providedCallSpec.ParallelLoop))
				Expect(actualOpPath).To(Equal(providedOpPath))
				Expect(*actualParentID).To(Equal(providedParentID))
				Expect(actualRootCallID).To(Equal(providedRootCallID))
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
				providedRootCallID := "dummyRootCallID"

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
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
					providedRootCallID,
				)

				/* assert */
				_,
					actualCallID,
					actualScope,
					actualRootCallID,
					actualOpPath,
					actualCallSpec := fakeSerialCaller.CallArgsForCall(0)

				Expect(actualCallID).To(Equal(providedCallID))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualRootCallID).To(Equal(providedRootCallID))
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
					SerialLoop: &model.SerialLoopCallSpec{
						Range: []interface{}{},
					},
				}
				providedOpPath := "providedOpPath"
				providedRootCallID := "dummyRootCallID"
				providedParentID := "providedParentID"

				fakePubSub := new(FakePubSub)
				// ensure eventChan closed so call exits
				fakePubSub.SubscribeReturns(closedEventChan, nil)

				objectUnderTest := _caller{
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
					providedRootCallID,
				)

				/* assert */
				_,
					actualID,
					actualScope,
					actualSerialLoopCallSpec,
					actualOpPath,
					actualParentID,
					actualRootCallID := fakeSerialLoopCaller.CallArgsForCall(0)

				Expect(actualID).To(Equal(providedCallID))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualSerialLoopCallSpec).To(Equal(*providedCallSpec.SerialLoop))
				Expect(actualOpPath).To(Equal(providedOpPath))
				Expect(*actualParentID).To(Equal(providedParentID))
				Expect(actualRootCallID).To(Equal(providedRootCallID))
			})
		})
	})
})
