package pubsub

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/dgraph-io/badger/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("pubSub", func() {
	dbDir, err := ioutil.TempDir("", "")
	if nil != err {
		panic(err)
	}

	db, err := badger.Open(
		badger.DefaultOptions(dbDir).WithLogger(nil),
	)
	if nil != err {
		panic(err)
	}

	Context("New", func() {
		It("should return PubSub", func() {
			/* arrange/act/assert */

			Expect(New(db)).To(Not(BeNil()))
		})
	})
	Context("Publish", func() {
		Context("subscription exist", func() {
			Context("is subscribed", func() {
				It("receives event", func() {
					/* arrange */
					db.DropAll()

					expectedEvent := model.Event{
						CallStarted: &model.CallStarted{
							Call: model.Call{
								RootID: "rootID",
							},
						},
					}

					objectUnderTest := New(db)

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
					subscriberEventFilter := model.EventFilter{Roots: []string{"notPublishedRootCallID"}}

					publishedEvent := model.Event{
						CallStarted: &model.CallStarted{
							Call: model.Call{
								RootID: "rootID",
							},
						},
					}

					objectUnderTest := New(db)

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
		Context("one publish has occurred", func() {
			Context("no filter", func() {
				It("should receive published event", func() {
					/* arrange */
					db.DropAll()

					expectedEvent := model.Event{
						CallStarted: &model.CallStarted{
							Call: model.Call{
								ID: "id",
							},
						},
					}

					objectUnderTest := New(db)
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
				It("should receive published event", func() {
					/* arrange */
					db.DropAll()

					expectedEvent := model.Event{
						CallStarted: &model.CallStarted{
							Call: model.Call{
								RootID: "rootId",
							},
						},
					}

					providedFilter := model.EventFilter{
						Roots: []string{
							expectedEvent.CallStarted.Call.RootID,
						},
					}

					objectUnderTest := New(db)
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
		Context("two publishes have occurred", func() {
			Context("no filter", func() {
				It("should receive published events", func() {
					/* arrange */
					db.DropAll()

					expectedEvent1 := model.Event{
						CallStarted: &model.CallStarted{
							Call: model.Call{
								ID: "id",
							},
						},
						Timestamp: time.Now().Add(-time.Second),
					}

					expectedEvent2 := model.Event{
						CallStarted: &model.CallStarted{
							Call: model.Call{
								RootID: "rootID",
							},
						},
						Timestamp: time.Now(),
					}

					notExpectedEvent := model.Event{
						CallStarted: &model.CallStarted{
							Call: model.Call{
								RootID: "rootID",
							},
						},
						Timestamp: time.Now().Add(time.Second),
					}

					objectUnderTest := New(db)
					objectUnderTest.Publish(expectedEvent1)
					objectUnderTest.Publish(notExpectedEvent)
					objectUnderTest.Publish(expectedEvent2)

					/* act */
					eventChannel, err := objectUnderTest.Subscribe(context.TODO(), model.EventFilter{})

					/* assert */
					Expect(err).To(BeNil())

					var actualEvent1 model.Event
					Eventually(eventChannel).Should(Receive(&actualEvent1))
					// ignore timestamp
					actualEvent1.Timestamp = expectedEvent1.Timestamp
					Expect(actualEvent1).To(Equal(expectedEvent1))

					var actualEvent2 model.Event
					Eventually(eventChannel).Should(Receive(&actualEvent2))
					// ignore timestamp
					actualEvent2.Timestamp = expectedEvent2.Timestamp
					Expect(actualEvent2).To(Equal(expectedEvent2))
				})
			})
		})
	})
})
