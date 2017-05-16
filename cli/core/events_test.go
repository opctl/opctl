package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/nodeprovider"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api/client"
)

var _ = Context("events", func() {
	Context("Execute", func() {
		It("should call client.GetEventStream", func() {
			/* arrange */
			fakeCliExiter := new(cliexiter.Fake)

			fakeOpspecNodeAPIClient := new(client.Fake)
			eventChannel := make(chan model.Event)
			close(eventChannel)
			fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

			objectUnderTest := _core{
				opspecNodeAPIClient: fakeOpspecNodeAPIClient,
				cliExiter:           fakeCliExiter,
				nodeProvider:        new(nodeprovider.Fake),
			}

			/* act */
			objectUnderTest.Events()

			/* assert */
			Expect(fakeOpspecNodeAPIClient.GetEventStreamCallCount()).To(Equal(1))

		})
		Context("client.GetEventStream errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeCliExiter := new(cliexiter.Fake)
				returnedError := errors.New("dummyError")

				fakeOpspecNodeAPIClient := new(client.Fake)
				fakeOpspecNodeAPIClient.GetEventStreamReturns(nil, returnedError)

				objectUnderTest := _core{
					opspecNodeAPIClient: fakeOpspecNodeAPIClient,
					cliExiter:           fakeCliExiter,
					nodeProvider:        new(nodeprovider.Fake),
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

					fakeOpspecNodeAPIClient := new(client.Fake)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeOpspecNodeAPIClient.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _core{
						opspecNodeAPIClient: fakeOpspecNodeAPIClient,
						cliExiter:           fakeCliExiter,
						nodeProvider:        new(nodeprovider.Fake),
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
