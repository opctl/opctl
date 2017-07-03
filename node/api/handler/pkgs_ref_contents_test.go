package handler

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ijson"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api"
	"github.com/opspec-io/sdk-golang/node/core"
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
			fakeCore.ListPkgContentsReturns(nil, errors.New("dummyError"))

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
		It("should call core.ListPkgContents w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef%2F"
			expectedPkgRef, err := url.PathUnescape(providedPkgRef)
			if nil != err {
				Fail(err.Error())
			}

			fakeCore := new(core.Fake)
			// error to trigger immediate return
			fakeCore.ListPkgContentsReturns(nil, errors.New("dummyError"))

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
			actualPkgRef := fakeCore.ListPkgContentsArgsForCall(0)
			Expect(actualPkgRef).To(Equal(expectedPkgRef))
		})
	})
	Context("core.ListPkgContents errs", func() {
		It("should return expected result", func() {
			/* arrange */
			expectedBody := "dummyErrorMsg"

			fakeCore := new(core.Fake)
			// error to trigger immediate return
			fakeCore.ListPkgContentsReturns(nil, errors.New(expectedBody))

			objectUnderTest := New(fakeCore)
			recorder := httptest.NewRecorder()

			httpReq, err := http.NewRequest(
				http.MethodGet,
				"/pkgs/dummy/contents",
				nil,
			)
			if nil != err {
				Fail(err.Error())
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
	Context("core.ListPkgContents doesn't err", func() {
		It("should call json.Marshal w/ expected args", func() {
			/* arrange */
			expectedPkgList := []*model.PkgContent{
				{Path: "dummyPath"},
				{Path: "dummyPath2"},
			}

			fakeCore := new(core.Fake)
			fakeCore.ListPkgContentsReturns(expectedPkgList, nil)

			fakeJSON := new(ijson.Fake)
			// error to trigger immediate return
			fakeJSON.MarshalReturns(nil, errors.New("dummyError"))

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
				Fail(err.Error())
			}

			/* act */
			objectUnderTest.ServeHTTP(recorder, httpReq)

			/* assert */
			actualPkgList := fakeJSON.MarshalArgsForCall(0)
			Expect(actualPkgList).To(Equal(expectedPkgList))
		})
		Context("json.Marshal errs", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedBody := "dummyErrorMsg"

				fakeJSON := new(ijson.Fake)
				fakeJSON.MarshalReturns(nil, errors.New(expectedBody))

				// construct objectUnderTest
				router := mux.NewRouter()
				objectUnderTest := _handler{
					core:   new(core.Fake),
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
					Fail(err.Error())
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
		Context("json.Marshal doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedBody := "dummyJSON"
				fakeJSON := new(ijson.Fake)
				fakeJSON.MarshalReturns([]byte(expectedBody), nil)

				// construct objectUnderTest
				router := mux.NewRouter()
				objectUnderTest := _handler{
					core:   new(core.Fake),
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
					Fail(err.Error())
				}

				/* act */
				objectUnderTest.ServeHTTP(recorder, httpReq)

				/* assert */
				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.HeaderMap.Get("Content-Type")).To(Equal("application/json"))
				actualBody := strings.TrimSpace(recorder.Body.String())
				Expect(actualBody).To(Equal(expectedBody))
			})
		})
	})
})
