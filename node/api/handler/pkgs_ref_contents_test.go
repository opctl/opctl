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
	"strings"
)

var _ = Context("GET /pkgs/{ref}/contents", func() {
	It("should call core.ListPkgContents w/ expected args", func() {
		/* arrange */
		providedPkgRef := "dummyPkgRef"

		fakeCore := new(core.Fake)
		// error to trigger immediate return
		fakeCore.ListPkgContentsReturns(nil, errors.New("dummyError"))

		objectUnderTest := New(fakeCore)
		recorder := httptest.NewRecorder()

		httpReq, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/pkgs/%v/contents", providedPkgRef),
			nil,
		)
		if nil != err {
			Fail(err.Error())
		}

		/* act */
		objectUnderTest.ServeHTTP(recorder, httpReq)

		/* assert */
		actualPkgRef := fakeCore.ListPkgContentsArgsForCall(0)
		Expect(actualPkgRef).To(Equal(providedPkgRef))
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
