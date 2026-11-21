package stream

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/fakes"
)

func TestOrderedQuery(t *testing.T) {
	for _, test := range []struct {
		query   string
		ordered bool
	}{{"", false}, {"?ordered", true}, {"?ordered=true", true}, {"?ordered=false", false}} {
		t.Run(test.query, func(t *testing.T) {
			requests := make(chan bool, 1)
			core := new(fakes.FakeNode)
			core.GetEventStreamStub = func(_ context.Context, req *model.GetEventStreamReq) (<-chan model.Event, error) {
				requests <- req.Filter.Ordered
				events := make(chan model.Event)
				close(events)
				return events, nil
			}
			server := httptest.NewServer(http.HandlerFunc(NewHandler(core).Handle))
			defer server.Close()
			conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+test.query, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			select {
			case ordered := <-requests:
				if ordered != test.ordered {
					t.Fatalf("Ordered = %v, want %v", ordered, test.ordered)
				}
			case <-time.After(time.Second):
				t.Fatal("stream handler did not subscribe")
			}
		})
	}
}

func TestOrderedStreamRejectsClientData(t *testing.T) {
	for _, test := range []struct {
		name       string
		size, code int
	}{
		{"unexpected", 1, websocket.CloseUnsupportedData},
		{"oversized", 513, websocket.CloseMessageTooBig},
	} {
		t.Run(test.name, func(t *testing.T) {
			done := make(chan struct{})
			core := new(fakes.FakeNode)
			core.GetEventStreamStub = func(ctx context.Context, _ *model.GetEventStreamReq) (<-chan model.Event, error) {
				events := make(chan model.Event)
				go func() {
					<-ctx.Done()
					close(events)
					close(done)
				}()
				return events, nil
			}
			server := httptest.NewServer(http.HandlerFunc(NewHandler(core).Handle))
			defer server.Close()
			conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"?ordered", nil)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
				t.Fatal(err)
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, make([]byte, test.size)); err != nil {
				t.Fatal(err)
			}
			if _, _, err := conn.ReadMessage(); !websocket.IsCloseError(err, test.code) {
				t.Fatalf("wanted close code %d, got %v", test.code, err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("stream subscription remained open after websocket closed")
			}
		})
	}
}

func TestOrderedAckStreamEndsAfterFinalAckAndDisconnect(t *testing.T) {
	subscriptionDone := make(chan struct{})
	handlerDone := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	core := new(fakes.FakeNode)
	core.GetEventStreamStub = func(ctx context.Context, _ *model.GetEventStreamReq) (<-chan model.Event, error) {
		events := make(chan model.Event, 1)
		events <- model.Event{CallEnded: &model.CallEnded{Call: model.Call{RootID: "root"}}}
		go func() {
			<-ctx.Done()
			close(events)
			close(subscriptionDone)
		}()
		return events, nil
	}
	handler := NewHandler(core)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer close(handlerDone)
		handler.Handle(w, r.WithContext(ctx))
	}))
	defer server.Close()

	conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"?ordered&ack&roots=root", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	deadline := time.Now().Add(time.Second)
	if err := conn.SetReadDeadline(deadline); err != nil {
		t.Fatal(err)
	}
	if err := conn.SetWriteDeadline(deadline); err != nil {
		t.Fatal(err)
	}
	var event model.Event
	if err := conn.ReadJSON(&event); err != nil {
		t.Fatal(err)
	}
	if event.CallEnded == nil || event.CallEnded.Call.RootID != "root" {
		t.Fatalf("expected final event for root, got %+v", event)
	}
	if err := conn.WriteMessage(websocket.TextMessage, nil); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}

	timeout := time.After(time.Second)
	select {
	case <-handlerDone:
	case <-timeout:
		t.Fatal("stream handler remained open after final ACK and disconnect")
	}
	select {
	case <-subscriptionDone:
	case <-timeout:
		t.Fatal("stream subscription remained open after final ACK and disconnect")
	}
}
