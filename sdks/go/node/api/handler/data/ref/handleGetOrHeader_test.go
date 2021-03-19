package ref

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golang-interfaces/ihttp"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	coreFakes "github.com/opctl/opctl/sdks/go/node/core/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("HandleGetOrHeader", func() {

	Context("req.Method not GET or HEAD", func() {
		It("should return expected result", func() {
			/* arrange */
			objectUnderTest := _handleGetOrHeader{}
			providedHTTPResp := httptest.NewRecorder()

			providedHTTPReq, err := http.NewRequest("notGETorHEAD", "", nil)
			if err != nil {
				panic(err.Error())
			}

			/* act */
			objectUnderTest.HandleGetOrHead(
				new(modelFakes.FakeDataHandle),
				providedHTTPResp,
				providedHTTPReq,
			)

			/* assert */
			Expect(providedHTTPResp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})
	Context("req path non empty", func() {
		Describe("Not found", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _handleGetOrHeader{}
				providedHTTPResp := httptest.NewRecorder()

				providedHTTPReq, err := http.NewRequest(http.MethodGet, "dummyPath", nil)
				if err != nil {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.HandleGetOrHead(
					new(modelFakes.FakeDataHandle),
					providedHTTPResp,
					providedHTTPReq,
				)

				/* assert */
				Expect(providedHTTPResp.Code).To(Equal(http.StatusNotFound))
			})
		})

		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		Describe("ref + path is dir", func() {

			It("should call handle.ListDescendants", func() {
				/* arrange */
				fakeDataHandle := new(modelFakes.FakeDataHandle)
				fakeDataHandle.PathReturns(&wd)
				// error to trigger immediate return
				fakeDataHandle.ListDescendantsReturns(nil, errors.New("dummyError"))

				fakeCore := new(coreFakes.FakeCore)
				fakeCore.ResolveDataReturns(fakeDataHandle, nil)

				objectUnderTest := _handleGetOrHeader{
					core: fakeCore,
				}
				providedHTTPResp := httptest.NewRecorder()

				providedHTTPReq, err := http.NewRequest(
					http.MethodGet,
					"",
					nil,
				)
				if err != nil {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.HandleGetOrHead(fakeDataHandle, providedHTTPResp, providedHTTPReq)

				/* assert */
				Expect(fakeDataHandle.ListDescendantsCallCount()).To(Equal(1))
			})
			Context("handle.ListDescendants errs", func() {
				It("should return expected result", func() {
					/* arrange */
					expectedBody := "dummyErrorMsg"

					fakeDataHandle := new(modelFakes.FakeDataHandle)
					fakeDataHandle.PathReturns(&wd)
					// error to trigger immediate return
					fakeDataHandle.ListDescendantsReturns(nil, errors.New(expectedBody))

					fakeCore := new(coreFakes.FakeCore)
					fakeCore.ResolveDataReturns(fakeDataHandle, nil)

					objectUnderTest := _handleGetOrHeader{
						core: fakeCore,
					}
					providedHTTPResp := httptest.NewRecorder()

					providedHTTPReq, err := http.NewRequest(
						http.MethodGet,
						"",
						nil,
					)
					if err != nil {
						panic(err.Error())
					}

					/* act */
					objectUnderTest.HandleGetOrHead(fakeDataHandle, providedHTTPResp, providedHTTPReq)

					/* assert */
					Expect(providedHTTPResp.Code).To(Equal(http.StatusInternalServerError))
					Expect(providedHTTPResp.HeaderMap.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
					actualBody := strings.TrimSpace(providedHTTPResp.Body.String())
					Expect(actualBody).To(Equal(expectedBody))
				})
			})
			Context("handle.ListDescendants doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */

					fakeCore := new(coreFakes.FakeCore)

					fakeDataHandle := new(modelFakes.FakeDataHandle)
					fakeDataHandle.PathReturns(&wd)
					contentsList := []*model.DirEntry{
						{Path: "dummyPath"},
					}
					fakeDataHandle.ListDescendantsReturns(contentsList, nil)

					expectedBodyBytes, err := json.Marshal(contentsList)
					if err != nil {
						panic(err)
					}

					objectUnderTest := _handleGetOrHeader{
						core: fakeCore,
					}

					providedHTTPResp := httptest.NewRecorder()

					providedHTTPReq, err := http.NewRequest(
						http.MethodGet,
						"",
						nil,
					)
					if err != nil {
						panic(err.Error())
					}

					/* act */
					objectUnderTest.HandleGetOrHead(fakeDataHandle, providedHTTPResp, providedHTTPReq)

					/* assert */
					Expect(providedHTTPResp.Code).To(Equal(http.StatusOK))
					Expect(providedHTTPResp.HeaderMap.Get("Content-Type")).To(Equal("application/vnd.opspec.0.1.6.dir+json; charset=UTF-8"))
					actualBody := strings.TrimSpace(providedHTTPResp.Body.String())
					Expect(actualBody).To(Equal(string(expectedBodyBytes)))
				})
			})
		})
		Describe("ref + path is file", func() {
			currentFilePath := path.Join(wd, "handleGetOrHeader_test.go")

			It("should call handle.GetContent w/ expected args", func() {
				/* arrange */
				fakeDataHandle := new(modelFakes.FakeDataHandle)
				fakeDataHandle.PathReturns(&currentFilePath)
				// error to trigger immediate return
				fakeDataHandle.GetContentReturns(nil, errors.New("dummyError"))

				fakeCore := new(coreFakes.FakeCore)
				fakeCore.ResolveDataReturns(fakeDataHandle, nil)

				objectUnderTest := _handleGetOrHeader{
					core: fakeCore,
				}
				providedHTTPResp := httptest.NewRecorder()

				providedHTTPReq, err := http.NewRequest(http.MethodGet, "", nil)
				if err != nil {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.HandleGetOrHead(
					fakeDataHandle,
					providedHTTPResp,
					providedHTTPReq,
				)

				/* assert */
				_,
					actualOpRef := fakeDataHandle.GetContentArgsForCall(0)

				Expect(actualOpRef).To(BeEmpty())
			})
			Context("handle.GetContent errors", func() {
				It("should return expected result", func() {
					/* arrange */
					expectedBody := "dummyErrorMsg"

					fakeDataHandle := new(modelFakes.FakeDataHandle)
					fakeDataHandle.PathReturns(&currentFilePath)
					fakeDataHandle.GetContentReturns(nil, errors.New(expectedBody))

					fakeCore := new(coreFakes.FakeCore)
					fakeCore.ResolveDataReturns(fakeDataHandle, nil)

					objectUnderTest := _handleGetOrHeader{
						core: fakeCore,
					}
					providedHTTPResp := httptest.NewRecorder()

					providedHTTPReq, err := http.NewRequest(
						http.MethodGet,
						"",
						nil,
					)
					if err != nil {
						panic(err.Error())
					}

					/* act */
					objectUnderTest.HandleGetOrHead(
						fakeDataHandle,
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

					expectedReadSeeker, err := ioutil.TempFile("", "")
					defer expectedReadSeeker.Close()

					fakeDataHandle := new(modelFakes.FakeDataHandle)
					fakeDataHandle.PathReturns(&currentFilePath)
					fakeDataHandle.GetContentReturns(expectedReadSeeker, nil)

					fakeCore := new(coreFakes.FakeCore)
					fakeCore.ResolveDataReturns(fakeDataHandle, nil)

					fakeHTTP := new(ihttp.Fake)

					objectUnderTest := _handleGetOrHeader{
						core: fakeCore,
						http: fakeHTTP,
					}

					providedHTTPResp := httptest.NewRecorder()

					providedRequest, err := http.NewRequest(
						http.MethodGet,
						"",
						nil,
					)
					if err != nil {
						panic(err.Error())
					}
					providedRequest = providedRequest.WithContext(context.TODO())

					/* act */
					objectUnderTest.HandleGetOrHead(
						fakeDataHandle,
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
					Expect(actualName).To(Equal(path.Base(fakeDataHandle.Ref())))
					Expect(actualTime).To(Equal(time.Time{}))
					Expect(actualReadSeeker).To(Equal(expectedReadSeeker))
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
