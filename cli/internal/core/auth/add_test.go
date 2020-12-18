package auth

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

var _ = Context("Adder", func() {
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
			expectedReq := model.AddAuthReq{
				Resources: "Resources",
				Creds: model.Creds{
					Username: "username",
					Password: "password",
				},
			}

			objectUnderTest := _adder{
				nodeProvider: fakeNodeProvider,
			}

			/* act */
			err := objectUnderTest.Add(
				expectedCtx,
				expectedReq.Resources,
				expectedReq.Username,
				expectedReq.Password,
			)

			/* assert */
			actualCtx, actualReq := fakeAPIClient.AddAuthArgsForCall(0)
			Expect(err).To(BeNil())
			Expect(actualCtx).To(Equal(expectedCtx))
			Expect(actualReq).To(BeEquivalentTo(expectedReq))
		})
		Context("apiClient.Invoke errors", func() {
			It("should return expected error", func() {
				/* arrange */
				fakeAPIClient := new(clientFakes.FakeClient)
				expectedError := errors.New("dummyError")
				fakeAPIClient.AddAuthReturns(expectedError)

				fakeNodeHandle := new(modelFakes.FakeNodeHandle)
				fakeNodeHandle.APIClientReturns(fakeAPIClient)

				fakeNodeProvider := new(nodeprovider.Fake)
				fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

				objectUnderTest := _adder{
					nodeProvider: fakeNodeProvider,
				}

				/* act */
				err := objectUnderTest.Add(context.TODO(), "", "", "")

				/* assert */
				Expect(err).To(MatchError(expectedError))
			})
		})
	})
})
