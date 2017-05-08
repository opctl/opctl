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

var _ = Context("killOpHandler", func() {
	Context("ServeHTTP() method", func() {
		Describe("json.Decoder.Decode errors", func() {
			It("should return StatusCode of 400", func() {

				/* arrange */
				objectUnderTest := killOpHandler{}
				recorder := httptest.NewRecorder()
				m := mux.NewRouter()
				m.Handle(api.Ops_KillsURLTpl, objectUnderTest)

				httpReq, err := http.NewRequest(http.MethodPost, api.Ops_KillsURLTpl, bytes.NewReader([]byte{}))
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
