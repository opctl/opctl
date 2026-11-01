package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Describe("PlanOp client", func() {
	It("posts to the describe endpoint and reads the interface", func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Expect(r.Method).To(Equal("POST"))
			Expect(r.URL.Path).To(Equal("/api/ops/plans"))
			var req model.PlanOpReq
			Expect(json.NewDecoder(r.Body).Decode(&req)).To(Succeed())
			Expect(req.Op.Ref).To(Equal("/ops/hello"))
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"hello","inputs":{"env":{"string":{"default":"dev"}}}}`))
		}))
		defer server.Close()
		base, err := url.Parse(server.URL + "/api")
		Expect(err).NotTo(HaveOccurred())
		spec, err := New(*base, nil).PlanOp(context.Background(), model.PlanOpReq{Op: model.StartOpReqOp{Ref: "/ops/hello"}})
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Inputs["env"].String.Default).To(Equal("dev"))
	})
	It("returns API and malformed response errors", func() {
		for _, status := range []int{200, 400} {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				w.Write([]byte("bad declaration"))
			}))
			base, err := url.Parse(server.URL)
			Expect(err).NotTo(HaveOccurred())
			_, err = New(*base, nil).PlanOp(context.Background(), model.PlanOpReq{})
			Expect(err).To(HaveOccurred())
			server.Close()
		}
	})
})
