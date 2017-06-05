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

			It("should call pkg.ParseRef w/ expected args", func() {
				/* arrange */
				providedPkgRef := "dummyPkgName"
				expectedPkgRef := providedPkgRef

				fakePkg := new(pkg.Fake)
				// error to trigger immediate return
				fakePkg.ParseRefReturns(nil, errors.New("dummyError"))

				objectUnderTest := _core{
					pkg:       fakePkg,
					cliExiter: new(cliexiter.Fake),
					os:        new(ios.Fake),
				}

				/* act */
				objectUnderTest.PkgValidate(providedPkgRef)

				/* assert */
				actualPkgRef := fakePkg.ParseRefArgsForCall(0)
				Expect(actualPkgRef).To(Equal(expectedPkgRef))

			})
			Context("pkg.ParseRef errors", func() {

				It("should call exiter w/ expected args", func() {
					/* arrange */
					expectedError := errors.New("dummyError")

					fakePkg := new(pkg.Fake)
					fakePkg.ParseRefReturns(nil, expectedError)

					fakeCliExiter := new(cliexiter.Fake)

					objectUnderTest := _core{
						pkg:       fakePkg,
						cliExiter: fakeCliExiter,
						os:        new(ios.Fake),
					}

					/* act */
					objectUnderTest.PkgValidate("dummyPkgRef")

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
				})
			})
			Context("pkg.ParseRef doesn't error", func() {

				It("should call pkg.Resolve w/ expected args", func() {
					/* arrange */
					providedPkgRef := "dummyPkgRef"

					expectedLookPaths := []string{"dummyWorkDir"}
					expectedPkgRef := &pkg.PkgRef{
						FullyQualifiedName: "dummyName",
						Version:            "dummyVersion",
					}

					fakePkg := new(pkg.Fake)
					fakePkg.ParseRefReturns(expectedPkgRef, nil)

					fakeCliExiter := new(cliexiter.Fake)

					fakeIOS := new(ios.Fake)
					fakeIOS.GetwdReturns(expectedLookPaths[0], nil)

					objectUnderTest := _core{
						pkg:       fakePkg,
						cliExiter: fakeCliExiter,
						os:        fakeIOS,
					}

					/* act */
					objectUnderTest.PkgValidate(providedPkgRef)

					/* assert */
					actualPkgRef, actualLookPaths := fakePkg.ResolveArgsForCall(0)
					Expect(actualLookPaths).To(Equal(expectedLookPaths))
					Expect(actualPkgRef).To(Equal(expectedPkgRef))
				})
				Context("pkg.Resolve fails", func() {

					It("should call exiter w/ expected args", func() {
						/* arrange */
						providedPkgRef := "dummyPkgRef"
						wdReturnedFromIOS := "dummyWorkDir"

						fakeIOS := new(ios.Fake)
						fakeIOS.GetwdReturns(wdReturnedFromIOS, nil)

						expectedMsg := fmt.Sprintf("Unable to resolve package '%v' from '%v'", providedPkgRef, wdReturnedFromIOS)

						fakePkg := new(pkg.Fake)
						fakePkg.ResolveReturns("", false)

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
				Context("pkg.Resolve succeeds", func() {
					It("should call pkg.Validate w/ expected args", func() {
						/* arrange */
						pkgPath := "dummyPkgName"
						wdReturnedFromIOS := "dummyWorkDir"

						fakeIOS := new(ios.Fake)
						fakeIOS.GetwdReturns(wdReturnedFromIOS, nil)

						fakePkg := new(pkg.Fake)
						fakePkg.ResolveReturns(pkgPath, true)

						objectUnderTest := _core{
							pkg:       fakePkg,
							cliExiter: new(cliexiter.Fake),
							os:        fakeIOS,
						}

						/* act */
						objectUnderTest.PkgValidate("dummyPkgRef")

						/* assert */
						Expect(fakePkg.ValidateArgsForCall(0)).To(Equal(pkgPath))
					})
					Context("pkg.Validate returns errors", func() {
						It("should call cliExiter.Exit w/ expected args", func() {
							/* arrange */
							fakePkg := new(pkg.Fake)
							fakePkg.ResolveReturns("dummyPkgName", true)

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
							pkgPath := "dummyPkgName"
							wdReturnedFromIOS := "dummyWorkDir"

							fakeIOS := new(ios.Fake)
							fakeIOS.GetwdReturns(wdReturnedFromIOS, nil)

							fakePkg := new(pkg.Fake)
							fakePkg.ResolveReturns(pkgPath, true)
							fakePkg.ValidateReturns([]error{})

							expectedExitReq := cliexiter.ExitReq{
								Message: fmt.Sprintf("%v is valid", pkgPath),
							}

							fakeCliExiter := new(cliexiter.Fake)

							objectUnderTest := _core{
								pkg:       fakePkg,
								cliExiter: fakeCliExiter,
								os:        fakeIOS,
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
})
