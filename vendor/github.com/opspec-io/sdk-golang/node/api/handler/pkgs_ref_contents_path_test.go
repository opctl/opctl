package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-interfaces/ihttp"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"github.com/opspec-io/sdk-golang/node/core"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strings"
	"time"
)

var _ = Context("GET /pkgs/{ref}/contents/{path}", func() {
	Context("has basic auth header", func() {
		It("should call core.ResolveData w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef%2F"
			expectedPkgRef, err := url.PathUnescape(providedPkgRef)
			if nil != err {
				panic(err.Error())
			}

			providedUsername := "dummyUsername"
			providedPassword := "dummyPassword"

			fakeCore := new(core.Fake)
			// error to trigger immediate return
			fakeCore.ResolveDataReturns(nil, errors.New("dummyError"))

			objectUnderTest := New(fakeCore)
			recorder := httptest.NewRecorder()

			httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/pkgs/%v/contents/dummyPath", providedPkgRef), nil)
			if nil != err {
				panic(err.Error())
			}
			httpReq.SetBasicAuth(providedUsername, providedPassword)

			/* act */
			objectUnderTest.ServeHTTP(recorder, httpReq)

			/* assert */
			_,
				actualPkgRef,
				actualPullCreds := fakeCore.ResolveDataArgsForCall(0)

			Expect(actualPkgRef).To(Equal(expectedPkgRef))
			Expect(*actualPullCreds).To(Equal(model.PullCreds{
				Username: providedUsername,
				Password: providedPassword,
			}))
		})
	})
	Context("doesn't have basic auth header", func() {
		It("should call core.ResolveData w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef%2F"
			expectedPkgRef, err := url.PathUnescape(providedPkgRef)
			if nil != err {
				panic(err.Error())
			}

			fakeCore := new(core.Fake)
			// error to trigger immediate return
			fakeCore.ResolveDataReturns(nil, errors.New("dummyError"))

			objectUnderTest := New(fakeCore)
			recorder := httptest.NewRecorder()

			httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/pkgs/%v/contents/dummyPath", providedPkgRef), nil)
			if nil != err {
				panic(err.Error())
			}

			/* act */
			objectUnderTest.ServeHTTP(recorder, httpReq)

			/* assert */
			_,
				actualPkgRef,
				actualPullCreds := fakeCore.ResolveDataArgsForCall(0)

			Expect(actualPkgRef).To(Equal(expectedPkgRef))
			Expect(actualPullCreds).To(BeNil())
		})
	})
	Context("core.ResolveData errs", func() {
		Context("err is model.ErrDataProviderAuthentication", func() {
			It("should return expected result", func() {
				/* arrange */
				pkgPullAuthenticationErr := model.ErrDataProviderAuthentication{}

				fakeCore := new(core.Fake)
				// error to trigger immediate return
				fakeCore.ResolveDataReturns(nil, pkgPullAuthenticationErr)

				objectUnderTest := New(fakeCore)
				recorder := httptest.NewRecorder()

				httpReq, err := http.NewRequest(
					http.MethodGet,
					"/pkgs/dummyRef/contents/dummyPath",
					nil,
				)
				if nil != err {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.ServeHTTP(recorder, httpReq)

				/* assert */
				Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
				Expect(recorder.HeaderMap.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
				actualBody := strings.TrimSpace(recorder.Body.String())
				Expect(actualBody).To(Equal(pkgPullAuthenticationErr.Error()))
			})
		})
		Context("err is model.ErrDataProviderAuthorization", func() {
			It("should return expected result", func() {
				/* arrange */
				pkgPullAuthorizationErr := model.ErrDataProviderAuthorization{}

				fakeCore := new(core.Fake)
				// error to trigger immediate return
				fakeCore.ResolveDataReturns(nil, pkgPullAuthorizationErr)

				objectUnderTest := New(fakeCore)
				recorder := httptest.NewRecorder()

				httpReq, err := http.NewRequest(
					http.MethodGet,
					"/pkgs/dummyRef/contents/dummyPath",
					nil,
				)
				if nil != err {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.ServeHTTP(recorder, httpReq)

				/* assert */
				Expect(recorder.Code).To(Equal(http.StatusForbidden))
				Expect(recorder.HeaderMap.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
				actualBody := strings.TrimSpace(recorder.Body.String())
				Expect(actualBody).To(Equal(pkgPullAuthorizationErr.Error()))
			})
		})
		Context("err is model.ErrDataRefResolution", func() {
			It("should return expected result", func() {
				/* arrange */
				pkgNotFoundErr := model.ErrDataRefResolution{}

				fakeCore := new(core.Fake)
				// error to trigger immediate return
				fakeCore.ResolveDataReturns(nil, pkgNotFoundErr)

				objectUnderTest := New(fakeCore)
				recorder := httptest.NewRecorder()

				httpReq, err := http.NewRequest(
					http.MethodGet,
					"/pkgs/dummyRef/contents/dummyPath",
					nil,
				)
				if nil != err {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.ServeHTTP(recorder, httpReq)

				/* assert */
				Expect(recorder.Code).To(Equal(http.StatusNotFound))
				Expect(recorder.HeaderMap.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
				actualBody := strings.TrimSpace(recorder.Body.String())
				Expect(actualBody).To(Equal(pkgNotFoundErr.Error()))
			})
		})
		Context("err isn't known type", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedBody := "dummyErrorMsg"

				fakeCore := new(core.Fake)
				// error to trigger immediate return
				fakeCore.ResolveDataReturns(nil, errors.New(expectedBody))

				objectUnderTest := New(fakeCore)
				recorder := httptest.NewRecorder()

				httpReq, err := http.NewRequest(
					http.MethodGet,
					"/pkgs/dummyRef/contents/dummyPath",
					nil,
				)
				if nil != err {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.ServeHTTP(recorder, httpReq)

				/* assert */
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				Expect(recorder.HeaderMap.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
				actualBody := strings.TrimSpace(recorder.Body.String())
				Expect(actualBody).To(Equal(expectedBody))
			})
		})
	})
	Context("core.GetPkgContent doesn't err", func() {
		It("should call handle.GetContent w/ expected args", func() {
			/* arrange */
			providedContentPath := "dummyPkgRef%2F"
			expectedContentPath, err := url.PathUnescape(providedContentPath)
			if nil != err {
				panic(err.Error())
			}

			fakeDataHandle := new(data.FakeHandle)
			// error to trigger immediate return
			fakeDataHandle.GetContentReturns(nil, errors.New("dummyError"))

			fakeCore := new(core.Fake)
			fakeCore.ResolveDataReturns(fakeDataHandle, nil)

			objectUnderTest := New(fakeCore)
			recorder := httptest.NewRecorder()

			httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/pkgs/dummyPkgRef/contents/%v", providedContentPath), nil)
			if nil != err {
				panic(err.Error())
			}

			/* act */
			objectUnderTest.ServeHTTP(recorder, httpReq)

			/* assert */
			_,
				actualPkgRef := fakeDataHandle.GetContentArgsForCall(0)

			Expect(actualPkgRef).To(Equal(expectedContentPath))
		})
		Context("handle.GetContent errors", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedBody := "dummyErrorMsg"

				fakeDataHandle := new(data.FakeHandle)
				fakeDataHandle.GetContentReturns(nil, errors.New(expectedBody))

				fakeCore := new(core.Fake)
				fakeCore.ResolveDataReturns(fakeDataHandle, nil)

				objectUnderTest := New(fakeCore)
				recorder := httptest.NewRecorder()

				httpReq, err := http.NewRequest(
					http.MethodGet,
					"/pkgs/dummyRef/contents/dummyPath",
					nil,
				)
				if nil != err {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.ServeHTTP(recorder, httpReq)

				/* assert */
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				Expect(recorder.HeaderMap.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
				actualBody := strings.TrimSpace(recorder.Body.String())
				Expect(actualBody).To(Equal(expectedBody))
			})
		})
		Context("handle.GetContent doesn't error", func() {
			It("should call http.ServeContent w/ expected args", func() {
				/* arrange */
				providedPath := "dummyPath"

				expectedReadSeeker, err := ioutil.TempFile("", "")
				defer expectedReadSeeker.Close()

				fakeDataHandle := new(data.FakeHandle)
				fakeDataHandle.GetContentReturns(expectedReadSeeker, nil)

				fakeCore := new(core.Fake)
				fakeCore.ResolveDataReturns(fakeDataHandle, nil)

				fakeHTTP := new(ihttp.Fake)

				// construct objectUnderTest
				router := mux.NewRouter()
				objectUnderTest := _handler{
					core:   fakeCore,
					http:   fakeHTTP,
					Router: router,
				}
				router.HandleFunc(api.URLPkgs_Ref_Contents_Path, objectUnderTest.pkgs_ref_contents_path).Methods(http.MethodGet)

				recorder := httptest.NewRecorder()

				providedRequest, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/pkgs/dummyRef/contents/%v", providedPath),
					nil,
				)
				if nil != err {
					panic(err.Error())
				}
				providedRequest = providedRequest.WithContext(context.TODO())

				/* act */
				objectUnderTest.ServeHTTP(recorder, providedRequest)

				/* assert */
				actualResponseWriter,
					actualRequest,
					actualName,
					actualTime,
					actualReadSeeker := fakeHTTP.ServeContentArgsForCall(0)

				// ignore context
				actualRequest = actualRequest.WithContext(providedRequest.Context())

				Expect(actualResponseWriter).To(Equal(recorder))
				Expect(*actualRequest).To(Equal(*providedRequest))
				Expect(actualName).To(Equal(path.Base(providedPath)))
				Expect(actualTime).To(Equal(time.Time{}))
				Expect(actualReadSeeker).To(Equal(expectedReadSeeker))
			})
		})
	})
})
