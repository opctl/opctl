package test_all

import (
	"testing"

	multiplex "github.com/whyrusleeping/go-smux-multiplex"
	multistream "github.com/whyrusleeping/go-smux-multistream"
	muxado "github.com/whyrusleeping/go-smux-muxado"
	spdy "github.com/whyrusleeping/go-smux-spdystream"
	yamux "github.com/whyrusleeping/go-smux-yamux"

	ttest "github.com/libp2p/go-stream-muxer/test"
)

func TestYamux(t *testing.T) {
	ttest.SubtestAll(t, yamux.DefaultTransport)
}

func TestMultistream(t *testing.T) {
	tpt := multistream.NewBlankTransport()
	tpt.AddTransport("/yamux", yamux.DefaultTransport)
	ttest.SubtestAll(t, tpt)
}

func TestMuxado(t *testing.T) {
	ttest.SubtestAll(t, muxado.Transport)
}

func TestMultiplex(t *testing.T) {
	ttest.SubtestAll(t, multiplex.DefaultTransport)
}

func TestSpdystream(t *testing.T) {
	ttest.SubtestAll(t, spdy.Transport)
}
