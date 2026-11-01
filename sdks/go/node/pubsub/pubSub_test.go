package pubsub

import (
	"context"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("pubSub", func() {
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
		It("removes a cancelled subscription without another publish", func() {
			Expect(db.DropAll()).To(Succeed())
			objectUnderTest := New(db).(*pubSub)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			eventChannel, err := objectUnderTest.Subscribe(ctx, model.EventFilter{})
			Expect(err).To(BeNil())
			subscriptionCount := func() int {
				objectUnderTest.subscriptionsMutex.RLock()
				defer objectUnderTest.subscriptionsMutex.RUnlock()
				return len(objectUnderTest.subscriptions)
			}
			Eventually(subscriptionCount).Should(Equal(1))

			cancel()

			Eventually(subscriptionCount).Should(Equal(0))
			Eventually(eventChannel).Should(BeClosed())
		})

		It("releases blocked publishers when a subscription is cancelled", func() {
			Expect(db.DropAll()).To(Succeed())
			objectUnderTest := New(db).(*pubSub)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			eventChannel, err := objectUnderTest.Subscribe(ctx, model.EventFilter{})
			Expect(err).To(BeNil())
			var publishEventChannel chan model.Event
			var info subscriptionInfo
			Eventually(func() bool {
				objectUnderTest.subscriptionsMutex.RLock()
				defer objectUnderTest.subscriptionsMutex.RUnlock()
				for channel, subscription := range objectUnderTest.subscriptions {
					publishEventChannel, info = channel, subscription
					return true
				}
				return false
			}).Should(BeTrue())

			const publisherCount = 10
			published := make(chan struct{}, publisherCount)
			for i := 0; i < publisherCount; i++ {
				go func() {
					objectUnderTest.publishToSubscription(publishEventChannel, info, model.Event{})
					published <- struct{}{}
				}()
			}
			// The subscription can accept one event before blocking on its unread output.
			Eventually(published).Should(Receive())
			Consistently(published).ShouldNot(Receive())

			cancel()

			Eventually(info.Done).Should(BeClosed())
			Eventually(eventChannel).Should(BeClosed())
			for i := 1; i < publisherCount; i++ {
				Eventually(published).Should(Receive())
			}
			objectUnderTest.subscriptionsMutex.RLock()
			subscriptionCount := len(objectUnderTest.subscriptions)
			objectUnderTest.subscriptionsMutex.RUnlock()
			Expect(subscriptionCount).To(Equal(0))
		})

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
						Timestamp: time.Now(),
					}

					expectedEvent2 := model.Event{
						CallStarted: &model.CallStarted{
							Call: model.Call{
								RootID: "rootID",
							},
						},
						Timestamp: time.Now().Add(time.Second),
					}

					objectUnderTest := New(db)
					objectUnderTest.Publish(expectedEvent1)
					objectUnderTest.Publish(expectedEvent2)

					/* act */
					eventChannel, _ := objectUnderTest.Subscribe(context.TODO(), model.EventFilter{})

					/* assert */
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
