package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/opctl/opctl/sdk/go/node/api/handler/data"
	"github.com/opctl/opctl/sdk/go/node/api/handler/events"
	"github.com/opctl/opctl/sdk/go/node/api/handler/liveness"
	"github.com/opctl/opctl/sdk/go/node/api/handler/ops"
	"github.com/opctl/opctl/sdk/go/node/api/handler/pkgs"
	"github.com/opctl/opctl/sdk/go/node/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("Handler", func() {
	Context("New", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(New(new(core.Fake))).Should(Not(BeNil()))
		})
	})
	Context("Handle", func() {
		Context("next URL path segment is empty", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _handler{}
				providedHTTPResp := httptest.NewRecorder()

				providedHTTPReq, err := http.NewRequest("dummyMethod", "", nil)
				if nil != err {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.ServeHTTP(providedHTTPResp, providedHTTPReq)

				/* assert */
				Expect(providedHTTPResp.Code).To(Equal(http.StatusNotFound))
			})
		})
		Context("next URL path segment is data", func() {
			It("should call dataHandler.Handle w/ expected args", func() {
				/* arrange */
				fakeDataHandler := new(data.FakeHandler)

				objectUnderTest := _handler{
					dataHandler: fakeDataHandler,
				}

				providedPath := "data/dummy"
				providedHTTPReq, err := http.NewRequest("dummyMethod", providedPath, nil)
				if nil != err {
					panic(err.Error())
				}

				expectedURLPath := strings.SplitN(providedPath, "/", 2)[1]

				/* act */
				objectUnderTest.ServeHTTP(httptest.NewRecorder(), providedHTTPReq)

				/* assert */
				_, actualHTTPReq := fakeDataHandler.HandleArgsForCall(0)

				Expect(actualHTTPReq.URL.Path).To(Equal(expectedURLPath))

				// this works because our URL path set mutates the httpRequest
				Expect(actualHTTPReq).To(Equal(providedHTTPReq))
			})
		})
		Context("next URL path segment is events", func() {
			It("should call eventsHandler.Handle w/ expected args", func() {
				/* arrange */
				fakeEventsHandler := new(events.FakeHandler)

				objectUnderTest := _handler{
					eventsHandler: fakeEventsHandler,
				}

				providedPath := "events/dummy"
				providedHTTPReq, err := http.NewRequest("dummyMethod", providedPath, nil)
				if nil != err {
					panic(err.Error())
				}

				expectedURLPath := strings.SplitN(providedPath, "/", 2)[1]

				/* act */
				objectUnderTest.ServeHTTP(httptest.NewRecorder(), providedHTTPReq)

				/* assert */
				_, actualHTTPReq := fakeEventsHandler.HandleArgsForCall(0)

				Expect(actualHTTPReq.URL.Path).To(Equal(expectedURLPath))

				// this works because our URL path set mutates the httpRequest
				Expect(actualHTTPReq).To(Equal(providedHTTPReq))
			})
		})
		Context("next URL path segment is liveness", func() {
			It("should call livenessHandler.Handle w/ expected args", func() {
				/* arrange */
				fakeLivenessHandler := new(liveness.FakeHandler)

				objectUnderTest := _handler{
					livenessHandler: fakeLivenessHandler,
				}

				providedPath := "liveness/dummy"
				providedHTTPReq, err := http.NewRequest("dummyMethod", providedPath, nil)
				if nil != err {
					panic(err.Error())
				}

				expectedURLPath := strings.SplitN(providedPath, "/", 2)[1]

				/* act */
				objectUnderTest.ServeHTTP(httptest.NewRecorder(), providedHTTPReq)

				/* assert */
				_, actualHTTPReq := fakeLivenessHandler.HandleArgsForCall(0)

				Expect(actualHTTPReq.URL.Path).To(Equal(expectedURLPath))

				// this works because our URL path set mutates the httpRequest
				Expect(actualHTTPReq).To(Equal(providedHTTPReq))
			})
		})
		Context("next URL path segment is ops", func() {
			It("should call livenessHandler.Handle w/ expected args", func() {
				/* arrange */
				fakeOpsHandler := new(ops.FakeHandler)

				objectUnderTest := _handler{
					opsHandler: fakeOpsHandler,
				}

				providedPath := "ops/dummy"
				providedHTTPReq, err := http.NewRequest("dummyMethod", providedPath, nil)
				if nil != err {
					panic(err.Error())
				}

				expectedURLPath := strings.SplitN(providedPath, "/", 2)[1]

				/* act */
				objectUnderTest.ServeHTTP(httptest.NewRecorder(), providedHTTPReq)

				/* assert */
				_, actualHTTPReq := fakeOpsHandler.HandleArgsForCall(0)

				Expect(actualHTTPReq.URL.Path).To(Equal(expectedURLPath))

				// this works because our URL path set mutates the httpRequest
				Expect(actualHTTPReq).To(Equal(providedHTTPReq))
			})
		})
		Context("next URL path segment is pkgs", func() {
			It("should call livenessHandler.Handle w/ expected args", func() {
				/* arrange */
				fakePkgsHandler := new(pkgs.FakeHandler)

				objectUnderTest := _handler{
					pkgsHandler: fakePkgsHandler,
				}

				providedPath := "pkgs/dummy"
				providedHTTPReq, err := http.NewRequest("dummyMethod", providedPath, nil)
				if nil != err {
					panic(err.Error())
				}

				expectedURLPath := strings.SplitN(providedPath, "/", 2)[1]

				/* act */
				objectUnderTest.ServeHTTP(httptest.NewRecorder(), providedHTTPReq)

				/* assert */
				_, actualHTTPReq := fakePkgsHandler.HandleArgsForCall(0)

				Expect(actualHTTPReq.URL.Path).To(Equal(expectedURLPath))

				// this works because our URL path set mutates the httpRequest
				Expect(actualHTTPReq).To(Equal(providedHTTPReq))
			})
		})
	})
})
