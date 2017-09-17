package handler

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/node/api"
	"github.com/opspec-io/sdk-golang/node/core"
	"net/http"
	"net/http/httptest"
)

var _ = Context("POST /ops/kills", func() {
	Context("json.Decoder.Decode errors", func() {
		It("should return StatusCode of 400", func() {

			/* arrange */
			objectUnderTest := New(new(core.Fake))
			recorder := httptest.NewRecorder()

			httpReq, err := http.NewRequest(http.MethodPost, api.URLOps_Kills, bytes.NewReader([]byte{}))
			if nil != err {
				panic(err.Error())
			}

			/* act */
			objectUnderTest.ServeHTTP(recorder, httpReq)

			/* assert */
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))

		})
	})
})
