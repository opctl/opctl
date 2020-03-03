package core

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	cliexiterFakes "github.com/opctl/opctl/cli/internal/cliexiter/fakes"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/cli/internal/updater"
)

var _ = Context("_selfUpdateInvoker", func() {

	Context("SelfUpdate", func() {
		Context("invalid channel", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)
				providedReleaseChannel := "invalidChannel"

				objectUnderTest := _selfUpdateInvoker{
					cliExiter: fakeCliExiter,
				}

				/* act */
				objectUnderTest.SelfUpdate(providedReleaseChannel)

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: fmt.Sprintf(
						"%v is not an available release channel. "+
							"Available release channels are 'alpha', 'beta', and 'stable'. \n", providedReleaseChannel), Code: 1}))
			})
		})
		Context("valid channel", func() {
			It("should call updater.GetUpdateIfExists w/ expected args", func() {
				/* arrange */
				fakeUpdater := new(updater.Fake)

				objectUnderTest := _selfUpdateInvoker{
					updater:   fakeUpdater,
					cliExiter: new(cliexiterFakes.FakeCliExiter),
				}

				providedChannel := "beta"

				/* act */
				objectUnderTest.SelfUpdate(providedChannel)

				/* assert */
				Expect(fakeUpdater.GetUpdateIfExistsArgsForCall(0)).To(Equal(providedChannel))
			})
			Context("updater.GetUpdateIfExists errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeCliExiter := new(cliexiterFakes.FakeCliExiter)
					returnedError := errors.New("dummyError")

					fakeUpdater := new(updater.Fake)
					fakeUpdater.GetUpdateIfExistsReturns(&updater.Update{}, returnedError)

					objectUnderTest := _selfUpdateInvoker{
						updater:   fakeUpdater,
						cliExiter: fakeCliExiter,
					}

					/* act */
					objectUnderTest.SelfUpdate("beta")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
				})
			})
			Context("updater.GetUpdateIfExists doesn't error", func() {
				Context("update doesn't exist", func() {
					It("should call exiter w/ expected args", func() {
						/* arrange */
						fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

						fakeUpdater := new(updater.Fake)
						fakeUpdater.GetUpdateIfExistsReturns(nil, nil)

						objectUnderTest := _selfUpdateInvoker{
							updater:   fakeUpdater,
							cliExiter: fakeCliExiter,
						}

						/* act */
						objectUnderTest.SelfUpdate("beta")

						/* assert */
						Expect(fakeCliExiter.ExitArgsForCall(0)).
							To(Equal(cliexiter.ExitReq{Message: "No update available, already at the latest version!", Code: 0}))
					})
				})
				Context("update exists", func() {
					It("should call updater.ApplyUpdate w/ expected args", func() {
						/* arrange */
						fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

						fakeUpdater := new(updater.Fake)
						returnedUpdate := &updater.Update{Version: "dummyVersion"}

						fakeUpdater.GetUpdateIfExistsReturns(returnedUpdate, nil)

						objectUnderTest := _selfUpdateInvoker{
							updater:      fakeUpdater,
							cliExiter:    fakeCliExiter,
							nodeProvider: new(nodeprovider.Fake),
						}

						/* act */
						objectUnderTest.SelfUpdate("beta")

						/* assert */
						Expect(fakeUpdater.ApplyUpdateArgsForCall(0)).
							To(Equal(returnedUpdate))
					})
					Context("updater.ApplyUpdate errors", func() {
						It("should call exiter w/ expected args", func() {
							/* arrange */
							fakeCliExiter := new(cliexiterFakes.FakeCliExiter)
							returnedError := errors.New("dummyError")

							fakeUpdater := new(updater.Fake)

							fakeUpdater.GetUpdateIfExistsReturns(&updater.Update{Version: "dummyVersion"}, nil)

							fakeUpdater.ApplyUpdateReturns(returnedError)

							objectUnderTest := _selfUpdateInvoker{
								updater:   fakeUpdater,
								cliExiter: fakeCliExiter,
							}

							/* act */
							objectUnderTest.SelfUpdate("beta")

							/* assert */
							Expect(fakeCliExiter.ExitArgsForCall(0)).
								To(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
						})
					})
					Context("updater.ApplyUpdate doesn't error", func() {
						It("should call nodeProvider.KillNodeIfExists", func() {
							/* arrange */
							fakeNodeProvider := new(nodeprovider.Fake)

							fakeUpdater := new(updater.Fake)
							returnedUpdate := &updater.Update{Version: "dummyVersion"}

							fakeUpdater.GetUpdateIfExistsReturns(returnedUpdate, nil)

							objectUnderTest := _selfUpdateInvoker{
								updater:      fakeUpdater,
								cliExiter:    new(cliexiterFakes.FakeCliExiter),
								nodeProvider: fakeNodeProvider,
							}

							/* act */
							objectUnderTest.SelfUpdate("beta")

							/* assert */
							Expect(fakeNodeProvider.KillNodeIfExistsCallCount()).To(Equal(1))
						})
						Context("nodeProvider.KillNodeIfExists errors", func() {
							It("should call exiter w/ expected args", func() {
								/* arrange */
								fakeCliExiter := new(cliexiterFakes.FakeCliExiter)
								returnedError := errors.New("dummyError")

								fakeNodeProvider := new(nodeprovider.Fake)
								fakeNodeProvider.KillNodeIfExistsReturns(returnedError)

								expectedExitMsg :=
									fmt.Sprintf("Unable to kill running node; run `node kill` to complete the update. Error was: %v", returnedError.Error())

								fakeUpdater := new(updater.Fake)

								fakeUpdater.GetUpdateIfExistsReturns(&updater.Update{Version: "dummyVersion"}, nil)

								objectUnderTest := _selfUpdateInvoker{
									nodeProvider: fakeNodeProvider,
									updater:      fakeUpdater,
									cliExiter:    fakeCliExiter,
								}

								/* act */
								objectUnderTest.SelfUpdate("beta")

								/* assert */
								Expect(fakeCliExiter.ExitArgsForCall(0)).
									To(Equal(cliexiter.ExitReq{Message: expectedExitMsg, Code: 1}))
							})
						})
						Context("nodeProvider.KillNodeIfExists doesn't error", func() {
							It("should call exiter w/ expected args", func() {
								/* arrange */
								fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

								fakeUpdater := new(updater.Fake)
								returnedUpdate := &updater.Update{Version: "dummyVersion"}

								fakeUpdater.GetUpdateIfExistsReturns(returnedUpdate, nil)

								objectUnderTest := _selfUpdateInvoker{
									updater:      fakeUpdater,
									cliExiter:    fakeCliExiter,
									nodeProvider: new(nodeprovider.Fake),
								}

								/* act */
								objectUnderTest.SelfUpdate("beta")

								/* assert */
								Expect(fakeCliExiter.ExitArgsForCall(0)).
									To(Equal(cliexiter.ExitReq{
										Message: fmt.Sprintf("Updated to new version: %s!\n", returnedUpdate.Version),
										Code:    0,
									}))
							})
						})
					})
				})
			})
		})
	})
})
