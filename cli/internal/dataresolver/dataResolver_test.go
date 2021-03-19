package dataresolver

import (
	"context"
	"errors"
	"os"
	"path"

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
			new(nodeFakes.FakeNode),
		)).NotTo(BeNil())
	})
	Context("Resolve", func() {
		Context("data.Resolve errs", func() {
			Context("data.ErrDataProviderAuthorization", func() {
				It("should call cliParamSatisfier.Satisfy w/ expected args", func() {
					/* arrange */
					fakeCore := new(nodeFakes.FakeNode)

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

					expectedInputs := map[string]*model.Param{
						usernameInputName: {
							String: &model.StringParam{
								Description: "Username for github.com.",
								Constraints: map[string]interface{}{
									"MinLength": 1,
								},
							},
						},
						passwordInputName: {
							String: &model.StringParam{
								Description: "Personal access token for github.com with 'Repo' permissions.",
								Constraints: map[string]interface{}{
									"MinLength": 1,
								},
								IsSecret: true,
							},
						},
					}

					objectUnderTest := _dataResolver{
						cliParamSatisfier: fakeCliParamSatisfier,
						node:              fakeCore,
					}

					/* act */
					objectUnderTest.Resolve(context.Background(), "github.com/opctl/opctl/.opspec/build", &model.Creds{})

					/* assert */
					_, actualInputs := fakeCliParamSatisfier.SatisfyArgsForCall(0)
					Expect(actualInputs).To(BeEquivalentTo(expectedInputs))
				})
			})
			Context("not data.ErrAuthenticationFailed", func() {
				It("should return expected error", func() {
					/* arrange */
					providedDataRef := "dummyDataRef"

					fakeCore := new(nodeFakes.FakeNode)
					fakeCore.ListDescendantsReturns(nil, errors.New("expectedErr"))

					objectUnderTest := _dataResolver{
						node: fakeCore,
					}

					/* act */
					response, err := objectUnderTest.Resolve(context.Background(), providedDataRef, &model.Creds{})

					/* assert */
					Expect(response).To(BeNil())
					Expect(err).To(MatchError(`unable to resolve op 'dummyDataRef':` + " " + `
- filesystem:` + " " + `
  - path /src/cli/internal/dataresolver/.opspec/dummyDataRef not found
  - path /src/cli/internal/dataresolver/dummyDataRef not found
- opctl node: expectedErr`))
				})
			})
		})
		Context("data.Resolve doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				wd, err := os.Getwd()
				if err != nil {
					panic(err)
				}

				fakeCore := new(nodeFakes.FakeNode)
				providedDataRef := "testdata/dummy-op"

				objectUnderTest := _dataResolver{
					node: fakeCore,
				}

				/* act */
				actualPkgHandle, err := objectUnderTest.Resolve(
					context.Background(),
					providedDataRef,
					&model.Creds{},
				)

				/* assert */
				Expect(err).To(BeNil())
				Expect(actualPkgHandle.Ref()).To(Equal(path.Join(wd, providedDataRef)))
			})
		})
	})
})
