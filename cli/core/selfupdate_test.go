package core

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/nodeprovider"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/updater"
)

var _ = Context("selfUpdate", func() {

	Context("Execute", func() {
		Context("invalid channel", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeCliExiter := new(cliexiter.Fake)
				providedReleaseChannel := "invalidChannel"

				objectUnderTest := _core{
					cliExiter: fakeCliExiter,
				}

				/* act */
				objectUnderTest.SelfUpdate(providedReleaseChannel)

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: fmt.Sprintf(
						"%v is not an available release channel. "+
							"Available release channels are 'alpha', 'beta', and 'stable'. \n", providedReleaseChannel), Code: 1}))
			})
		})
		Context("valid channel", func() {
			It("should call updater.GetUpdateIfExists w/ expected args", func() {
				/* arrange */
				fakeUpdater := new(updater.Fake)

				objectUnderTest := _core{
					updater:   fakeUpdater,
					cliExiter: new(cliexiter.Fake),
				}

				providedChannel := "beta"

				/* act */
				objectUnderTest.SelfUpdate(providedChannel)

				/* assert */
				Expect(fakeUpdater.GetUpdateIfExistsArgsForCall(0)).Should(Equal(providedChannel))
			})
			Context("updater.GetUpdateIfExists errors", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeCliExiter := new(cliexiter.Fake)
					returnedError := errors.New("dummyError")

					fakeUpdater := new(updater.Fake)
					fakeUpdater.GetUpdateIfExistsReturns(&updater.Update{}, returnedError)

					objectUnderTest := _core{
						updater:   fakeUpdater,
						cliExiter: fakeCliExiter,
					}

					/* act */
					objectUnderTest.SelfUpdate("beta")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
				})
			})
			Context("updater.GetUpdateIfExists doesn't error", func() {
				Context("update doesn't exist", func() {
					It("should call exiter w/ expected args", func() {
						/* arrange */
						fakeCliExiter := new(cliexiter.Fake)

						fakeUpdater := new(updater.Fake)
						fakeUpdater.GetUpdateIfExistsReturns(nil, nil)

						objectUnderTest := _core{
							updater:   fakeUpdater,
							cliExiter: fakeCliExiter,
						}

						/* act */
						objectUnderTest.SelfUpdate("beta")

						/* assert */
						Expect(fakeCliExiter.ExitArgsForCall(0)).
							Should(Equal(cliexiter.ExitReq{Message: "No update available, already at the latest version!", Code: 0}))
					})
				})
				Context("update exists", func() {
					It("should call updater.ApplyUpdate w/ expected args", func() {
						/* arrange */
						fakeCliExiter := new(cliexiter.Fake)

						fakeUpdater := new(updater.Fake)
						returnedUpdate := &updater.Update{Version: "dummyVersion"}

						fakeUpdater.GetUpdateIfExistsReturns(returnedUpdate, nil)

						objectUnderTest := _core{
							updater:      fakeUpdater,
							cliExiter:    fakeCliExiter,
							nodeProvider: new(nodeprovider.Fake),
						}

						/* act */
						objectUnderTest.SelfUpdate("beta")

						/* assert */
						Expect(fakeUpdater.ApplyUpdateArgsForCall(0)).
							Should(Equal(returnedUpdate))
					})
					Context("updater.ApplyUpdate errors", func() {
						It("should call exiter w/ expected args", func() {
							/* arrange */
							fakeCliExiter := new(cliexiter.Fake)
							returnedError := errors.New("dummyError")

							fakeUpdater := new(updater.Fake)

							fakeUpdater.GetUpdateIfExistsReturns(&updater.Update{Version: "dummyVersion"}, nil)

							fakeUpdater.ApplyUpdateReturns(returnedError)

							objectUnderTest := _core{
								updater:   fakeUpdater,
								cliExiter: fakeCliExiter,
							}

							/* act */
							objectUnderTest.SelfUpdate("beta")

							/* assert */
							Expect(fakeCliExiter.ExitArgsForCall(0)).
								Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
						})
					})
					Context("updater.ApplyUpdate doesn't error", func() {
						It("should call nodeProvider.KillNodeIfExists", func() {
							/* arrange */
							fakeNodeProvider := new(nodeprovider.Fake)

							fakeUpdater := new(updater.Fake)
							returnedUpdate := &updater.Update{Version: "dummyVersion"}

							fakeUpdater.GetUpdateIfExistsReturns(returnedUpdate, nil)

							objectUnderTest := _core{
								updater:      fakeUpdater,
								cliExiter:    new(cliexiter.Fake),
								nodeProvider: fakeNodeProvider,
							}

							/* act */
							objectUnderTest.SelfUpdate("beta")

							/* assert */
							Expect(fakeNodeProvider.KillNodeIfExistsCallCount()).To(Equal(1))
						})
						It("should call nodeProvider.CreateNode", func() {
							/* arrange */
							fakeNodeProvider := new(nodeprovider.Fake)

							fakeUpdater := new(updater.Fake)
							returnedUpdate := &updater.Update{Version: "dummyVersion"}

							fakeUpdater.GetUpdateIfExistsReturns(returnedUpdate, nil)

							objectUnderTest := _core{
								updater:      fakeUpdater,
								cliExiter:    new(cliexiter.Fake),
								nodeProvider: fakeNodeProvider,
							}

							/* act */
							objectUnderTest.SelfUpdate("beta")

							/* assert */
							Expect(fakeNodeProvider.CreateNodeCallCount()).To(Equal(1))
						})
						It("should call exiter w/ expected args", func() {
							/* arrange */
							fakeCliExiter := new(cliexiter.Fake)

							fakeUpdater := new(updater.Fake)
							returnedUpdate := &updater.Update{Version: "dummyVersion"}

							fakeUpdater.GetUpdateIfExistsReturns(returnedUpdate, nil)

							objectUnderTest := _core{
								updater:      fakeUpdater,
								cliExiter:    fakeCliExiter,
								nodeProvider: new(nodeprovider.Fake),
							}

							/* act */
							objectUnderTest.SelfUpdate("beta")

							/* assert */
							Expect(fakeCliExiter.ExitArgsForCall(0)).
								Should(Equal(cliexiter.ExitReq{
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
