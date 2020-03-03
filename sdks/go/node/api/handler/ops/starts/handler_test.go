package starts

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/node/api"
	. "github.com/opctl/opctl/sdks/go/node/core/fakes"
	"net/http"
	"net/http/httptest"
)

var _ = Context("Handler", func() {
	Context("NewHandler", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewHandler(new(FakeCore))).Should(Not(BeNil()))
		})
	})
	Context("Handle", func() {
		Context("json.Decoder.Decode errors", func() {
			It("should return StatusCode of 400", func() {

				/* arrange */
				objectUnderTest := _handler{
					core: new(FakeCore),
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
