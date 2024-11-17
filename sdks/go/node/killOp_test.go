package node

import (
	"context"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

var _ = Context("core", func() {
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

	Context("KillOp", func() {
		It("should call callKiller.Kill w/ expected args", func() {

			/* arrange */
			providedReq := model.KillOpReq{
				OpID: "dummyOpID",
			}

			expectedEvent := model.Event{
				CallKillRequested: &model.CallKillRequested{
					Request: providedReq,
				},
				Timestamp: time.Now().UTC(),
			}

			pubSub := pubsub.New(db)

			eventChannel, err := pubSub.Subscribe(
				context.Background(),
				model.EventFilter{},
			)
			if err != nil {
				panic(err)
			}

			objectUnderTest := core{
				pubSub: pubSub,
			}

			/* act */
			objectUnderTest.KillOp(
				context.Background(),
				providedReq,
			)

			/* assert */

			for {
				actualEvent := <-eventChannel

				// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
				Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
				// set temporal fields to expected vals since they're already asserted
				actualEvent.Timestamp = expectedEvent.Timestamp

				Expect(actualEvent).To(Equal(expectedEvent))

				break
			}
		})
	})
})
