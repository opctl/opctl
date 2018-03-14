package libp2p

import (
	"context"
	"fmt"
	"testing"

	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
)

func TestNewHost(t *testing.T) {
	_, err := makeRandomHost(t, 9000)
	if err != nil {
		t.Fatal(err)
	}
}

func makeRandomHost(t *testing.T, port int) (host.Host, error) {
	ctx := context.Background()
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		t.Fatal(err)
	}

	opts := []Option{
		ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)),
		Identity(priv),
		Muxer(DefaultMuxer()),
		NATPortMap(),
	}

	return New(ctx, opts...)
}
