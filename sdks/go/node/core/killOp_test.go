package core

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("core", func() {
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

			fakePubSub := new(FakePubSub)

			objectUnderTest := core{
				pubSub: fakePubSub,
			}

			/* act */
			objectUnderTest.KillOp(
				context.Background(),
				providedReq,
			)

			/* assert */
			actualEvent := fakePubSub.PublishArgsForCall(0)

			// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
			Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
			// set temporal fields to expected vals since they're already asserted
			actualEvent.Timestamp = expectedEvent.Timestamp

			Expect(actualEvent).To(Equal(expectedEvent))
		})
	})
})
