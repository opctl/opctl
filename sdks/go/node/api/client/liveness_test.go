package client

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-interfaces/ihttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api"
)

var _ = Context("Liveness", func() {

	It("should call httpClient.Do() w/ expected args & return result", func() {

		/* arrange */
		providedCtx := context.TODO()

		expectedReqURL := url.URL{}
		expectedReqURL.Path = api.URLLiveness

		expectedHTTPReq, _ := http.NewRequest(
			"GET",
			expectedReqURL.String(),
			nil,
		)

		fakeHttpClient := new(ihttp.FakeClient)
		fakeHttpClient.DoReturns(
			&http.Response{
				Body:       ioutil.NopCloser(strings.NewReader("[]")),
				StatusCode: http.StatusOK,
			},
			nil,
		)

		objectUnderTest := APIClient{
			httpClient: fakeHttpClient,
		}

		/* act */
		objectUnderTest.Liveness(providedCtx)

		/* assert */
		actualHTTPReq := fakeHttpClient.DoArgsForCall(0)

		Expect(actualHTTPReq.URL).To(Equal(expectedHTTPReq.URL))
		Expect(actualHTTPReq.Body).To(BeNil())
		Expect(actualHTTPReq.Header).To(Equal(expectedHTTPReq.Header))
		Expect(actualHTTPReq.Context()).To(Equal(providedCtx))

	})
	Context("StatusCode == 200", func() {

		It("should return expected result", func() {

			/* arrange */
			httpResp := &http.Response{
				Body:       ioutil.NopCloser(strings.NewReader("[]")),
				StatusCode: http.StatusOK,
			}

			fakeHttpClient := new(ihttp.FakeClient)
			fakeHttpClient.DoReturns(httpResp, nil)

			objectUnderTest := APIClient{
				httpClient: fakeHttpClient,
			}

			/* act */
			actualErr := objectUnderTest.Liveness(
				context.TODO(),
			)

			/* assert */
			Expect(actualErr).To(BeNil())

		})
	})
	Context("StatusCode != 200", func() {
		It("should return expected result", func() {
			/* arrange */
			expectedErr := errors.New("dummyMsg")
			httpResp := &http.Response{
				Body:       ioutil.NopCloser(strings.NewReader(expectedErr.Error())),
				StatusCode: http.StatusInternalServerError,
			}

			fakeHttpClient := new(ihttp.FakeClient)
			fakeHttpClient.DoReturns(httpResp, nil)

			objectUnderTest := APIClient{
				httpClient: fakeHttpClient,
			}

			/* act */
			_, actualErr := objectUnderTest.ListDescendants(
				context.TODO(),
				model.ListDescendantsReq{},
			)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))

		})
	})
})
