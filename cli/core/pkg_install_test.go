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
		It("should call pkg.Resolve w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"
			providedUsername := "dummyUsername"
			providedPassword := "dummyPassword"

			fakePkg := new(pkg.Fake)
			// err to trigger immediate return
			fakePkg.ResolveReturns(nil, errors.New("dummyError"))

			objectUnderTest := _core{
				pkg:       fakePkg,
				cliExiter: new(cliexiter.Fake),
			}

			/* act */
			objectUnderTest.PkgInstall("dummyPath", providedPkgRef, providedUsername, providedPassword)

			/* assert */
			actualPkgRef, actualResolveOpts := fakePkg.ResolveArgsForCall(0)

			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualResolveOpts).To(Equal(&pkg.ResolveOpts{
				PullCreds: &pkg.PullCreds{
					Username: providedUsername,
					Password: providedPassword,
				},
			}))
		})
		Context("pkg.Resolve errors", func() {
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
					}

					/* act */
					objectUnderTest.PkgInstall("", "", "", "")

					/* assert */
					_, actualInputs := fakeCliParamSatisfier.SatisfyArgsForCall(0)
					Expect(actualInputs).To(Equal(credsPromptInputs))
				})
				It("should retry pkg.Resolve w/ expected args", func() {
					/* arrange */
					expectedPkgRef := "dummyPkgRef"

					fakePkg := new(pkg.Fake)
					expectedError := pkg.ErrAuthenticationFailed{}
					fakePkg.ResolveReturnsOnCall(0, nil, expectedError)

					expectedResolveOpts := &pkg.ResolveOpts{
						PullCreds: &pkg.PullCreds{
							Username: "dummyUsername",
							Password: "dummyPassword",
						},
					}

					fakeCliParamSatisfier := new(cliparamsatisfier.Fake)
					fakeCliParamSatisfier.SatisfyReturns(
						map[string]*model.Value{
							usernameInputName: {String: &(expectedResolveOpts).PullCreds.Username},
							passwordInputName: {String: &(expectedResolveOpts.PullCreds.Password)},
						},
					)

					objectUnderTest := _core{
						pkg:               fakePkg,
						cliParamSatisfier: fakeCliParamSatisfier,
						cliExiter:         new(cliexiter.Fake),
					}

					/* act */
					objectUnderTest.PkgInstall(
						"dummyPath",
						"dummyPkgRef",
						"dummyUsername",
						"dummyPassword",
					)

					/* assert */
					actualPkgRef, actualResolveOpts := fakePkg.ResolveArgsForCall(1)

					Expect(actualPkgRef).To(Equal(expectedPkgRef))
					Expect(actualResolveOpts).To(Equal(expectedResolveOpts))
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
