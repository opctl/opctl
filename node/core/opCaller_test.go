package core

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/validate"
	"github.com/pkg/errors"
	"path/filepath"
	"time"
)

var _ = Context("opCaller", func() {
	Context("newOpCaller", func() {
		It("should return opCaller", func() {
			/* arrange/act/assert */
			Expect(newOpCaller(
				new(pkg.Fake),
				new(pubsub.Fake),
				newDCGNodeRepo(),
				new(fakeCaller),
				new(uniquestring.Fake),
				new(validate.Fake),
				"dummyRootFSPath",
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		Context("deprecated pkg format", func() {
			It("should call pkg.Get w/ expected args", func() {
				/* arrange */
				providedPkgBasePath := "dummyPkgBasePath"
				providedSCGOpCall := &model.SCGOpCall{
					Ref: "dummySCGOpCallPkgRef",
				}

				expectedPkgRef := filepath.Join(providedPkgBasePath, providedSCGOpCall.Ref)

				fakePkg := new(pkg.Fake)
				fakePkg.GetReturns(nil, errors.New("dummyError"))

				objectUnderTest := newOpCaller(
					fakePkg,
					new(pubsub.Fake),
					new(fakeDCGNodeRepo),
					new(fakeCaller),
					new(uniquestring.Fake),
					new(validate.Fake),
					"dummyRootFSPath",
				)

				/* act */
				objectUnderTest.Call(
					map[string]*model.Data{},
					"dummyOpId",
					providedPkgBasePath,
					"dummyRootOpId",
					providedSCGOpCall,
				)

				/* assert */
				Expect(fakePkg.GetArgsForCall(0)).To(Equal(expectedPkgRef))
			})
		})
		It("should call pkg.ParseRef w/ expected args", func() {
			/* arrange */
			providedPkgBasePath := "dummyPkgBasePath"
			providedSCGOpCall := &model.SCGOpCall{
				Pkg: &model.SCGOpCallPkg{
					Ref: "dummySCGOpCallPkgRef",
				},
			}

			expectedPkgRef := providedSCGOpCall.Pkg.Ref

			fakePkg := new(pkg.Fake)
			// error to trigger immediate return
			fakePkg.ParseRefReturns(nil, errors.New("dummyError"))

			objectUnderTest := newOpCaller(
				fakePkg,
				new(pubsub.Fake),
				new(fakeDCGNodeRepo),
				new(fakeCaller),
				new(uniquestring.Fake),
				new(validate.Fake),
				"dummyRootFSPath",
			)

			/* act */
			objectUnderTest.Call(
				map[string]*model.Data{},
				"dummyOpId",
				providedPkgBasePath,
				"dummyRootOpId",
				providedSCGOpCall,
			)

			/* assert */
			Expect(fakePkg.ParseRefArgsForCall(0)).To(Equal(expectedPkgRef))
		})
		Context("pkg.ParseRef errors", func() {
			It("should return error", func() {

				/* arrange */
				providedPkgBasePath := "dummyPkgBasePath"
				providedSCGOpCall := &model.SCGOpCall{
					Pkg: &model.SCGOpCallPkg{
						Ref: "dummySCGOpCallPkgRef",
					},
				}

				expectedError := errors.New("dummyError")

				fakePkg := new(pkg.Fake)
				fakePkg.ParseRefReturns(nil, expectedError)

				objectUnderTest := newOpCaller(
					fakePkg,
					new(pubsub.Fake),
					new(fakeDCGNodeRepo),
					new(fakeCaller),
					new(uniquestring.Fake),
					new(validate.Fake),
					"dummyRootFSPath",
				)

				/* act */
				actualError := objectUnderTest.Call(
					map[string]*model.Data{},
					"dummyOpId",
					providedPkgBasePath,
					"dummyRootOpId",
					providedSCGOpCall,
				)

				/* assert */
				Expect(actualError).To(Equal(actualError))
			})
		})
		Context("pkg.ParseRef doesn't error", func() {
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

				objectUnderTest := newOpCaller(
					fakePkg,
					new(pubsub.Fake),
					new(fakeDCGNodeRepo),
					new(fakeCaller),
					new(uniquestring.Fake),
					new(validate.Fake),
					providedRootFSPath,
				)

				/* act */
				objectUnderTest.Call(
					map[string]*model.Data{},
					"dummyOpId",
					providedPkgBasePath,
					"dummyRootOpId",
					&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
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

					expectedPath := filepath.Join(providedRootFSPath, "pkgs")
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

					objectUnderTest := newOpCaller(
						fakePkg,
						new(pubsub.Fake),
						new(fakeDCGNodeRepo),
						new(fakeCaller),
						new(uniquestring.Fake),
						new(validate.Fake),
						providedRootFSPath,
					)

					/* act */
					objectUnderTest.Call(
						map[string]*model.Data{},
						"dummyOpId",
						providedPkgBasePath,
						"dummyRootOpId",
						providedSCGOpCall,
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
						fakePkg.PullReturns(expectedErr)

						objectUnderTest := newOpCaller(
							fakePkg,
							new(pubsub.Fake),
							new(fakeDCGNodeRepo),
							new(fakeCaller),
							new(uniquestring.Fake),
							new(validate.Fake),
							"dummyRootFSPath",
						)

						/* act */
						actualErr := objectUnderTest.Call(
							map[string]*model.Data{},
							"dummyOpId",
							"dummyPkgBasePath",
							"dummyRootOpId",
							&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
			})
			Context("pkg.resolve succeeds", func() {
				It("should call dcgNodeRepo.add w/ expected args", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{}
					providedOpId := "dummyOpId"
					providedRootOpId := "dummyRootOpId"
					providedSCGOpCall := &model.SCGOpCall{
						Pkg: &model.SCGOpCallPkg{
							Ref: "dummyPkgRef",
						}}

					resolvedPkgPath := "dummyResolvedPkgPath"

					expectedDCGNodeDescriptor := &dcgNodeDescriptor{
						Id:       providedOpId,
						PkgRef:   resolvedPkgPath,
						RootOpId: providedRootOpId,
						Op:       &dcgOpDescriptor{},
					}

					fakePkg := new(pkg.Fake)
					fakePkg.ResolveReturns(resolvedPkgPath, true)
					fakePkg.GetReturns(&model.PkgManifest{}, nil)

					fakeDCGNodeRepo := new(fakeDCGNodeRepo)

					objectUnderTest := newOpCaller(
						fakePkg,
						new(pubsub.Fake),
						fakeDCGNodeRepo,
						new(fakeCaller),
						new(uniquestring.Fake),
						new(validate.Fake),
						"dummyRootFSPath",
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOpId,
						resolvedPkgPath,
						providedRootOpId,
						providedSCGOpCall,
					)

					/* assert */
					Expect(fakeDCGNodeRepo.AddArgsForCall(0)).To(Equal(expectedDCGNodeDescriptor))
				})
				It("should call pkg.Get w/ expected args", func() {
					/* arrange */
					providedUsernameString := "name1Value"
					providedInboundScope := map[string]*model.Data{
						"username": {String: &providedUsernameString},
					}
					providedOpId := "dummyOpId"
					providedPkgBasePath := "dummyPkgBasePath"
					providedRootOpId := "dummyRootOpId"
					providedSCGOpCall := &model.SCGOpCall{
						Pkg: &model.SCGOpCallPkg{
							Ref: "dummySCGOpCallPkgRef",
							PullCreds: &model.SCGPullCreds{
								Username: "$(username)",
								Password: "dummyPassword",
							},
						},
					}

					resolvedPkgRef := "dummyResolvedPkgRef"
					expectedPkgRef := resolvedPkgRef

					fakePkg := new(pkg.Fake)
					fakePkg.ResolveReturns(resolvedPkgRef, true)
					fakePkg.GetReturns(nil, errors.New("dummyError"))

					objectUnderTest := newOpCaller(
						fakePkg,
						new(pubsub.Fake),
						new(fakeDCGNodeRepo),
						new(fakeCaller),
						new(uniquestring.Fake),
						new(validate.Fake),
						"dummyRootFSPath",
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOpId,
						providedPkgBasePath,
						providedRootOpId,
						providedSCGOpCall,
					)

					/* assert */
					Expect(fakePkg.GetArgsForCall(0)).To(Equal(expectedPkgRef))
				})
				Context("pkg.Get errors", func() {
					It("should call pubSub.Publish w/ expected args", func() {
						/* arrange */
						providedInboundScope := map[string]*model.Data{}
						providedOpId := "dummyOpId"
						providedPkgBasePath := "dummyPkgBasePath"
						providedRootOpId := "dummyRootOpId"
						providedSCGOpCall := &model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}}

						resolvedPkgRef := "dummyResolvedPkgRef"

						expectedEvent := &model.Event{
							Timestamp: time.Now().UTC(),
							OpErred: &model.OpErredEvent{
								Msg:      "dummyError",
								OpId:     providedOpId,
								PkgRef:   resolvedPkgRef,
								RootOpId: providedRootOpId,
							},
						}

						fakePkg := new(pkg.Fake)
						fakePkg.ResolveReturns(resolvedPkgRef, true)
						fakePkg.GetReturns(
							&model.PkgManifest{},
							errors.New(expectedEvent.OpErred.Msg),
						)

						fakeDCGNodeRepo := new(fakeDCGNodeRepo)
						fakePkg.ResolveReturns(resolvedPkgRef, true)
						fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

						fakePubSub := new(pubsub.Fake)

						objectUnderTest := newOpCaller(
							fakePkg,
							fakePubSub,
							fakeDCGNodeRepo,
							new(fakeCaller),
							new(uniquestring.Fake),
							new(validate.Fake),
							"dummyRootFSPath",
						)

						/* act */
						objectUnderTest.Call(
							providedInboundScope,
							providedOpId,
							providedPkgBasePath,
							providedRootOpId,
							providedSCGOpCall,
						)

						/* assert */
						actualEvent := fakePubSub.PublishArgsForCall(0)

						// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
						Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
						// set temporal fields to expected vals since they're already asserted
						actualEvent.Timestamp = expectedEvent.Timestamp

						Expect(actualEvent).To(Equal(expectedEvent))
					})
				})
				Context("pkg.Get doesn't error", func() {
					It("should call validate.Param w/ expected args", func() {
						/* arrange */
						providedInboundVar1String := "val1"
						providedInboundVarFile := "val2"
						providedInboundVar3Dir := "val3"
						providedInboundVar4Socket := "val4"
						providedInboundVar5Number := float64(5)
						providedInboundScope := map[string]*model.Data{
							"name1": {String: &providedInboundVar1String},
							"name2": {File: &providedInboundVarFile},
							"name3": {Dir: &providedInboundVar3Dir},
							"name4": {Socket: &providedInboundVar4Socket},
							"name5": {Number: &providedInboundVar5Number},
						}
						providedOpId := "dummyOpId"
						providedPkgBasePath := "dummyPkgBasePath"
						providedRootOpId := "dummyRootOpId"
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
						fakePkg.ResolveReturns("", true)
						fakePkg.GetReturns(returnedPkg, nil)

						expectedCalls := map[model.Data]*model.Param{
							// from scope
							*providedInboundScope["name1"]: returnedPkg.Inputs["name1"],
							*providedInboundScope["name2"]: returnedPkg.Inputs["name2"],
							*providedInboundScope["name3"]: returnedPkg.Inputs["name3"],
							*providedInboundScope["name4"]: returnedPkg.Inputs["name4"],
							*providedInboundScope["name5"]: returnedPkg.Inputs["name5"],
							// from defaults
							model.Data{
								Number: returnedPkg.Inputs["name6"].Number.Default,
							}: returnedPkg.Inputs["name6"],
							model.Data{
								String: returnedPkg.Inputs["name7"].String.Default,
							}: returnedPkg.Inputs["name7"],
						}

						fakeValidate := new(validate.Fake)

						objectUnderTest := newOpCaller(
							fakePkg,
							new(pubsub.Fake),
							new(fakeDCGNodeRepo),
							new(fakeCaller),
							new(uniquestring.Fake),
							fakeValidate,
							"dummyRootFSPath",
						)

						/* act */
						objectUnderTest.Call(
							providedInboundScope,
							providedOpId,
							providedPkgBasePath,
							providedRootOpId,
							providedSCGOpCall,
						)

						/* assert */
						actualCalls := map[model.Data]*model.Param{}
						for i := 0; i < fakeValidate.ParamCallCount(); i++ {
							actualVarData, actualParam := fakeValidate.ParamArgsForCall(i)
							actualCalls[*actualVarData] = actualParam
						}
						Expect(actualCalls).To(Equal(expectedCalls))
					})
					Context("validate.Param errors", func() {
						It("should call pubSub.Publish w/ expected args", func() {
							/* arrange */
							providedInboundScope := map[string]*model.Data{}
							providedOpId := "dummyOpId"
							providedRootOpId := "dummyRootOpId"
							providedSCGOpCall := &model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}}

							fakeDCGNodeRepo := new(fakeDCGNodeRepo)
							fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

							resolvedPkgRef := "dummyResolvedPkgRef"

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
							fakePkg.ResolveReturns(resolvedPkgRef, true)
							fakePkg.GetReturns(opReturnedFromPkg, nil)

							fakeValidate := new(validate.Fake)

							errorReturnedFromValidate := "validationError0"
							fakeValidate.ParamReturns([]error{errors.New(errorReturnedFromValidate)})

							expectedMsg := fmt.Sprintf(`
-
  validation of the following input failed:

  Name: %v
  Value: %v
  Error(s):
    - %v

-`, "dummyVar1Name", "************", errorReturnedFromValidate)

							fakePubSub := new(pubsub.Fake)
							expectedEvent := &model.Event{
								Timestamp: time.Now().UTC(),
								OpErred: &model.OpErredEvent{
									Msg:      expectedMsg,
									OpId:     providedOpId,
									PkgRef:   resolvedPkgRef,
									RootOpId: providedRootOpId,
								},
							}

							objectUnderTest := newOpCaller(
								fakePkg,
								fakePubSub,
								fakeDCGNodeRepo,
								new(fakeCaller),
								new(uniquestring.Fake),
								fakeValidate,
								"dummyRootFSPath",
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOpId,
								"dummyPkgBasePath",
								providedRootOpId,
								providedSCGOpCall,
							)

							/* assert */
							actualEvent := fakePubSub.PublishArgsForCall(0)

							// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
							Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
							// set temporal fields to expected vals since they're already asserted
							actualEvent.Timestamp = expectedEvent.Timestamp

							Expect(actualEvent).To(Equal(expectedEvent))
						})
					})
					Context("validate.Param doesn't error", func() {
						It("should call pubSub.Publish w/ expected args", func() {
							/* arrange */
							providedInboundScope := map[string]*model.Data{}
							providedOpId := "dummyOpId"
							providedRootOpId := "dummyRootOpId"
							providedSCGOpCall := &model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}}

							resolvedPkgRef := "dummyPkgRef"

							expectedEvent := &model.Event{
								Timestamp: time.Now().UTC(),
								OpStarted: &model.OpStartedEvent{
									OpId:     providedOpId,
									PkgRef:   resolvedPkgRef,
									RootOpId: providedRootOpId,
								},
							}

							fakePkg := new(pkg.Fake)
							fakePkg.ResolveReturns(resolvedPkgRef, true)
							fakePkg.GetReturns(&model.PkgManifest{}, nil)

							fakePubSub := new(pubsub.Fake)

							objectUnderTest := newOpCaller(
								fakePkg,
								fakePubSub,
								new(fakeDCGNodeRepo),
								new(fakeCaller),
								new(uniquestring.Fake),
								new(validate.Fake),
								"dummyRootFSPath",
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOpId,
								"dummyPkgBasePath",
								providedRootOpId,
								providedSCGOpCall,
							)

							/* assert */
							actualEvent := fakePubSub.PublishArgsForCall(0)

							// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
							Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
							// set temporal fields to expected vals since they're already asserted
							actualEvent.Timestamp = expectedEvent.Timestamp

							Expect(actualEvent).To(Equal(expectedEvent))
						})
						It("should call caller.Call w/ expected args", func() {
							/* arrange */
							providedInboundScope := map[string]*model.Data{}
							providedOpId := "dummyOpId"
							providedRootOpId := "dummyRootOpId"
							providedSCGOpCall := &model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}}

							resolvedPkgRef := "dummyResolvedPkgRef"

							opReturnedFromPkg := &model.PkgManifest{
								Run: &model.SCG{
									Parallel: []*model.SCG{
										{
											Container: &model.SCGContainerCall{},
										},
									},
								},
							}
							fakePkg := new(pkg.Fake)
							fakePkg.ResolveReturns(resolvedPkgRef, true)
							fakePkg.GetReturns(opReturnedFromPkg, nil)

							fakeUniqueStringFactory := new(uniquestring.Fake)
							expectedNodeId := "dummyNodeId"
							fakeUniqueStringFactory.ConstructReturns(expectedNodeId)

							fakeCaller := new(fakeCaller)

							objectUnderTest := newOpCaller(
								fakePkg,
								new(pubsub.Fake),
								new(fakeDCGNodeRepo),
								fakeCaller,
								fakeUniqueStringFactory,
								new(validate.Fake),
								"dummyRootFSPath",
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOpId,
								"dummyPkgBasePath",
								providedRootOpId,
								providedSCGOpCall,
							)

							/* assert */
							actualNodeId,
								actualInboundScope,
								actualSCG,
								actualPkgRef,
								actualRootOpId := fakeCaller.CallArgsForCall(0)

							Expect(actualNodeId).To(Equal(expectedNodeId))
							Expect(actualInboundScope).To(Equal(providedInboundScope))
							Expect(actualSCG).To(Equal(opReturnedFromPkg.Run))
							Expect(actualPkgRef).To(Equal(resolvedPkgRef))
							Expect(actualRootOpId).To(Equal(providedRootOpId))
						})
						It("should call dcgNodeRepo.GetIfExists w/ expected args", func() {
							/* arrange */
							providedInboundScope := map[string]*model.Data{}
							providedOpId := "dummyOpId"
							providedPkgBasePath := "dummyPkgBasePath"
							providedRootOpId := "dummyRootOpId"
							providedSCGOpCall := &model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}}

							fakePkg := new(pkg.Fake)
							fakePkg.ResolveReturns("", true)
							fakePkg.GetReturns(&model.PkgManifest{}, nil)

							fakeDCGNodeRepo := new(fakeDCGNodeRepo)

							objectUnderTest := newOpCaller(
								fakePkg,
								new(pubsub.Fake),
								fakeDCGNodeRepo,
								new(fakeCaller),
								new(uniquestring.Fake),
								new(validate.Fake),
								"dummyRootFSPath",
							)

							/* act */
							objectUnderTest.Call(
								providedInboundScope,
								providedOpId,
								providedPkgBasePath,
								providedRootOpId,
								providedSCGOpCall,
							)

							/* assert */
							Expect(fakeDCGNodeRepo.GetIfExistsArgsForCall(0)).To(Equal(providedRootOpId))
						})
						Context("dcgNodeRepo.GetIfExists returns nil", func() {
							It("should call pubSub.Publish w/ expected args", func() {
								/* arrange */
								providedInboundScope := map[string]*model.Data{}
								providedOpId := "dummyOpId"
								providedRootOpId := "dummyRootOpId"
								providedSCGOpCall := &model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}}

								resolvedPkgRef := "dummyResolvedPkgRef"

								expectedEvent := &model.Event{
									Timestamp: time.Now().UTC(),
									OpEnded: &model.OpEndedEvent{
										OpId:     providedOpId,
										Outcome:  model.OpOutcomeKilled,
										RootOpId: providedRootOpId,
										PkgRef:   resolvedPkgRef,
									},
								}

								fakePkg := new(pkg.Fake)
								fakePkg.ResolveReturns(resolvedPkgRef, true)
								fakePkg.GetReturns(&model.PkgManifest{}, nil)

								fakePubSub := new(pubsub.Fake)

								objectUnderTest := newOpCaller(
									fakePkg,
									fakePubSub,
									new(fakeDCGNodeRepo),
									new(fakeCaller),
									new(uniquestring.Fake),
									new(validate.Fake),
									"dummyRootFSPath",
								)

								/* act */
								objectUnderTest.Call(
									providedInboundScope,
									providedOpId,
									"dummyPkgBasePath",
									providedRootOpId,
									providedSCGOpCall,
								)

								/* assert */
								actualEvent := fakePubSub.PublishArgsForCall(1)

								// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
								Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
								// set temporal fields to expected vals since they're already asserted
								actualEvent.Timestamp = expectedEvent.Timestamp

								Expect(actualEvent).To(Equal(expectedEvent))
							})
						})
						Context("dcgNodeRepo.GetIfExists doesn't return nil", func() {
							It("should call dcgNodeRepo.DeleteIfExists w/ expected args", func() {
								/* arrange */
								providedInboundScope := map[string]*model.Data{}
								providedOpId := "dummyOpId"
								providedRootOpId := "dummyRootOpId"
								providedSCGOpCall := &model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}}

								fakePkg := new(pkg.Fake)
								fakePkg.ResolveReturns("", true)
								fakePkg.GetReturns(&model.PkgManifest{}, nil)

								fakeDCGNodeRepo := new(fakeDCGNodeRepo)
								fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

								objectUnderTest := newOpCaller(
									fakePkg,
									new(pubsub.Fake),
									fakeDCGNodeRepo,
									new(fakeCaller),
									new(uniquestring.Fake),
									new(validate.Fake),
									"dummyRootFSPath",
								)

								/* act */
								objectUnderTest.Call(
									providedInboundScope,
									providedOpId,
									"dummyPkgBasePath",
									providedRootOpId,
									providedSCGOpCall,
								)

								/* assert */
								Expect(fakeDCGNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedOpId))
							})
							Context("caller.Call errored", func() {
								It("should call pubSub.Publish w/ expected args", func() {
									/* arrange */
									providedInboundScope := map[string]*model.Data{}
									providedOpId := "dummyOpId"
									providedRootOpId := "dummyRootOpId"
									providedSCGOpCall := &model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}}

									resolvedPkgRef := "dummyResolvedPkgRef"
									callErr := errors.New("dummyError")

									expectedEvent := &model.Event{
										Timestamp: time.Now().UTC(),
										OpErred: &model.OpErredEvent{
											Msg:      callErr.Error(),
											OpId:     providedOpId,
											PkgRef:   resolvedPkgRef,
											RootOpId: providedRootOpId,
										},
									}

									fakePkg := new(pkg.Fake)
									fakePkg.ResolveReturns(resolvedPkgRef, true)
									fakePkg.GetReturns(&model.PkgManifest{}, nil)

									fakeDCGNodeRepo := new(fakeDCGNodeRepo)
									fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

									fakePubSub := new(pubsub.Fake)

									fakeCaller := new(fakeCaller)
									fakeCaller.CallReturns(
										callErr,
									)

									objectUnderTest := newOpCaller(
										fakePkg,
										fakePubSub,
										fakeDCGNodeRepo,
										fakeCaller,
										new(uniquestring.Fake),
										new(validate.Fake),
										"dummyRootFSPath",
									)

									/* act */
									objectUnderTest.Call(
										providedInboundScope,
										providedOpId,
										"dummyPkgBasePath",
										providedRootOpId,
										providedSCGOpCall,
									)

									/* assert */
									actualEvent := fakePubSub.PublishArgsForCall(1)

									// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
									Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
									// set temporal fields to expected vals since they're already asserted
									actualEvent.Timestamp = expectedEvent.Timestamp

									Expect(actualEvent).To(Equal(expectedEvent))
								})
								It("should call pubSub.Publish w/ expected args", func() {
									/* arrange */
									providedInboundScope := map[string]*model.Data{}
									providedOpId := "dummyOpId"
									providedRootOpId := "dummyRootOpId"
									providedSCGOpCall := &model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}}

									resolvedPkgRef := "dummyResolvedPkgRef"

									expectedEvent := &model.Event{
										Timestamp: time.Now().UTC(),
										OpEnded: &model.OpEndedEvent{
											OpId:     providedOpId,
											PkgRef:   resolvedPkgRef,
											Outcome:  model.OpOutcomeFailed,
											RootOpId: providedRootOpId,
										},
									}

									fakePkg := new(pkg.Fake)
									fakePkg.ResolveReturns(resolvedPkgRef, true)
									fakePkg.GetReturns(&model.PkgManifest{}, nil)

									fakeDCGNodeRepo := new(fakeDCGNodeRepo)
									fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

									fakePubSub := new(pubsub.Fake)

									fakeCaller := new(fakeCaller)
									fakeCaller.CallReturns(
										errors.New("dummyError"),
									)

									objectUnderTest := newOpCaller(
										fakePkg,
										fakePubSub,
										fakeDCGNodeRepo,
										fakeCaller,
										new(uniquestring.Fake),
										new(validate.Fake),
										"dummyRootFSPath",
									)

									/* act */
									objectUnderTest.Call(
										providedInboundScope,
										providedOpId,
										"dummyPkgBasePath",
										providedRootOpId,
										providedSCGOpCall,
									)

									/* assert */
									actualEvent := fakePubSub.PublishArgsForCall(2)

									// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
									Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
									// set temporal fields to expected vals since they're already asserted
									actualEvent.Timestamp = expectedEvent.Timestamp

									Expect(actualEvent).To(Equal(expectedEvent))
								})
							})
							Context("caller.Call didn't error", func() {
								It("should call pubSub.Publish w/ expected args", func() {
									/* arrange */
									providedInboundScope := map[string]*model.Data{}
									providedOpId := "dummyOpId"
									providedRootOpId := "dummyRootOpId"
									providedSCGOpCall := &model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}}

									resolvedPkgRef := "dummyResolvedPkgRef"

									expectedEvent := &model.Event{
										Timestamp: time.Now().UTC(),
										OpEnded: &model.OpEndedEvent{
											OpId:     providedOpId,
											PkgRef:   resolvedPkgRef,
											Outcome:  model.OpOutcomeSucceeded,
											RootOpId: providedRootOpId,
										},
									}

									fakePkg := new(pkg.Fake)
									fakePkg.ResolveReturns(resolvedPkgRef, true)
									fakePkg.GetReturns(&model.PkgManifest{}, nil)

									fakeDCGNodeRepo := new(fakeDCGNodeRepo)
									fakeDCGNodeRepo.GetIfExistsReturns(&dcgNodeDescriptor{})

									fakePubSub := new(pubsub.Fake)

									objectUnderTest := newOpCaller(
										fakePkg,
										fakePubSub,
										fakeDCGNodeRepo,
										new(fakeCaller),
										new(uniquestring.Fake),
										new(validate.Fake),
										"dummyRootFSPath",
									)

									/* act */
									objectUnderTest.Call(
										providedInboundScope,
										providedOpId,
										"dummyPkgBasePath",
										providedRootOpId,
										providedSCGOpCall,
									)

									/* assert */
									actualEvent := fakePubSub.PublishArgsForCall(1)

									// @TODO: implement/use VTime (similar to IOS & VFS) so we don't need custom assertions on temporal fields
									Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
									// set temporal fields to expected vals since they're already asserted
									actualEvent.Timestamp = expectedEvent.Timestamp

									Expect(actualEvent).To(Equal(expectedEvent))
								})
							})
						})
					})
				})
			})
		})
	})
})
