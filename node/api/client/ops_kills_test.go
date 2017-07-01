package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang-interfaces/ihttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"io/ioutil"
	"net/http"
	"net/url"
)

var _ = Context("KillOp", func() {

	It("should call httpClient.Do() with expected args", func() {

		/* arrange */
		providedCtx := context.TODO()
		providedReq := model.KillOpReq{
			OpId: "dummyRootOpId",
		}

		expectedReqUrl := url.URL{}
		expectedReqUrl.Path = api.URLOps_Kills

		expectedBytes, _ := json.Marshal(providedReq)

		expectedHttpReq, _ := http.NewRequest(
			"POST",
			expectedReqUrl.String(),
			bytes.NewBuffer(expectedBytes),
		)
		expectedHttpReq.WithContext(providedCtx)

		fakeHttpClient := new(ihttp.FakeClient)
		fakeHttpClient.DoReturns(&http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte{}))}, nil)

		objectUnderTest := client{
			httpClient: fakeHttpClient,
		}

		/* act */
		objectUnderTest.KillOp(providedCtx, providedReq)

		/* assert */
		Expect(expectedHttpReq).To(Equal(fakeHttpClient.DoArgsForCall(0)))

	})
})
