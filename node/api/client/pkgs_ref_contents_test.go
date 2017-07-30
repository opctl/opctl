package client

import (
	"bytes"
	"context"
	"github.com/golang-interfaces/ihttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var _ = Context("ListPkgContents", func() {

	It("should call httpClient.Do() w/ expected args & return result", func() {

		/* arrange */
		providedCtx := context.TODO()
		providedReq := model.ListPkgContentsReq{
			PkgRef: "dummyPkgRef",
			PullCreds: &model.PullCreds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			},
		}

		expectedReqUrl := url.URL{}
		path := strings.Replace(api.URLPkgs_Ref_Contents, "{ref}", url.PathEscape(providedReq.PkgRef), 1)
		expectedReqUrl.Path = path

		httpResp := &http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte("[]")))}

		expectedHttpReq, _ := http.NewRequest(
			"GET",
			expectedReqUrl.String(),
			nil,
		)

		fakeHttpClient := new(ihttp.FakeClient)
		fakeHttpClient.DoReturns(httpResp, nil)

		objectUnderTest := client{
			httpClient: fakeHttpClient,
		}

		/* act */
		actualContentsList, _ := objectUnderTest.ListPkgContents(providedCtx, providedReq)

		/* assert */
		actualHttpReq := fakeHttpClient.DoArgsForCall(0)

		Expect(actualHttpReq.URL).To(Equal(expectedHttpReq.URL))
		Expect(actualHttpReq.Body).To(BeNil())
		Expect(actualHttpReq.Header).To(Equal(expectedHttpReq.Header))
		Expect(actualHttpReq.Context()).To(Equal(providedCtx))

		Expect(actualContentsList).To(BeNil())

	})
})
