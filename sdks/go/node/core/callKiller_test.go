package core

import (
	"context"
	"os"
	"time"

	"github.com/dgraph-io/badger/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/fakes"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

var _ = Context("_callKiller", func() {
	Context("Kill", func() {
		Context("stateStore.ListWithParentID returns nodes", func() {
			It("should call pubsub.Publish for each", func() {
				/* arrange */
				providedCallID := "providedCallID"
				providedRootCallID := "providedRootCallID"

				expectedChildCallIDs := []string{
					"childCallID1",
					"childCallID2",
					"childCallID3",
				}

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
				eventChannel, err := pubSub.Subscribe(
					context.Background(),
					model.EventFilter{},
				)
				if err != nil {
					panic(err)
				}

				stateStore := newStateStore(context.Background(), db, pubSub)

				// seed call
				pubSub.Publish(model.Event{
					CallStarted: &model.CallStarted{
						Call: model.Call{
							ID:     providedCallID,
							RootID: providedRootCallID,
						},
					},
					Timestamp: time.Now().UTC(),
				})

				// seed child calls
				for _, childCallID := range expectedChildCallIDs {
					pubSub.Publish(model.Event{
						CallStarted: &model.CallStarted{
							Call: model.Call{
								ID:       childCallID,
								ParentID: &providedCallID,
								RootID:   providedRootCallID,
							},
						},
						Timestamp: time.Now().UTC(),
					})
				}

				// give stateStore time to receive & apply events
				time.Sleep(time.Second)

				objectUnderTest := newCallKiller(
					stateStore,
					new(FakeContainerRuntime),
					pubSub,
				)

				/* act */
				objectUnderTest.Kill(
					context.Background(),
					providedCallID,
					providedRootCallID,
				)

				/* assert */
				actualChildCallIDs := []string{}
				go func() {
					for event := range eventChannel {
						if event.CallKillRequested != nil {
							actualChildCallIDs = append(actualChildCallIDs, event.CallKillRequested.Request.OpID)
						}
					}
				}()

				Eventually(
					func() []string { return actualChildCallIDs },
				).Should(
					ContainElements(expectedChildCallIDs),
				)
			})
		})
	})
})
