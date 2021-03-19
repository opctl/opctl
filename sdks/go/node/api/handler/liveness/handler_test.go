package liveness

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("Handler", func() {
	Context("Handle", func() {
		Context("next URL path segment is empty", func() {
			It("should return expected result", func() {
				/* arrange */

				providedHTTPResp := httptest.NewRecorder()

				objectUnderTest := _handler{}

				providedHTTPReq, err := http.NewRequest("dummyMethod", "", nil)
				if err != nil {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.Handle(providedHTTPResp, providedHTTPReq)

				/* assert */
				Expect(providedHTTPResp.Code).To(Equal(http.StatusOK))
			})
		})
		Context("next URL path segment not empty", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _handler{}
				providedHTTPResp := httptest.NewRecorder()

				providedHTTPReq, err := http.NewRequest("dummyMethod", "blah", nil)
				if err != nil {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.Handle(providedHTTPResp, providedHTTPReq)

				/* assert */
				Expect(providedHTTPResp.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
