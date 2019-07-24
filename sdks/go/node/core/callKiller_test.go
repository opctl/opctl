package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/types"
	"github.com/opctl/opctl/sdks/go/util/pubsub"
)

var _ = Context("core", func() {
	Context("Kill", func() {
		It("should call callStore.SetIsKilled w/ expected args", func() {
			/* arrange */
			providedCallID := "providedCallID"

			fakeCallStore := new(fakeCallStore)

			objectUnderTest := _callKiller{
				containerRuntime: new(containerruntime.Fake),
				callStore:        fakeCallStore,
				eventPublisher:   new(pubsub.FakeEventPublisher),
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

			fakeCallStore := new(fakeCallStore)

			objectUnderTest := _callKiller{
				containerRuntime: new(containerruntime.Fake),
				callStore:        fakeCallStore,
				eventPublisher:   new(pubsub.FakeEventPublisher),
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

				nodesReturnedFromCallStore := []*types.DCG{
					{Id: "dummyNode1Id"},
					{Id: "dummyNode2Id"},
					{Id: "dummyNode3Id"},
				}
				fakeCallStore := new(fakeCallStore)
				fakeCallStore.ListWithParentIDReturnsOnCall(0, nodesReturnedFromCallStore)

				// use map so order ignored; calls happen in parallel so all ordering bets are off
				expectedCalls := map[string]bool{
					providedCallID:                   true,
					nodesReturnedFromCallStore[0].Id: true,
					nodesReturnedFromCallStore[1].Id: true,
					nodesReturnedFromCallStore[2].Id: true,
				}

				objectUnderTest := _callKiller{
					containerRuntime: new(containerruntime.Fake),
					callStore:        fakeCallStore,
					eventPublisher:   new(pubsub.FakeEventPublisher),
				}

				/* act */
				objectUnderTest.Kill(
					providedCallID,
					"rootCallID",
				)

				/* assert */
				actualCalls := map[string]bool{}
				callIndex := 0
				for callIndex < fakeCallStore.SetIsKilledCallCount() {
					actualCalls[fakeCallStore.SetIsKilledArgsForCall(callIndex)] = true
					callIndex++
				}

				Expect(actualCalls).To(Equal(expectedCalls))
			})
			It("should call containerRuntime.DeleteContainerIfExists w/ expected args", func() {
				/* arrange */
				providedCallID := "providedCallID"

				nodesReturnedFromCallStore := []*types.DCG{
					{Id: "dummyNode1Id"},
					{Id: "dummyNode2Id"},
					{Id: "dummyNode3Id"},
				}
				fakeCallStore := new(fakeCallStore)
				fakeCallStore.ListWithParentIDReturnsOnCall(0, nodesReturnedFromCallStore)

				// use map so order ignored; calls happen in parallel so all ordering bets are off
				expectedCalls := map[string]bool{
					providedCallID:                   true,
					nodesReturnedFromCallStore[0].Id: true,
					nodesReturnedFromCallStore[1].Id: true,
					nodesReturnedFromCallStore[2].Id: true,
				}

				fakeContainerRuntime := new(containerruntime.Fake)

				objectUnderTest := &_callKiller{
					containerRuntime: fakeContainerRuntime,
					callStore:        fakeCallStore,
					eventPublisher:   new(pubsub.FakeEventPublisher),
				}

				/* act */
				objectUnderTest.Kill(
					providedCallID,
					"rootCallID",
				)

				/* assert */
				actualCalls := map[string]bool{}
				callIndex := 0
				for callIndex < fakeContainerRuntime.DeleteContainerIfExistsCallCount() {
					actualCalls[fakeContainerRuntime.DeleteContainerIfExistsArgsForCall(callIndex)] = true
					callIndex++
				}

				Expect(actualCalls).To(Equal(expectedCalls))
			})
		})
	})
})
