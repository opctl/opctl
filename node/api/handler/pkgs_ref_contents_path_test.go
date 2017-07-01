package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-interfaces/ihttp"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/node/api"
	"github.com/opspec-io/sdk-golang/node/core"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"time"
)

var _ = Context("GET /pkgs/{ref}/contents/{path}", func() {
	It("should call core.GetPkgContent w/ expected args", func() {
		/* arrange */
		providedPkgRef := "dummyPkgRef"
		providedPkgPath := "dummyPkgPath"

		fakeCore := new(core.Fake)
		// error to trigger immediate return
		fakeCore.GetPkgContentReturns(nil, errors.New("dummyError"))

		objectUnderTest := New(fakeCore)
		recorder := httptest.NewRecorder()

		httpReq, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/pkgs/%v/contents/%v", providedPkgRef, providedPkgPath),
			nil,
		)
		if nil != err {
			Fail(err.Error())
		}

		/* act */
		objectUnderTest.ServeHTTP(recorder, httpReq)

		/* assert */
		actualPkgRef, actualPkgPath := fakeCore.GetPkgContentArgsForCall(0)
		Expect(actualPkgRef).To(Equal(providedPkgRef))
		Expect(actualPkgPath).To(Equal(providedPkgPath))
	})
	Context("core.GetPkgContent errs", func() {
		It("should return expected result", func() {
			/* arrange */
			expectedErrMsg := "dummyErrorMsg"

			fakeCore := new(core.Fake)
			// error to trigger immediate return
			fakeCore.GetPkgContentReturns(nil, errors.New(expectedErrMsg))

			objectUnderTest := New(fakeCore)
			recorder := httptest.NewRecorder()

			httpReq, err := http.NewRequest(
				http.MethodGet,
				"/pkgs/dummyRef/contents/dummyPath",
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
			actualErrMsg := strings.TrimSpace(recorder.Body.String())
			Expect(actualErrMsg).To(Equal(expectedErrMsg))
		})
	})
	Context("core.GetPkgContent doesn't err", func() {
		It("should call http.ServeContent w/ expected args", func() {
			/* arrange */
			providedPath := "dummyPath"

			expectedReadSeeker, err := ioutil.TempFile("", "")
			defer expectedReadSeeker.Close()

			fakeCore := new(core.Fake)
			// error to trigger immediate return
			fakeCore.GetPkgContentReturns(expectedReadSeeker, nil)

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
				Fail(err.Error())
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
