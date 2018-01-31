package node

import (
	"context"
	"fmt"

	"bufio"
	"github.com/libp2p/go-libp2p"
	libP2PNet "github.com/libp2p/go-libp2p-net"
	"github.com/opctl/opctl/util/clicolorer"
	"net/http"
	"strings"
)

/**
newLibP2PListener returns a listener which exposes opctl over libP2P
*/
func newLibP2PListener() listener {
	return _libP2PListener{
		cliColorer: clicolorer.New(),
		listenAddr: "/ip4/0.0.0.0/tcp/42225/ws",
	}
}

type _libP2PListener struct {
	cliColorer clicolorer.CliColorer
	listenAddr string
}

func (hf _libP2PListener) Listen(
	ctx context.Context,
) <-chan error {
	errChan := make(chan error, 1)

	h, err := libp2p.New(
		ctx,
		libp2p.ListenAddrStrings(hf.listenAddr),
	)
	if nil != err {
		errChan <- err
		close(errChan)
		return errChan
	}

	h.SetStreamHandler("/opspec/0.1.5", func(stream libP2PNet.Stream) {
		defer stream.Close()

		// Create a new buffered reader, as ReadRequest needs one.
		// The buffered reader reads from our stream, on which we
		// have sent the HTTP request (see ServeHTTP())
		buf := bufio.NewReader(stream)
		// Read the HTTP request from the buffer
		req, err := http.ReadRequest(buf)
		if err != nil {
			stream.Reset()
			return
		}
		defer req.Body.Close()

		// We need to reset these fields in the request
		// URL as they are not maintained.
		req.URL.Scheme = "http"
		hp := strings.Split(req.Host, ":")
		if len(hp) > 1 && hp[1] == "443" {
			req.URL.Scheme = "https"
		} else {
			req.URL.Scheme = "http"
		}
		req.URL.Host = req.Host

		outreq := new(http.Request)
		*outreq = *req

		// We now make the request
		fmt.Printf("Making request to %s\n", req.URL)
		resp, err := http.DefaultTransport.RoundTrip(outreq)
		if err != nil {
			stream.Reset()
			return
		}

		// resp.Write writes whatever response we obtained for our
		// request back to the stream.
		resp.Write(stream)

	})

	fmt.Println(
		hf.cliColorer.Info(
			"libP2PListener bound to %v w/ peer id %s",
			hf.listenAddr,
			h.ID().Pretty(),
		),
	)
	return errChan
}
