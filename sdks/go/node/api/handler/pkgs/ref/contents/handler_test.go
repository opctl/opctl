package contents

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	coreFakes "github.com/opctl/opctl/sdks/go/node/core/fakes"
)

var _ = Context("Handler", func() {
	Context("Handle", func() {
		It("should call handle.ListDescendants", func() {
			/* arrange */
			fakeDataHandle := new(modelFakes.FakeDataHandle)
			// error to trigger immediate return
			fakeDataHandle.ListDescendantsReturns(nil, errors.New("dummyError"))

			fakeCore := new(coreFakes.FakeCore)
			fakeCore.ResolveDataReturns(fakeDataHandle, nil)

			objectUnderTest := _handler{
				node: fakeCore,
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
			objectUnderTest.Handle(fakeDataHandle, providedHTTPResp, providedHTTPReq)

			/* assert */
			Expect(fakeDataHandle.ListDescendantsCallCount()).To(Equal(1))
		})
		Context("handle.ListDescendants errs", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedBody := "dummyErrorMsg"

				fakeDataHandle := new(modelFakes.FakeDataHandle)
				// error to trigger immediate return
				fakeDataHandle.ListDescendantsReturns(nil, errors.New(expectedBody))

				fakeCore := new(coreFakes.FakeCore)
				fakeCore.ResolveDataReturns(fakeDataHandle, nil)

				objectUnderTest := _handler{
					node: fakeCore,
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
				objectUnderTest.Handle(fakeDataHandle, providedHTTPResp, providedHTTPReq)

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

				fakeHandle := new(modelFakes.FakeDataHandle)
				contentsList := []*model.DirEntry{
					{Path: "dummyPath"},
				}
				fakeHandle.ListDescendantsReturns(contentsList, nil)

				expectedBodyBytes, err := json.Marshal(contentsList)
				if nil != err {
					panic(err)
				}

				objectUnderTest := _handler{
					node: fakeCore,
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
				objectUnderTest.Handle(fakeHandle, providedHTTPResp, providedHTTPReq)

				/* assert */
				Expect(providedHTTPResp.Code).To(Equal(http.StatusOK))
				Expect(providedHTTPResp.HeaderMap.Get("Content-Type")).To(Equal("application/json; charset=UTF-8"))
				actualBody := strings.TrimSpace(providedHTTPResp.Body.String())
				Expect(actualBody).To(Equal(string(expectedBodyBytes)))
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
