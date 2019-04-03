package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/api/client"
)

var _ = Context("events", func() {
	Context("Execute", func() {
		It("should call client.GetEventStream", func() {
			/* arrange */
			fakeCliExiter := new(cliexiter.Fake)

			fakeAPIClient := new(client.Fake)
			eventChannel := make(chan model.Event)
			close(eventChannel)
			fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

			objectUnderTest := _core{
				apiClient:               fakeAPIClient,
				cliExiter:               fakeCliExiter,
				nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
			}

			/* act */
			objectUnderTest.Events()

			/* assert */
			Expect(fakeAPIClient.GetEventStreamCallCount()).To(Equal(1))

		})
		Context("client.GetEventStream errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeCliExiter := new(cliexiter.Fake)
				returnedError := errors.New("dummyError")

				fakeAPIClient := new(client.Fake)
				fakeAPIClient.GetEventStreamReturns(nil, returnedError)

				objectUnderTest := _core{
					apiClient:               fakeAPIClient,
					cliExiter:               fakeCliExiter,
					nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
				}

				/* act */
				objectUnderTest.Events()

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
			})
		})
		Context("client.GetEventStream doesn't error", func() {
			Context("channel closes unexpectedly", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeCliExiter := new(cliexiter.Fake)

					fakeAPIClient := new(client.Fake)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _core{
						apiClient:               fakeAPIClient,
						cliExiter:               fakeCliExiter,
						nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
					}

					/* act */
					objectUnderTest.Events()

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: "Connection to event stream lost", Code: 1}))
				})
			})
		})
	})
})
