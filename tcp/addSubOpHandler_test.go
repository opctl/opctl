package tcp

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "net/http/httptest"
  "net/http"
  "strings"
  "github.com/chrisdostert/mux"
  "github.com/opctl/engine/core/models"
  "bytes"
  "encoding/json"
)

var _ = Describe("addSubOpHandler", func() {

  Context("ServeHTTP() method", func() {
    It("should return StatusCode of 400 if projectUrl is malformed in Request", func() {

      /* arrange */
      objectUnderTest := addSubOpHandler{}
      recorder := httptest.NewRecorder()
      m := mux.NewRouter()
      m.Handle(addSubOpRelUrlTemplate, objectUnderTest)

      providedProjectUrl := "%%invalidProjectUrl%%"
      providedAddSubOpReqJson, err := json.Marshal(models.AddSubOpReq{})
      if (nil != err) {
        Fail(err.Error())
      }

      httpReq, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(providedAddSubOpReqJson))
      if (nil != err) {
        Fail(err.Error())
      }

      // brute force a request with malformed projectUrl
      httpReq.URL.Path = strings.Replace(addSubOpRelUrlTemplate, "{projectUrl}", providedProjectUrl, 1)

      /* act */
      m.ServeHTTP(recorder, httpReq)

      /* assert */
      Expect(recorder.Code).To(Equal(http.StatusBadRequest))

    })
  })
})
