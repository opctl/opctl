package core

import (
	"context"
	"io"
	"os"

	"github.com/dgraph-io/badger/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	containerRuntimeFakes "github.com/opctl/opctl/sdks/go/node/core/containerruntime/fakes"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	"github.com/opctl/opctl/sdks/go/pubsub"
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
					caller: fakeCaller,
					pubSub: new(FakePubSub),
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

		Context("iteration spec invalid", func() {

			It("should return expected results", func() {
				/* arrange */
				dbDir, err := os.MkdirTemp("", "")
				if err != nil {
					panic(err)
				}

				db, err := badger.Open(
					badger.DefaultOptions(dbDir).WithLogger(nil),
				)
				if err != nil {
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
							providedCtx,
							db,
							pubSub,
						),
					),
					dbDir,
					pubSub,
				)

				objectUnderTest := _parallelLoopCaller{
					caller: caller,
					pubSub: pubSub,
				}

				/* act */
				actualOutputs, actualErr := objectUnderTest.Call(
					providedCtx,
					"id",
					providedScope,
					model.ParallelLoopCallSpec{
						Range: model.Value{
							Array: &[]interface{}{0},
						},
						Run: model.CallSpec{
							Container: &model.ContainerCallSpec{},
						},
					},
					"opPath",
					new(string),
					"rootCallID",
				)

				/* assert */
				Expect(actualErr.Error()).To(Equal("child call failed"))
				Expect(actualOutputs).To(BeNil())
			})
		})

		It("should start each child as expected", func() {
			/* arrange */
			dbDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			db, err := badger.Open(
				badger.DefaultOptions(dbDir).WithLogger(nil),
			)
			if err != nil {
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
			if err != nil {
				panic(err)
			}

			objectUnderTest := _parallelLoopCaller{
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
				model.ParallelLoopCallSpec{
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
					if event.CallStarted != nil && event.CallStarted.Call.Container != nil {
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
				ContainElements(
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
