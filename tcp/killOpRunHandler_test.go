package tcp

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "net/http/httptest"
  "net/http"
  "github.com/gorilla/mux"
  "bytes"
)

var _ = Describe("killOpRunHandler", func() {

  Context("ServeHTTP() method", func() {
    It("should return StatusCode of 400 if body of request is malformed", func() {

      /* arrange */
      objectUnderTest := killOpRunHandler{}
      recorder := httptest.NewRecorder()
      m := mux.NewRouter()
      m.Handle(killOpRunRelUrlTemplate, objectUnderTest)

      httpReq, err := http.NewRequest(http.MethodPost, killOpRunRelUrlTemplate, bytes.NewReader([]byte{}))
      if (nil != err) {
        Fail(err.Error())
      }

      /* act */
      m.ServeHTTP(recorder, httpReq)

      /* assert */
      Expect(recorder.Code).To(Equal(http.StatusBadRequest))

    })
  })
})
