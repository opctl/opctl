package pkgs

import (
	"github.com/opspec-io/sdk-golang/node/core"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/opspec-io/sdk-golang/node/api/handler/pkgs/ref"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("Handler", func() {
	Context("NewHandler", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewHandler(new(core.Fake))).Should(Not(BeNil()))
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
				objectUnderTest.Handle(providedHTTPResp, providedHTTPReq)

				/* assert */
				Expect(providedHTTPResp.Code).To(Equal(http.StatusNotFound))
			})
		})
		Context("next URL path segment isnt empty", func() {
			It("should call refHandler.Handle w/ expected args", func() {
				/* arrange */
				fakeRefHandler := new(ref.FakeHandler)

				objectUnderTest := _handler{
					refHandler: fakeRefHandler,
				}

				providedPath := "ref/dummy"
				providedHTTPReq, err := http.NewRequest("dummyMethod", providedPath, nil)
				if nil != err {
					panic(err.Error())
				}

				expectedURLPath := strings.SplitN(providedPath, "/", 2)[1]

				/* act */
				objectUnderTest.Handle(httptest.NewRecorder(), providedHTTPReq)

				/* assert */
				_, actualHTTPReq := fakeRefHandler.HandleArgsForCall(0)

				Expect(actualHTTPReq.URL.Path).To(Equal(expectedURLPath))

				// this works because our URL path set mutates the httpRequest
				Expect(actualHTTPReq).To(Equal(providedHTTPReq))
			})
		})
	})
})
