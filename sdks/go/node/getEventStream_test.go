package node

import (
	"context"
	"os"

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

	Context("GetEventStream", func() {
		It("should return subscription & not error", func() {
			/* arrange */
			providedCtx := context.TODO()
			providedReq := &model.GetEventStreamReq{
				Filter: model.EventFilter{
					Roots: []string{
						"dummyRootCallID",
					},
				},
			}

			objectUnderTest := core{
				pubSub: pubsub.New(db),
			}

			/* act */
			actualSubscription, actualErr := objectUnderTest.GetEventStream(
				providedCtx,
				providedReq,
			)

			/* assert */

			Expect(actualErr).To(BeNil())
			Expect(actualSubscription).NotTo(BeClosed())
		})
	})
})
