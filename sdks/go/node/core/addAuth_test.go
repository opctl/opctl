package core

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/dgraph-io/badger/v3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

var _ = Context("core", func() {
	Context("AddAuth", func() {
		It("should call opAdder.Add w/ expected args", func() {

			/* arrange */
			providedReq := model.AddAuthReq{
				Creds: model.Creds{
					Username: "username",
					Password: "password",
				},
				Resources: "resources",
			}

			dbDir, err := ioutil.TempDir("", "")
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

			expectedEvent := model.Event{
				AuthAdded: &model.AuthAdded{
					Auth: model.Auth{
						Creds:     providedReq.Creds,
						Resources: providedReq.Resources,
					},
				},
				Timestamp: time.Now().UTC(),
			}

			objectUnderTest := core{
				pubSub: pubSub,
			}

			/* act */
			objectUnderTest.AddAuth(
				context.Background(),
				providedReq,
			)

			/* assert */
			var actualEvent model.Event
			go func() {
				for event := range eventChannel {
					if event.AuthAdded != nil {
						// ignore timestamp from assertion
						event.Timestamp = expectedEvent.Timestamp
						actualEvent = event
					}
				}
			}()

			Eventually(
				func() model.Event { return actualEvent },
			).Should(
				Equal(expectedEvent),
			)
		})
	})
})
