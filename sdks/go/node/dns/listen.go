package dns

import (
	"context"
	"net"
	"runtime"
	"time"

	miekgdns "github.com/miekg/dns"
)

var (
	nsIPAddress = ""
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
	// var nsPort string
	nsIPAddress, _, err = net.SplitHostPort(address)
	if err != nil {
		return err
	}

	dnsServer := miekgdns.Server{
		Addr:    address,
		Handler: newHandler(),
		Net:     "udp",
	}

	if runtime.GOOS == "darwin" {
		// on Mac OS, mDNSResponder binds to *:53. This is a workaround
		dnsServer.ReusePort = true
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
