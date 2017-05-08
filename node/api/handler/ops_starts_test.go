package handler

import (
	"bytes"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/node/api"
	"net/http"
	"net/http/httptest"
)

var _ = Context("startOpHandler", func() {
	Context("ServeHTTP() method", func() {
		Describe("json.Decoder.Decode errors", func() {
			It("should return StatusCode of 400", func() {

				/* arrange */
				objectUnderTest := startOpHandler{}
				recorder := httptest.NewRecorder()
				m := mux.NewRouter()
				m.Handle(api.Ops_StartsURLTpl, objectUnderTest)

				httpReq, err := http.NewRequest(http.MethodPost, api.Ops_StartsURLTpl, bytes.NewReader([]byte{}))
				if nil != err {
					Fail(err.Error())
				}

				/* act */
				m.ServeHTTP(recorder, httpReq)

				/* assert */
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))

			})
		})
	})
})
