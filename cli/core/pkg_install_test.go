package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/cliparamsatisfier"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
)

var _ = Context("core", func() {
	Context("PkgInstall", func() {
		It("should call pkg.ParseRef w/ expected args", func() {
			/* arrange */
			fakePkg := new(pkg.Fake)

			providedPkgRef := "dummyPkgRef"
			expectedPkgRef := providedPkgRef

			// error to trigger immediate return
			fakePkg.ParseRefReturns(nil, errors.New("dummyError"))

			objectUnderTest := _core{
				pkg:       fakePkg,
				cliExiter: new(cliexiter.Fake),
			}

			/* act */
			objectUnderTest.PkgInstall("", providedPkgRef, "", "")

			/* assert */
			Expect(fakePkg.ParseRefArgsForCall(0)).To(Equal(expectedPkgRef))
		})
		Context("pkg.ParseRef errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)
				expectedError := errors.New("dummyError")
				fakePkg.ParseRefReturns(nil, expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:       fakePkg,
					cliExiter: fakeCliExiter,
				}

				/* act */
				objectUnderTest.PkgInstall("", "", "", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

			})
		})
		Context("pkg.ParseRef doesn't error", func() {
			It("should call pkg.Pull w/ expected args", func() {
				/* arrange */
				providedPath := "dummyPath"
				providedPkgRef := "dummyPkgRef"
				providedUsername := "dummyUsername"
				providedPassword := "dummyPassword"

				expectedPath := providedPath
				expectedPkgRef := &pkg.PkgRef{
					FullyQualifiedName: "dummyFQName",
					Version:            "dummyVersion",
				}

				fakePkg := new(pkg.Fake)
				fakePkg.ParseRefReturns(expectedPkgRef, nil)

				expectedPullOpts := &pkg.PullOpts{
					Username: providedUsername,
					Password: providedPassword,
				}

				objectUnderTest := _core{
					pkg: fakePkg,
				}

				/* act */
				objectUnderTest.PkgInstall(providedPath, providedPkgRef, providedUsername, providedPassword)

				/* assert */
				actualPath,
					actualPkgRef,
					actualPullOpts := fakePkg.PullArgsForCall(0)

				Expect(actualPath).To(Equal(expectedPath))
				Expect(actualPkgRef).To(Equal(expectedPkgRef))
				Expect(actualPullOpts).To(Equal(expectedPullOpts))
			})
			Context("pkg.Pull errors", func() {
				Context("pkg.ErrAuthenticationFailed", func() {
					It("should call cliParamSatisfier.Satisfy w/ expected args", func() {
						/* arrange */
						fakePkg := new(pkg.Fake)
						expectedError := pkg.ErrAuthenticationFailed{}
						fakePkg.PullReturnsOnCall(0, expectedError)

						username := "dummyUsername"
						password := "dummyPassword"

						fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
						fakeCliParamSatisfier.SatisfyReturns(
							map[string]*model.Data{
								usernameInputName: {String: &username},
								passwordInputName: {String: &password},
							},
						)

						fakeCliExiter := new(cliexiter.Fake)

						objectUnderTest := _core{
							pkg:               fakePkg,
							cliParamSatisfier: fakeCliParamSatisfier,
							cliExiter:         fakeCliExiter,
						}

						/* act */
						objectUnderTest.PkgInstall("", "", "", "")

						/* assert */
						_, actualInputs := fakeCliParamSatisfier.SatisfyArgsForCall(0)
						Expect(actualInputs).To(Equal(credsPromptInputs))
					})
					It("should retry pkg.Pull w/ expected args", func() {
						/* arrange */
						providedPath := "dummyPath"

						expectedPath := providedPath
						expectedPkgRef := &pkg.PkgRef{
							FullyQualifiedName: "dummyFQName",
							Version:            "dummyVersion",
						}

						fakePkg := new(pkg.Fake)
						fakePkg.ParseRefReturns(expectedPkgRef, nil)
						expectedError := pkg.ErrAuthenticationFailed{}
						fakePkg.PullReturnsOnCall(0, expectedError)

						expectedPullOpts := &pkg.PullOpts{
							Username: "dummyUsername",
							Password: "dummyPassword",
						}

						fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
						fakeCliParamSatisfier.SatisfyReturns(
							map[string]*model.Data{
								usernameInputName: {String: &(expectedPullOpts.Username)},
								passwordInputName: {String: &(expectedPullOpts.Password)},
							},
						)

						fakeCliExiter := new(cliexiter.Fake)

						objectUnderTest := _core{
							pkg:               fakePkg,
							cliParamSatisfier: fakeCliParamSatisfier,
							cliExiter:         fakeCliExiter,
						}

						/* act */
						objectUnderTest.PkgInstall(providedPath, "dummyPkgRef", "dummyUsername", "dummyPassword")

						/* assert */
						actualPath,
							actualPkgRef,
							actualPullOpts := fakePkg.PullArgsForCall(1)

						Expect(actualPath).To(Equal(expectedPath))
						Expect(actualPkgRef).To(Equal(expectedPkgRef))
						Expect(actualPullOpts).To(Equal(expectedPullOpts))
					})
				})
				Context("not pkg.ErrAuthenticationFailed", func() {
					It("should call exiter w/ expected args", func() {
						/* arrange */
						fakePkg := new(pkg.Fake)
						expectedError := errors.New("dummyError")
						fakePkg.PullReturns(expectedError)

						fakeCliExiter := new(cliexiter.Fake)

						objectUnderTest := _core{
							pkg:       fakePkg,
							cliExiter: fakeCliExiter,
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
})
