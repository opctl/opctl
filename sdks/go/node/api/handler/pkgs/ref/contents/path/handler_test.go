package path

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"time"

	"github.com/golang-interfaces/ihttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	. "github.com/opctl/opctl/sdks/go/node/core/fakes"
)

var _ = Context("Handler", func() {
	Context("Handle", func() {
		It("should call handle.GetContent w/ expected args", func() {
			/* arrange */
			providedContentPath := "dummyOpRef"

			fakeDataHandle := new(modelFakes.FakeDataHandle)
			// error to trigger immediate return
			fakeDataHandle.GetContentReturns(nil, errors.New("dummyError"))

			fakeCore := new(FakeCore)
			fakeCore.ResolveDataReturns(fakeDataHandle, nil)

			objectUnderTest := _handler{
				core: fakeCore,
			}
			providedHTTPResp := httptest.NewRecorder()

			providedHTTPReq, err := http.NewRequest(http.MethodGet, "", nil)
			if nil != err {
				panic(err.Error())
			}

			/* act */
			objectUnderTest.Handle(
				fakeDataHandle,
				providedContentPath,
				providedHTTPResp,
				providedHTTPReq,
			)

			/* assert */
			_,
				actualOpRef := fakeDataHandle.GetContentArgsForCall(0)

			Expect(actualOpRef).To(Equal(providedContentPath))
		})
		Context("handle.GetContent errors", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedBody := "dummyErrorMsg"

				fakeDataHandle := new(modelFakes.FakeDataHandle)
				fakeDataHandle.GetContentReturns(nil, errors.New(expectedBody))

				fakeCore := new(FakeCore)
				fakeCore.ResolveDataReturns(fakeDataHandle, nil)

				objectUnderTest := _handler{
					core: fakeCore,
				}
				providedHTTPResp := httptest.NewRecorder()

				providedHTTPReq, err := http.NewRequest(
					http.MethodGet,
					"",
					nil,
				)
				if nil != err {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.Handle(
					fakeDataHandle,
					"dummyPath",
					providedHTTPResp,
					providedHTTPReq,
				)

				/* assert */
				Expect(providedHTTPResp.Code).To(Equal(http.StatusInternalServerError))
				Expect(providedHTTPResp.HeaderMap.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
				actualBody := strings.TrimSpace(providedHTTPResp.Body.String())
				Expect(actualBody).To(Equal(expectedBody))
			})
		})
		Context("handle.GetContent doesn't error", func() {
			It("should call http.ServeContent w/ expected args", func() {
				/* arrange */
				providedPath := "dummyPath"

				expectedReadSeeker, err := ioutil.TempFile("", "")
				defer expectedReadSeeker.Close()

				fakeDataHandle := new(modelFakes.FakeDataHandle)
				fakeDataHandle.GetContentReturns(expectedReadSeeker, nil)

				fakeCore := new(FakeCore)
				fakeCore.ResolveDataReturns(fakeDataHandle, nil)

				fakeHTTP := new(ihttp.Fake)

				objectUnderTest := _handler{
					core: fakeCore,
					http: fakeHTTP,
				}

				providedHTTPResp := httptest.NewRecorder()

				providedRequest, err := http.NewRequest(
					http.MethodGet,
					"",
					nil,
				)
				if nil != err {
					panic(err.Error())
				}
				providedRequest = providedRequest.WithContext(context.TODO())

				/* act */
				objectUnderTest.Handle(
					fakeDataHandle,
					providedPath,
					providedHTTPResp,
					providedRequest,
				)

				/* assert */
				actualResponseWriter,
					actualRequest,
					actualName,
					actualTime,
					actualReadSeeker := fakeHTTP.ServeContentArgsForCall(0)

				// ignore context
				actualRequest = actualRequest.WithContext(providedRequest.Context())

				Expect(actualResponseWriter).To(Equal(providedHTTPResp))
				Expect(*actualRequest).To(Equal(*providedRequest))
				Expect(actualName).To(Equal(path.Base(providedPath)))
				Expect(actualTime).To(Equal(time.Time{}))
				Expect(actualReadSeeker).To(Equal(expectedReadSeeker))
			})
		})
	})
})
