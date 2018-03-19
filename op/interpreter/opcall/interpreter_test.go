package opcall

import (
	"context"
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/dotyml"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression"
	"github.com/opspec-io/sdk-golang/op/interpreter/opcall/inputs"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter("")).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		Context("called w/ opspec test-suite scenarios", func() {
			It("should return result fulfilling scenario.call.expect", func() {
				tempDir, err := ioutil.TempDir("", "")
				if nil != err {
					panic(err)
				}
				rootPath := "../../../github.com/opspec-io/test-suite/scenarios/pkg"

				filepath.Walk(rootPath,
					func(path string, info os.FileInfo, err error) error {
						if info.IsDir() {
							scenariosDotYmlFilePath := filepath.Join(path, "scenarios.yml")
							if _, err := os.Stat(scenariosDotYmlFilePath); nil == err {
								/* arrange */
								absPkgPath, err := filepath.Abs(path)
								if nil != err {
									panic(fmt.Errorf("error getting absPkgPath %v; error was %v", path, err))
								}

								data := data.New()
								opDirHandle, err := data.Resolve(
									context.Background(),
									absPkgPath,
									data.NewFSProvider(),
								)
								if nil != err {
									panic(fmt.Errorf("error resolving pkg for %v; error was %v", path, err))
								}

								scenariosDotYmlBytes, err := ioutil.ReadFile(scenariosDotYmlFilePath)
								if nil != err {
									panic(err)
								}

								var scenarioDotYml []struct {
									Name      string
									Interpret *struct {
										Expect string
										Scope  map[string]*model.Value
									}
								}
								if err := yaml.Unmarshal(scenariosDotYmlBytes, &scenarioDotYml); nil != err {
									panic(fmt.Errorf("error unmarshalling scenario.yml for %v; error was %v", path, err))
								}

								for _, scenario := range scenarioDotYml {
									if nil != scenario.Interpret {
										scgOpCall := &model.SCGOpCall{
											Pkg: &model.SCGOpCallPkg{
												Ref: absPkgPath,
											},
											Inputs: map[string]interface{}{},
										}

										for name := range scenario.Interpret.Scope {
											// map as passed
											scgOpCall.Inputs[name] = ""
										}

										/* act */
										objectUnderTest := NewInterpreter(tempDir)
										_, actualErr := objectUnderTest.Interpret(
											scenario.Interpret.Scope,
											scgOpCall,
											"",
											opDirHandle,
											"",
										)

										/* assert */
										description := fmt.Sprintf("pkg: '%v'\nscenario: '%v'", path, scenario.Name)
										switch expect := scenario.Interpret.Expect; expect {
										case "success":
											Expect(actualErr).To(BeNil(), description)
										case "failure":
											Expect(actualErr).To(Not(BeNil()), description)
										}
									}
								}
							}
						}
						return nil
					})
			})
		})
		It("should call pkg.NewFSProvider w/ expected args", func() {
			/* arrange */
			providedParentOpDirHandle := new(data.FakeHandle)
			providedParentOpDirHandle.PathReturns(new(string))

			fakeData := new(data.Fake)
			// error to trigger immediate return
			fakeData.ResolveReturns(nil, errors.New("dummyError"))

			objectUnderTest := _interpreter{
				expression: new(expression.Fake),
				data:       fakeData,
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
				providedParentOpDirHandle,
				"dummyRootOpId",
			)

			/* assert */
			Expect(fakeData.NewFSProviderArgsForCall(0)).To(ConsistOf(filepath.Dir(providedParentOpDirHandle.Ref())))
		})
		Context("scgOpCall.Pkg.PullCreds is nil", func() {
			It("should call pkg.NewGitProvider w/ expected args", func() {
				/* arrange */
				providedParentOpDirHandle := new(data.FakeHandle)
				providedParentOpDirHandle.PathReturns(new(string))

				providedPkgCachePath := "dummyPkgCachePath"

				fakeData := new(data.Fake)
				// error to trigger immediate return
				fakeData.ResolveReturns(nil, errors.New("dummyError"))

				objectUnderTest := _interpreter{
					data:         fakeData,
					pkgCachePath: providedPkgCachePath,
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
					providedParentOpDirHandle,
					"dummyRootOpId",
				)

				/* assert */
				actualBasePath,
					actualPullCreds := fakeData.NewGitProviderArgsForCall(0)

				Expect(actualBasePath).To(Equal(providedPkgCachePath))
				Expect(actualPullCreds).To(BeNil())
			})
		})
		Context("scgOpCall.Pkg.PullCreds isn't nil", func() {
			Context("string.Interpret errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeExpression := new(expression.Fake)
					interpretError := errors.New("dummyError")
					fakeExpression.EvalToStringReturns(nil, interpretError)

					objectUnderTest := _interpreter{
						expression: fakeExpression,
					}

					/* act */
					_, actualError := objectUnderTest.Interpret(
						map[string]*model.Value{},
						&model.SCGOpCall{
							Pkg: &model.SCGOpCallPkg{
								PullCreds: &model.SCGPullCreds{},
							},
						},
						"dummyOpId",
						new(data.FakeHandle),
						"dummyRootOpId",
					)

					/* assert */
					Expect(actualError).To(Equal(interpretError))
				})
			})
			Context("string.Interpret doesn't err", func() {
				It("should call pkg.NewGitProvider w/ expected args", func() {
					/* arrange */
					providedParentOpDirHandle := new(data.FakeHandle)
					providedParentOpDirHandle.PathReturns(new(string))

					providedPkgCachePath := "dummyPkgCachePath"

					fakeExpression := new(expression.Fake)
					expectedPullCreds := &model.PullCreds{Username: "dummyUsername", Password: "dummyPassword"}
					fakeExpression.EvalToStringReturnsOnCall(0, &model.Value{String: &expectedPullCreds.Username}, nil)
					fakeExpression.EvalToStringReturnsOnCall(1, &model.Value{String: &expectedPullCreds.Password}, nil)

					fakeData := new(data.Fake)
					// error to trigger immediate return
					fakeData.ResolveReturns(nil, errors.New("dummyError"))

					objectUnderTest := _interpreter{
						expression:   fakeExpression,
						data:         fakeData,
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
						providedParentOpDirHandle,
						"dummyRootOpId",
					)

					/* assert */
					actualBasePath,
						actualPullCreds := fakeData.NewGitProviderArgsForCall(0)

					Expect(actualBasePath).To(Equal(providedPkgCachePath))
					Expect(actualPullCreds).To(Equal(expectedPullCreds))
				})
			})
		})
		It("should call pkg.Resolve w/ expected args", func() {
			/* arrange */
			providedParentOpDirHandle := new(data.FakeHandle)
			providedParentOpDirHandle.PathReturns(new(string))

			providedRootFSPath := "dummyRootFSPath"
			providedSCGOpCall := &model.SCGOpCall{
				Pkg: &model.SCGOpCallPkg{
					Ref: "dummyPkgRef",
				},
			}

			expectedPkgRef := providedSCGOpCall.Pkg.Ref

			fakeData := new(data.Fake)

			expectedPkgProviders := []data.Provider{
				new(data.FakeProvider),
				new(data.FakeProvider),
			}
			fakeData.NewFSProviderReturns(expectedPkgProviders[0])
			fakeData.NewGitProviderReturns(expectedPkgProviders[1])

			// error to trigger immediate return
			fakeData.ResolveReturns(nil, errors.New("dummyError"))

			objectUnderTest := _interpreter{
				data:         fakeData,
				pkgCachePath: filepath.Join(providedRootFSPath, "pkgs"),
			}

			/* act */
			objectUnderTest.Interpret(
				map[string]*model.Value{},
				providedSCGOpCall,
				"dummyOpId",
				providedParentOpDirHandle,
				"dummyRootOpId",
			)

			/* assert */
			actualCtx,
				actualPkgRef,
				actualPkgProviders := fakeData.ResolveArgsForCall(0)

			Expect(actualCtx).To(Equal(context.TODO()))
			Expect(actualPkgRef).To(Equal(expectedPkgRef))
			Expect(actualPkgProviders).To(Equal(expectedPkgProviders))
		})
		Context("pkg.Resolve errs", func() {
			It("should return err", func() {
				/* arrange */
				providedParentOpDirHandle := new(data.FakeHandle)
				providedParentOpDirHandle.PathReturns(new(string))

				expectedErr := errors.New("dummyError")
				fakeData := new(data.Fake)
				fakeData.ResolveReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					data:                fakeData,
					uniqueStringFactory: new(uniquestring.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
					"dummyOpId",
					providedParentOpDirHandle,
					"dummyRootOpId",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("pkg.Resolve doesn't err", func() {
			It("should call pkg.GetManifest w/ expected args", func() {
				/* arrange */
				providedParentOpDirHandle := new(data.FakeHandle)
				providedParentOpDirHandle.PathReturns(new(string))

				fakeDataHandle := new(data.FakeHandle)

				fakeData := new(data.Fake)
				fakeData.ResolveReturns(fakeDataHandle, nil)

				fakeDotYmlGetter := new(dotyml.FakeGetter)
				expectedErr := errors.New("dummyError")
				// err to trigger immediate return
				fakeDotYmlGetter.GetReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					data:                fakeData,
					dotYmlGetter:        fakeDotYmlGetter,
					uniqueStringFactory: new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
					"dummyOpId",
					providedParentOpDirHandle,
					"dummyRootOpId",
				)

				/* assert */
				actualCtx,
					actualHandle := fakeDotYmlGetter.GetArgsForCall(0)

				Expect(actualCtx).To(Equal(context.TODO()))
				Expect(actualHandle).To(Equal(fakeDataHandle))
			})
			Context("pkg.GetManifest errs", func() {
				It("should return err", func() {
					/* arrange */
					providedParentOpDirHandle := new(data.FakeHandle)
					providedParentOpDirHandle.PathReturns(new(string))

					expectedErr := errors.New("dummyError")
					fakeDotYmlGetter := new(dotyml.FakeGetter)
					fakeDotYmlGetter.GetReturns(nil, expectedErr)

					objectUnderTest := _interpreter{
						data:                new(data.Fake),
						dotYmlGetter:        fakeDotYmlGetter,
						uniqueStringFactory: new(uniquestring.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
						"dummyOpId",
						providedParentOpDirHandle,
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

					providedOpId := "dummyOpId"

					providedParentOpDirHandle := new(data.FakeHandle)
					parentOpDirPath := "dummyParentOpDirPath"
					providedParentOpDirHandle.PathReturns(&parentOpDirPath)

					fakeDataHandle := new(data.FakeHandle)
					pkgPath := "dummyPkgPath"
					fakeDataHandle.PathReturns(&pkgPath)

					fakeData := new(data.Fake)
					fakeData.ResolveReturns(fakeDataHandle, nil)

					expectedInputParams := map[string]*model.Param{
						"dummyParam1Name": {String: &model.StringParam{}},
					}

					fakeDotYmlGetter := new(dotyml.FakeGetter)
					returnedManifest := &model.PkgManifest{
						Inputs: expectedInputParams,
					}
					fakeDotYmlGetter.GetReturns(returnedManifest, nil)

					fakeInputsInterpreter := new(inputs.FakeInterpreter)

					dcgScratchDir := "dummyDCGScratchDir"

					objectUnderTest := _interpreter{
						dcgScratchDir:       dcgScratchDir,
						data:                fakeData,
						dotYmlGetter:        fakeDotYmlGetter,
						uniqueStringFactory: new(uniquestring.Fake),
						inputsInterpreter:   fakeInputsInterpreter,
					}

					/* act */
					objectUnderTest.Interpret(
						providedScope,
						providedSCGOpCall,
						providedOpId,
						providedParentOpDirHandle,
						"dummyRootOpId",
					)

					/* assert */
					actualSCGArgs,
						actualSCGInputs,
						actualParentOpDirHandle,
						actualPkgRef,
						actualScope,
						actualOpScratchDir := fakeInputsInterpreter.InterpretArgsForCall(0)

					Expect(actualScope).To(Equal(expectedScope))
					Expect(actualSCGArgs).To(Equal(expectedInputArgs))
					Expect(actualParentOpDirHandle).To(Equal(providedParentOpDirHandle))
					Expect(actualPkgRef).To(Equal(pkgPath))
					Expect(actualSCGInputs).To(Equal(expectedInputParams))
					Expect(actualOpScratchDir).To(Equal(filepath.Join(dcgScratchDir, providedOpId)))

				})
			})
		})
	})
})
