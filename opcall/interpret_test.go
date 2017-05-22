package opcall

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/satori-go.uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall/paramvalidator"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
	"path/filepath"
)

var _ = Context("OpCall", func() {
	Context("Interpret", func() {
	})
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
			map[string]*model.Data{},
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
			map[string]*model.Data{},
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
				map[string]*model.Data{},
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
					map[string]*model.Data{},
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
				map[string]*model.Data{},
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
			It("should call paramvalidator.Validate w/ expected args", func() {
				/* arrange */
				providedScopeVar1String := "val1"
				providedScopeVarFile := "val2"
				providedScopeVar3Dir := "val3"
				providedScopeVar4Socket := "val4"
				providedScopeVar5Number := float64(5)
				providedScope := map[string]*model.Data{
					"name1": {String: &providedScopeVar1String},
					"name2": {File: &providedScopeVarFile},
					"name3": {Dir: &providedScopeVar3Dir},
					"name4": {Socket: &providedScopeVar4Socket},
					"name5": {Number: &providedScopeVar5Number},
				}
				providedSCGOpCall := &model.SCGOpCall{
					Pkg: &model.SCGOpCallPkg{},
					Inputs: map[string]string{
						"name1": "",
						"name2": "",
						"name3": "",
						"name4": "",
						"name5": "",
						"name6": "",
						"name7": "",
					},
				}

				returnedPkgInput6Default := float64(6)
				returnedPkgInput7Default := "seven"
				returnedPkg := &model.PkgManifest{
					Inputs: map[string]*model.Param{
						"name1": {String: &model.StringParam{}},
						"name2": {File: &model.FileParam{}},
						"name3": {Dir: &model.DirParam{}},
						"name4": {Socket: &model.SocketParam{}},
						"name5": {Number: &model.NumberParam{}},
						"name6": {Number: &model.NumberParam{Default: &returnedPkgInput6Default}},
						"name7": {String: &model.StringParam{Default: &returnedPkgInput7Default}},
					},
				}
				fakePkg := new(pkg.Fake)
				fakePkg.ParseRefReturns(&pkg.PkgRef{}, nil)
				fakePkg.ResolveReturns("", true)

				fakeManifest := new(manifest.Fake)
				fakeManifest.UnmarshalReturns(returnedPkg, nil)

				expectedCalls := map[model.Data]*model.Param{
					// from scope
					*providedScope["name1"]: returnedPkg.Inputs["name1"],
					*providedScope["name2"]: returnedPkg.Inputs["name2"],
					*providedScope["name3"]: returnedPkg.Inputs["name3"],
					*providedScope["name4"]: returnedPkg.Inputs["name4"],
					*providedScope["name5"]: returnedPkg.Inputs["name5"],
					// from defaults
					model.Data{
						Number: returnedPkg.Inputs["name6"].Number.Default,
					}: returnedPkg.Inputs["name6"],
					model.Data{
						String: returnedPkg.Inputs["name7"].String.Default,
					}: returnedPkg.Inputs["name7"],
				}

				fakeParamValidator := new(paramvalidator.Fake)
				// error to trigger immediate return
				fakeParamValidator.ValidateReturns([]error{errors.New("dummyError")})

				objectUnderTest := _OpCall{
					manifest: fakeManifest,
					pkg:      fakePkg,
					uuid:     new(iuuid.Fake),
					validate: fakeParamValidator,
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
				actualCalls := map[model.Data]*model.Param{}
				for i := 0; i < fakeParamValidator.ValidateCallCount(); i++ {
					actualVarData, actualParam := fakeParamValidator.ValidateArgsForCall(i)
					actualCalls[*actualVarData] = actualParam
				}
				Expect(actualCalls).To(Equal(expectedCalls))
			})

			Context("paramvalidator.Validate errors", func() {
				It("should return error", func() {
					/* arrange */
					opReturnedFromPkg := &model.PkgManifest{
						Inputs: map[string]*model.Param{
							"dummyVar1Name": {
								String: &model.StringParam{
									IsSecret: true,
								},
							},
						},
					}
					fakePkg := new(pkg.Fake)
					fakePkg.ParseRefReturns(&pkg.PkgRef{}, nil)
					fakePkg.ResolveReturns("", true)

					fakeManifest := new(manifest.Fake)
					fakeManifest.UnmarshalReturns(opReturnedFromPkg, nil)

					fakeParamValidator := new(paramvalidator.Fake)

					errorReturnedFromValidate := "validationError0"
					fakeParamValidator.ValidateReturns([]error{errors.New(errorReturnedFromValidate)})

					expectedErr := fmt.Errorf(`
-
  validation of the following input failed:

  Name: %v
  Value: %v
  Error(s):
    - %v

-`, "dummyVar1Name", "************", errorReturnedFromValidate)

					objectUnderTest := _OpCall{
						manifest: fakeManifest,
						pkg:      fakePkg,
						uuid:     new(iuuid.Fake),
						validate: fakeParamValidator,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						map[string]*model.Data{},
						&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
						"dummyOpId",
						"dummyPkgBasePath",
						"dummyRootOpId",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("paramvalidator.Validate doesn't error", func() {
				Context("It should return expected result", func() {

				})
			})
		})
	})
})
