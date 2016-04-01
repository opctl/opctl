package rest

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "net/http/httptest"
  "net/http"
  "strings"
  "github.com/chrisdostert/mux"
  "github.com/dev-op-spec/engine/core/models"
  "bytes"
  "encoding/json"
)

var _ = Describe("setDescriptionOfPipelineHandler", func() {

  Context("ServeHTTP() method", func() {
    It("should return StatusCode of 400 if projectUrl is malformed in Request", func() {

      /* arrange */
      objectUnderTest := setDescriptionOfPipelineHandler{}
      recorder := httptest.NewRecorder()
      m := mux.NewRouter()
      m.Handle(setDescriptionOfPipelineRelUrlTemplate, objectUnderTest)

      providedProjectUrl := "%%invalidProjectUrl%%"
      providedPipelineName := "validPipelineName"
      providedSetDescriptionOfPipelineReqJson, err := json.Marshal(models.SetDescriptionOfPipelineReq{})
      if (nil != err) {
        Fail(err.Error())
      }

      httpReq, err := http.NewRequest(http.MethodGet, "", bytes.NewReader(providedSetDescriptionOfPipelineReqJson))
      if (nil != err) {
        Fail(err.Error())
      }

      // brute force a request with malformed projectUrl
      httpReq.URL.Path = strings.Replace(setDescriptionOfPipelineRelUrlTemplate, "{projectUrl}", providedProjectUrl, 1)
      httpReq.URL.Path = strings.Replace(httpReq.URL.Path, "{pipelineName}", providedPipelineName, 1)

      /* act */
      m.ServeHTTP(recorder, httpReq)

      /* assert */
      Expect(recorder.Code).To(Equal(http.StatusBadRequest))

    })
  })
})