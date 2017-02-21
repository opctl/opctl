package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/pkg/nodeprovider"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg/apiclient"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Context("streamEvents", func() {
	Context("Execute", func() {
		It("should call bundle.GetEventStream", func() {
			/* arrange */
			fakeCliExiter := new(cliexiter.Fake)

			fakeApiClient := new(apiclient.Fake)
			eventChannel := make(chan model.Event)
			close(eventChannel)
			fakeApiClient.GetEventStreamReturns(eventChannel, nil)

			objectUnderTest := _core{
				apiClient:    fakeApiClient,
				cliExiter:    fakeCliExiter,
				nodeProvider: new(nodeprovider.Fake),
			}

			/* act */
			objectUnderTest.StreamEvents()

			/* assert */
			Expect(fakeApiClient.GetEventStreamCallCount()).Should(Equal(1))

		})
		Context("bundle.GetEventStream errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeCliExiter := new(cliexiter.Fake)
				returnedError := errors.New("dummyError")

				fakeApiClient := new(apiclient.Fake)
				fakeApiClient.GetEventStreamReturns(nil, returnedError)

				objectUnderTest := _core{
					apiClient:    fakeApiClient,
					cliExiter:    fakeCliExiter,
					nodeProvider: new(nodeprovider.Fake),
				}

				/* act */
				objectUnderTest.StreamEvents()

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
			})
		})
		Context("bundle.GetEventStream doesn't error", func() {
			Context("channel closes unexpectedly", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeCliExiter := new(cliexiter.Fake)

					fakeApiClient := new(apiclient.Fake)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeApiClient.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _core{
						apiClient:    fakeApiClient,
						cliExiter:    fakeCliExiter,
						nodeProvider: new(nodeprovider.Fake),
					}

					/* act */
					objectUnderTest.StreamEvents()

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: "Connection to event stream lost", Code: 1}))
				})
			})
		})
	})
})
