package node

import (
	"context"
	"os"
	"time"

	"github.com/dgraph-io/badger/v3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

var _ = Context("stateStore", func() {
	Context("TryGetAuth", func() {
		Context("AuthAdded", func() {
			It("should return expected auth", func() {

				/* arrange */
				providedReq := model.AddAuthReq{
					Creds: model.Creds{
						Username: "username",
						Password: "password",
					},
					Resources: "resources",
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

				expectedAuth := model.Auth{
					Creds:     providedReq.Creds,
					Resources: providedReq.Resources,
				}
				// seed auth
				pubSub.Publish(model.Event{
					AuthAdded: &model.AuthAdded{
						Auth: expectedAuth,
					},
					Timestamp: time.Now().UTC(),
				})

				objectUnderTest := newStateStore(
					context.Background(),
					db,
					pubSub,
				)

				// give stateStore time to receive & apply events
				time.Sleep(time.Second)

				/* act */
				objectUnderTest.TryGetAuth(
					expectedAuth.Resources,
				)

				/* assert */
				var actualAuth model.Auth
				go func() {
					for event := range eventChannel {
						if event.AuthAdded != nil {
							actualAuth = event.AuthAdded.Auth
						}
					}
				}()

				Eventually(
					func() model.Auth { return actualAuth },
				).Should(
					Equal(expectedAuth),
				)
			})
		})
	})
})
