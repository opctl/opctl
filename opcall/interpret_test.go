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
	"github.com/opspec-io/sdk-golang/pkg/manifest"
	"path/filepath"
)

var _ = Context("OpCall", func() {
	Context("Interpret", func() {
		It("should call pkg.ParseRef w/ expected args & return errors", func() {
			/* arrange */
			providedSCGOpCall := &model.SCGOpCall{
				Pkg: &model.SCGOpCallPkg{
					Ref: "dummyPkgRef",
				},
			}

			expectedErr := errors.New("dummyError")
			fakePkg := new(pkg.Fake)
			fakePkg.ParseRefReturns(nil, expectedErr)

			objectUnderTest := _OpCall{
				pkg: fakePkg,
			}

			/* act */
			_, actualError := objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedSCGOpCall,
				"dummyOpId",
				"dummyPkgBasePath",
				"dummyRootOpId",
			)

			/* assert */
			Expect(actualError).To(Equal(expectedErr))
		})
		It("should call pkg.Resolve w/ expected args", func() {
			/* arrange */

			providedRootFSPath := "dummyRootFSPath"
			providedPkgBasePath := "dummyPkgBasePath"

			expectedPkgRef := &pkg.PkgRef{
				FullyQualifiedName: "dummyFQName",
				Version:            "dummyVersion",
			}
			expectedLookPaths := []string{providedPkgBasePath, filepath.Join(providedRootFSPath, "pkgs")}

			fakePkg := new(pkg.Fake)
			fakePkg.ParseRefReturns(expectedPkgRef, nil)

			objectUnderTest := _OpCall{
				manifest:     new(manifest.Fake),
				pkg:          fakePkg,
				pkgCachePath: filepath.Join(providedRootFSPath, "pkgs"),
			}

			/* act */
			objectUnderTest.Interpret(
				map[string]*model.Value{},
				&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
				"dummyOpId",
				providedPkgBasePath,
				"dummyRootOpId",
			)

			/* assert */
			actualPkgRef, actualLookPaths := fakePkg.ResolveArgsForCall(0)
			Expect(actualLookPaths).To(Equal(expectedLookPaths))
			Expect(actualPkgRef).To(Equal(expectedPkgRef))
		})
		Context("pkg.resolve fails", func() {
			It("should call pkg.pull w/ expected args", func() {
				/* arrange */
				providedRootFSPath := "dummyRootFSPath"
				providedPkgBasePath := "dummyPkgBasePath"
				providedSCGOpCall := &model.SCGOpCall{
					Pkg: &model.SCGOpCallPkg{
						Ref: "dummyPkgRef",
						PullCreds: &model.SCGPullCreds{
							Username: "dummyUsername",
							Password: "dummyPassword",
						},
					}}

				pkgCachePath := filepath.Join(providedRootFSPath, "pkgs")

				expectedPath := pkgCachePath
				expectedPkgRef := &pkg.PkgRef{
					FullyQualifiedName: "dummyFQName",
					Version:            "dummyVersion",
				}
				expectedPullOpts := &pkg.PullOpts{
					Username: providedSCGOpCall.Pkg.PullCreds.Username,
					Password: providedSCGOpCall.Pkg.PullCreds.Password,
				}

				fakePkg := new(pkg.Fake)
				fakePkg.ParseRefReturns(expectedPkgRef, nil)
				// error to trigger immediate return
				fakePkg.PullReturns(errors.New("dummyErr"))

				objectUnderTest := _OpCall{
					interpolater: interpolater.New(),
					manifest:     new(manifest.Fake),
					pkg:          fakePkg,
					pkgCachePath: pkgCachePath,
				}

				/* act */
				objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCGOpCall,
					"dummyOpId",
					providedPkgBasePath,
					"dummyRootOpId",
				)

				/* assert */
				actualPath, actualPkgRef, actualPullOpts := fakePkg.PullArgsForCall(0)
				Expect(actualPath).To(Equal(expectedPath))
				Expect(actualPkgRef).To(Equal(expectedPkgRef))
				Expect(actualPullOpts).To(Equal(expectedPullOpts))
			})
			Context("pkg.pull errors", func() {
				It("should return err", func() {

					/* arrange */
					expectedErr := errors.New("dummyError")

					fakePkg := new(pkg.Fake)
					fakePkg.ParseRefReturns(&pkg.PkgRef{}, nil)
					fakePkg.PullReturns(expectedErr)

					objectUnderTest := _OpCall{
						manifest: new(manifest.Fake),
						pkg:      fakePkg,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
						"dummyOpId",
						"dummyPkgBasePath",
						"dummyRootOpId",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
		})
		Context("pkg.resolve succeeds", func() {
			It("should call manifest.Unmarshal w/ expected args & return errors", func() {
				/* arrange */
				pkgPath := "dummyResolvedPkgRef"
				expectedPath := filepath.Join(pkgPath, pkg.OpDotYmlFileName)

				fakePkg := new(pkg.Fake)
				fakePkg.ParseRefReturns(&pkg.PkgRef{}, nil)
				fakePkg.ResolveReturns(pkgPath, true)
				fakeManifest := new(manifest.Fake)

				expectedErr := errors.New("dummyError")
				fakeManifest.UnmarshalReturns(nil, expectedErr)

				objectUnderTest := _OpCall{
					manifest: fakeManifest,
					pkg:      fakePkg,
					uuid:     new(iuuid.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
					"dummyOpId",
					"dummyPkgBasePath",
					"dummyRootOpId",
				)

				/* assert */
				Expect(fakeManifest.UnmarshalArgsForCall(0)).To(Equal(expectedPath))
				Expect(actualErr).To(Equal(expectedErr))
			})
			Context("manifest.Unmarshal succeeds", func() {
				It("should call inputs.Interpret w/ expected inputs", func() {
					/* arrange */
					providedScope := map[string]*model.Value{
						"dummyScopeRef1Name": {String: new(string)},
					}
					expectedScope := providedScope

					expectedInputArgs := map[string]string{"dummySCGOpCallInputName": "dummyScgOpCallInputValue"}
					providedSCGOpCall := &model.SCGOpCall{
						Inputs: expectedInputArgs,
						Pkg:    &model.SCGOpCallPkg{},
					}

					resolvedPkgPath := "resolvedPkgPath"
					expectedPkgPath := resolvedPkgPath

					fakePkg := new(pkg.Fake)
					fakePkg.ParseRefReturns(&pkg.PkgRef{}, nil)
					fakePkg.ResolveReturns(resolvedPkgPath, true)
					fakeManifest := new(manifest.Fake)

					expectedInputParams := map[string]*model.Param{
						"dummyParam1Name": {String: &model.StringParam{}},
					}

					returnedManifest := &model.PkgManifest{
						Inputs: expectedInputParams,
					}
					fakeManifest.UnmarshalReturns(returnedManifest, nil)

					fakeArgs := new(inputs.Fake)

					objectUnderTest := _OpCall{
						manifest: fakeManifest,
						pkg:      fakePkg,
						uuid:     new(iuuid.Fake),
						inputs:   fakeArgs,
					}

					/* act */
					objectUnderTest.Interpret(
						providedScope,
						providedSCGOpCall,
						"dummyOpId",
						"dummyPkgBasePath",
						"dummyRootOpId",
					)

					/* assert */
					actualSCGArgs, actualSCGInputs, actualPkgPath, actualScope := fakeArgs.InterpretArgsForCall(0)
					Expect(actualScope).To(Equal(expectedScope))
					Expect(actualSCGArgs).To(Equal(expectedInputArgs))
					Expect(actualPkgPath).To(Equal(expectedPkgPath))
					Expect(actualSCGInputs).To(Equal(expectedInputParams))
				})
			})
		})
	})
})
