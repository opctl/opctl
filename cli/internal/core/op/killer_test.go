package op

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

var _ = Context("Killer", func() {
	Context("Invoke", func() {
		It("should call apiClient.Invoke w/ expected args", func() {
			/* arrange */
			fakeAPIClient := new(clientFakes.FakeClient)
			fakeNodeHandle := new(modelFakes.FakeNodeHandle)
			fakeNodeHandle.APIClientReturns(fakeAPIClient)

			fakeNodeProvider := new(nodeprovider.Fake)
			fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

			providedCtx := context.TODO()

			expectedCtx := providedCtx
			expectedReq := model.KillOpReq{
				OpID:       "dummyOpID",
				RootCallID: "dummyOpID",
			}

			objectUnderTest := _killer{
				nodeProvider: fakeNodeProvider,
			}

			/* act */
			err := objectUnderTest.Kill(expectedCtx, expectedReq.OpID)

			/* assert */
			actualCtx, actualReq := fakeAPIClient.KillOpArgsForCall(0)
			Expect(err).To(BeNil())
			Expect(actualCtx).To(Equal(expectedCtx))
			Expect(actualReq).To(BeEquivalentTo(expectedReq))
		})
		Context("apiClient.Invoke errors", func() {
			It("should return expected error", func() {
				/* arrange */
				fakeAPIClient := new(clientFakes.FakeClient)
				expectedError := errors.New("dummyError")
				fakeAPIClient.KillOpReturns(expectedError)

				fakeNodeHandle := new(modelFakes.FakeNodeHandle)
				fakeNodeHandle.APIClientReturns(fakeAPIClient)

				fakeNodeProvider := new(nodeprovider.Fake)
				fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

				objectUnderTest := _killer{
					nodeProvider: fakeNodeProvider,
				}

				/* act */
				err := objectUnderTest.Kill(context.TODO(), "")

				/* assert */
				Expect(err).To(MatchError(expectedError))
			})
		})
	})
})
