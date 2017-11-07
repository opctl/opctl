package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"time"
)

var _ = Context("core", func() {
	Context("KillOp", func() {
		It("should call pubSub.Publish w/ expected args", func() {
			/* arrange */
			providedReq := model.KillOpReq{RootOpId: "dummyRootOpId"}

			fakePubSub := new(pubsub.Fake)

			objectUnderTest := _core{
				pubSub: fakePubSub,
			}

			expectedEvent := model.Event{
				Timestamp: time.Now().UTC(),
				OpKilled: &model.OpKilledEvent{
					RootOpId: providedReq.RootOpId,
				},
			}

			/* act */
			objectUnderTest.KillOp(providedReq)

			/* assert */
			actualEvent := fakePubSub.PublishArgsForCall(0)

			// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
			Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
			// set temporal fields to expected vals since they're already asserted
			actualEvent.Timestamp = expectedEvent.Timestamp

			Expect(*actualEvent).To(Equal(expectedEvent))
		})
	})
})
