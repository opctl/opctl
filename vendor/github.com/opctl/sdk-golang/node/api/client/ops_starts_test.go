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

var _ = Context("StartOp", func() {

	It("should call httpClient.Do() with expected args", func() {

		/* arrange */
		providedCtx := context.TODO()
		providedReq := model.StartOpReq{
			Args: map[string]*model.Value{},
			Op: model.StartOpReqOp{
				Ref: "dummyOpRef",
				PullCreds: &model.PullCreds{
					Username: "dummyUsername",
					Password: "dummyPassword",
				},
			},
		}

		expectedReqURL := url.URL{}
		expectedReqURL.Path = api.URLOps_Starts

		expectedReqBytes, _ := json.Marshal(providedReq)
		expectedResult := "dummyOpID"

		expectedHTTPReq, _ := http.NewRequest(
			"POST",
			expectedReqURL.String(),
			bytes.NewBuffer(expectedReqBytes),
		)

		fakeHttpClient := new(ihttp.FakeClient)
		fakeHttpClient.DoReturns(&http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte(expectedResult)))}, nil)

		objectUnderTest := client{
			httpClient: fakeHttpClient,
		}

		/* act */
		actualResult, _ := objectUnderTest.StartOp(providedCtx, providedReq)

		/* assert */
		actualHTTPReq := fakeHttpClient.DoArgsForCall(0)

		Expect(actualHTTPReq.URL).To(Equal(expectedHTTPReq.URL))
		Expect(actualHTTPReq.Body).To(Equal(expectedHTTPReq.Body))
		Expect(actualHTTPReq.Header).To(Equal(expectedHTTPReq.Header))
		Expect(actualHTTPReq.Context()).To(Equal(providedCtx))

		Expect(expectedResult).To(Equal(actualResult))

	})
})
