package docker

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

var dnsIPByHostname = map[string]map[string]interface{}{}

func StartDNSServer() error {
	// start server
	port := 53
	server := &dns.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: &handler{},
		Net:     "udp",
	}

	log.Printf("Starting at %d\n", port)

	err := server.ListenAndServe()
	defer server.Shutdown()

	return err
}

type handler struct{}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	log.Println("DNS handler called")
	msg := dns.Msg{}
	msg.SetReply(r)
	for _, q := range r.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			ips := dnsIPByHostname[strings.TrimSuffix(q.Name, ".")]
			fmt.Printf("Ips: %+v", ips)
			if len(ips) > 0 {
				// randomize returned ip for load balancing
				for ip := range ips {
					msg.Answer = append(
						msg.Answer,
						&dns.A{
							Hdr: dns.RR_Header{
								Name:   q.Name,
								Rrtype: dns.TypeA,
								Class:  dns.ClassINET,
								Ttl:    5,
							},
							A: net.ParseIP(ip),
						},
					)
				}
			} else {
				// refuse to answer (triggers clients to failover)
				msg.SetRcode(r, 5)
			}
		}
	}
	w.WriteMsg(&msg)
}
