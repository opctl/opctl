package client

import (
	"context"
	"errors"
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

		expectedHttpReq, _ := http.NewRequest(
			"GET",
			expectedReqUrl.String(),
			nil,
		)

		expectedHttpReq.SetBasicAuth(
			providedReq.PullCreds.Username,
			providedReq.PullCreds.Password,
		)

		fakeHttpClient := new(ihttp.FakeClient)
		fakeHttpClient.DoReturns(
			&http.Response{
				Body:       ioutil.NopCloser(strings.NewReader("[]")),
				StatusCode: http.StatusOK,
			},
			nil,
		)

		objectUnderTest := client{
			httpClient: fakeHttpClient,
		}

		/* act */
		objectUnderTest.ListPkgContents(providedCtx, providedReq)

		/* assert */
		actualHttpReq := fakeHttpClient.DoArgsForCall(0)

		Expect(actualHttpReq.URL).To(Equal(expectedHttpReq.URL))
		Expect(actualHttpReq.Body).To(BeNil())
		Expect(actualHttpReq.Header).To(Equal(expectedHttpReq.Header))
		Expect(actualHttpReq.Context()).To(Equal(providedCtx))

	})
	Context("StatusCode < 400", func() {

		It("should return expected result", func() {

			/* arrange */
			httpResp := &http.Response{
				Body:       ioutil.NopCloser(strings.NewReader("[]")),
				StatusCode: http.StatusOK,
			}

			fakeHttpClient := new(ihttp.FakeClient)
			fakeHttpClient.DoReturns(httpResp, nil)

			objectUnderTest := client{
				httpClient: fakeHttpClient,
			}

			/* act */
			actualContentsList, actualErr := objectUnderTest.ListPkgContents(
				context.TODO(),
				model.ListPkgContentsReq{},
			)

			/* assert */
			Expect(actualContentsList).To(Equal([]*model.PkgContent{}))
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
				_, actualErr := objectUnderTest.ListPkgContents(
					context.TODO(),
					model.ListPkgContentsReq{},
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
				_, actualErr := objectUnderTest.ListPkgContents(
					context.TODO(),
					model.ListPkgContentsReq{},
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
				_, actualErr := objectUnderTest.ListPkgContents(
					context.TODO(),
					model.ListPkgContentsReq{},
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
				_, actualErr := objectUnderTest.ListPkgContents(
					context.TODO(),
					model.ListPkgContentsReq{},
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))

			})

		})
	})
})
