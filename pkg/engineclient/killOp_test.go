package engineclient

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/engineprovider/providers/fake"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/http"
	netHttp "net/http"
)

var _ = Describe("KillOp", func() {

	It("should call httpClient.Do() with expected args", func() {

		/* arrange */
		providedKillOpReq := model.KillOpReq{
			OpGraphId: "dummyOpGraphId",
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
			httpClient:     fakeHttpClient,
			engineProvider: fakeEngineProvider,
			jsonFormat:     format.NewJsonFormat(),
		}

		/* act */
		objectUnderTest.KillOp(providedKillOpReq)

		/* assert */
		Expect(expectedHttpReq).To(Equal(fakeHttpClient.DoArgsForCall(0)))

	})
})
