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
	"io/ioutil"
	netHttp "net/http"
)

var _ = Describe("StartOp", func() {

	It("should call httpClient.Do() with expected args", func() {

		/* arrange */
		providedStartOpReq := model.StartOpReq{
			Args:  map[string]*model.Arg{},
			OpRef: "dummyOpRef",
		}

		expectedReqBytes, _ := format.NewJsonFormat().From(providedStartOpReq)
		engineProtocolRelativeBaseUrl := "dummyEngineProtocolBaseUrl"
		expectedResult := "dummyOpId"

		expectedHttpReq, _ := netHttp.NewRequest(
			"POST",
			fmt.Sprintf("http:%v/instances/starts", engineProtocolRelativeBaseUrl),
			bytes.NewBuffer(expectedReqBytes),
		)

		fakeEngineProvider := new(fake.EngineProvider)
		fakeEngineProvider.GetEngineProtocolRelativeBaseUrlReturns(engineProtocolRelativeBaseUrl, nil)

		fakeHttpClient := new(http.FakeClient)
		fakeHttpClient.DoReturns(&netHttp.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte(expectedResult)))}, nil)

		objectUnderTest := _engineClient{
			httpClient:     fakeHttpClient,
			engineProvider: fakeEngineProvider,
			jsonFormat:     format.NewJsonFormat(),
		}

		/* act */
		actualResult, _ := objectUnderTest.StartOp(providedStartOpReq)

		/* assert */
		Expect(expectedHttpReq).To(Equal(fakeHttpClient.DoArgsForCall(0)))
		Expect(expectedResult).To(Equal(actualResult))

	})
})
