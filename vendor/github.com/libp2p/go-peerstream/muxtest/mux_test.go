package muxtest

import (
	"testing"

	multiplex "github.com/whyrusleeping/go-smux-multiplex"
	multistream "github.com/whyrusleeping/go-smux-multistream"
	yamux "github.com/whyrusleeping/go-smux-yamux"
)

func TestYamuxTransport(t *testing.T) {
	SubtestAll(t, yamux.DefaultTransport)
}

func TestMultiplexTransport(t *testing.T) {
	SubtestAll(t, multiplex.DefaultTransport)
}

func TestMultistreamTransport(t *testing.T) {
	tpt := multistream.NewBlankTransport()
	tpt.AddTransport("/yamux", yamux.DefaultTransport)
	tpt.AddTransport("/mplex", multiplex.DefaultTransport)
	SubtestAll(t, tpt)
}
