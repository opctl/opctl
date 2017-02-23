package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/containerprovider"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Context("core", func() {
	Context("GetEventStream", func() {
		It("should call pubSub.RegisterSubscriber w/ expected args", func() {
			/* arrange */
			providedReq := &model.GetEventStreamReq{
				Filter: &model.EventFilter{
					OpGraphIds: []string{
						"dummyOpGraphId",
					},
				},
			}

			providedEventStream := make(chan *model.Event)

			fakePubSub := new(pubsub.Fake)

			objectUnderTest := _core{
				containerProvider:   new(containerprovider.Fake),
				pubSub:              fakePubSub,
				opCaller:            new(fakeOpCaller),
				dcgNodeRepo:         new(fakeDcgNodeRepo),
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
