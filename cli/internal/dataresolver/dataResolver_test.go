package dataresolver

import (
	"errors"
	"fmt"

	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cliparamsatisfierFakes "github.com/opctl/opctl/cli/internal/cliparamsatisfier/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
)

var _ = Context("dataResolver", func() {
	It("Can be constructed", func() {
		Expect(New(
			new(cliparamsatisfierFakes.FakeCLIParamSatisfier),
			new(nodeFakes.FakeOpNode),
		)).NotTo(BeNil())
	})
	Context("Resolve", func() {
		Context("data.Resolve errs", func() {
			Context("data.ErrDataProviderAuthorization", func() {
				It("should call cliParamSatisfier.Satisfy w/ expected args", func() {
					/* arrange */
					fakeCore := new(nodeFakes.FakeOpNode)

					fakeCore.ListDescendantsReturnsOnCall(0, nil, model.ErrDataProviderAuthorization{})
					fakeCore.ListDescendantsReturnsOnCall(1, nil, errors.New(""))

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
						opNode:            fakeCore,
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
					fakeCore := new(nodeFakes.FakeOpNode)
					fakeCore.ListDescendantsReturns(nil, errors.New(expectedErr))

					objectUnderTest := _dataResolver{
						os:     new(ios.Fake),
						opNode: fakeCore,
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
				fakeCore := new(nodeFakes.FakeOpNode)

				objectUnderTest := _dataResolver{
					os:     new(ios.Fake),
					opNode: fakeCore,
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
