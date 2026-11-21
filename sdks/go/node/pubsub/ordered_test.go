package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Describe("Ordered event subscriptions", func() {
	It("delivers replay and a burst of live events to a slow subscriber once in order", func() {
		db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ps := New(db)
		now := time.Now().UTC()
		publish := func(i int) {
			ps.Publish(model.Event{
				Timestamp: now.Add(time.Duration(i)),
				ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{
					Data: []byte(fmt.Sprint(i)), RootCallID: "root",
				},
			})
		}
		publish(0)
		events, err := ps.Subscribe(ctx, model.EventFilter{Roots: []string{"root"}, Ordered: true})
		Expect(err).NotTo(HaveOccurred())
		for i := 1; i < 1100; i++ {
			publish(i)
		}
		deadline := time.After(time.Second)
		for i := 0; i < 1100; i++ {
			select {
			case event, open := <-events:
				Expect(open).To(BeTrue(), "stream closed at event %d", i)
				Expect(string(event.ContainerStdOutWrittenTo.Data)).To(Equal(fmt.Sprint(i)))
			case <-deadline:
				Fail(fmt.Sprintf("event %d did not arrive", i))
			}
		}
		Consistently(events).ShouldNot(Receive())
		cancel()
		Eventually(events).Should(BeClosed())
	})

	It("unregisters when canceled while replaying events for other roots", func() {
		db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		ps := New(db).(*pubSub)
		now := time.Now().UTC()
		for i := 0; i < 1024; i++ {
			Expect(ps.eventStore.Add(model.Event{
				Timestamp: now.Add(time.Duration(i)),
				ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{
					RootCallID: "other",
				},
			})).To(Succeed())
		}

		ctx, cancel := context.WithCancel(context.Background())
		replayCtx := &replayCheckpoint{Context: ctx, reached: make(chan struct{})}
		events, err := ps.Subscribe(replayCtx, model.EventFilter{Roots: []string{"root"}, Ordered: true})
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			cancel()
			<-events
		}()

		Eventually(replayCtx.reached).Should(BeClosed(), "replay never checked for cancellation while scanning other roots")
		cancel()
		Eventually(events).Should(BeClosed())
		Expect(replayCtx.checks).To(Equal(8), "replay continued after cancellation")
		ps.orderedMu.RLock()
		defer ps.orderedMu.RUnlock()
		Expect(ps.ordered).To(BeEmpty())
	})

	It("zeroes delivered batch slots while later events are pending", func() {
		db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()

		ps := New(db).(*pubSub)
		now := time.Now().UTC()
		publish := func(i int) model.Event {
			event := model.Event{
				Timestamp: now.Add(time.Duration(i)),
				ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{
					Data: []byte{byte(i)}, RootCallID: "root",
				},
			}
			ps.Publish(event)
			return event
		}
		replay := publish(0)
		ctx, cancel := context.WithCancel(context.Background())
		events, err := ps.Subscribe(ctx, model.EventFilter{Roots: []string{"root"}, Ordered: true})
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			cancel()
			<-events
		}()

		// With replay blocked, the live events all queue in a single batch.
		first, second, third := publish(1), publish(2), publish(3)
		publish(4)
		var batch []model.Event
		ps.orderedMu.RLock()
		for sub := range ps.ordered {
			batch = sub.events
		}
		ps.orderedMu.RUnlock()
		Expect(batch).To(HaveLen(4))

		for _, event := range []model.Event{replay, first, second, third} {
			Eventually(events).Should(Receive(Equal(event)))
		}
		// Receiving the third live event synchronizes with clearing the first two.
		Expect(batch[:2]).To(Equal(make([]model.Event, 2)))
	})
})

type replayCheckpoint struct {
	context.Context
	checks  int
	reached chan struct{}
}

func (ctx *replayCheckpoint) Err() error {
	ctx.checks++
	if ctx.checks == 8 {
		close(ctx.reached)
		<-ctx.Done()
	}
	return ctx.Context.Err()
}
