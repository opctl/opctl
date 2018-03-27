package starts

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/node/api"
	"github.com/opspec-io/sdk-golang/node/core"
	"net/http"
	"net/http/httptest"
)

var _ = Context("Handler", func() {
	Context("Handle", func() {
		Context("json.Decoder.Decode errors", func() {
			It("should return StatusCode of 400", func() {

				/* arrange */
				objectUnderTest := _handler{
					core: new(core.Fake),
				}
				providedHTTPResp := httptest.NewRecorder()

				providedHTTPReq, err := http.NewRequest(http.MethodPost, api.URLOps_Starts, bytes.NewReader([]byte{}))
				if nil != err {
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
