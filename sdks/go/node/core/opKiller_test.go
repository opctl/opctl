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

var _ = Context("_opKiller", func() {
	Context("Kill", func() {
		It("should call stateStore.ListWithParentID w/ expected args", func() {
			/* arrange */
			providedCallID := "providedCallID"

			fakeStateStore := new(FakeStateStore)

			objectUnderTest := _opKiller{
				containerRuntime: new(FakeContainerRuntime),
				stateStore:       fakeStateStore,
				eventPublisher:   new(FakeEventPublisher),
			}

			/* act */
			objectUnderTest.Kill(
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
					{Id: "dummyNode1Id"},
					{Id: "dummyNode2Id"},
					{Id: "dummyNode3Id"},
				}
				fakeStateStore := new(FakeStateStore)
				fakeStateStore.ListWithParentIDReturnsOnCall(0, nodesReturnedFromStateStore)

				db, err := badger.Open(
					badger.DefaultOptions(os.TempDir()).WithLogger(nil),
				)
				if err != nil {
					panic(err)
				}

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

				objectUnderTest := _opKiller{
					containerRuntime: new(FakeContainerRuntime),
					stateStore:       fakeStateStore,
					eventPublisher:   pubSub,
				}

				/* act */
				objectUnderTest.Kill(
					providedCallID,
					"rootCallID",
				)

				/* assert */
				cancelEventChan()

				Expect(len(actualCalls)).To(Equal(len(nodesReturnedFromStateStore) + 1))
			})
		})
	})
})
