package dns

import (
	"net"
	"strings"
	"sync"
	"time"

	miekgDNS "github.com/miekg/dns"
)

var (
	ipsByHostname   = map[string]map[string]interface{}{}
	ipsByHostnameMu sync.RWMutex
)

// upstreamTimeout bounds each upstream exchange so dead upstreams can't stall clients
const upstreamTimeout = 2 * time.Second

/*
newHandler returns an opctl dns.Handler which authoritatively answers queries for
names registered via RegisterName and forwards all other queries to upstreams.

ownAddress is excluded from upstreams to prevent forwarding loops.
*/
func newHandler(
	ownAddress string,
	upstreams []string,
) miekgDNS.Handler {
	filteredUpstreams := []string{}
	for _, upstream := range upstreams {
		if upstream != ownAddress {
			filteredUpstreams = append(filteredUpstreams, upstream)
		}
	}

	return &handler{
		upstreams: filteredUpstreams,
	}
}

type handler struct {
	// immutable after construction
	upstreams []string
}

func (hdlr *handler) ServeDNS(w miekgDNS.ResponseWriter, r *miekgDNS.Msg) {
	if len(r.Question) == 0 {
		msg := miekgDNS.Msg{}
		msg.SetReply(r)
		w.WriteMsg(&msg)
		return
	}

	q := r.Question[0]

	ipsByHostnameMu.RLock()
	ips, isKnownName := ipsByHostname[strings.TrimSuffix(q.Name, ".")]
	ipList := make([]string, 0, len(ips))
	// randomize returned ip for load balancing
	for ip := range ips {
		ipList = append(ipList, ip)
	}
	ipsByHostnameMu.RUnlock()

	if isKnownName {
		msg := miekgDNS.Msg{}
		msg.SetReply(r)
		msg.Authoritative = true
		msg.RecursionAvailable = true
		if q.Qtype == miekgDNS.TypeA {
			for _, ip := range ipList {
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
		}
		w.WriteMsg(&msg)
		return
	}

	if resp := hdlr.forward(r, w.RemoteAddr().Network()); resp != nil {
		resp.Id = r.Id
		w.WriteMsg(resp)
		return
	}

	// not our zone & no upstream answered; refuse so clients look elsewhere
	msg := miekgDNS.Msg{}
	msg.SetRcode(r, miekgDNS.RcodeRefused)
	w.WriteMsg(&msg)
}

/*
forward relays r to the first responsive upstream and returns its response
verbatim, or nil if no upstream responded.
*/
func (hdlr *handler) forward(
	r *miekgDNS.Msg,
	clientNet string,
) *miekgDNS.Msg {
	proto := "udp"
	if strings.HasPrefix(clientNet, "tcp") {
		proto = "tcp"
	}

	for _, upstream := range hdlr.upstreams {
		resp, err := exchange(r, upstream, proto)
		if err != nil {
			continue
		}

		if proto == "udp" && resp.Truncated {
			if tcpResp, err := exchange(r, upstream, "tcp"); err == nil {
				resp = tcpResp
			}
		}

		return resp
	}

	return nil
}

func exchange(
	r *miekgDNS.Msg,
	upstream,
	proto string,
) (*miekgDNS.Msg, error) {
	client := miekgDNS.Client{
		Net:     proto,
		Timeout: upstreamTimeout,
	}

	resp, _, err := client.Exchange(r, upstream)
	return resp, err
}
