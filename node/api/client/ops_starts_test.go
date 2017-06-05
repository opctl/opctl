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

var _ = Describe("StartOp", func() {

	It("should call httpClient.Do() with expected args", func() {

		/* arrange */
		providedCtx := context.TODO()
		providedReq := model.StartOpReq{
			Args: map[string]*model.Value{},
			Pkg: &model.DCGOpCallPkg{
				Ref: "dummyPkgRef",
				PullCreds: &model.DCGPullCreds{
					Username: "dummyUsername",
					Password: "dummyPassword",
				},
			},
		}

		expectedReqUrl := url.URL{}
		expectedReqUrl.Path = api.Ops_StartsURLTpl

		expectedReqBytes, _ := json.Marshal(providedReq)
		expectedResult := "dummyOpId"

		expectedHttpReq, _ := http.NewRequest(
			"POST",
			expectedReqUrl.String(),
			bytes.NewBuffer(expectedReqBytes),
		)
		expectedHttpReq.WithContext(providedCtx)

		fakeHttpClient := new(ihttp.Fake)
		fakeHttpClient.DoReturns(&http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte(expectedResult)))}, nil)

		objectUnderTest := client{
			httpClient: fakeHttpClient,
		}

		/* act */
		actualResult, _ := objectUnderTest.StartOp(providedCtx, providedReq)

		/* assert */
		Expect(expectedHttpReq).To(Equal(fakeHttpClient.DoArgsForCall(0)))
		Expect(expectedResult).To(Equal(actualResult))

	})
})
