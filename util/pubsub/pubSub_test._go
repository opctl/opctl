package pubsub

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"io/ioutil"
)

var _ = Context("pubSub", func() {
	tempFilePath, err := ioutil.TempFile("", "")
	if nil != err {
		panic(err)
	}
	tempEventRepo := NewEventRepo(tempFilePath.Name())

	Context("New", func() {
		It("should return PubSub", func() {
			/* arrange/act/assert */

			Expect(New(tempEventRepo)).To(Not(BeNil()))
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
							RootOpId: "dummyRootOpId",
							OpId:     "dummyOpId",
							PkgRef:   "dummyPkgRef",
						},
					}

					objectUnderTest := New(tempEventRepo)

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
					subscriberEventFilter := &model.EventFilter{Roots: []string{"notPublishedRootOpId"}}
					subscriberChannel := make(chan *model.Event, 1000)

					publishedEvent := &model.Event{
						OpStarted: &model.OpStartedEvent{
							RootOpId: "dummyRootOpId",
							OpId:     "dummyOpId",
							PkgRef:   "dummyPkgRef",
						},
					}

					objectUnderTest := New(tempEventRepo)

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
							RootOpId: "dummyRootOpId",
							OpId:     "dummyOpId",
							PkgRef:   "dummyPkgRef",
						},
					}

					objectUnderTest := New(tempEventRepo)
					objectUnderTest.Publish(expectedEvent)

					/* act */
					objectUnderTest.Subscribe(nil, subscriberChannel)

					/* assert */
					var actualEvent *model.Event
					Eventually(subscriberChannel).Should(Receive(&actualEvent))
					// ignore timestamp
					actualEvent.Timestamp = expectedEvent.Timestamp
					Expect(actualEvent).To(Equal(expectedEvent))
				})
			})
			Context("filter allows previous publish", func() {
				It("should receive previous event", func() {
					/* arrange */
					subscriberChannel := make(chan *model.Event, 1000)

					expectedEvent := &model.Event{
						OpStarted: &model.OpStartedEvent{
							RootOpId: "dummyRootOpId",
							OpId:     "dummyOpId",
							PkgRef:   "dummyPkgRef",
						},
					}

					providedFilter := &model.EventFilter{
						Roots: []string{
							expectedEvent.OpStarted.RootOpId,
						},
					}

					objectUnderTest := New(tempEventRepo)
					objectUnderTest.Publish(expectedEvent)

					/* act */
					objectUnderTest.Subscribe(providedFilter, subscriberChannel)

					/* assert */
					var actualEvent *model.Event
					Eventually(subscriberChannel).Should(Receive(&actualEvent))
					// ignore timestamp
					actualEvent.Timestamp = expectedEvent.Timestamp
					Expect(actualEvent).To(Equal(expectedEvent))
				})
			})
		})
	})
})
