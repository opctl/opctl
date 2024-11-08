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

	dnsServer := miekgdns.Server{
		Addr:    address,
		Handler: newHandler(),
		Net:     "udp",
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// little hammer
		dnsServer.ShutdownContext(ctx)
	}()

	return dnsServer.ListenAndServe()
}
