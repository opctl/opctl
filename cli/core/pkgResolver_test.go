package core

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/cliparamsatisfier"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"net/url"
)

var _ = Context("pkgResolver", func() {
	Context("Resolve", func() {
		It("should call nodeReachabilityEnsurer.EnsureNodeReachable", func() {
			/* arrange */
			fakeNodeReachabilityEnsurer := new(fakeNodeReachabilityEnsurer)

			fakeIOS := new(ios.Fake)
			// err to trigger immediate return
			fakeIOS.GetwdReturns("", errors.New("dummyError"))

			objectUnderTest := _pkgResolver{
				cliExiter:               new(cliexiter.Fake),
				nodeReachabilityEnsurer: fakeNodeReachabilityEnsurer,
				os: fakeIOS,
			}

			/* act */
			objectUnderTest.Resolve("dummyPkgRef", &model.PullCreds{})

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

				objectUnderTest := _pkgResolver{
					cliExiter: fakeCliExiter,
					os:        fakeIOS,
					nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
				}

				/* act */
				objectUnderTest.Resolve(
					"dummyPkgRef",
					nil,
				)

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("os.Getwd doesn't err", func() {
			It("should call pkg.NewFSProvider w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)
				fakeFSProvider := new(pkg.FakeProvider)
				fakePkg.NewFSProviderReturns(fakeFSProvider)

				// error to trigger immediate return
				fakePkg.ResolveReturns(nil, errors.New("dummyError"))

				fakeIOS := new(ios.Fake)
				workDir := "dummyWorkDir"
				fakeIOS.GetwdReturns(workDir, nil)

				objectUnderTest := _pkgResolver{
					pkg:                     fakePkg,
					cliExiter:               new(cliexiter.Fake),
					nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
					os: fakeIOS,
				}

				/* act */
				objectUnderTest.Resolve("dummyPkgRef", &model.PullCreds{})

				/* assert */
				Expect(fakePkg.NewFSProviderArgsForCall(0)).To(ConsistOf(workDir))
			})
			It("should call pkg.NewNodeProvider w/ expected args", func() {
				/* arrange */
				providedPullCreds := &model.PullCreds{
					Username: "dummyUsername",
					Password: "dummyPassword",
				}

				fakePkg := new(pkg.Fake)
				fakeNodeProvider := new(pkg.FakeProvider)
				fakePkg.NewNodeProviderReturns(fakeNodeProvider)

				// err to trigger immediate return
				fakePkg.ResolveReturns(nil, errors.New("dummyError"))

				nodeURL := url.URL{Path: "dummyNodeURL"}

				objectUnderTest := _pkgResolver{
					pkg:                     fakePkg,
					cliExiter:               new(cliexiter.Fake),
					nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
					nodeURL:                 nodeURL,
					os:                      new(ios.Fake),
				}

				/* act */
				objectUnderTest.Resolve(
					"dummyPkgRef",
					providedPullCreds,
				)

				/* assert */
				actualNodeURL,
					actualPullCreds := fakePkg.NewNodeProviderArgsForCall(0)
				Expect(actualNodeURL).To(Equal(nodeURL))
				Expect(actualPullCreds).To(Equal(providedPullCreds))
			})
			It("should call pkg.Resolve w/ expected args", func() {
				/* arrange */
				providedPkgRef := "dummyPkgRef"

				fakePkg := new(pkg.Fake)

				fakeFSProvider := new(pkg.FakeProvider)
				fakePkg.NewFSProviderReturns(fakeFSProvider)

				fakeNodeProvider := new(pkg.FakeProvider)
				fakePkg.NewNodeProviderReturns(fakeNodeProvider)

				// err to trigger immediate return
				fakePkg.ResolveReturns(nil, errors.New("dummyError"))

				objectUnderTest := _pkgResolver{
					pkg:                     fakePkg,
					cliExiter:               new(cliexiter.Fake),
					nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
					os: new(ios.Fake),
				}

				/* act */
				objectUnderTest.Resolve(
					providedPkgRef,
					&model.PullCreds{},
				)

				/* assert */
				_,
					actualPkgRef,
					actualProviders := fakePkg.ResolveArgsForCall(0)

				Expect(actualPkgRef).To(Equal(providedPkgRef))
				Expect(actualProviders).To(ConsistOf(fakeFSProvider, fakeNodeProvider))
			})
			Context("pkg.Resolve errs", func() {
				Context("pkg.ErrPkgPullAuthorization", func() {
					It("should call cliParamSatisfier.Satisfy w/ expected args", func() {
						/* arrange */
						fakePkg := new(pkg.Fake)
						expectedError := model.ErrPkgPullAuthorization{}
						fakePkg.ResolveReturnsOnCall(0, nil, expectedError)

						username := "dummyUsername"
						password := "dummyPassword"

						fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
						fakeCliParamSatisfier.SatisfyReturns(
							map[string]*model.Value{
								usernameInputName: {String: &username},
								passwordInputName: {String: &password},
							},
						)

						objectUnderTest := _pkgResolver{
							pkg:                     fakePkg,
							cliParamSatisfier:       fakeCliParamSatisfier,
							cliExiter:               new(cliexiter.Fake),
							nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
							os: new(ios.Fake),
						}

						/* act */
						objectUnderTest.Resolve("dummyPkgRef", &model.PullCreds{})

						/* assert */
						_, actualInputs := fakeCliParamSatisfier.SatisfyArgsForCall(0)
						Expect(actualInputs).To(Equal(credsPromptInputs))
					})
					It("should call pkg.NewNodeProvider w/ expected args", func() {
						/* arrange */
						fakePkg := new(pkg.Fake)
						fakeNodeProvider := new(pkg.FakeProvider)
						fakePkg.NewNodeProviderReturns(fakeNodeProvider)

						expectedError := model.ErrPkgPullAuthentication{}
						fakePkg.ResolveReturnsOnCall(0, nil, expectedError)

						pullCreds := &model.PullCreds{
							Username: "dummyUsername",
							Password: "dummyPassword",
						}

						fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
						fakeCliParamSatisfier.SatisfyReturns(
							map[string]*model.Value{
								usernameInputName: {String: &pullCreds.Username},
								passwordInputName: {String: &pullCreds.Password},
							},
						)

						nodeURL := url.URL{Path: "dummyPath"}

						objectUnderTest := _pkgResolver{
							pkg:                     fakePkg,
							cliParamSatisfier:       fakeCliParamSatisfier,
							cliExiter:               new(cliexiter.Fake),
							nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
							nodeURL:                 nodeURL,
							os:                      new(ios.Fake),
						}

						/* act */
						objectUnderTest.Resolve(
							"dummyPkgRef",
							&model.PullCreds{},
						)

						/* assert */
						actualNodeURL,
							actualPullCreds := fakePkg.NewNodeProviderArgsForCall(1)
						Expect(actualNodeURL).To(Equal(nodeURL))
						Expect(actualPullCreds).To(Equal(pullCreds))
					})
					It("should call pkg.Resolve w/ expected args", func() {
						/* arrange */
						providedPkgRef := "dummyPkgRef"

						fakePkg := new(pkg.Fake)

						fakeFSProvider := new(pkg.FakeProvider)
						fakePkg.NewFSProviderReturns(fakeFSProvider)

						fakeNodeProvider := new(pkg.FakeProvider)
						fakePkg.NewNodeProviderReturns(fakeNodeProvider)

						expectedError := model.ErrPkgPullAuthentication{}
						fakePkg.ResolveReturnsOnCall(0, nil, expectedError)

						fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
						fakeCliParamSatisfier.SatisfyReturns(
							map[string]*model.Value{
								usernameInputName: {String: new(string)},
								passwordInputName: {String: new(string)},
							},
						)

						objectUnderTest := _pkgResolver{
							pkg:                     fakePkg,
							cliParamSatisfier:       fakeCliParamSatisfier,
							cliExiter:               new(cliexiter.Fake),
							nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
							os: new(ios.Fake),
						}

						/* act */
						objectUnderTest.Resolve(
							providedPkgRef,
							&model.PullCreds{},
						)

						/* assert */
						_,
							actualPkgRef,
							actualProviders := fakePkg.ResolveArgsForCall(1)

						Expect(actualPkgRef).To(Equal(providedPkgRef))
						Expect(actualProviders).To(ConsistOf(fakeFSProvider, fakeNodeProvider))
					})
				})
				Context("not pkg.ErrAuthenticationFailed", func() {
					It("should call exiter w/ expected args", func() {
						/* arrange */
						providedPkgRef := "dummyPkgRef"

						fakePkg := new(pkg.Fake)
						resolveError := errors.New("dummyError")
						fakePkg.ResolveReturns(nil, resolveError)

						expectedError := fmt.Sprintf(
							"Unable to resolve pkg '%v'; error was %v",
							providedPkgRef,
							resolveError.Error(),
						)

						fakeCliExiter := new(cliexiter.Fake)

						objectUnderTest := _pkgResolver{
							pkg:                     fakePkg,
							cliExiter:               fakeCliExiter,
							nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
							os: new(ios.Fake),
						}

						/* act */
						objectUnderTest.Resolve("dummyPkgRef", &model.PullCreds{})

						/* assert */
						Expect(fakeCliExiter.ExitArgsForCall(0)).
							To(Equal(cliexiter.ExitReq{Message: expectedError, Code: 1}))

					})
				})
			})
			Context("pkg.Resolve doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakePkg := new(pkg.Fake)
					fakeHandle := new(pkg.FakeHandle)

					// err to trigger immediate return
					fakePkg.ResolveReturns(fakeHandle, nil)

					objectUnderTest := _pkgResolver{
						pkg: fakePkg,
						nodeReachabilityEnsurer: new(fakeNodeReachabilityEnsurer),
						os: new(ios.Fake),
					}

					/* act */
					actualPkgHandle := objectUnderTest.Resolve(
						"dummyPkgRef",
						&model.PullCreds{},
					)

					/* assert */
					Expect(actualPkgHandle).To(Equal(fakeHandle))
				})
			})
		})
	})
})
