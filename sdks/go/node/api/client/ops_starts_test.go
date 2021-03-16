package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/golang-interfaces/ihttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api"
)

var _ = Context("StartOp", func() {

	It("should call httpClient.Do() with expected args", func() {

		/* arrange */
		providedCtx := context.TODO()
		providedReq := model.StartOpReq{
			Args: map[string]*model.Value{},
			Op: model.StartOpReqOp{
				Ref: "dummyOpRef",
				PullCreds: &model.Creds{
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

		fakeHTTPClient := new(ihttp.FakeClient)
		fakeHTTPClient.DoReturns(
			&http.Response{
				Body:       io.NopCloser(bytes.NewReader([]byte(expectedResult))),
				StatusCode: http.StatusCreated,
			},
			nil,
		)

		objectUnderTest := apiClient{
			httpClient: fakeHTTPClient,
		}

		/* act */
		actualResult, _ := objectUnderTest.StartOp(providedCtx, providedReq)

		/* assert */
		actualHTTPReq := fakeHTTPClient.DoArgsForCall(0)

		Expect(actualHTTPReq.URL).To(Equal(expectedHTTPReq.URL))
		Expect(actualHTTPReq.Body).To(Equal(expectedHTTPReq.Body))
		Expect(actualHTTPReq.Header).To(Equal(expectedHTTPReq.Header))
		Expect(actualHTTPReq.Context()).To(Equal(providedCtx))

		Expect(expectedResult).To(Equal(actualResult))

	})
})
