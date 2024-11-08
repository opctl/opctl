package dns

import (
	"context"
	"runtime"
	"time"

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
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// little hammer
		dnsServer.ShutdownContext(ctx)
	}()

	if runtime.GOOS == "linux" {
		err := registerServer(address)
		if nil != err {
			return err
		}
	}

	return dnsServer.ListenAndServe()
}
