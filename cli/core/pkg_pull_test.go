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
	Context("PkgPull", func() {
		It("should call pkg.Pull w/ expected args", func() {
			/* arrange */
			fakePkg := new(pkg.Fake)

			providedPkgRef := "dummyPkgRef"
			providedUsername := "dummyUsername"
			providedPassword := "dummyPassword"

			expectedPkgRef := providedPkgRef
			expectedPullOpts := &pkg.PullOpts{
				Username: providedUsername,
				Password: providedPassword,
			}

			objectUnderTest := _core{
				pkg: fakePkg,
			}

			/* act */
			objectUnderTest.PkgPull(providedPkgRef, providedUsername, providedPassword)

			/* assert */
			actualPkgRef,
				actualPullOpts := fakePkg.PullArgsForCall(0)

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
					objectUnderTest.PkgPull("", "", "")

					/* assert */
					_, actualInputs := fakeCliParamSatisfier.SatisfyArgsForCall(0)
					Expect(actualInputs).To(Equal(pullAuthPromptInputs))
				})
				It("should retry pkg.Pull w/ expected args", func() {
					/* arrange */
					providedPkgRef := "dummyPkgRef"
					expectedPkgRef := providedPkgRef

					fakePkg := new(pkg.Fake)
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
					objectUnderTest.PkgPull(providedPkgRef, "", "")

					/* assert */
					actualPkgRef,
						actualPullOpts := fakePkg.PullArgsForCall(1)

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
					objectUnderTest.PkgPull("", "", "")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

				})
			})
		})
	})
})
