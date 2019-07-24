package core

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/util/cliexiter"
	"github.com/opctl/opctl/cli/util/cliparamsatisfier"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/types"
	"net/url"
)

var _ = Context("dataResolver", func() {
	Context("Resolve", func() {
		It("should call nodeReachabilityEnsurer.EnsureNodeReachable", func() {
			/* arrange */
			fakeNodeReachabilityEnsurer := new(fakeNodeReachabilityEnsurer)

			fakeIOS := new(ios.Fake)
			// err to trigger immediate return
			fakeIOS.GetwdReturns("", errors.New("dummyError"))

			objectUnderTest := _dataResolver{
				cliExiter:               new(cliexiter.Fake),
				nodeReachabilityEnsurer: fakeNodeReachabilityEnsurer,
				os: fakeIOS,
			}

			/* act */
			objectUnderTest.Resolve("dummyDataRef", &types.PullCreds{})

			/* assert */
			Expect(fakeNodeReachabilityEnsurer.EnsureNodeReachableCallCount()).To(Equal(1))
		})
		Context("os.Getwd errs", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeIOS := new(ios.Fake)
				expectedError := errors.New("dummyError")
				fakeIOS.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _dataResolver{
					cliExiter: fakeCliExiter,
					os:        fakeIOS,
					nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
				}

				/* act */
				objectUnderTest.Resolve(
					"dummyDataRef",
					nil,
				)

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("os.Getwd doesn't err", func() {
			It("should call data.NewFSProvider w/ expected args", func() {
				/* arrange */
				fakeData := new(data.Fake)
				fakeFSProvider := new(data.FakeProvider)
				fakeData.NewFSProviderReturns(fakeFSProvider)

				// error to trigger immediate return
				fakeData.ResolveReturns(nil, errors.New("dummyError"))

				fakeIOS := new(ios.Fake)
				workDir := "dummyWorkDir"
				fakeIOS.GetwdReturns(workDir, nil)

				objectUnderTest := _dataResolver{
					data:                    fakeData,
					cliExiter:               new(cliexiter.Fake),
					nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
					os: fakeIOS,
				}

				/* act */
				objectUnderTest.Resolve("dummyDataRef", &types.PullCreds{})

				/* assert */
				Expect(fakeData.NewFSProviderArgsForCall(0)).To(ConsistOf(workDir))
			})
			It("should call data.NewNodeProvider w/ expected args", func() {
				/* arrange */
				providedPullCreds := &types.PullCreds{
					Username: "dummyUsername",
					Password: "dummyPassword",
				}

				fakeData := new(data.Fake)
				fakeNodeProvider := new(data.FakeProvider)
				fakeData.NewNodeProviderReturns(fakeNodeProvider)

				// err to trigger immediate return
				fakeData.ResolveReturns(nil, errors.New("dummyError"))

				nodeURL := url.URL{Path: "dummyNodeURL"}

				objectUnderTest := _dataResolver{
					data:                    fakeData,
					cliExiter:               new(cliexiter.Fake),
					nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
					nodeURL:                 nodeURL,
					os:                      new(ios.Fake),
				}

				/* act */
				objectUnderTest.Resolve(
					"dummyDataRef",
					providedPullCreds,
				)

				/* assert */
				actualNodeURL,
					actualPullCreds := fakeData.NewNodeProviderArgsForCall(0)
				Expect(actualNodeURL).To(Equal(nodeURL))
				Expect(actualPullCreds).To(Equal(providedPullCreds))
			})
			It("should call data.Resolve w/ expected args", func() {
				/* arrange */
				providedDataRef := "dummyDataRef"

				fakeData := new(data.Fake)

				fakeFSProvider := new(data.FakeProvider)
				fakeData.NewFSProviderReturns(fakeFSProvider)

				fakeNodeProvider := new(data.FakeProvider)
				fakeData.NewNodeProviderReturns(fakeNodeProvider)

				// err to trigger immediate return
				fakeData.ResolveReturns(nil, errors.New("dummyError"))

				objectUnderTest := _dataResolver{
					data:                    fakeData,
					cliExiter:               new(cliexiter.Fake),
					nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
					os: new(ios.Fake),
				}

				/* act */
				objectUnderTest.Resolve(
					providedDataRef,
					&types.PullCreds{},
				)

				/* assert */
				_,
					actualDataRef,
					actualProviders := fakeData.ResolveArgsForCall(0)

				Expect(actualDataRef).To(Equal(providedDataRef))
				Expect(actualProviders).To(ConsistOf(fakeFSProvider, fakeNodeProvider))
			})
			Context("data.Resolve errs", func() {
				Context("data.ErrDataProviderAuthorization", func() {
					It("should call cliParamSatisfier.Satisfy w/ expected args", func() {
						/* arrange */
						fakeData := new(data.Fake)
						expectedError := types.ErrDataProviderAuthorization{}
						fakeData.ResolveReturnsOnCall(0, nil, expectedError)

						username := "dummyUsername"
						password := "dummyPassword"

						fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
						fakeCliParamSatisfier.SatisfyReturns(
							map[string]*types.Value{
								usernameInputName: {String: &username},
								passwordInputName: {String: &password},
							},
						)

						objectUnderTest := _dataResolver{
							data:                    fakeData,
							cliParamSatisfier:       fakeCliParamSatisfier,
							cliExiter:               new(cliexiter.Fake),
							nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
							os: new(ios.Fake),
						}

						/* act */
						objectUnderTest.Resolve("dummyDataRef", &types.PullCreds{})

						/* assert */
						_, actualInputs := fakeCliParamSatisfier.SatisfyArgsForCall(0)
						Expect(actualInputs).To(Equal(credsPromptInputs))
					})
					It("should call data.NewNodeProvider w/ expected args", func() {
						/* arrange */
						fakeData := new(data.Fake)
						fakeNodeProvider := new(data.FakeProvider)
						fakeData.NewNodeProviderReturns(fakeNodeProvider)

						expectedError := types.ErrDataProviderAuthentication{}
						fakeData.ResolveReturnsOnCall(0, nil, expectedError)

						pullCreds := &types.PullCreds{
							Username: "dummyUsername",
							Password: "dummyPassword",
						}

						fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
						fakeCliParamSatisfier.SatisfyReturns(
							map[string]*types.Value{
								usernameInputName: {String: &pullCreds.Username},
								passwordInputName: {String: &pullCreds.Password},
							},
						)

						nodeURL := url.URL{Path: "dummyPath"}

						objectUnderTest := _dataResolver{
							data:                    fakeData,
							cliParamSatisfier:       fakeCliParamSatisfier,
							cliExiter:               new(cliexiter.Fake),
							nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
							nodeURL:                 nodeURL,
							os:                      new(ios.Fake),
						}

						/* act */
						objectUnderTest.Resolve(
							"dummyDataRef",
							&types.PullCreds{},
						)

						/* assert */
						actualNodeURL,
							actualPullCreds := fakeData.NewNodeProviderArgsForCall(1)
						Expect(actualNodeURL).To(Equal(nodeURL))
						Expect(actualPullCreds).To(Equal(pullCreds))
					})
					It("should call data.Resolve w/ expected args", func() {
						/* arrange */
						providedDataRef := "dummyDataRef"

						fakeData := new(data.Fake)

						fakeFSProvider := new(data.FakeProvider)
						fakeData.NewFSProviderReturns(fakeFSProvider)

						fakeNodeProvider := new(data.FakeProvider)
						fakeData.NewNodeProviderReturns(fakeNodeProvider)

						expectedError := types.ErrDataProviderAuthentication{}
						fakeData.ResolveReturnsOnCall(0, nil, expectedError)

						fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
						fakeCliParamSatisfier.SatisfyReturns(
							map[string]*types.Value{
								usernameInputName: {String: new(string)},
								passwordInputName: {String: new(string)},
							},
						)

						objectUnderTest := _dataResolver{
							data:                    fakeData,
							cliParamSatisfier:       fakeCliParamSatisfier,
							cliExiter:               new(cliexiter.Fake),
							nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
							os: new(ios.Fake),
						}

						/* act */
						objectUnderTest.Resolve(
							providedDataRef,
							&types.PullCreds{},
						)

						/* assert */
						_,
							actualDataRef,
							actualProviders := fakeData.ResolveArgsForCall(1)

						Expect(actualDataRef).To(Equal(providedDataRef))
						Expect(actualProviders).To(ConsistOf(fakeFSProvider, fakeNodeProvider))
					})
				})
				Context("not data.ErrAuthenticationFailed", func() {
					It("should call exiter w/ expected args", func() {
						/* arrange */
						providedDataRef := "dummyDataRef"

						fakeData := new(data.Fake)
						resolveError := errors.New("dummyError")
						fakeData.ResolveReturns(nil, resolveError)

						expectedError := fmt.Sprintf(
							"Unable to resolve pkg '%v'; error was %v",
							providedDataRef,
							resolveError.Error(),
						)

						fakeCliExiter := new(cliexiter.Fake)

						objectUnderTest := _dataResolver{
							data:                    fakeData,
							cliExiter:               fakeCliExiter,
							nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
							os: new(ios.Fake),
						}

						/* act */
						objectUnderTest.Resolve("dummyDataRef", &types.PullCreds{})

						/* assert */
						Expect(fakeCliExiter.ExitArgsForCall(0)).
							To(Equal(cliexiter.ExitReq{Message: expectedError, Code: 1}))

					})
				})
			})
			Context("data.Resolve doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeData := new(data.Fake)
					fakeHandle := new(data.FakeHandle)

					// err to trigger immediate return
					fakeData.ResolveReturns(fakeHandle, nil)

					objectUnderTest := _dataResolver{
						data: fakeData,
						nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
						os: new(ios.Fake),
					}

					/* act */
					actualPkgHandle := objectUnderTest.Resolve(
						"dummyDataRef",
						&types.PullCreds{},
					)

					/* assert */
					Expect(actualPkgHandle).To(Equal(fakeHandle))
				})
			})
		})
	})
})
