package dns

import (
	"context"
	"net"
	"time"

	miekgdns "github.com/miekg/dns"
)

var (
	nsIPAddress = ""
	nsPort      = ""
)

/*
*
Listen for DNS requests
*/
func Listen(
	ctx context.Context,
	address string,
) error {

	var err error
	nsIPAddress, nsPort, err = net.SplitHostPort(address)
	if err != nil {
		return err
	}

	// snapshot upstreams before RegisterName gets a chance to modify resolv.conf
	handler := newHandler(
		address,
		snapshotUpstreams(
			resolvConfPath,
			nsIPAddress,
			nsPort,
		),
	)

	reusePort, err := shouldReusePort(
		ctx,
		nsPort,
	)
	if err != nil {
		return err
	}

	udpServer := miekgdns.Server{
		Addr:      address,
		Handler:   handler,
		Net:       "udp",
		ReusePort: reusePort,
	}

	// also serve TCP so clients can retry truncated responses
	tcpServer := miekgdns.Server{
		Addr:      address,
		Handler:   handler,
		Net:       "tcp",
		ReusePort: reusePort,
	}

	errChan := make(chan error, 2)
	go func() {
		errChan <- udpServer.ListenAndServe()
	}()
	go func() {
		errChan <- tcpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		err = nil
	case err = <-errChan:
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// little hammer
	udpServer.ShutdownContext(shutdownCtx)
	tcpServer.ShutdownContext(shutdownCtx)

	return err
}
