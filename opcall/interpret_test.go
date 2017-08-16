package opcall

import (
	"errors"
	"github.com/golang-interfaces/satori-go.uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall/inputs"
	"github.com/opspec-io/sdk-golang/pkg"
	"path/filepath"
)

var _ = Context("OpCall", func() {
	Context("Interpret", func() {
		It("should call pkg.Resolve w/ expected args", func() {
			/* arrange */
			providedRootFSPath := "dummyRootFSPath"
			providedSCGOpCall := &model.SCGOpCall{
				Pkg: &model.SCGOpCallPkg{
					Ref: "dummyPkgRef",
					PullCreds: &model.SCGPullCreds{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				},
			}

			expectedPkgRef := providedSCGOpCall.Pkg.Ref

			fakeInterpolater := new(interpolater.Fake)
			fakeInterpolater.InterpolateReturnsOnCall(0, providedSCGOpCall.Pkg.PullCreds.Username)
			fakeInterpolater.InterpolateReturnsOnCall(1, providedSCGOpCall.Pkg.PullCreds.Password)

			fakePkg := new(pkg.Fake)

			expectedPkgProviders := []pkg.Provider{
				new(pkg.FakeProvider),
				new(pkg.FakeProvider),
			}
			fakePkg.NewFSProviderReturns(expectedPkgProviders[0])
			fakePkg.NewGitProviderReturns(expectedPkgProviders[1])

			// error to trigger immediate return
			fakePkg.ResolveReturns(nil, errors.New("dummyError"))

			objectUnderTest := _OpCall{
				interpolater: fakeInterpolater,
				pkg:          fakePkg,
				pkgCachePath: filepath.Join(providedRootFSPath, "pkgs"),
			}

			/* act */
			objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedSCGOpCall,
				"dummyOpId",
				new(pkg.FakeHandle),
				"dummyRootOpId",
			)

			/* assert */
			actualPkgRef, actualPkgProviders := fakePkg.ResolveArgsForCall(0)
			Expect(actualPkgRef).To(Equal(expectedPkgRef))
			Expect(actualPkgProviders).To(Equal(expectedPkgProviders))
		})
		It("should call pkg.NewFSProvider w/ expected args", func() {
			/* arrange */
			providedParentPkgHandle := new(pkg.FakeHandle)

			fakePkg := new(pkg.Fake)
			// error to trigger immediate return
			fakePkg.ResolveReturns(nil, errors.New("dummyError"))

			objectUnderTest := _OpCall{
				interpolater: new(interpolater.Fake),
				pkg:          fakePkg,
			}

			/* act */
			objectUnderTest.Interpret(
				map[string]*model.Value{},
				&model.SCGOpCall{
					Pkg: &model.SCGOpCallPkg{
						Ref: "dummyPkgRef",
					},
				},
				"dummyOpId",
				providedParentPkgHandle,
				"dummyRootOpId",
			)

			/* assert */
			Expect(fakePkg.NewFSProviderArgsForCall(0)).To(ConsistOf(filepath.Dir(providedParentPkgHandle.Ref())))
		})
		It("should call pkg.NewGitProvider w/ expected args", func() {
			/* arrange */
			providedPkgCachePath := "dummyPkgCachePath"

			fakeInterpolater := new(interpolater.Fake)
			expectedPullCreds := &model.PullCreds{Username: "dummyUsername", Password: "dummyPassword"}
			fakeInterpolater.InterpolateReturnsOnCall(0, expectedPullCreds.Username)
			fakeInterpolater.InterpolateReturnsOnCall(1, expectedPullCreds.Password)

			fakePkg := new(pkg.Fake)
			// error to trigger immediate return
			fakePkg.ResolveReturns(nil, errors.New("dummyError"))

			objectUnderTest := _OpCall{
				interpolater: fakeInterpolater,
				pkg:          fakePkg,
				pkgCachePath: providedPkgCachePath,
			}

			/* act */
			objectUnderTest.Interpret(
				map[string]*model.Value{},
				&model.SCGOpCall{
					Pkg: &model.SCGOpCallPkg{
						Ref:       "dummyPkgRef",
						PullCreds: &model.SCGPullCreds{},
					},
				},
				"dummyOpId",
				new(pkg.FakeHandle),
				"dummyRootOpId",
			)

			/* assert */
			actualBasePath,
				actualPullCreds := fakePkg.NewGitProviderArgsForCall(0)

			Expect(actualBasePath).To(Equal(providedPkgCachePath))
			Expect(actualPullCreds).To(Equal(expectedPullCreds))
		})
		Context("pkg.Resolve errs", func() {
			It("should return err", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakePkg := new(pkg.Fake)
				fakePkg.ResolveReturns(nil, expectedErr)

				objectUnderTest := _OpCall{
					pkg:  fakePkg,
					uuid: new(iuuid.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
					"dummyOpId",
					new(pkg.FakeHandle),
					"dummyRootOpId",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("pkg.Resolve doesn't err", func() {
			It("should call pkg.GetManifest w/ expected args", func() {
				/* arrange */
				fakePkgHandle := new(pkg.FakeHandle)

				fakePkg := new(pkg.Fake)
				fakePkg.ResolveReturns(fakePkgHandle, nil)

				expectedErr := errors.New("dummyError")
				// err to trigger immediate return
				fakePkg.GetManifestReturns(nil, expectedErr)

				objectUnderTest := _OpCall{
					pkg:  fakePkg,
					uuid: new(iuuid.Fake),
				}

				/* act */
				objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
					"dummyOpId",
					new(pkg.FakeHandle),
					"dummyRootOpId",
				)

				/* assert */
				actualHandle := fakePkg.GetManifestArgsForCall(0)
				Expect(actualHandle).To(Equal(fakePkgHandle))
			})
			Context("pkg.GetManifest errs", func() {
				It("should return err", func() {
					/* arrange */
					expectedErr := errors.New("dummyError")
					fakePkg := new(pkg.Fake)
					fakePkg.GetManifestReturns(nil, expectedErr)

					objectUnderTest := _OpCall{
						pkg:  fakePkg,
						uuid: new(iuuid.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
						"dummyOpId",
						new(pkg.FakeHandle),
						"dummyRootOpId",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("pkg.GetManifest doesn't err", func() {
				It("should call inputs.Interpret w/ expected inputs", func() {
					/* arrange */
					providedScope := map[string]*model.Value{
						"dummyScopeRef1Name": {String: new(string)},
					}
					expectedScope := providedScope

					expectedInputArgs := map[string]interface{}{"dummySCGOpCallInputName": "dummyScgOpCallInputValue"}

					providedSCGOpCall := &model.SCGOpCall{
						Inputs: expectedInputArgs,
						Pkg:    &model.SCGOpCallPkg{},
					}

					fakePkgHandle := new(pkg.FakeHandle)

					fakePkg := new(pkg.Fake)
					fakePkg.ResolveReturns(fakePkgHandle, nil)

					expectedInputParams := map[string]*model.Param{
						"dummyParam1Name": {String: &model.StringParam{}},
					}

					returnedManifest := &model.PkgManifest{
						Inputs: expectedInputParams,
					}
					fakePkg.GetManifestReturns(returnedManifest, nil)

					fakeArgs := new(inputs.Fake)

					objectUnderTest := _OpCall{
						pkg:    fakePkg,
						uuid:   new(iuuid.Fake),
						inputs: fakeArgs,
					}

					/* act */
					objectUnderTest.Interpret(
						providedScope,
						providedSCGOpCall,
						"dummyOpId",
						fakePkgHandle,
						"dummyRootOpId",
					)

					/* assert */
					actualSCGArgs, actualSCGInputs, actualPkgRef, actualScope := fakeArgs.InterpretArgsForCall(0)
					Expect(actualScope).To(Equal(expectedScope))
					Expect(actualSCGArgs).To(Equal(expectedInputArgs))
					Expect(actualPkgRef).To(Equal(fakePkgHandle.Ref()))
					Expect(actualSCGInputs).To(Equal(expectedInputParams))
				})
			})
		})
	})
})
