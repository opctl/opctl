package core

import (
	"context"
	"os"

	"github.com/dgraph-io/badger/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/fakes"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	"github.com/opctl/opctl/sdks/go/pubsub"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("_callKiller", func() {
	Context("Kill", func() {
		It("should call stateStore.ListWithParentID w/ expected args", func() {
			/* arrange */
			providedCallID := "providedCallID"

			fakeStateStore := new(FakeStateStore)

			objectUnderTest := _callKiller{
				containerRuntime: new(FakeContainerRuntime),
				stateStore:       fakeStateStore,
				eventPublisher:   new(FakeEventPublisher),
			}

			/* act */
			objectUnderTest.Kill(
				context.Background(),
				providedCallID,
				"rootCallID",
			)

			/* assert */
			Expect(fakeStateStore.ListWithParentIDArgsForCall(0)).To(Equal(providedCallID))
		})
		Context("stateStore.ListWithParentID returns nodes", func() {
			It("should call pubsub.Publish for each", func() {
				/* arrange */
				providedCallID := "providedCallID"

				nodesReturnedFromStateStore := []*model.Call{
					{ID: "dummyNode1ID"},
					{ID: "dummyNode2ID"},
					{ID: "dummyNode3ID"},
				}
				fakeStateStore := new(FakeStateStore)
				fakeStateStore.ListWithParentIDReturnsOnCall(0, nodesReturnedFromStateStore)

				db, err := badger.Open(
					badger.DefaultOptions(os.TempDir()).WithLogger(nil),
				)
				if err != nil {
					panic(err)
				}
				db.DropAll()

				pubSub := pubsub.New(
					db,
				)

				eventChanCtx, cancelEventChan := context.WithCancel(context.Background())
				eventChannel, _ := pubSub.Subscribe(
					eventChanCtx,
					model.EventFilter{},
				)

				actualCalls := []model.Event{}
				go func() {
					for event := range eventChannel {
						actualCalls = append(actualCalls, event)
					}
				}()

				objectUnderTest := _callKiller{
					containerRuntime: new(FakeContainerRuntime),
					stateStore:       fakeStateStore,
					eventPublisher:   pubSub,
				}

				/* act */
				objectUnderTest.Kill(
					context.Background(),
					providedCallID,
					"rootCallID",
				)

				/* assert */
				cancelEventChan()
				Eventually(func() int {
					return len(actualCalls)
				}).Should(Equal(len(nodesReturnedFromStateStore)))
			})
		})
	})
})
