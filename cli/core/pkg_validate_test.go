package core

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
)

var _ = Context("pkgValidate", func() {
	Context("Execute", func() {
		Context("ios.Getwd errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeIOS := new(ios.Fake)
				expectedError := errors.New("dummyError")
				fakeIOS.GetwdReturns("", expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:       new(pkg.Fake),
					cliExiter: fakeCliExiter,
					os:        fakeIOS,
				}

				/* act */
				objectUnderTest.PkgValidate("dummyPkgRef")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
		Context("ios.Getwd doesn't error", func() {
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

				objectUnderTest := _core{
					pkg:       fakePkg,
					cliExiter: new(cliexiter.Fake),
					os:        fakeIOS,
				}

				/* act */
				objectUnderTest.PkgValidate("dummyPkgRef")

				/* assert */
				Expect(fakePkg.NewFSProviderArgsForCall(0)).To(ConsistOf(workDir))
			})
			It("should call pkg.Resolve w/ expected args", func() {
				/* arrange */
				providedPkgRef := "dummyPkgRef"

				fakePkg := new(pkg.Fake)
				fakeFSProvider := new(pkg.FakeProvider)
				fakePkg.NewFSProviderReturns(fakeFSProvider)

				// error to trigger immediate return
				fakePkg.ResolveReturns(nil, errors.New("dummyError"))

				objectUnderTest := _core{
					pkg:       fakePkg,
					cliExiter: new(cliexiter.Fake),
					os:        new(ios.Fake),
				}

				/* act */
				objectUnderTest.PkgValidate(providedPkgRef)

				/* assert */
				actualPkgRef, actualProviders := fakePkg.ResolveArgsForCall(0)
				Expect(actualPkgRef).To(Equal(providedPkgRef))
				Expect(actualProviders).To(ConsistOf(fakeFSProvider))
			})
			Context("pkg.Resolve errs", func() {

				It("should call exiter w/ expected args", func() {
					/* arrange */
					providedPkgRef := "dummyPkgRef"
					wdReturnedFromIOS := "dummyWorkDir"

					fakeIOS := new(ios.Fake)
					fakeIOS.GetwdReturns(wdReturnedFromIOS, nil)

					resolveError := errors.New("dummyError")
					expectedMsg := fmt.Sprintf(
						"Unable to resolve package '%v' from '%v'; error was: %v",
						providedPkgRef,
						wdReturnedFromIOS,
						resolveError.Error(),
					)

					fakePkg := new(pkg.Fake)
					fakePkg.ResolveReturns(nil, resolveError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						pkg:       fakePkg,
						cliExiter: fakeCliExiter,
						os:        fakeIOS,
					}

					/* act */
					objectUnderTest.PkgValidate(providedPkgRef)

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: expectedMsg, Code: 1}))
				})
			})
			Context("pkg.Resolve doesn't error", func() {
				It("should call pkg.Validate w/ expected args", func() {
					/* arrange */

					fakePkgHandle := new(pkg.FakeHandle)

					fakePkg := new(pkg.Fake)
					fakePkg.ResolveReturns(fakePkgHandle, nil)

					objectUnderTest := _core{
						pkg:       fakePkg,
						cliExiter: new(cliexiter.Fake),
						os:        new(ios.Fake),
					}

					/* act */
					objectUnderTest.PkgValidate("dummyPkgRef")

					/* assert */
					Expect(fakePkg.ValidateArgsForCall(0)).To(Equal(fakePkgHandle))
				})
				Context("pkg.Validate returns errors", func() {
					It("should call cliExiter.Exit w/ expected args", func() {
						/* arrange */
						fakePkg := new(pkg.Fake)

						errsReturnedFromValidate := []error{errors.New("dummyError")}
						fakePkg.ValidateReturns(errsReturnedFromValidate)

						expectedExitReq := cliexiter.ExitReq{
							Message: fmt.Sprintf(`
-
  Error(s):
    - %v
-`, errsReturnedFromValidate[0]),
							Code: 1,
						}

						fakeCliExiter := new(cliexiter.Fake)

						objectUnderTest := _core{
							pkg:       fakePkg,
							cliExiter: fakeCliExiter,
							os:        new(ios.Fake),
						}

						/* act */
						objectUnderTest.PkgValidate("dummyPkgRef")

						/* assert */

						Expect(fakeCliExiter.ExitArgsForCall(0)).To(Equal(expectedExitReq))
					})
				})
				Context("pkg.Validate doesn't return errors", func() {
					It("should call cliExiter.Exit w/ expected args", func() {
						/* arrange */
						fakePkgHandle := new(pkg.FakeHandle)
						pkgRef := "dummyPkgRef"
						fakePkgHandle.RefReturns(pkgRef)

						fakePkg := new(pkg.Fake)
						fakePkg.ResolveReturns(fakePkgHandle, nil)

						expectedExitReq := cliexiter.ExitReq{
							Message: fmt.Sprintf("%v is valid", pkgRef),
						}

						fakeCliExiter := new(cliexiter.Fake)

						objectUnderTest := _core{
							pkg:       fakePkg,
							cliExiter: fakeCliExiter,
							os:        new(ios.Fake),
						}

						/* act */
						objectUnderTest.PkgValidate("dummyPkgRef")

						/* assert */

						Expect(fakeCliExiter.ExitArgsForCall(0)).To(Equal(expectedExitReq))
					})
				})
			})
		})
	})
})
