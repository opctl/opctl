package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/opcall"
	"github.com/opspec-io/sdk-golang/opcall/outputs"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"time"
)

var _ = Context("opCaller", func() {
	Context("newOpCaller", func() {
		It("should return opCaller", func() {
			/* arrange/act/assert */
			Expect(newOpCaller(
				new(pubsub.Fake),
				new(fakeCaller),
				"",
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call opCall.Construct w/ expected args & return errors", func() {
			/* arrange */
			providedScope := map[string]*model.Value{}
			providedOpId := "dummyOpId"
			providedRootOpId := "dummyRootOpId"
			providedPkgHandle := new(pkg.FakeHandle)
			providedSCGOpCall := &model.SCGOpCall{
				Pkg: &model.SCGOpCallPkg{
					Ref: "dummyPkgRef",
				},
			}

			fakeOpCall := new(opcall.Fake)

			// error to trigger immediate return
			expectedErr := errors.New("dummyError")
			fakeOpCall.InterpretReturns(
				&model.DCGOpCall{
					DCGBaseCall: &model.DCGBaseCall{},
				},
				expectedErr,
			)

			objectUnderTest := _opCaller{
				opCall: fakeOpCall,
				pubSub: new(pubsub.Fake),
				caller: new(fakeCaller),
			}

			/* act */
			actualErr := objectUnderTest.Call(
				providedScope,
				providedOpId,
				providedPkgHandle,
				providedRootOpId,
				providedSCGOpCall,
			)

			/* assert */
			actualScope,
				actualSCGOpCall,
				actualOpId,
				actualPkgHandle,
				actualRootOpId := fakeOpCall.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualSCGOpCall).To(Equal(providedSCGOpCall))
			Expect(actualOpId).To(Equal(providedOpId))
			Expect(actualPkgHandle).To(Equal(providedPkgHandle))
			Expect(actualRootOpId).To(Equal(providedRootOpId))

			Expect(actualErr).To(Equal(expectedErr))
		})
		Context("opCall.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				providedOpId := "dummyOpId"
				providedRootOpId := "dummyRootOpId"
				providedSCGOpCall := &model.SCGOpCall{
					Pkg: &model.SCGOpCallPkg{
						Ref: "dummyPkgRef",
					},
				}

				expectedEvent := &model.Event{
					Timestamp: time.Now().UTC(),
					OpErred: &model.OpErredEvent{
						Msg:      "dummyError",
						OpId:     providedOpId,
						PkgRef:   providedSCGOpCall.Pkg.Ref,
						RootOpId: providedRootOpId,
					},
				}

				fakeOpCall := new(opcall.Fake)
				fakeOpCall.InterpretReturns(
					nil,
					errors.New(expectedEvent.OpErred.Msg),
				)

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := _opCaller{
					opCall: fakeOpCall,
					pubSub: fakePubSub,
					caller: new(fakeCaller),
				}

				/* act */
				objectUnderTest.Call(
					map[string]*model.Value{},
					providedOpId,
					new(pkg.FakeHandle),
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
		Context("opCall.Construct doesn't error", func() {
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedOpId := "dummyOpId"
				providedRootOpId := "dummyRootOpId"
				providedSCGOpCall := &model.SCGOpCall{
					Pkg: &model.SCGOpCallPkg{
						Ref: "dummyPkgRef",
					},
				}

				expectedEvent := &model.Event{
					Timestamp: time.Now().UTC(),
					OpStarted: &model.OpStartedEvent{
						OpId:     providedOpId,
						PkgRef:   providedSCGOpCall.Pkg.Ref,
						RootOpId: providedRootOpId,
					},
				}

				fakeOpCall := new(opcall.Fake)
				fakeOpCall.InterpretReturns(
					&model.DCGOpCall{
						DCGBaseCall: &model.DCGBaseCall{},
					},
					nil,
				)

				fakePubSub := new(pubsub.Fake)
				fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
					// close eventChannel to trigger immediate return
					close(eventChannel)
				}

				fakeCaller := new(fakeCaller)
				// err to trigger immediate return
				fakeCaller.CallReturns(errors.New("dummyError"))

				objectUnderTest := _opCaller{
					opCall: fakeOpCall,
					pubSub: fakePubSub,
					caller: fakeCaller,
				}

				/* act */
				objectUnderTest.Call(
					map[string]*model.Value{},
					providedOpId,
					new(pkg.FakeHandle),
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
				providedRootOpId := "dummyRootOpId"

				dummyString := "dummyString"
				dcgOpCall := model.DCGOpCall{
					DCGBaseCall: &model.DCGBaseCall{
						PkgHandle: new(pkg.FakeHandle),
					},
					ChildCallId: "dummyChildCallId",
					ChildCallSCG: &model.SCG{
						Parallel: []*model.SCG{
							{
								Container: &model.SCGContainerCall{},
							},
						},
					},
					Inputs: map[string]*model.Value{
						"dummyScopeName": {String: &dummyString},
					},
				}
				fakeOpCall := new(opcall.Fake)
				fakeOpCall.InterpretReturns(
					&dcgOpCall,
					nil,
				)

				fakePubSub := new(pubsub.Fake)
				fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
					// close eventChannel to trigger immediate return
					close(eventChannel)
				}

				fakeCaller := new(fakeCaller)
				// err to trigger immediate return
				fakeCaller.CallReturns(errors.New("dummyErr"))

				objectUnderTest := _opCaller{
					opCall: fakeOpCall,
					pubSub: fakePubSub,
					caller: fakeCaller,
				}

				/* act */
				objectUnderTest.Call(
					map[string]*model.Value{},
					"dummyOpId",
					new(pkg.FakeHandle),
					providedRootOpId,
					&model.SCGOpCall{Pkg: &model.SCGOpCallPkg{}},
				)

				/* assert */
				actualChildCallId,
					actualChildCallScope,
					actualChildSCG,
					actualPkgRef,
					actualRootOpId := fakeCaller.CallArgsForCall(0)

				Expect(actualChildCallId).To(Equal(dcgOpCall.ChildCallId))
				Expect(actualChildCallScope).To(Equal(dcgOpCall.Inputs))
				Expect(actualChildSCG).To(Equal(dcgOpCall.ChildCallSCG))
				Expect(actualPkgRef).To(Equal(dcgOpCall.PkgHandle))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
			})
			Context("caller.Call errs", func() {
				It("should call pubSub.Publish w/ expected args", func() {
					/* arrange */
					providedOpId := "dummyOpId"
					providedRootOpId := "dummyRootOpId"
					providedSCGOpCall := &model.SCGOpCall{
						Pkg: &model.SCGOpCallPkg{
							Ref: "dummyPkgRef",
						},
					}

					expectedEvent := &model.Event{
						Timestamp: time.Now().UTC(),
						OpEnded: &model.OpEndedEvent{
							OpId:     providedOpId,
							PkgRef:   providedSCGOpCall.Pkg.Ref,
							Outcome:  model.OpOutcomeFailed,
							RootOpId: providedRootOpId,
						},
					}

					fakeOpCall := new(opcall.Fake)
					fakeOpCall.InterpretReturns(
						&model.DCGOpCall{
							DCGBaseCall: &model.DCGBaseCall{},
						},
						nil,
					)

					fakePubSub := new(pubsub.Fake)
					fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
						// close eventChannel to trigger immediate return
						close(eventChannel)
					}

					fakeCaller := new(fakeCaller)
					fakeCaller.CallReturns(
						errors.New("dummyError"),
					)

					objectUnderTest := _opCaller{
						pubSub: fakePubSub,
						opCall: fakeOpCall,
						caller: fakeCaller,
					}

					/* act */
					objectUnderTest.Call(
						map[string]*model.Value{},
						providedOpId,
						new(pkg.FakeHandle),
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
					providedOpId := "dummyOpId"
					providedRootOpId := "dummyRootOpId"
					providedSCGOpCall := &model.SCGOpCall{
						Pkg: &model.SCGOpCallPkg{
							Ref: "dummyPkgRef",
						},
					}

					fakeOutputs := new(outputs.Fake)
					interpretedOutputs := map[string]*model.Value{"dummyOutputName": new(model.Value)}
					fakeOutputs.InterpretReturns(interpretedOutputs, nil)

					expectedEvent := &model.Event{
						Timestamp: time.Now().UTC(),
						OpEnded: &model.OpEndedEvent{
							OpId:     providedOpId,
							PkgRef:   providedSCGOpCall.Pkg.Ref,
							Outcome:  model.OpOutcomeSucceeded,
							RootOpId: providedRootOpId,
							Outputs:  interpretedOutputs,
						},
					}

					fakePkgHandle := new(pkg.FakeHandle)
					fakePkgHandle.PathReturns(new(string))

					fakeOpCall := new(opcall.Fake)
					fakeOpCall.InterpretReturns(
						&model.DCGOpCall{
							DCGBaseCall: &model.DCGBaseCall{
								PkgHandle: fakePkgHandle,
							},
						},
						nil,
					)

					fakePkg := new(pkg.Fake)
					fakePkg.GetManifestReturns(&model.PkgManifest{}, nil)

					fakePubSub := new(pubsub.Fake)
					fakePubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
						// close eventChannel to trigger immediate return
						close(eventChannel)
					}

					objectUnderTest := _opCaller{
						pkg:     fakePkg,
						pubSub:  fakePubSub,
						opCall:  fakeOpCall,
						outputs: fakeOutputs,
						caller:  new(fakeCaller),
					}

					/* act */
					objectUnderTest.Call(
						map[string]*model.Value{},
						providedOpId,
						new(pkg.FakeHandle),
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
