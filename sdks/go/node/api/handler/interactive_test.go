package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/api/client"
	"github.com/opctl/opctl/sdks/go/node/fakes"
)

type testContainerInput struct {
	data      chan string
	sizes     chan model.TermSize
	done      chan struct{}
	eof       chan struct{}
	writing   chan struct{}
	writeOnce sync.Once
}

func (input *testContainerInput) Write(data []byte) (int, error) {
	if input.writing != nil {
		input.writeOnce.Do(func() { close(input.writing) })
	}
	select {
	case input.data <- string(data):
		return len(data), nil
	case <-input.eof:
		return 0, io.ErrClosedPipe
	}
}

func (input *testContainerInput) Resize(rows, cols uint) error {
	input.sizes <- model.TermSize{Rows: rows, Cols: cols}
	return nil
}

func (input *testContainerInput) Done() <-chan struct{} { return input.done }

func (input *testContainerInput) Close() error {
	close(input.eof)
	return nil
}

var _ = Describe("Interactive containers", func() {
	It("round-trips Interactive and initial terminal size through the API client and starts handler", func() {
		fakeCore := new(fakes.FakeCore)
		fakeCore.StartOpReturns("root", nil)
		server := httptest.NewServer(http.StripPrefix("/api/", New(fakeCore)))
		defer server.Close()
		serverURL, err := url.Parse(server.URL + "/api")
		Expect(err).NotTo(HaveOccurred())
		apiClient := client.New(*serverURL, nil)

		size := &model.TermSize{Rows: 33, Cols: 95}
		id, err := apiClient.StartOp(context.Background(), model.StartOpReq{Interactive: true, TermSize: size})
		Expect(err).NotTo(HaveOccurred())
		Expect(id).To(Equal("root"))
		_, req := fakeCore.StartOpArgsForCall(0)
		Expect(req.Interactive).To(BeTrue())
		Expect(req.TermSize).To(Equal(size))
	})

	It("forwards binary stdin and text resize, then closes stdin on disconnect", func() {
		input := &testContainerInput{
			data: make(chan string, 1), sizes: make(chan model.TermSize, 1),
			done: make(chan struct{}), eof: make(chan struct{}),
		}
		fakeCore := new(fakes.FakeCore)
		fakeCore.ClaimContainerInputReturns(input, func() {}, nil)
		server := httptest.NewServer(http.StripPrefix("/api/", New(fakeCore)))
		defer server.Close()
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/containers/container-id/attach"
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		Expect(err).NotTo(HaveOccurred())
		defer conn.Close()

		Expect(conn.WriteMessage(websocket.BinaryMessage, []byte("hello\r"))).To(Succeed())
		Expect(conn.WriteMessage(websocket.TextMessage, []byte(`{"rows":34,"cols":112}`))).To(Succeed())
		Eventually(input.data).Should(Receive(Equal("hello\r")))
		Eventually(input.sizes).Should(Receive(Equal(model.TermSize{Rows: 34, Cols: 112})))
		_, id := fakeCore.ClaimContainerInputArgsForCall(0)
		Expect(id).To(Equal("container-id"))

		Expect(conn.Close()).To(Succeed())
		Eventually(input.eof).Should(BeClosed())
	})

	It("closes the websocket when the container ends, even with a blocked stdin write", func() {
		input := &testContainerInput{
			data: make(chan string), done: make(chan struct{}), eof: make(chan struct{}),
			writing: make(chan struct{}),
		}
		fakeCore := new(fakes.FakeCore)
		fakeCore.ClaimContainerInputReturns(input, func() {}, nil)
		server := httptest.NewServer(http.StripPrefix("/api/", New(fakeCore)))
		defer server.Close()
		serverURL, err := url.Parse(server.URL + "/api")
		Expect(err).NotTo(HaveOccurred())

		conn, err := client.New(*serverURL, nil).AttachContainer(context.Background(), "container-id")
		Expect(err).NotTo(HaveOccurred())
		defer conn.Close()
		_, err = conn.Write([]byte("blocked"))
		Expect(err).NotTo(HaveOccurred())
		Eventually(input.writing).Should(BeClosed())

		close(input.done)
		Eventually(conn.Done()).Should(BeClosed())
		Eventually(input.eof).Should(BeClosed())
	})

	for _, scenario := range []struct {
		name             string
		messages         int
		closeAfterWrites bool
	}{
		{"writes accepted input before an orderly close", 5, true},
		{"keeps a slow input attachment open during a burst", 24, false},
	} {
		It(scenario.name, func() {
			input := &testContainerInput{
				data: make(chan string), done: make(chan struct{}), eof: make(chan struct{}),
				writing: make(chan struct{}),
			}
			fakeCore := new(fakes.FakeCore)
			fakeCore.ClaimContainerInputReturns(input, func() {}, nil)
			server := httptest.NewServer(http.StripPrefix("/api/", New(fakeCore)))
			defer server.Close()
			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/containers/container-id/attach"
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			Expect(err).NotTo(HaveOccurred())
			defer conn.Close()

			Expect(conn.WriteMessage(websocket.BinaryMessage, []byte("0"))).To(Succeed())
			Eventually(input.writing).Should(BeClosed())
			for i := 1; i < scenario.messages; i++ {
				Expect(conn.WriteMessage(websocket.BinaryMessage, []byte(fmt.Sprint(i)))).To(Succeed())
			}
			if scenario.closeAfterWrites {
				Expect(conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))).To(Succeed())
			}
			Consistently(input.eof, 50*time.Millisecond).ShouldNot(BeClosed())

			for i := 0; i < scenario.messages; i++ {
				Eventually(input.data).Should(Receive(Equal(fmt.Sprint(i))))
			}
			if !scenario.closeAfterWrites {
				Consistently(input.eof, 50*time.Millisecond).ShouldNot(BeClosed())
				Expect(conn.Close()).To(Succeed())
			}
			Eventually(input.eof).Should(BeClosed())
		})
	}

	It("splits a large SDK write into accepted stdin messages in order", func() {
		input := &testContainerInput{
			data: make(chan string, 3), done: make(chan struct{}), eof: make(chan struct{}),
		}
		fakeCore := new(fakes.FakeCore)
		fakeCore.ClaimContainerInputReturns(input, func() {}, nil)
		server := httptest.NewServer(http.StripPrefix("/api/", New(fakeCore)))
		defer server.Close()
		serverURL, err := url.Parse(server.URL + "/api")
		Expect(err).NotTo(HaveOccurred())
		conn, err := client.New(*serverURL, nil).AttachContainer(context.Background(), "container-id")
		Expect(err).NotTo(HaveOccurred())
		defer conn.Close()

		first := strings.Repeat("a", 64*1024)
		second := strings.Repeat("b", 64*1024)
		payload := first + second + "c"
		n, err := conn.Write([]byte(payload))
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(len(payload)))
		for _, want := range []string{first, second, "c"} {
			var got string
			Eventually(input.data).Should(Receive(&got))
			Expect(len(got) == len(want) && got == want).To(BeTrue())
		}
	})

	It("rejects an oversized stdin message at the websocket", func() {
		input := &testContainerInput{
			data: make(chan string, 1), done: make(chan struct{}), eof: make(chan struct{}),
		}
		fakeCore := new(fakes.FakeCore)
		fakeCore.ClaimContainerInputReturns(input, func() {}, nil)
		server := httptest.NewServer(http.StripPrefix("/api/", New(fakeCore)))
		defer server.Close()
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/containers/container-id/attach"
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		Expect(err).NotTo(HaveOccurred())
		defer conn.Close()

		Expect(conn.SetReadDeadline(time.Now().Add(time.Second))).To(Succeed())
		// The reader can send its close frame before the client finishes writing.
		_ = conn.WriteMessage(websocket.BinaryMessage, make([]byte, 64*1024+1))
		_, _, err = conn.ReadMessage()
		Expect(websocket.IsCloseError(err, websocket.CloseMessageTooBig)).To(BeTrue(), "%v", err)
		Eventually(input.eof).Should(BeClosed())
		Expect(input.data).NotTo(Receive())
	})

	It("releases the claim when websocket upgrade fails", func() {
		input := &testContainerInput{done: make(chan struct{}), eof: make(chan struct{})}
		released := make(chan struct{})
		fakeCore := new(fakes.FakeCore)
		fakeCore.ClaimContainerInputReturns(input, func() { close(released) }, nil)
		server := httptest.NewServer(http.StripPrefix("/api/", New(fakeCore)))
		defer server.Close()

		resp, err := server.Client().Get(server.URL + "/api/containers/container-id/attach")
		Expect(err).NotTo(HaveOccurred())
		resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		Eventually(released).Should(BeClosed())
		Expect(input.eof).NotTo(BeClosed())
	})

	It("maps missing and already attached sessions to 404 and 409", func() {
		fakeCore := new(fakes.FakeCore)
		fakeCore.ClaimContainerInputReturnsOnCall(0, nil, nil, node.ErrContainerInputNotFound)
		fakeCore.ClaimContainerInputReturnsOnCall(1, nil, nil, node.ErrContainerInputAttached)
		server := httptest.NewServer(http.StripPrefix("/api/", New(fakeCore)))
		defer server.Close()
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/containers/container-id/attach"
		dialer := websocket.Dialer{HandshakeTimeout: time.Second}

		for _, status := range []int{http.StatusNotFound, http.StatusConflict} {
			_, response, err := dialer.Dial(wsURL, nil)
			Expect(err).To(HaveOccurred())
			Expect(response.StatusCode).To(Equal(status))
			response.Body.Close()
		}
	})
})
