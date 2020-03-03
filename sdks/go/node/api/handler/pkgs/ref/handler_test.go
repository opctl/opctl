package ref

import (
	"errors"
	"fmt"
	contentsFakes "github.com/opctl/opctl/sdks/go/node/api/handler/pkgs/ref/contents/fakes"
	"net/http"
	"net/http/httptest"
	"strings"

	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"

	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/fakes"

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
				objectUnderTest.Handle(
					"dummyDataRef",
					providedHTTPResp,
					providedHTTPReq,
				)

				/* assert */
				Expect(providedHTTPResp.Code).To(Equal(http.StatusNotFound))
			})
		})
		Context("next URL path segment is 'contents'", func() {
			Context("req has BasicAuth", func() {
				It("should call core.ResolveData w/ expected args", func() {
					/* arrange */
					providedDataRef := "dummyDataRef"
					providedUsername := "dummyUsername"
					providedPassword := "dummyPassword"

					fakeCore := new(FakeCore)
					// err to trigger immediate return
					fakeCore.ResolveDataReturns(nil, errors.New("dummyErr"))

					objectUnderTest := _handler{
						core: fakeCore,
					}

					providedHTTPReq, err := http.NewRequest("dummyHttpMethod", "contents", nil)
					if nil != err {
						panic(err.Error())
					}

					providedHTTPReq.SetBasicAuth(
						providedUsername,
						providedPassword,
					)

					expectedPullCreds := model.PullCreds{
						Username: providedUsername,
						Password: providedPassword,
					}

					/* act */
					objectUnderTest.Handle(
						providedDataRef,
						httptest.NewRecorder(),
						providedHTTPReq,
					)

					/* assert */
					_,
						actualRef,
						actualPullCreds := fakeCore.ResolveDataArgsForCall(0)

					Expect(actualRef).To(Equal(providedDataRef))
					Expect(*actualPullCreds).To(Equal(expectedPullCreds))
				})
			})
			It("should call core.ResolveData w/ expected args", func() {
				/* arrange */
				providedDataRef := "dummyDataRef"

				fakeCore := new(FakeCore)
				// err to trigger immediate return
				fakeCore.ResolveDataReturns(nil, errors.New("dummyErr"))

				objectUnderTest := _handler{
					core: fakeCore,
				}

				providedHTTPReq, err := http.NewRequest("dummyHttpMethod", "contents", nil)
				if nil != err {
					panic(err.Error())
				}

				/* act */
				objectUnderTest.Handle(
					providedDataRef,
					httptest.NewRecorder(),
					providedHTTPReq,
				)

				/* assert */
				_,
					actualRef,
					actualPullCreds := fakeCore.ResolveDataArgsForCall(0)

				Expect(actualRef).To(Equal(providedDataRef))
				Expect(actualPullCreds).To(BeNil())
			})
			Context("core.ResolveData errs", func() {
				Context("err is ErrDataProviderAuthentication", func() {
					It("should return expected result", func() {
						/* arrange */
						dataRefSegment1 := "dataRefSegment1"
						providedDataRef := strings.Join([]string{dataRefSegment1, "dataRefSegment2"}, "/")

						fakeCore := new(FakeCore)
						fakeCore.ResolveDataReturns(nil, model.ErrDataProviderAuthentication{})

						objectUnderTest := _handler{
							core: fakeCore,
						}
						providedHTTPResp := httptest.NewRecorder()

						providedHTTPReq, err := http.NewRequest("dummyHttpMethod", "contents", nil)
						if nil != err {
							panic(err.Error())
						}

						/* act */
						objectUnderTest.Handle(
							providedDataRef,
							providedHTTPResp,
							providedHTTPReq,
						)

						/* assert */
						Expect(providedHTTPResp.Code).To(Equal(http.StatusUnauthorized))
						Expect(providedHTTPResp.Header().Get("WWW-Authenticate")).To(Equal(fmt.Sprintf(`Basic realm="%s"`, dataRefSegment1)))
					})
				})
				Context("err is ErrDataProviderAuthorization", func() {
					It("should return expected result", func() {
						/* arrange */
						dataRefSegment1 := "dataRefSegment1"
						providedDataRef := strings.Join([]string{dataRefSegment1, "dataRefSegment2"}, "/")

						fakeCore := new(FakeCore)
						fakeCore.ResolveDataReturns(nil, model.ErrDataProviderAuthorization{})

						objectUnderTest := _handler{
							core: fakeCore,
						}
						providedHTTPResp := httptest.NewRecorder()

						providedHTTPReq, err := http.NewRequest("dummyHttpMethod", "contents", nil)
						if nil != err {
							panic(err.Error())
						}

						/* act */
						objectUnderTest.Handle(
							providedDataRef,
							providedHTTPResp,
							providedHTTPReq,
						)

						/* assert */
						Expect(providedHTTPResp.Code).To(Equal(http.StatusForbidden))
						Expect(providedHTTPResp.Header().Get("WWW-Authenticate")).To(Equal(fmt.Sprintf(`Basic realm="%s"`, dataRefSegment1)))
					})
				})
				Context("err is ErrDataRefResolution", func() {
					It("should return expected result", func() {
						/* arrange */
						fakeCore := new(FakeCore)
						fakeCore.ResolveDataReturns(nil, model.ErrDataRefResolution{})

						objectUnderTest := _handler{
							core: fakeCore,
						}
						providedHTTPResp := httptest.NewRecorder()

						providedHTTPReq, err := http.NewRequest("dummyHttpMethod", "contents", nil)
						if nil != err {
							panic(err.Error())
						}

						/* act */
						objectUnderTest.Handle(
							"dummyDataRef",
							providedHTTPResp,
							providedHTTPReq,
						)

						/* assert */
						Expect(providedHTTPResp.Code).To(Equal(http.StatusNotFound))
					})
				})
			})
			Context("core.ResolveData doesn't err", func() {
				It("should call handleGetOrHeader.HandleGetOrHead w/ expected args", func() {
					/* arrange */
					providedDataRef := "dummyDataRef"

					fakeCore := new(FakeCore)
					fakeDataHandle := new(modelFakes.FakeDataHandle)
					fakeCore.ResolveDataReturns(fakeDataHandle, nil)

					fakeContentsHandler := new(contentsFakes.FakeHandler)

					objectUnderTest := _handler{
						core:            fakeCore,
						contentsHandler: fakeContentsHandler,
					}

					providedHTTPReq, err := http.NewRequest("dummyHttpMethod", "contents", nil)
					if nil != err {
						panic(err.Error())
					}

					/* act */
					objectUnderTest.Handle(
						providedDataRef,
						httptest.NewRecorder(),
						providedHTTPReq,
					)

					/* assert */
					actualDataHandle,
						_,
						actualHTTPReq := fakeContentsHandler.HandleArgsForCall(0)

					Expect(actualDataHandle).To(Equal(fakeDataHandle))
					Expect(actualHTTPReq.URL.Path).To(BeEmpty())

					// this works because our URL path set mutates the request
					Expect(actualHTTPReq).To(Equal(providedHTTPReq))
				})
			})
		})
	})
})
