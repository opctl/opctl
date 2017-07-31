package core

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/cliparamsatisfier"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"net/url"
)

var _ = Context("core", func() {
	Context("PkgInstall", func() {
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

			objectUnderTest := _core{
				pkg:       fakePkg,
				cliExiter: new(cliexiter.Fake),
				nodeURL:   nodeURL,
			}

			/* act */
			objectUnderTest.PkgInstall(
				"dummyPath",
				"dummyPkgRef",
				providedPullCreds.Username,
				providedPullCreds.Password,
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
			fakeNodeProvider := new(pkg.FakeProvider)
			fakePkg.NewNodeProviderReturns(fakeNodeProvider)

			// err to trigger immediate return
			fakePkg.ResolveReturns(nil, errors.New("dummyError"))

			objectUnderTest := _core{
				pkg:       fakePkg,
				cliExiter: new(cliexiter.Fake),
				nodeURL:   url.URL{},
			}

			/* act */
			objectUnderTest.PkgInstall(
				"dummyPath",
				providedPkgRef,
				"dummyUsername",
				"dummyPassword",
			)

			/* assert */
			actualPkgRef, actualProviders := fakePkg.ResolveArgsForCall(0)

			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualProviders).To(ConsistOf(fakeNodeProvider))
		})
		Context("pkg.Resolve errs", func() {
			Context("pkg.ErrAuthenticationFailed", func() {
				It("should call cliParamSatisfier.Satisfy w/ expected args", func() {
					/* arrange */
					fakePkg := new(pkg.Fake)
					expectedError := pkg.ErrAuthenticationFailed{}
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

					objectUnderTest := _core{
						pkg:               fakePkg,
						cliParamSatisfier: fakeCliParamSatisfier,
						cliExiter:         new(cliexiter.Fake),
						nodeURL:           url.URL{},
					}

					/* act */
					objectUnderTest.PkgInstall("", "", "", "")

					/* assert */
					_, actualInputs := fakeCliParamSatisfier.SatisfyArgsForCall(0)
					Expect(actualInputs).To(Equal(credsPromptInputs))
				})
				It("should call pkg.NewNodeProvider w/ expected args", func() {
					/* arrange */
					fakePkg := new(pkg.Fake)
					fakeNodeProvider := new(pkg.FakeProvider)
					fakePkg.NewNodeProviderReturns(fakeNodeProvider)

					expectedError := pkg.ErrAuthenticationFailed{}
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

					objectUnderTest := _core{
						pkg:               fakePkg,
						cliParamSatisfier: fakeCliParamSatisfier,
						cliExiter:         new(cliexiter.Fake),
						nodeURL:           nodeURL,
					}

					/* act */
					objectUnderTest.PkgInstall(
						"dummyPath",
						"dummyPkgRef",
						"",
						"",
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
					fakeNodeProvider := new(pkg.FakeProvider)
					fakePkg.NewNodeProviderReturns(fakeNodeProvider)

					expectedError := pkg.ErrAuthenticationFailed{}
					fakePkg.ResolveReturnsOnCall(0, nil, expectedError)

					fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
					fakeCliParamSatisfier.SatisfyReturns(
						map[string]*model.Value{
							usernameInputName: {String: new(string)},
							passwordInputName: {String: new(string)},
						},
					)

					objectUnderTest := _core{
						pkg:               fakePkg,
						cliParamSatisfier: fakeCliParamSatisfier,
						cliExiter:         new(cliexiter.Fake),
						nodeURL:           url.URL{},
					}

					/* act */
					objectUnderTest.PkgInstall(
						"dummyPath",
						providedPkgRef,
						"",
						"",
					)

					/* assert */
					actualPkgRef, actualProviders := fakePkg.ResolveArgsForCall(1)

					Expect(actualPkgRef).To(Equal(providedPkgRef))
					Expect(actualProviders).To(ConsistOf(fakeNodeProvider))
				})
			})
			Context("not pkg.ErrAuthenticationFailed", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakePkg := new(pkg.Fake)
					expectedError := errors.New("dummyError")
					fakePkg.ResolveReturns(nil, expectedError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						pkg:       fakePkg,
						cliExiter: fakeCliExiter,
						nodeURL:   url.URL{},
					}

					/* act */
					objectUnderTest.PkgInstall("", "", "", "")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

				})
			})
		})
		Context("pkg.Resolve doesn't err", func() {
			It("should call pkg.Install w/ expected args", func() {
				/* arrange */
				providedPath := "dummyPath"

				fakePkg := new(pkg.Fake)
				fakeHandle := new(pkg.FakeHandle)

				// err to trigger immediate return
				fakePkg.ResolveReturns(fakeHandle, nil)

				objectUnderTest := _core{
					pkg:     fakePkg,
					nodeURL: url.URL{},
				}

				/* act */
				objectUnderTest.PkgInstall(
					providedPath,
					"dummyPkgRef",
					"dummyUsername",
					"dummyPassword",
				)

				/* assert */
				actualContext, actualPath, actualHandle := fakePkg.InstallArgsForCall(0)

				Expect(actualContext).To(Equal(context.TODO()))
				Expect(actualPath).To(Equal(providedPath))
				Expect(actualHandle).To(Equal(fakeHandle))
			})
			Context("pkg.Install errs", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakePkg := new(pkg.Fake)
					fakeHandle := new(pkg.FakeHandle)

					fakePkg.ResolveReturns(fakeHandle, nil)

					expectedError := errors.New("dummyError")
					fakePkg.InstallReturns(expectedError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						pkg:       fakePkg,
						cliExiter: fakeCliExiter,
						nodeURL:   url.URL{},
					}

					/* act */
					objectUnderTest.PkgInstall("", "", "", "")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

				})
			})
		})
	})
})
