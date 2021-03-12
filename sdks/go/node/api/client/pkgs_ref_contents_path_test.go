package client

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-interfaces/ihttp"
	"github.com/jfbus/httprs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api"
)

var _ = Context("GetData", func() {

	It("should call httpClient.Do() with expected args & return result", func() {

		/* arrange */
		providedCtx := context.TODO()
		providedReq := model.GetDataReq{
			ContentPath: "dummy/content/path",
			PkgRef:      "dummyOpRef",
			PullCreds: &model.Creds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			},
		}

		path := strings.Replace(api.URLPkgs_Ref_Contents_Path, "{ref}", url.PathEscape(providedReq.PkgRef), 1)
		path = strings.Replace(path, "{path}", url.PathEscape(providedReq.ContentPath), 1)

		expectedReqURL, err := url.Parse(path)
		if nil != err {
			panic(err)
		}

		expectedHTTPReq, _ := http.NewRequest(
			"GET",
			expectedReqURL.String(),
			nil,
		)

		expectedHTTPReq.SetBasicAuth(
			providedReq.PullCreds.Username,
			providedReq.PullCreds.Password,
		)

		fakeHttpClient := new(ihttp.FakeClient)
		fakeHttpClient.DoReturns(
			&http.Response{
				Body:       ioutil.NopCloser(strings.NewReader("dummyBody")),
				StatusCode: http.StatusOK,
				Request:    expectedHTTPReq,
			},
			nil,
		)

		objectUnderTest := apiClient{
			httpClient: fakeHttpClient,
		}

		/* act */
		objectUnderTest.GetData(providedCtx, providedReq)

		/* assert */
		actualHTTPReq := fakeHttpClient.DoArgsForCall(0)

		Expect(actualHTTPReq.URL).To(Equal(expectedHTTPReq.URL))
		Expect(actualHTTPReq.Body).To(BeNil())
		Expect(actualHTTPReq.Header).To(Equal(expectedHTTPReq.Header))
		Expect(actualHTTPReq.Context()).To(Equal(providedCtx))

	})
	Context("StatusCode < 400", func() {

		It("should return expected result", func() {

			/* arrange */
			httpResp := &http.Response{
				Body:       ioutil.NopCloser(strings.NewReader("dummyBody")),
				StatusCode: http.StatusOK,
				Request:    &http.Request{},
			}

			expectedReadSeekCloser := httprs.NewHttpReadSeeker(httpResp)

			fakeHttpClient := new(ihttp.FakeClient)
			fakeHttpClient.DoReturns(httpResp, nil)

			objectUnderTest := apiClient{
				httpClient: fakeHttpClient,
			}

			/* act */
			actualReadSeekCloser, actualErr := objectUnderTest.GetData(
				context.TODO(),
				model.GetDataReq{},
			)

			/* assert */
			Expect(actualReadSeekCloser).To(Equal(expectedReadSeekCloser))
			Expect(actualErr).To(BeNil())

		})
	})
	Context("StatusCode >= 400", func() {
		Context("401", func() {
			It("should return expected result", func() {

				/* arrange */
				httpResp := &http.Response{
					Body:       ioutil.NopCloser(strings.NewReader("response")),
					StatusCode: http.StatusUnauthorized,
				}

				fakeHttpClient := new(ihttp.FakeClient)
				fakeHttpClient.DoReturns(httpResp, nil)

				objectUnderTest := apiClient{
					httpClient: fakeHttpClient,
				}

				/* act */
				_, actualErr := objectUnderTest.GetData(
					context.TODO(),
					model.GetDataReq{},
				)

				/* assert */
				Expect(actualErr).To(MatchError(model.ErrDataProviderAuthentication{}))
				Expect(actualErr.Error()).To(ContainSubstring("response"))

			})
		})
		Context("403", func() {
			It("should return expected result", func() {

				/* arrange */
				httpResp := &http.Response{
					Body:       ioutil.NopCloser(strings.NewReader("response")),
					StatusCode: http.StatusForbidden,
				}

				fakeHttpClient := new(ihttp.FakeClient)
				fakeHttpClient.DoReturns(httpResp, nil)

				objectUnderTest := apiClient{
					httpClient: fakeHttpClient,
				}

				/* act */
				_, actualErr := objectUnderTest.GetData(
					context.TODO(),
					model.GetDataReq{},
				)

				/* assert */
				Expect(actualErr).To(MatchError(model.ErrDataProviderAuthorization{}))
				Expect(actualErr.Error()).To(ContainSubstring("response"))

			})

		})
		Context("404", func() {
			It("should return expected result", func() {

				/* arrange */
				httpResp := &http.Response{
					Body:       ioutil.NopCloser(strings.NewReader("response")),
					StatusCode: http.StatusNotFound,
				}

				fakeHttpClient := new(ihttp.FakeClient)
				fakeHttpClient.DoReturns(httpResp, nil)

				objectUnderTest := apiClient{
					httpClient: fakeHttpClient,
				}

				/* act */
				_, actualErr := objectUnderTest.GetData(
					context.TODO(),
					model.GetDataReq{},
				)

				/* assert */
				Expect(actualErr).To(Equal(model.ErrDataRefResolution{}))

			})

		})
		Context("500", func() {
			It("should return expected result", func() {

				/* arrange */
				expectedErr := errors.New("dummyMsg")
				httpResp := &http.Response{
					Body:       ioutil.NopCloser(strings.NewReader(expectedErr.Error())),
					StatusCode: http.StatusInternalServerError,
				}

				fakeHttpClient := new(ihttp.FakeClient)
				fakeHttpClient.DoReturns(httpResp, nil)

				objectUnderTest := apiClient{
					httpClient: fakeHttpClient,
				}

				/* act */
				_, actualErr := objectUnderTest.GetData(
					context.TODO(),
					model.GetDataReq{},
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))

			})

		})
	})
})
