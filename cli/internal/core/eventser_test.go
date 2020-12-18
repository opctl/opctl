package core

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	modelFakes "github.com/opctl/opctl/cli/internal/model/fakes"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
	clientFakes "github.com/opctl/opctl/sdks/go/node/api/client/fakes"
)

var _ = Context("Eventser", func() {
	Context("Events", func() {
		It("should call nodeHandle.APIClient().GetEventStream w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()

			fakeAPIClient := new(clientFakes.FakeClient)
			fakeNodeHandle := new(modelFakes.FakeNodeHandle)
			fakeNodeHandle.APIClientReturns(fakeAPIClient)

			fakeNodeProvider := new(nodeprovider.Fake)
			fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

			eventChannel := make(chan model.Event)
			close(eventChannel)
			fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

			objectUnderTest := _eventser{
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
			It("should return expected error", func() {
				/* arrange */
				returnedError := errors.New("dummyError")

				fakeAPIClient := new(clientFakes.FakeClient)
				fakeNodeHandle := new(modelFakes.FakeNodeHandle)
				fakeNodeHandle.APIClientReturns(fakeAPIClient)

				fakeNodeProvider := new(nodeprovider.Fake)
				fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

				fakeAPIClient.GetEventStreamReturns(nil, returnedError)

				objectUnderTest := _eventser{
					nodeProvider: fakeNodeProvider,
				}

				/* act */
				err := objectUnderTest.Events(context.Background())

				/* assert */
				Expect(err).To(MatchError(returnedError))
			})
		})
		Context("client.GetEventStream doesn't error", func() {
			Context("channel closes unexpectedly", func() {
				It("should return expected error", func() {
					/* arrange */
					fakeAPIClient := new(clientFakes.FakeClient)
					fakeNodeHandle := new(modelFakes.FakeNodeHandle)
					fakeNodeHandle.APIClientReturns(fakeAPIClient)

					fakeNodeProvider := new(nodeprovider.Fake)
					fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeAPIClient.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _eventser{
						nodeProvider: fakeNodeProvider,
					}

					/* act */
					err := objectUnderTest.Events(context.Background())

					/* assert */
					Expect(err).To(MatchError(err))
				})
			})
		})
	})
})
