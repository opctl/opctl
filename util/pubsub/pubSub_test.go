package pubsub

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Context("pubSub", func() {
	Context("New", func() {
		It("should return PubSub", func() {
			/* arrange/act/assert */
			Expect(New()).Should(Not(BeNil()))
		})
	})
	Context("Publish", func() {
		Context("subscription exist", func() {
			Context("is subscribed", func() {
				It("receives event", func() {
					/* arrange */
					subscriberChannel := make(chan *model.Event, 1000)

					expectedEvent := &model.Event{
						OpStarted: &model.OpStartedEvent{
							OpGraphId: "dummyOpGraphId",
							OpId:      "dummyOpId",
							OpRef:     "dummyOpRef",
						},
					}

					objectUnderTest := New()

					objectUnderTest.Subscribe(nil, subscriberChannel)

					/* act */
					objectUnderTest.Publish(expectedEvent)

					/* assert */
					var actualEvent *model.Event
					Eventually(subscriberChannel).Should(Receive(&actualEvent))
					Expect(actualEvent).To(Equal(expectedEvent))
				})
			})
			Context("isn't subscribed", func() {
				It("doesn't receive event", func() {
					/* arrange */
					subscriberEventFilter := &model.EventFilter{OpGraphIds: []string{"notPublishedOgid"}}
					subscriberChannel := make(chan *model.Event, 1000)

					publishedEvent := &model.Event{
						OpStarted: &model.OpStartedEvent{
							OpGraphId: "dummyOpGraphId",
							OpId:      "dummyOpId",
							OpRef:     "dummyOpRef",
						},
					}

					objectUnderTest := New()

					objectUnderTest.Subscribe(subscriberEventFilter, subscriberChannel)

					/* act */
					objectUnderTest.Publish(publishedEvent)

					/* assert */
					Consistently(subscriberChannel).ShouldNot(Receive())
				})
			})
		})
	})
	Context("Subscribe", func() {
		Context("previous publish exists", func() {
			Context("no filter", func() {
				It("should receive previous event", func() {
					/* arrange */
					subscriberChannel := make(chan *model.Event, 1000)

					expectedEvent := &model.Event{
						OpStarted: &model.OpStartedEvent{
							OpGraphId: "dummyOpGraphId",
							OpId:      "dummyOpId",
							OpRef:     "dummyOpRef",
						},
					}

					objectUnderTest := New()
					objectUnderTest.Publish(expectedEvent)

					/* act */
					objectUnderTest.Subscribe(nil, subscriberChannel)

					/* assert */
					var actualEvent *model.Event
					Eventually(subscriberChannel).Should(Receive(&actualEvent))
					Expect(actualEvent).To(Equal(expectedEvent))
				})
			})
			Context("filter allows previous publish", func() {
				It("should receive previous event", func() {
					/* arrange */
					subscriberChannel := make(chan *model.Event, 1000)

					expectedEvent := &model.Event{
						OpStarted: &model.OpStartedEvent{
							OpGraphId: "dummyOpGraphId",
							OpId:      "dummyOpId",
							OpRef:     "dummyOpRef",
						},
					}

					providedFilter := &model.EventFilter{
						OpGraphIds: []string{
							expectedEvent.OpStarted.OpGraphId,
						},
					}

					objectUnderTest := New()
					objectUnderTest.Publish(expectedEvent)

					/* act */
					objectUnderTest.Subscribe(providedFilter, subscriberChannel)

					/* assert */
					var actualEvent *model.Event
					Eventually(subscriberChannel).Should(Receive(&actualEvent))
					Expect(actualEvent).To(Equal(expectedEvent))
				})
			})
		})
	})
})
