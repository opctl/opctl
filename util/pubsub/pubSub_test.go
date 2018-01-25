package pubsub

import (
	"context"
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
	tempEventRepo := NewBuntDBEventRepo(tempFilePath.Name())

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
					expectedEvent := model.Event{
						OpStarted: &model.OpStartedEvent{
							RootOpId: "dummyRootOpId",
							OpId:     "dummyOpId",
							PkgRef:   "dummyPkgRef",
						},
					}

					objectUnderTest := New(tempEventRepo)

					eventChannel, _ := objectUnderTest.Subscribe(context.TODO(), model.EventFilter{})

					/* act */
					objectUnderTest.Publish(expectedEvent)

					/* assert */
					var actualEvent model.Event
					Eventually(eventChannel).Should(Receive(&actualEvent))
					Expect(actualEvent).To(Equal(expectedEvent))
				})
			})
			Context("isn't subscribed", func() {
				It("doesn't receive event", func() {
					/* arrange */
					subscriberEventFilter := model.EventFilter{Roots: []string{"notPublishedRootOpId"}}

					publishedEvent := model.Event{
						OpStarted: &model.OpStartedEvent{
							RootOpId: "dummyRootOpId",
							OpId:     "dummyOpId",
							PkgRef:   "dummyPkgRef",
						},
					}

					objectUnderTest := New(tempEventRepo)

					eventChannel, _ := objectUnderTest.Subscribe(context.TODO(), subscriberEventFilter)

					/* act */
					objectUnderTest.Publish(publishedEvent)

					/* assert */
					Consistently(eventChannel).ShouldNot(Receive())
				})
			})
		})
	})
	Context("Subscribe", func() {
		Context("previous publish exists", func() {
			Context("no filter", func() {
				It("should receive previous event", func() {
					/* arrange */
					expectedEvent := model.Event{
						OpStarted: &model.OpStartedEvent{
							RootOpId: "dummyRootOpId",
							OpId:     "dummyOpId",
							PkgRef:   "dummyPkgRef",
						},
					}

					objectUnderTest := New(tempEventRepo)
					objectUnderTest.Publish(expectedEvent)

					/* act */
					eventChannel, _ := objectUnderTest.Subscribe(context.TODO(), model.EventFilter{})

					/* assert */
					var actualEvent model.Event
					Eventually(eventChannel).Should(Receive(&actualEvent))
					// ignore timestamp
					actualEvent.Timestamp = expectedEvent.Timestamp
					Expect(actualEvent).To(Equal(expectedEvent))
				})
			})
			Context("filter allows previous publish", func() {
				It("should receive previous event", func() {
					/* arrange */
					expectedEvent := model.Event{
						OpStarted: &model.OpStartedEvent{
							RootOpId: "dummyRootOpId",
							OpId:     "dummyOpId",
							PkgRef:   "dummyPkgRef",
						},
					}

					providedFilter := model.EventFilter{
						Roots: []string{
							expectedEvent.OpStarted.RootOpId,
						},
					}

					objectUnderTest := New(tempEventRepo)
					objectUnderTest.Publish(expectedEvent)

					/* act */
					eventChannel, _ := objectUnderTest.Subscribe(context.TODO(), providedFilter)

					/* assert */
					var actualEvent model.Event
					Eventually(eventChannel).Should(Receive(&actualEvent))
					// ignore timestamp
					actualEvent.Timestamp = expectedEvent.Timestamp
					Expect(actualEvent).To(Equal(expectedEvent))
				})
			})
		})
	})
})
