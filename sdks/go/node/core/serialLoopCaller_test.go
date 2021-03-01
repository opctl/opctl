package core

import (
	"context"
	"errors"
	"io/ioutil"

	"github.com/dgraph-io/badger/v2"
	containerRuntimeFakes "github.com/opctl/opctl/sdks/go/node/core/containerruntime/fakes"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	"github.com/opctl/opctl/sdks/go/pubsub"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
					caller: fakeCaller,
					pubSub: new(FakePubSub),
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
					caller: fakeCaller,
					pubSub: new(FakePubSub),
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

			Context("iteration spec invalid", func() {

				It("should return expected results", func() {
					/* arrange */
					dbDir, err := ioutil.TempDir("", "")
					if nil != err {
						panic(err)
					}

					db, err := badger.Open(
						badger.DefaultOptions(dbDir).WithLogger(nil),
					)
					if nil != err {
						panic(err)
					}
					pubSub := pubsub.New(db)

					providedCtx := context.Background()
					providedScope := map[string]*model.Value{}

					caller := newCaller(
						newContainerCaller(
							new(containerRuntimeFakes.FakeContainerRuntime),
							pubSub,
							newStateStore(
								context.Background(),
								db,
								pubSub,
							),
						),
						dbDir,
						pubSub,
					)

					objectUnderTest := _serialLoopCaller{
						caller: caller,
						pubSub: pubSub,
					}

					/* act */
					actualOutputs, actualErr := objectUnderTest.Call(
						providedCtx,
						"id",
						providedScope,
						model.SerialLoopCallSpec{
							Run: model.CallSpec{
								Container: new(model.ContainerCallSpec),
							},
							Vars: &model.LoopVarsSpec{
								Index: new(string),
							},
						},
						"opPath",
						new(string),
						"rootCallID",
					)

					/* assert */
					Expect(actualErr).To(Equal(errors.New("image required")))
					Expect(actualOutputs).To(BeNil())
				})
			})

			It("should start each child as expected", func() {
				/* arrange */
				dbDir, err := ioutil.TempDir("", "")
				if nil != err {
					panic(err)
				}

				db, err := badger.Open(
					badger.DefaultOptions(dbDir).WithLogger(nil),
				)
				if nil != err {
					panic(err)
				}
				pubSub := pubsub.New(db)

				providedOpRef := "providedOpRef"
				providedParentID := "providedParentID"
				providedRootID := "providedRootID"
				imageRef := "docker.io/library/alpine"

				ctx := context.Background()

				fakeContainerRuntime := new(containerRuntimeFakes.FakeContainerRuntime)
				fakeContainerRuntime.RunContainerStub = func(
					ctx context.Context,
					req *model.ContainerCall,
					rootCallID string,
					eventPublisher pubsub.EventPublisher,
					stdOut io.WriteCloser,
					stdErr io.WriteCloser,
				) (*int64, error) {

					stdErr.Close()
					stdOut.Close()

					return nil, nil
				}

				eventChannel, err := pubSub.Subscribe(
					ctx,
					model.EventFilter{},
				)
				if nil != err {
					panic(err)
				}

				objectUnderTest := _serialLoopCaller{
					caller: newCaller(
						newContainerCaller(
							fakeContainerRuntime,
							pubSub,
							newStateStore(
								ctx,
								db,
								pubSub,
							),
						),
						dbDir,
						pubSub,
					),
					pubSub: pubSub,
				}

				/* act */
				_, actualErr := objectUnderTest.Call(
					ctx,
					"",
					map[string]*model.Value{},
					model.SerialLoopCallSpec{
						Range: model.Value{
							Array: &[]interface{}{0, 1},
						},
						Run: model.CallSpec{
							Container: &model.ContainerCallSpec{
								Image: &model.ContainerCallImageSpec{
									Ref: imageRef,
								},
							},
						},
					},
					providedOpRef,
					&providedParentID,
					providedRootID,
				)

				/* assert */
				Expect(actualErr).To(BeNil())

				actualChildCalls := []model.CallStarted{}
				go func() {
					for event := range eventChannel {
						if nil != event.CallStarted && nil != event.CallStarted.Call.Container {
							// ignore props we can't readily assert
							event.CallStarted.Call.Container.ContainerID = ""
							event.CallStarted.Call.ID = ""

							actualChildCalls = append(actualChildCalls, *event.CallStarted)
						}
					}
				}()

				Eventually(
					func() []model.CallStarted { return actualChildCalls },
				).Should(
					ConsistOf(
						[]model.CallStarted{
							{
								Call: model.Call{
									Container: &model.ContainerCall{
										BaseCall: model.BaseCall{
											OpPath: providedOpRef,
										},
										Cmd:   []string{},
										Dirs:  map[string]string{},
										Files: map[string]string{},
										Image: &model.ContainerCallImage{
											Ref: &imageRef,
										},
										Sockets: map[string]string{},
									},
									ParentID: &providedParentID,
									RootID:   providedRootID,
								},
								Ref: providedOpRef,
							},
							{
								Call: model.Call{
									Container: &model.ContainerCall{
										BaseCall: model.BaseCall{
											OpPath: providedOpRef,
										},
										Cmd:   []string{},
										Dirs:  map[string]string{},
										Files: map[string]string{},
										Image: &model.ContainerCallImage{
											Ref: &imageRef,
										},
										Sockets: map[string]string{},
									},
									ParentID: &providedParentID,
									RootID:   providedRootID,
								},
								Ref: providedOpRef,
							},
						},
					),
				)
			})
		})
	})
})
