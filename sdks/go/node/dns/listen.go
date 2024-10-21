package dns

import (
	"context"
	"runtime"

	miekgdns "github.com/miekg/dns"
)

/*
*
Listen for DNS requests
*/
func Listen(
	ctx context.Context,
	address string,
) error {

	dnsServer := miekgdns.Server{
		Addr:    address,
		Handler: newHandler(),
		Net:     "udp",
	}

	go func() {
		<-ctx.Done()
		// little hammer
		dnsServer.Shutdown()
	}()

	if runtime.GOOS == "linux" {
		err := registerServer(address)
		if nil != err {
			return err
		}
	}

	return dnsServer.ListenAndServe()
}
