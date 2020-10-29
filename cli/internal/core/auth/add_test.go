package auth

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	cliexiterFakes "github.com/opctl/opctl/cli/internal/cliexiter/fakes"
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
			objectUnderTest.Add(
				expectedCtx,
				expectedReq.Resources,
				expectedReq.Username,
				expectedReq.Password,
			)

			/* assert */
			actualCtx, actualReq := fakeAPIClient.AddAuthArgsForCall(0)
			Expect(actualCtx).To(Equal(expectedCtx))
			Expect(actualReq).To(BeEquivalentTo(expectedReq))
		})
		Context("apiClient.Invoke errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeAPIClient := new(clientFakes.FakeClient)
				expectedError := errors.New("dummyError")
				fakeAPIClient.AddAuthReturns(expectedError)

				fakeNodeHandle := new(modelFakes.FakeNodeHandle)
				fakeNodeHandle.APIClientReturns(fakeAPIClient)

				fakeNodeProvider := new(nodeprovider.Fake)
				fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

				objectUnderTest := _adder{
					cliExiter:    fakeCliExiter,
					nodeProvider: fakeNodeProvider,
				}

				/* act */
				objectUnderTest.Add(context.TODO(), "", "", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
