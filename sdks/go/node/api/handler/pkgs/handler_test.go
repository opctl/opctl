package pkgs

import (
	. "github.com/opctl/opctl/sdks/go/node/core/fakes"
	"net/http"
	"net/http/httptest"
	"strings"

	refFakes "github.com/opctl/opctl/sdks/go/node/api/handler/pkgs/ref/fakes"

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
				providedPath := "ref/dummy"
				providedPathParts := strings.SplitN(providedPath, "/", 2)

				providedHTTPRes := httptest.NewRecorder()
				providedHTTPReq, err := http.NewRequest("dummyMethod", providedPath, nil)
				if nil != err {
					panic(err.Error())
				}

				fakeRefHandler := new(refFakes.FakeHandler)

				objectUnderTest := _handler{
					refHandler: fakeRefHandler,
				}

				expectedDataRef := providedPathParts[0]
				expectedURLPath := providedPathParts[1]

				/* act */
				objectUnderTest.Handle(providedHTTPRes, providedHTTPReq)

				/* assert */
				actualDataRef,
					actualHTTPRes,
					actualHTTPReq := fakeRefHandler.HandleArgsForCall(0)

				Expect(actualDataRef).To(Equal(expectedDataRef))
				Expect(actualHTTPRes).To(Equal(providedHTTPRes))
				Expect(actualHTTPReq.URL.Path).To(Equal(expectedURLPath))

				// this works because our URL path set mutates the httpRequest
				Expect(actualHTTPReq).To(Equal(providedHTTPReq))
			})
		})
	})
})
