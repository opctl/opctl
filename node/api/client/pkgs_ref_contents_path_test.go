package client

import (
	"context"
	"errors"
	"github.com/golang-interfaces/ihttp"
	"github.com/jfbus/httprs"
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
			PkgRef:      "dummyOpRef",
			PullCreds: &model.PullCreds{
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

		objectUnderTest := client{
			httpClient: fakeHttpClient,
		}

		/* act */
		objectUnderTest.GetPkgContent(providedCtx, providedReq)

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

			objectUnderTest := client{
				httpClient: fakeHttpClient,
			}

			/* act */
			actualReadSeekCloser, actualErr := objectUnderTest.GetPkgContent(
				context.TODO(),
				model.GetPkgContentReq{},
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
					Body:       ioutil.NopCloser(nil),
					StatusCode: http.StatusUnauthorized,
				}

				fakeHttpClient := new(ihttp.FakeClient)
				fakeHttpClient.DoReturns(httpResp, nil)

				objectUnderTest := client{
					httpClient: fakeHttpClient,
				}

				/* act */
				_, actualErr := objectUnderTest.GetPkgContent(
					context.TODO(),
					model.GetPkgContentReq{},
				)

				/* assert */
				Expect(actualErr).To(Equal(model.ErrDataProviderAuthentication{}))

			})
		})
		Context("403", func() {
			It("should return expected result", func() {

				/* arrange */
				httpResp := &http.Response{
					Body:       ioutil.NopCloser(nil),
					StatusCode: http.StatusForbidden,
				}

				fakeHttpClient := new(ihttp.FakeClient)
				fakeHttpClient.DoReturns(httpResp, nil)

				objectUnderTest := client{
					httpClient: fakeHttpClient,
				}

				/* act */
				_, actualErr := objectUnderTest.GetPkgContent(
					context.TODO(),
					model.GetPkgContentReq{},
				)

				/* assert */
				Expect(actualErr).To(Equal(model.ErrDataProviderAuthorization{}))

			})

		})
		Context("404", func() {
			It("should return expected result", func() {

				/* arrange */
				httpResp := &http.Response{
					Body:       ioutil.NopCloser(nil),
					StatusCode: http.StatusNotFound,
				}

				fakeHttpClient := new(ihttp.FakeClient)
				fakeHttpClient.DoReturns(httpResp, nil)

				objectUnderTest := client{
					httpClient: fakeHttpClient,
				}

				/* act */
				_, actualErr := objectUnderTest.GetPkgContent(
					context.TODO(),
					model.GetPkgContentReq{},
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

				objectUnderTest := client{
					httpClient: fakeHttpClient,
				}

				/* act */
				_, actualErr := objectUnderTest.GetPkgContent(
					context.TODO(),
					model.GetPkgContentReq{},
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))

			})

		})
	})
})
