package starts

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/node/api"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
)

var _ = Context("Handler", func() {
	Context("NewHandler", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewHandler(new(nodeFakes.FakeNode))).Should(Not(BeNil()))
		})
	})
	Context("Handle", func() {
		Context("json.Decoder.Decode errors", func() {
			It("should return StatusCode of 400", func() {

				/* arrange */
				objectUnderTest := _handler{
					node: new(nodeFakes.FakeNode),
				}
				providedHTTPResp := httptest.NewRecorder()

				providedHTTPReq, err := http.NewRequest(http.MethodPost, api.URLOps_Starts, bytes.NewReader([]byte{}))
				if err != nil {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.Handle(providedHTTPResp, providedHTTPReq)

				/* assert */
				Expect(providedHTTPResp.Code).To(Equal(http.StatusBadRequest))

			})
		})
	})
})
