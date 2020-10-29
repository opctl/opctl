package core

import (
	"context"
	"os"

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
		It("should call callStore.SetIsKilled w/ expected args", func() {
			/* arrange */
			providedCallID := "providedCallID"

			fakeCallStore := new(FakeCallStore)

			objectUnderTest := _opKiller{
				containerRuntime: new(FakeContainerRuntime),
				callStore:        fakeCallStore,
				eventPublisher:   new(FakeEventPublisher),
			}

			/* act */
			objectUnderTest.Kill(
				providedCallID,
				"rootCallID",
			)

			/* assert */
			Expect(fakeCallStore.SetIsKilledArgsForCall(0)).To(Equal(providedCallID))
		})
		It("should call callStore.ListWithParentID w/ expected args", func() {
			/* arrange */
			providedCallID := "providedCallID"

			fakeCallStore := new(FakeCallStore)

			objectUnderTest := _opKiller{
				containerRuntime: new(FakeContainerRuntime),
				callStore:        fakeCallStore,
				eventPublisher:   new(FakeEventPublisher),
			}

			/* act */
			objectUnderTest.Kill(
				providedCallID,
				"rootCallID",
			)

			/* assert */
			Expect(fakeCallStore.ListWithParentIDArgsForCall(0)).To(Equal(providedCallID))
		})
		Context("callStore.ListWithParentID returns nodes", func() {
			It("should call callStore.SetIsKilled w/ expected args for each", func() {
				/* arrange */
				providedCallID := "providedCallID"

				nodesReturnedFromCallStore := []*model.Call{
					{Id: "dummyNode1Id"},
					{Id: "dummyNode2Id"},
					{Id: "dummyNode3Id"},
				}
				fakeCallStore := new(FakeCallStore)
				fakeCallStore.ListWithParentIDReturnsOnCall(0, nodesReturnedFromCallStore)

				tmpDir := os.TempDir()

				pubSub := pubsub.New(
					pubsub.NewBadgerDBEventStore(tmpDir),
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
					callStore:        fakeCallStore,
					eventPublisher:   pubSub,
				}

				/* act */
				objectUnderTest.Kill(
					providedCallID,
					"rootCallID",
				)

				/* assert */
				cancelEventChan()

				Expect(len(actualCalls)).To(Equal(len(nodesReturnedFromCallStore) + 1))
			})
		})
	})
})
