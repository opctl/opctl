package core

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/updater"
)

var _ = Describe("selfUpdate", func() {

	Context("Execute", func() {
		Context("invalid channel", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeExiter := new(fakeExiter)
				providedReleaseChannel := "invalidChannel"

				objectUnderTest := _core{
					exiter: fakeExiter,
				}

				/* act */
				objectUnderTest.SelfUpdate(providedReleaseChannel)

				/* assert */
				Expect(fakeExiter.ExitArgsForCall(0)).
					Should(Equal(ExitReq{Message: fmt.Sprintf(
						"%v is not an available release channel. "+
							"Available release channels are 'beta' 'stable'. \n", providedReleaseChannel), Code: 1}))
			})
		})
		Context("valid channel", func() {
			It("should call updater.GetUpdateIfExists w/ expected args", func() {
				/* arrange */
				fakeUpdater := new(updater.FakeUpdater)

				objectUnderTest := _core{
					updater: fakeUpdater,
					exiter:  new(fakeExiter),
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
					fakeExiter := new(fakeExiter)
					returnedError := errors.New("dummyError")

					fakeUpdater := new(updater.FakeUpdater)
					fakeUpdater.GetUpdateIfExistsReturns(&updater.Update{}, returnedError)

					objectUnderTest := _core{
						updater: fakeUpdater,
						exiter:  fakeExiter,
					}

					/* act */
					objectUnderTest.SelfUpdate("beta")

					/* assert */
					Expect(fakeExiter.ExitArgsForCall(0)).
						Should(Equal(ExitReq{Message: returnedError.Error(), Code: 1}))
				})
			})
			Context("updater.GetUpdateIfExists doesn't error", func() {
				Context("update doesn't exist", func() {
					It("should call exiter w/ expected args", func() {
						/* arrange */
						fakeExiter := new(fakeExiter)

						fakeUpdater := new(updater.FakeUpdater)
						fakeUpdater.GetUpdateIfExistsReturns(nil, nil)

						objectUnderTest := _core{
							updater: fakeUpdater,
							exiter:  fakeExiter,
						}

						/* act */
						objectUnderTest.SelfUpdate("beta")

						/* assert */
						Expect(fakeExiter.ExitArgsForCall(0)).
							Should(Equal(ExitReq{Message: "No update available, already at the latest version!", Code: 0}))
					})
				})
				Context("update exists", func() {
					It("should call updater.ApplyUpdate w/ expected args", func() {
						/* arrange */
						fakeExiter := new(fakeExiter)

						fakeUpdater := new(updater.FakeUpdater)
						returnedUpdate := &updater.Update{Version: "dummyVersion"}

						fakeUpdater.GetUpdateIfExistsReturns(returnedUpdate, nil)

						objectUnderTest := _core{
							updater: fakeUpdater,
							exiter:  fakeExiter,
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
							fakeExiter := new(fakeExiter)
							returnedError := errors.New("dummyError")

							fakeUpdater := new(updater.FakeUpdater)

							fakeUpdater.GetUpdateIfExistsReturns(&updater.Update{Version: "dummyVersion"}, nil)

							fakeUpdater.ApplyUpdateReturns(returnedError)

							objectUnderTest := _core{
								updater: fakeUpdater,
								exiter:  fakeExiter,
							}

							/* act */
							objectUnderTest.SelfUpdate("beta")

							/* assert */
							Expect(fakeExiter.ExitArgsForCall(0)).
								Should(Equal(ExitReq{Message: returnedError.Error(), Code: 1}))
						})
					})
					Context("updater.ApplyUpdate doesn't error", func() {
						It("should call exiter w/ expected args", func() {
							/* arrange */
							fakeExiter := new(fakeExiter)

							fakeUpdater := new(updater.FakeUpdater)
							returnedUpdate := &updater.Update{Version: "dummyVersion"}

							fakeUpdater.GetUpdateIfExistsReturns(returnedUpdate, nil)

							objectUnderTest := _core{
								updater: fakeUpdater,
								exiter:  fakeExiter,
							}

							/* act */
							objectUnderTest.SelfUpdate("beta")

							/* assert */
							Expect(fakeExiter.ExitArgsForCall(0)).
								Should(Equal(ExitReq{
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
