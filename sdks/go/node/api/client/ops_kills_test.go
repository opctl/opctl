package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang-interfaces/ihttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/api"
	"io/ioutil"
	"net/http"
	"net/url"
)

var _ = Context("KillOp", func() {

	It("should call httpClient.Do() with expected args", func() {

		/* arrange */
		providedCtx := context.TODO()
		providedReq := model.KillOpReq{
			OpID: "dummyRootOpID",
		}

		expectedReqURL := url.URL{}
		expectedReqURL.Path = api.URLOps_Kills

		expectedBytes, _ := json.Marshal(providedReq)

		expectedHTTPReq, _ := http.NewRequest(
			"POST",
			expectedReqURL.String(),
			bytes.NewBuffer(expectedBytes),
		)

		fakeHttpClient := new(ihttp.FakeClient)
		fakeHttpClient.DoReturns(&http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte{}))}, nil)

		objectUnderTest := client{
			httpClient: fakeHttpClient,
		}

		/* act */
		objectUnderTest.KillOp(providedCtx, providedReq)

		/* assert */
		actualHTTPReq := fakeHttpClient.DoArgsForCall(0)

		Expect(actualHTTPReq.URL).To(Equal(expectedHTTPReq.URL))
		Expect(actualHTTPReq.Body).To(Equal(expectedHTTPReq.Body))
		Expect(actualHTTPReq.Header).To(Equal(expectedHTTPReq.Header))
		Expect(actualHTTPReq.Context()).To(Equal(providedCtx))

	})
})
