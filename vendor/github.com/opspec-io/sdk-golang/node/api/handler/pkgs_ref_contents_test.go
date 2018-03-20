package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"github.com/opspec-io/sdk-golang/node/core"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

var _ = Context("GET /pkgs/{ref}/contents", func() {
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

			httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/pkgs/%v/contents", providedPkgRef), nil)
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

			httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/pkgs/%v/contents", providedPkgRef), nil)
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
					"/pkgs/dummy/contents",
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
					"/pkgs/dummy/contents",
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
					"/pkgs/dummy/contents",
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
					"/pkgs/dummy/contents",
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
	Context("core.ResolveData doesn't err", func() {
		It("should call handle.ListContents", func() {
			/* arrange */
			fakeDataHandle := new(data.FakeHandle)
			// error to trigger immediate return
			fakeDataHandle.ListContentsReturns(nil, errors.New("dummyError"))

			fakeCore := new(core.Fake)
			fakeCore.ResolveDataReturns(fakeDataHandle, nil)

			objectUnderTest := New(fakeCore)
			recorder := httptest.NewRecorder()

			httpReq, err := http.NewRequest(
				http.MethodGet,
				"/pkgs/dummy/contents",
				nil,
			)
			if nil != err {
				panic(err.Error())
			}

			/* act */
			objectUnderTest.ServeHTTP(recorder, httpReq)

			/* assert */
			Expect(fakeDataHandle.ListContentsCallCount()).To(Equal(1))
		})
		Context("handle.ListContents errs", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedBody := "dummyErrorMsg"

				fakeDataHandle := new(data.FakeHandle)
				// error to trigger immediate return
				fakeDataHandle.ListContentsReturns(nil, errors.New(expectedBody))

				fakeCore := new(core.Fake)
				fakeCore.ResolveDataReturns(fakeDataHandle, nil)

				objectUnderTest := New(fakeCore)
				recorder := httptest.NewRecorder()

				httpReq, err := http.NewRequest(
					http.MethodGet,
					"/pkgs/dummy/contents",
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
		Context("handle.ListContents doesn't err", func() {
			Context("encoder.Encode errs", func() {
				It("should return expected result", func() {
					/* arrange */
					expectedBody := "dummyErrorMsg"

					fakeCore := new(core.Fake)
					fakeCore.ResolveDataReturns(new(data.FakeHandle), nil)

					fakeJSON := new(ijson.Fake)
					fakeJSON.NewEncoderReturns(json.NewEncoder(errWriter{Msg: expectedBody}))

					// construct objectUnderTest
					router := mux.NewRouter()
					objectUnderTest := _handler{
						core:   fakeCore,
						json:   fakeJSON,
						Router: router,
					}
					router.HandleFunc(api.URLPkgs_Ref_Contents, objectUnderTest.pkgs_ref_contents).Methods(http.MethodGet)

					recorder := httptest.NewRecorder()

					httpReq, err := http.NewRequest(
						http.MethodGet,
						"/pkgs/dummy/contents",
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
			Context("encoder.Encode doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */

					fakeCore := new(core.Fake)

					fakeHandle := new(data.FakeHandle)

					contentsList := []*model.PkgContent{
						{Path: "dummyPath"},
					}

					expectedBodyBytes, err := json.Marshal(contentsList)
					if nil != err {
						panic(err)
					}

					fakeHandle.ListContentsReturns(contentsList, nil)

					fakeCore.ResolveDataReturns(fakeHandle, nil)

					fakeJSON := new(ijson.Fake)
					fakeJSON.NewEncoderStub = func(w io.Writer) *json.Encoder {
						return json.NewEncoder(w)
					}

					// construct objectUnderTest
					router := mux.NewRouter()
					objectUnderTest := _handler{
						core:   fakeCore,
						json:   fakeJSON,
						Router: router,
					}
					router.HandleFunc(api.URLPkgs_Ref_Contents, objectUnderTest.pkgs_ref_contents).Methods(http.MethodGet)

					recorder := httptest.NewRecorder()

					httpReq, err := http.NewRequest(
						http.MethodGet,
						"/pkgs/dummy/contents",
						nil,
					)
					if nil != err {
						panic(err.Error())
					}

					/* act */
					objectUnderTest.ServeHTTP(recorder, httpReq)

					/* assert */
					Expect(recorder.Code).To(Equal(http.StatusOK))
					Expect(recorder.HeaderMap.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
					actualBody := strings.TrimSpace(recorder.Body.String())
					Expect(actualBody).To(Equal(string(expectedBodyBytes)))
				})
			})
		})
	})
})

// errWriter always errs
type errWriter struct {
	Msg string
}

func (this errWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New(this.Msg)
}
