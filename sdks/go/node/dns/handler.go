package dns

import (
	"net"
	"strings"

	miekgDNS "github.com/miekg/dns"
)

var ipsByHostname = map[string]map[string]interface{}{}

// New returns an opctl dns.Handler
func newHandler() miekgDNS.Handler {
	return &handler{}
}

type handler struct{}

func (hdlr *handler) ServeDNS(w miekgDNS.ResponseWriter, r *miekgDNS.Msg) {
	msg := miekgDNS.Msg{}
	msg.SetReply(r)
	msg.RecursionAvailable = true
	for _, q := range r.Question {
		switch q.Qtype {
		case miekgDNS.TypeA:
			msg.Authoritative = true
			ips := ipsByHostname[strings.TrimSuffix(q.Name, ".")]

			// randomize returned ip for load balancing
			for ip := range ips {
				msg.Answer = append(
					msg.Answer,
					&miekgDNS.A{
						Hdr: miekgDNS.RR_Header{
							Name:   q.Name,
							Rrtype: miekgDNS.TypeA,
							Class:  miekgDNS.ClassINET,
							Ttl:    60,
						},
						A: net.ParseIP(ip),
					},
				)
			}
			if len(msg.Answer) == 0 {
				// refuse to answer (triggers clients to failover)
				msg.SetRcode(r, 2)
			}
		}
	}
	w.WriteMsg(&msg)
}
