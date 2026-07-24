package dns

import (
	"net"
	"sync"
	"time"

	miekgDNS "github.com/miekg/dns"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// startTestServer serves h over both UDP & TCP on a random loopback port
func startTestServer(h miekgDNS.Handler) (string, func()) {
	packetConn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	address := packetConn.LocalAddr().String()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	udpServer := &miekgDNS.Server{PacketConn: packetConn, Handler: h}
	tcpServer := &miekgDNS.Server{Listener: listener, Handler: h}
	go udpServer.ActivateAndServe()
	go tcpServer.ActivateAndServe()

	return address, func() {
		udpServer.Shutdown()
		tcpServer.Shutdown()
	}
}

type fakeUpstream struct {
	address string
	stop    func()

	mux     sync.Mutex
	queries []miekgDNS.Question
}

func (fu *fakeUpstream) queryCount() int {
	fu.mux.Lock()
	defer fu.mux.Unlock()
	return len(fu.queries)
}

// startFakeUpstream serves serve on a random loopback port, recording every query
func startFakeUpstream(serve func(w miekgDNS.ResponseWriter, r *miekgDNS.Msg)) *fakeUpstream {
	fu := &fakeUpstream{}
	fu.address, fu.stop = startTestServer(
		miekgDNS.HandlerFunc(func(w miekgDNS.ResponseWriter, r *miekgDNS.Msg) {
			fu.mux.Lock()
			fu.queries = append(fu.queries, r.Question...)
			fu.mux.Unlock()
			serve(w, r)
		}),
	)
	return fu
}

func aRecord(name, ip string) *miekgDNS.A {
	return &miekgDNS.A{
		Hdr: miekgDNS.RR_Header{
			Name:   name,
			Rrtype: miekgDNS.TypeA,
			Class:  miekgDNS.ClassINET,
			Ttl:    60,
		},
		A: net.ParseIP(ip),
	}
}

func query(serverAddress, name string, qtype uint16) *miekgDNS.Msg {
	q := miekgDNS.Msg{}
	q.SetQuestion(miekgDNS.Fqdn(name), qtype)

	client := miekgDNS.Client{Timeout: 5 * time.Second}
	resp, _, err := client.Exchange(&q, serverAddress)
	Expect(err).To(BeNil())

	return resp
}

var _ = Context("handler", func() {
	Context("name isn't registered", func() {
		It("should forward A queries & relay the upstream answer verbatim", func() {
			/* arrange */
			upstream := startFakeUpstream(func(w miekgDNS.ResponseWriter, r *miekgDNS.Msg) {
				msg := miekgDNS.Msg{}
				msg.SetReply(r)
				msg.RecursionAvailable = true
				msg.Answer = append(msg.Answer, aRecord(r.Question[0].Name, "93.184.216.34"))
				w.WriteMsg(&msg)
			})
			defer upstream.stop()

			serverAddress, stop := startTestServer(newHandler("", []string{upstream.address}))
			defer stop()

			/* act */
			resp := query(serverAddress, "example.com", miekgDNS.TypeA)

			/* assert */
			Expect(resp.Rcode).To(Equal(miekgDNS.RcodeSuccess))
			Expect(resp.Answer).To(HaveLen(1))
			Expect(resp.Answer[0].(*miekgDNS.A).A.String()).To(Equal("93.184.216.34"))
			Expect(upstream.queryCount()).To(Equal(1))
		})

		It("should relay upstream NXDOMAIN", func() {
			/* arrange */
			upstream := startFakeUpstream(func(w miekgDNS.ResponseWriter, r *miekgDNS.Msg) {
				msg := miekgDNS.Msg{}
				msg.SetRcode(r, miekgDNS.RcodeNameError)
				w.WriteMsg(&msg)
			})
			defer upstream.stop()

			serverAddress, stop := startTestServer(newHandler("", []string{upstream.address}))
			defer stop()

			/* act */
			resp := query(serverAddress, "doesnotexist.example.com", miekgDNS.TypeA)

			/* assert */
			Expect(resp.Rcode).To(Equal(miekgDNS.RcodeNameError))
		})

		It("should forward AAAA queries instead of answering empty NOERROR", func() {
			/* arrange */
			upstream := startFakeUpstream(func(w miekgDNS.ResponseWriter, r *miekgDNS.Msg) {
				msg := miekgDNS.Msg{}
				msg.SetReply(r)
				msg.Answer = append(
					msg.Answer,
					&miekgDNS.AAAA{
						Hdr: miekgDNS.RR_Header{
							Name:   r.Question[0].Name,
							Rrtype: miekgDNS.TypeAAAA,
							Class:  miekgDNS.ClassINET,
							Ttl:    60,
						},
						AAAA: net.ParseIP("2606:2800:220:1::1"),
					},
				)
				w.WriteMsg(&msg)
			})
			defer upstream.stop()

			serverAddress, stop := startTestServer(newHandler("", []string{upstream.address}))
			defer stop()

			/* act */
			resp := query(serverAddress, "example.com", miekgDNS.TypeAAAA)

			/* assert */
			Expect(resp.Rcode).To(Equal(miekgDNS.RcodeSuccess))
			Expect(resp.Answer).To(HaveLen(1))
			Expect(resp.Answer[0].(*miekgDNS.AAAA).AAAA.String()).To(Equal("2606:2800:220:1::1"))
			Expect(upstream.queryCount()).To(Equal(1))
		})

		It("should answer REFUSED when all upstreams are unreachable", func() {
			/* arrange */
			// bind then immediately close to obtain a dead upstream address
			deadConn, err := net.ListenPacket("udp", "127.0.0.1:0")
			if err != nil {
				panic(err)
			}
			deadAddress := deadConn.LocalAddr().String()
			deadConn.Close()

			serverAddress, stop := startTestServer(newHandler("", []string{deadAddress}))
			defer stop()

			/* act */
			resp := query(serverAddress, "example.com", miekgDNS.TypeA)

			/* assert */
			Expect(resp.Rcode).To(Equal(miekgDNS.RcodeRefused))
		})

		It("should answer REFUSED when the only upstream is itself", func() {
			/* arrange */
			packetConn, err := net.ListenPacket("udp", "127.0.0.1:0")
			if err != nil {
				panic(err)
			}
			serverAddress := packetConn.LocalAddr().String()

			server := &miekgDNS.Server{
				PacketConn: packetConn,
				Handler:    newHandler(serverAddress, []string{serverAddress}),
			}
			go server.ActivateAndServe()
			defer server.Shutdown()

			/* act */
			resp := query(serverAddress, "example.com", miekgDNS.TypeA)

			/* assert */
			Expect(resp.Rcode).To(Equal(miekgDNS.RcodeRefused))
		})

		It("should retry over TCP & relay the full answer when the upstream truncates over UDP", func() {
			/* arrange */
			upstream := startFakeUpstream(func(w miekgDNS.ResponseWriter, r *miekgDNS.Msg) {
				msg := miekgDNS.Msg{}
				msg.SetReply(r)
				if w.LocalAddr().Network() == "udp" {
					msg.Truncated = true
				} else {
					msg.Answer = append(msg.Answer, aRecord(r.Question[0].Name, "93.184.216.34"))
				}
				w.WriteMsg(&msg)
			})
			defer upstream.stop()

			serverAddress, stop := startTestServer(newHandler("", []string{upstream.address}))
			defer stop()

			/* act */
			resp := query(serverAddress, "example.com", miekgDNS.TypeA)

			/* assert */
			Expect(resp.Truncated).To(BeFalse())
			Expect(resp.Answer).To(HaveLen(1))
			Expect(resp.Answer[0].(*miekgDNS.A).A.String()).To(Equal("93.184.216.34"))
		})
	})

	Context("name is registered", func() {
		BeforeEach(func() {
			// RegisterName also applies host resolver config (requires root) so seed
			// the registry via the same internal path it uses
			registerIP("testcontainer", "172.17.0.2")
		})

		AfterEach(func() {
			ipsByHostnameMu.Lock()
			delete(ipsByHostname, "testcontainer")
			ipsByHostnameMu.Unlock()
		})

		It("should answer A queries locally without querying upstreams", func() {
			/* arrange */
			upstream := startFakeUpstream(func(w miekgDNS.ResponseWriter, r *miekgDNS.Msg) {
				msg := miekgDNS.Msg{}
				msg.SetReply(r)
				w.WriteMsg(&msg)
			})
			defer upstream.stop()

			serverAddress, stop := startTestServer(newHandler("", []string{upstream.address}))
			defer stop()

			/* act */
			resp := query(serverAddress, "testcontainer", miekgDNS.TypeA)

			/* assert */
			Expect(resp.Rcode).To(Equal(miekgDNS.RcodeSuccess))
			Expect(resp.Authoritative).To(BeTrue())
			Expect(resp.Answer).To(HaveLen(1))
			Expect(resp.Answer[0].(*miekgDNS.A).A.String()).To(Equal("172.17.0.2"))
			Expect(upstream.queryCount()).To(Equal(0))
		})

		It("should answer AAAA queries with empty NOERROR without querying upstreams", func() {
			/* arrange */
			upstream := startFakeUpstream(func(w miekgDNS.ResponseWriter, r *miekgDNS.Msg) {
				msg := miekgDNS.Msg{}
				msg.SetReply(r)
				w.WriteMsg(&msg)
			})
			defer upstream.stop()

			serverAddress, stop := startTestServer(newHandler("", []string{upstream.address}))
			defer stop()

			/* act */
			resp := query(serverAddress, "testcontainer", miekgDNS.TypeAAAA)

			/* assert */
			Expect(resp.Rcode).To(Equal(miekgDNS.RcodeSuccess))
			Expect(resp.Authoritative).To(BeTrue())
			Expect(resp.Answer).To(BeEmpty())
			Expect(upstream.queryCount()).To(Equal(0))
		})
	})
})
