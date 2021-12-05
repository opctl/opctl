package docker

import (
	"sync"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

type mockPublisher struct {
	events []model.Event
	mu     sync.Mutex
	wg     sync.WaitGroup
}

func NewMockPublisher() *mockPublisher {
	return &mockPublisher{
		events: make([]model.Event, 0),
	}
}

func (p *mockPublisher) Publish(e model.Event) {
	defer p.wg.Done()
	p.mu.Lock()
	p.events = append(p.events, e)
	p.mu.Unlock()
}

func TestStdOutWriteCloser(t *testing.T) {
	g := NewGomegaWithT(t)

	/* arrange */
	eventPublisher := &mockPublisher{}

	objectUnderTest := NewStdOutWriteCloser(
		eventPublisher,
		"containerId",
		"rootCallId",
	)

	expectedEventLen := 2

	eventPublisher.wg.Add(expectedEventLen)

	/* act */
	_, err := objectUnderTest.Write([]byte("testing 1\ntesting 2"))
	if err != nil {
		panic(err)
	}
	objectUnderTest.Close()

	/* assert */
	eventPublisher.wg.Wait()
	eventPublisher.wg.Wait()
	g.Expect(eventPublisher.events).To(HaveLen(2))
}
