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

var _ = Context("GetPkgContent", func() {

	It("should call httpClient.Do() with expected args & return result", func() {

		/* arrange */
		providedCtx := context.TODO()
		providedReq := model.GetPkgContentReq{
			ContentPath: "dummy/content/path",
			PkgRef:      "dummyPkgRef",
			PullCreds: &model.PullCreds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			},
		}

		expectedReqUrl := url.URL{}
		path := strings.Replace(api.URLPkgs_Ref_Contents_Path, "{ref}", url.PathEscape(providedReq.PkgRef), 1)
		path = strings.Replace(path, "{path}", url.PathEscape(providedReq.ContentPath), 1)
		expectedReqUrl.Path = path

		pkgContent := "dummyPkgContent"
		expectedPkgContentReadCloser := ioutil.NopCloser(bytes.NewReader([]byte(pkgContent)))

		expectedHttpReq, _ := http.NewRequest(
			"GET",
			expectedReqUrl.String(),
			nil,
		)

		fakeHttpClient := new(ihttp.FakeClient)
		fakeHttpClient.DoReturns(&http.Response{Body: expectedPkgContentReadCloser}, nil)

		objectUnderTest := client{
			httpClient: fakeHttpClient,
		}

		/* act */
		actualPkgContentReadCloser, _ := objectUnderTest.GetPkgContent(providedCtx, providedReq)

		/* assert */
		actualHttpReq := fakeHttpClient.DoArgsForCall(0)

		Expect(actualHttpReq.URL).To(Equal(expectedHttpReq.URL))
		Expect(actualHttpReq.Body).To(BeNil())
		Expect(actualHttpReq.Header).To(Equal(expectedHttpReq.Header))
		Expect(actualHttpReq.Context()).To(Equal(providedCtx))

		Expect(actualPkgContentReadCloser).To(Equal(expectedPkgContentReadCloser))

	})
})
