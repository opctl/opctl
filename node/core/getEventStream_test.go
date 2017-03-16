package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/containerprovider"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("core", func() {
	Context("GetEventStream", func() {
		It("should call pubSub.RegisterSubscriber w/ expected args", func() {
			/* arrange */
			providedReq := &model.GetEventStreamReq{
				Filter: &model.EventFilter{
					RootOpIds: []string{
						"dummyRootOpId",
					},
				},
			}

			providedEventStream := make(chan *model.Event, 150)

			fakePubSub := new(pubsub.Fake)

			objectUnderTest := _core{
				containerProvider:   new(containerprovider.Fake),
				pubSub:              fakePubSub,
				opCaller:            new(fakeOpCaller),
				dcgNodeRepo:         new(fakeDCGNodeRepo),
				uniqueStringFactory: new(uniquestring.Fake),
			}

			/* act */
			objectUnderTest.GetEventStream(providedReq, providedEventStream)

			/* assert */

			actualFilter,
				actualEventChannel := fakePubSub.SubscribeArgsForCall(0)

			Expect(actualFilter).To(Equal(providedReq.Filter))
			Expect(actualEventChannel).To(Equal(providedEventStream))
		})
	})
})
