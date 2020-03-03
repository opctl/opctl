package ops

import (
	killsFakes "github.com/opctl/opctl/sdks/go/node/api/handler/ops/kills/fakes"
	startsFakes "github.com/opctl/opctl/sdks/go/node/api/handler/ops/starts/fakes"
	. "github.com/opctl/opctl/sdks/go/node/core/fakes"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("Handler", func() {
	Context("NewHandler", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewHandler(new(FakeCore))).Should(Not(BeNil()))
		})
	})
	Context("Handle", func() {
		Context("next URL path segment isn't starts or kills", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _handler{}
				providedHTTPResp := httptest.NewRecorder()

				providedHTTPReq, err := http.NewRequest("dummyMethod", "", nil)
				if nil != err {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.Handle(providedHTTPResp, providedHTTPReq)

				/* assert */
				Expect(providedHTTPResp.Code).To(Equal(http.StatusNotFound))
			})
		})
		Context("next URL path segment is starts", func() {
			It("should call refHandler.Handle w/ expected args", func() {
				/* arrange */
				fakeStartsHandler := new(startsFakes.FakeHandler)

				objectUnderTest := _handler{
					startsHandler: fakeStartsHandler,
				}

				providedPath := "starts/dummy"
				providedHTTPReq, err := http.NewRequest("dummyMethod", providedPath, nil)
				if nil != err {
					panic(err.Error())
				}

				expectedURLPath := strings.SplitN(providedPath, "/", 2)[1]

				/* act */
				objectUnderTest.Handle(httptest.NewRecorder(), providedHTTPReq)

				/* assert */
				_, actualHTTPReq := fakeStartsHandler.HandleArgsForCall(0)

				Expect(actualHTTPReq.URL.Path).To(Equal(expectedURLPath))

				// this works because our URL path set mutates the httpRequest
				Expect(actualHTTPReq).To(Equal(providedHTTPReq))
			})
		})
		Context("next URL path segment is kills", func() {
			It("should call refHandler.Handle w/ expected args", func() {
				/* arrange */
				fakeKillsHandler := new(killsFakes.FakeHandler)

				objectUnderTest := _handler{
					killsHandler: fakeKillsHandler,
				}

				providedPath := "kills/dummy"
				providedHTTPReq, err := http.NewRequest("dummyMethod", providedPath, nil)
				if nil != err {
					panic(err.Error())
				}

				expectedURLPath := strings.SplitN(providedPath, "/", 2)[1]

				/* act */
				objectUnderTest.Handle(httptest.NewRecorder(), providedHTTPReq)

				/* assert */
				_, actualHTTPReq := fakeKillsHandler.HandleArgsForCall(0)

				Expect(actualHTTPReq.URL.Path).To(Equal(expectedURLPath))

				// this works because our URL path set mutates the httpRequest
				Expect(actualHTTPReq).To(Equal(providedHTTPReq))
			})
		})
	})
})
