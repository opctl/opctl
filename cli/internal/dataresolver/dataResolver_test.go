package dataresolver

import (
	"errors"
	"fmt"

	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cliparamsatisfierFakes "github.com/opctl/opctl/cli/internal/cliparamsatisfier/fakes"
	modelFakes "github.com/opctl/opctl/cli/internal/model/fakes"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
	clientFakes "github.com/opctl/opctl/sdks/go/node/api/client/fakes"
)

var _ = Context("dataResolver", func() {
	Context("Resolve", func() {
		Context("data.Resolve errs", func() {
			Context("data.ErrDataProviderAuthorization", func() {
				It("should call cliParamSatisfier.Satisfy w/ expected args", func() {
					/* arrange */
					fakeAPIClient := new(clientFakes.FakeClient)

					fakeAPIClient.ListDescendantsReturnsOnCall(0, nil, model.ErrDataProviderAuthorization{})
					fakeAPIClient.ListDescendantsReturnsOnCall(1, nil, errors.New(""))

					fakeNodeHandle := new(modelFakes.FakeNodeHandle)
					fakeNodeHandle.APIClientReturns(fakeAPIClient)

					fakeNodeProvider := new(nodeprovider.Fake)
					fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

					username := "dummyUsername"
					password := "dummyPassword"

					fakeCliParamSatisfier := new(cliparamsatisfierFakes.FakeCLIParamSatisfier)
					fakeCliParamSatisfier.SatisfyReturns(
						map[string]*model.Value{
							usernameInputName: {String: &username},
							passwordInputName: {String: &password},
						},
						nil,
					)

					objectUnderTest := _dataResolver{
						cliParamSatisfier: fakeCliParamSatisfier,
						os:                new(ios.Fake),
						nodeProvider:      fakeNodeProvider,
					}

					/* act */
					objectUnderTest.Resolve("ref", &model.Creds{})

					/* assert */
					_, actualInputs := fakeCliParamSatisfier.SatisfyArgsForCall(0)
					Expect(actualInputs).To(Equal(credsPromptInputs))
				})
			})
			Context("not data.ErrAuthenticationFailed", func() {
				It("should return expected error", func() {
					/* arrange */
					providedDataRef := "dummyDataRef"

					expectedErr := "expectedErr"
					fakeAPIClient := new(clientFakes.FakeClient)
					fakeAPIClient.ListDescendantsReturns(nil, errors.New(expectedErr))

					fakeNodeHandle := new(modelFakes.FakeNodeHandle)
					fakeNodeHandle.APIClientReturns(fakeAPIClient)

					fakeNodeProvider := new(nodeprovider.Fake)
					fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

					objectUnderTest := _dataResolver{
						os:           new(ios.Fake),
						nodeProvider: fakeNodeProvider,
					}

					/* act */
					response, err := objectUnderTest.Resolve(providedDataRef, &model.Creds{})

					/* assert */
					Expect(response).To(BeNil())
					Expect(err.Error()).To(Equal(fmt.Sprintf("Unable to resolve pkg 'dummyDataRef'; error was %s", expectedErr)))
				})
			})
		})
		Context("data.Resolve doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeAPIClient := new(clientFakes.FakeClient)
				fakeNodeHandle := new(modelFakes.FakeNodeHandle)
				fakeNodeHandle.APIClientReturns(fakeAPIClient)

				fakeNodeProvider := new(nodeprovider.Fake)
				fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)

				objectUnderTest := _dataResolver{
					os:           new(ios.Fake),
					nodeProvider: fakeNodeProvider,
				}

				/* act */
				actualPkgHandle, err := objectUnderTest.Resolve(
					"testdata/dummy-op",
					&model.Creds{},
				)

				/* assert */
				Expect(err).To(BeNil())
				Expect(actualPkgHandle.Ref()).To(Equal("testdata/dummy-op"))
			})
		})
	})
})
