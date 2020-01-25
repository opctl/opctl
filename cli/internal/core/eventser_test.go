package core

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	cliModel "github.com/opctl/opctl/cli/internal/model"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

var _ = Context("Eventser", func() {
	Context("Events", func() {
		It("should call nodeHandle.APIClient().GetEventStream w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()

			fakeAPIClient := new(client.Fake)
			fakeNodeHandle := new(cliModel.FakeNodeHandle)
			fakeNodeHandle.APIClientReturns(fakeAPIClient)

			fakeNodeProvider := new(nodeprovider.Fake)
			fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

			fakeCliExiter := new(cliexiter.Fake)

			eventChannel := make(chan model.Event)
			close(eventChannel)
			fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

			objectUnderTest := _eventser{
				cliExiter:    fakeCliExiter,
				nodeProvider: fakeNodeProvider,
			}

			/* act */
			objectUnderTest.Events(
				providedCtx,
			)

			/* assert */
			actualCtx,
				actualGetEventStreamReq := fakeAPIClient.GetEventStreamArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(*actualGetEventStreamReq).To(Equal(model.GetEventStreamReq{}))

		})
		Context("client.GetEventStream errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeCliExiter := new(cliexiter.Fake)
				returnedError := errors.New("dummyError")

				fakeAPIClient := new(client.Fake)
				fakeNodeHandle := new(cliModel.FakeNodeHandle)
				fakeNodeHandle.APIClientReturns(fakeAPIClient)

				fakeNodeProvider := new(nodeprovider.Fake)
				fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

				fakeAPIClient.GetEventStreamReturns(nil, returnedError)

				objectUnderTest := _eventser{
					cliExiter:    fakeCliExiter,
					nodeProvider: fakeNodeProvider,
				}

				/* act */
				objectUnderTest.Events(
					context.Background(),
				)

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
					fakeNodeHandle := new(cliModel.FakeNodeHandle)
					fakeNodeHandle.APIClientReturns(fakeAPIClient)

					fakeNodeProvider := new(nodeprovider.Fake)
					fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _eventser{
						cliExiter:    fakeCliExiter,
						nodeProvider: fakeNodeProvider,
					}

					/* act */
					objectUnderTest.Events(
						context.Background(),
					)

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: "Connection to event stream lost", Code: 1}))
				})
			})
		})
	})
})
