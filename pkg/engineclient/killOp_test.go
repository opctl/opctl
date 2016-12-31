package engineclient

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/pkg/model"
  "github.com/opspec-io/sdk-golang/util/http"
  netHttp "net/http"
  "github.com/opspec-io/sdk-golang/pkg/engineprovider/providers/fake"
  "fmt"
  "bytes"
  "github.com/opspec-io/sdk-golang/util/format"
)

var _ = Describe("KillOp", func() {

  It("should call httpClient.Do() with expected args", func() {

    /* arrange */
    providedKillOpReq := model.KillOpReq{
      OpId:"dummyOpId",
    }

    expectedBytes, _ := format.NewJsonFormat().From(providedKillOpReq)
    engineProtocolRelativeBaseUrl := "dummyEngineProtocolBaseUrl"

    expectedHttpReq, _ := netHttp.NewRequest(
      "POST",
      fmt.Sprintf("http:%v/instances/kills", engineProtocolRelativeBaseUrl),
      bytes.NewBuffer(expectedBytes),
    )

    fakeEngineProvider := new(fake.EngineProvider)
    fakeEngineProvider.GetEngineProtocolRelativeBaseUrlReturns(engineProtocolRelativeBaseUrl, nil)

    fakeHttpClient := new(http.FakeClient)

    objectUnderTest := _engineClient{
      httpClient : fakeHttpClient,
      engineProvider: fakeEngineProvider,
      jsonFormat:format.NewJsonFormat(),
    }

    /* act */
    objectUnderTest.KillOp(providedKillOpReq)

    /* assert */
    Expect(expectedHttpReq).To(Equal(fakeHttpClient.DoArgsForCall(0)))

  })
})
