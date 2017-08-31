package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"github.com/opspec-io/sdk-golang/node/core"
	"github.com/opspec-io/sdk-golang/pkg"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

var _ = Context("GET /pkgs/{ref}/contents", func() {
	Context("pkgRef isn't valid URL encoded segment", func() {
		It("should return expected error", func() {
			/* arrange */
			invalidURLEncoding := "%%"
			providedPath := fmt.Sprintf("/pkgs/dummyRef%v/contents", invalidURLEncoding)
			expectedBody := fmt.Sprintf("invalid URL escape \"%v\"", invalidURLEncoding)

			fakeCore := new(core.Fake)
			// error to trigger immediate return
			fakeCore.ResolvePkgReturns(nil, errors.New("dummyError"))

			objectUnderTest := New(fakeCore)
			recorder := httptest.NewRecorder()

			// manually construct request so path doesn't get parsed
			httpReq := &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{Path: providedPath},
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
	Context("pkgRef is valid URL encoded segment", func() {
		Context("has basic auth header", func() {
			It("should call core.ResolvePkg w/ expected args", func() {
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
				fakeCore.ResolvePkgReturns(nil, errors.New("dummyError"))

				objectUnderTest := New(fakeCore)
				recorder := httptest.NewRecorder()

				// manually construct request so path doesn't get parsed
				httpReq := &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Path: fmt.Sprintf("/pkgs/%v/contents", providedPkgRef)},
					Header: map[string][]string{},
				}
				httpReq.SetBasicAuth(providedUsername, providedPassword)

				/* act */
				objectUnderTest.ServeHTTP(recorder, httpReq)

				/* assert */
				actualPkgRef, actualPullCreds := fakeCore.ResolvePkgArgsForCall(0)
				Expect(actualPkgRef).To(Equal(expectedPkgRef))
				Expect(*actualPullCreds).To(Equal(model.PullCreds{
					Username: providedUsername,
					Password: providedPassword,
				}))
			})
		})
	})
	Context("doesn't have basic auth header", func() {
		It("should call core.ResolvePkg w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef%2F"
			expectedPkgRef, err := url.PathUnescape(providedPkgRef)
			if nil != err {
				panic(err.Error())
			}

			fakeCore := new(core.Fake)
			// error to trigger immediate return
			fakeCore.ResolvePkgReturns(nil, errors.New("dummyError"))

			objectUnderTest := New(fakeCore)
			recorder := httptest.NewRecorder()

			// manually construct request so path doesn't get parsed
			httpReq := &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{Path: fmt.Sprintf("/pkgs/%v/contents", providedPkgRef)},
			}

			/* act */
			objectUnderTest.ServeHTTP(recorder, httpReq)

			/* assert */
			actualPkgRef, actualPullCreds := fakeCore.ResolvePkgArgsForCall(0)
			Expect(actualPkgRef).To(Equal(expectedPkgRef))
			Expect(actualPullCreds).To(BeNil())
		})
	})
	Context("core.ResolvePkg errs", func() {
		Context("err is model.ErrPkgPullAuthentication", func() {
			It("should return expected result", func() {
				/* arrange */
				pkgPullAuthenticationErr := model.ErrPkgPullAuthentication{}

				fakeCore := new(core.Fake)
				// error to trigger immediate return
				fakeCore.ResolvePkgReturns(nil, pkgPullAuthenticationErr)

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
		Context("err is model.ErrPkgPullAuthorization", func() {
			It("should return expected result", func() {
				/* arrange */
				pkgPullAuthorizationErr := model.ErrPkgPullAuthorization{}

				fakeCore := new(core.Fake)
				// error to trigger immediate return
				fakeCore.ResolvePkgReturns(nil, pkgPullAuthorizationErr)

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
		Context("err is model.ErrPkgNotFound", func() {
			It("should return expected result", func() {
				/* arrange */
				pkgNotFoundErr := model.ErrPkgNotFound{}

				fakeCore := new(core.Fake)
				// error to trigger immediate return
				fakeCore.ResolvePkgReturns(nil, pkgNotFoundErr)

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
				fakeCore.ResolvePkgReturns(nil, errors.New(expectedBody))

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
	Context("core.ResolvePkg doesn't err", func() {
		It("should call handle.ListContents", func() {
			/* arrange */
			fakePkgHandle := new(pkg.FakeHandle)
			// error to trigger immediate return
			fakePkgHandle.ListContentsReturns(nil, errors.New("dummyError"))

			fakeCore := new(core.Fake)
			fakeCore.ResolvePkgReturns(fakePkgHandle, nil)

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
			Expect(fakePkgHandle.ListContentsCallCount()).To(Equal(1))
		})
		Context("handle.ListContents errs", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedBody := "dummyErrorMsg"

				fakePkgHandle := new(pkg.FakeHandle)
				// error to trigger immediate return
				fakePkgHandle.ListContentsReturns(nil, errors.New(expectedBody))

				fakeCore := new(core.Fake)
				fakeCore.ResolvePkgReturns(fakePkgHandle, nil)

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
					fakeCore.ResolvePkgReturns(new(pkg.FakeHandle), nil)

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

					fakeHandle := new(pkg.FakeHandle)

					contentsList := []*model.PkgContent{
						{Path: "dummyPath"},
					}

					expectedBodyBytes, err := json.Marshal(contentsList)
					if nil != err {
						panic(err)
					}

					fakeHandle.ListContentsReturns(contentsList, nil)

					fakeCore.ResolvePkgReturns(fakeHandle, nil)

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
