package node

import (
	"sync"

	"github.com/opctl/opctl/sdks/go/model"
)

// newChildCallEndedWatcher drains run events while serial children execute synchronously.
func newChildCallEndedWatcher(eventChannel <-chan model.Event) func(string) <-chan *model.CallEnded {
	var mutex sync.Mutex
	pending := map[string]chan *model.CallEnded{}
	closed := false

	go func() {
		for event := range eventChannel {
			if event.CallEnded == nil {
				continue
			}
			mutex.Lock()
			if channel, ok := pending[event.CallEnded.Call.ID]; ok {
				channel <- event.CallEnded
				// Forget completed children so duplicate events cannot block the reader.
				delete(pending, event.CallEnded.Call.ID)
			}
			mutex.Unlock()
		}

		mutex.Lock()
		defer mutex.Unlock()
		closed = true
		for _, channel := range pending {
			close(channel)
		}
	}()

	return func(callID string) <-chan *model.CallEnded {
		channel := make(chan *model.CallEnded, 1)
		mutex.Lock()
		defer mutex.Unlock()
		if closed {
			close(channel)
		} else {
			pending[callID] = channel
		}
		return channel
	}
}
