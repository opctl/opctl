package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/containerprovider"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/pkg/errors"
	"time"
)

var _ = Context("containerCaller", func() {
	Context("newContainerCaller", func() {
		It("should return containerCaller", func() {
			/* arrange/act/assert */
			Expect(newContainerCaller(
				new(bundle.Fake),
				new(containerprovider.Fake),
				new(pubsub.Fake),
				new(fakeDcgNodeRepo),
			)).Should(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call bundle.GetOp w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedContainerId := "dummyContainerId"
			providedScgContainerCall := &model.ScgContainerCall{}
			providedOpRef := "dummyOpRef"
			providedOpGraphId := "dummyOpGraphId"

			fakeBundle := new(bundle.Fake)

			objectUnderTest := newContainerCaller(
				fakeBundle,
				new(containerprovider.Fake),
				new(pubsub.Fake),
				new(fakeDcgNodeRepo),
			)

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedContainerId,
				providedScgContainerCall,
				providedOpRef,
				providedOpGraphId,
			)

			/* assert */
			Expect(fakeBundle.GetOpArgsForCall(0)).To(Equal(providedOpRef))
		})
		Context("bundle.GetOp errors", func() {
			It("should return expected error", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedContainerId := "dummyContainerId"
				providedScgContainerCall := &model.ScgContainerCall{}
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"

				expectedError := errors.New("dummyError")

				fakeBundle := new(bundle.Fake)
				fakeBundle.GetOpReturns(model.OpView{}, expectedError)

				objectUnderTest := newContainerCaller(
					fakeBundle,
					new(containerprovider.Fake),
					new(pubsub.Fake),
					new(fakeDcgNodeRepo),
				)

				/* act */
				_, actualError := objectUnderTest.Call(
					providedInboundScope,
					providedContainerId,
					providedScgContainerCall,
					providedOpRef,
					providedOpGraphId,
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("bundle.GetOp doesn't error", func() {
			It("should call dcgNodeRepo.Add w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedContainerId := "dummyContainerId"
				providedScgContainerCall := &model.ScgContainerCall{}
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"

				expectedDcgNodeDescriptor := &dcgNodeDescriptor{
					Id:        providedContainerId,
					OpRef:     providedOpRef,
					OpGraphId: providedOpGraphId,
					Container: &dcgContainerDescriptor{},
				}

				fakeDcgNodeRepo := new(fakeDcgNodeRepo)

				objectUnderTest := newContainerCaller(
					new(bundle.Fake),
					new(containerprovider.Fake),
					new(pubsub.Fake),
					fakeDcgNodeRepo,
				)

				/* act */
				objectUnderTest.Call(
					providedInboundScope,
					providedContainerId,
					providedScgContainerCall,
					providedOpRef,
					providedOpGraphId,
				)

				/* assert */
				Expect(fakeDcgNodeRepo.AddArgsForCall(0)).To(Equal(expectedDcgNodeDescriptor))
			})
			It("should call pubSub.Publish w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedContainerId := "dummyContainerId"
				providedScgContainerCall := &model.ScgContainerCall{}
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"

				expectedEvent := &model.Event{
					Timestamp: time.Now().UTC(),
					ContainerStarted: &model.ContainerStartedEvent{
						ContainerId: providedContainerId,
						OpRef:       providedOpRef,
						OpGraphId:   providedOpGraphId,
					},
				}

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := newContainerCaller(
					new(bundle.Fake),
					new(containerprovider.Fake),
					fakePubSub,
					new(fakeDcgNodeRepo),
				)

				/* act */
				objectUnderTest.Call(
					providedInboundScope,
					providedContainerId,
					providedScgContainerCall,
					providedOpRef,
					providedOpGraphId,
				)

				/* assert */
				actualEvent := fakePubSub.PublishArgsForCall(0)

				// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
				Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
				// set temporal fields to expected vals since they're already asserted
				actualEvent.Timestamp = expectedEvent.Timestamp

				Expect(actualEvent).To(Equal(expectedEvent))
			})
			It("should call containerProvider.RunContainer w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedContainerId := "dummyContainerId"
				providedScgContainerCall := &model.ScgContainerCall{}
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"

				opViewReturnedFromBundle := model.OpView{
					Inputs: map[string]*model.Param{
						"dummyParam1Name": {
							String: &model.StringParam{},
						},
					},
				}
				fakeBundle := new(bundle.Fake)
				fakeBundle.GetOpReturns(opViewReturnedFromBundle, nil)

				expectedReq := newRunContainerReq(
					providedInboundScope,
					providedScgContainerCall,
					providedContainerId,
					opViewReturnedFromBundle.Inputs,
					providedOpGraphId)

				fakeContainerProvider := new(containerprovider.Fake)

				fakePubSub := new(pubsub.Fake)

				objectUnderTest := newContainerCaller(
					fakeBundle,
					fakeContainerProvider,
					fakePubSub,
					new(fakeDcgNodeRepo),
				)

				/* act */
				objectUnderTest.Call(
					providedInboundScope,
					providedContainerId,
					providedScgContainerCall,
					providedOpRef,
					providedOpGraphId,
				)

				/* assert */
				actualReq, actualEventPublisher := fakeContainerProvider.RunContainerArgsForCall(0)
				Expect(actualReq).To(Equal(expectedReq))
				Expect(actualEventPublisher).To(Equal(fakePubSub))
			})
			Context("containerProvider.RunContainer errors", func() {
				It("should return expected error", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{}
					providedContainerId := "dummyContainerId"
					providedScgContainerCall := &model.ScgContainerCall{}
					providedOpRef := "dummyOpRef"
					providedOpGraphId := "dummyOpGraphId"

					expectedError := errors.New("dummyError")

					fakeContainerProvider := new(containerprovider.Fake)
					fakeContainerProvider.RunContainerReturns(expectedError)

					objectUnderTest := newContainerCaller(
						new(bundle.Fake),
						fakeContainerProvider,
						new(pubsub.Fake),
						new(fakeDcgNodeRepo),
					)

					/* act */
					_, actualError := objectUnderTest.Call(
						providedInboundScope,
						providedContainerId,
						providedScgContainerCall,
						providedOpRef,
						providedOpGraphId,
					)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("containerProvider.RunContainer doesn't error", func() {
				It("should call containerProvider.InspectContainerIfExists w/ expected args", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{}
					providedContainerId := "dummyContainerId"
					providedScgContainerCall := &model.ScgContainerCall{}
					providedOpRef := "dummyOpRef"
					providedOpGraphId := "dummyOpGraphId"

					fakeContainerProvider := new(containerprovider.Fake)

					objectUnderTest := newContainerCaller(
						new(bundle.Fake),
						fakeContainerProvider,
						new(pubsub.Fake),
						new(fakeDcgNodeRepo),
					)

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedContainerId,
						providedScgContainerCall,
						providedOpRef,
						providedOpGraphId,
					)

					/* assert */
					Expect(fakeContainerProvider.InspectContainerIfExistsArgsForCall(0)).Should(Equal(providedContainerId))
				})
				Context("containerProvider.InspectContainerIfExists errors", func() {
					It("should return expected error", func() {
						/* arrange */
						providedInboundScope := map[string]*model.Data{}
						providedContainerId := "dummyContainerId"
						providedScgContainerCall := &model.ScgContainerCall{}
						providedOpRef := "dummyOpRef"
						providedOpGraphId := "dummyOpGraphId"

						expectedError := errors.New("dummyError")

						fakeContainerProvider := new(containerprovider.Fake)
						fakeContainerProvider.InspectContainerIfExistsReturns(&model.DcgContainerCall{}, expectedError)

						objectUnderTest := newContainerCaller(
							new(bundle.Fake),
							fakeContainerProvider,
							new(pubsub.Fake),
							new(fakeDcgNodeRepo),
						)

						/* act */
						_, actualError := objectUnderTest.Call(
							providedInboundScope,
							providedContainerId,
							providedScgContainerCall,
							providedOpRef,
							providedOpGraphId,
						)

						/* assert */
						Expect(actualError).To(Equal(expectedError))
					})
				})
				Context("containerProvider.InspectContainerIfExists doesn't error", func() {
					It("should return expected outboundScope", func() {
						/* arrange */
						providedScgContainerCall := &model.ScgContainerCall{
							Dirs: map[string]*model.ScgBinding{
								"dir1ContainerPath": {Bind: "dir1VarName"},
								"dir2ContainerPath": {Bind: "dir2VarName"},
							},
							Files: map[string]*model.ScgBinding{
								"file1ContainerPath": {Bind: "file1VarName"},
								"file2ContainerPath": {Bind: "file2VarName"},
							},
							Sockets: map[string]*model.ScgBinding{
								"/unixSocket1ContainerAddress": {Bind: "socket1VarName"},
								"/unixSocket2ContainerAddress": {Bind: "socket2VarName"},
							},
						}

						dcgContainerCallReturnedFromEngine := &model.DcgContainerCall{
							Dirs: map[string]string{
								"dir1ContainerPath": "dir1HostPath",
								"dir2ContainerPath": "dir2HostPath",
							},
							Files: map[string]string{
								"file1ContainerPath": "file1HostPath",
								"file2ContainerPath": "file2HostPath",
							},
							Sockets: map[string]string{
								"/unixSocket1ContainerAddress": "/unixSocket1HostAddress",
								"/unixSocket2ContainerAddress": "/unixSocket2HostAddress",
							},
						}
						fakeContainerProvider := new(containerprovider.Fake)
						fakeContainerProvider.InspectContainerIfExistsReturns(dcgContainerCallReturnedFromEngine, nil)

						expectedOutboundScope := map[string]*model.Data{
							// dirs
							"dir1VarName": {Dir: "dir1HostPath"},
							"dir2VarName": {Dir: "dir2HostPath"},

							// files
							"file1VarName": {File: "file1HostPath"},
							"file2VarName": {File: "file2HostPath"},

							// sockets
							"socket1VarName": {Socket: "/unixSocket1HostAddress"},
							"socket2VarName": {Socket: "/unixSocket2HostAddress"},
						}

						objectUnderTest := newContainerCaller(
							new(bundle.Fake),
							fakeContainerProvider,
							new(pubsub.Fake),
							new(fakeDcgNodeRepo),
						)

						/* act */
						actualOutboundScope, _ := objectUnderTest.Call(
							map[string]*model.Data{},
							"dummyContainerId",
							providedScgContainerCall,
							"dummyOpRef",
							"dummyOpGraphId",
						)

						/* assert */
						Expect(actualOutboundScope).To(Equal(expectedOutboundScope))
					})
				})
			})
		})
		It("should call dcgNodeRepo.DeleteIfExists w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedContainerId := "dummyContainerId"
			providedScgContainerCall := &model.ScgContainerCall{}
			providedOpRef := "dummyOpRef"
			providedOpGraphId := "dummyOpGraphId"

			fakeDcgNodeRepo := new(fakeDcgNodeRepo)

			objectUnderTest := newContainerCaller(
				new(bundle.Fake),
				new(containerprovider.Fake),
				new(pubsub.Fake),
				fakeDcgNodeRepo,
			)

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedContainerId,
				providedScgContainerCall,
				providedOpRef,
				providedOpGraphId,
			)

			/* assert */
			Expect(fakeDcgNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedContainerId))
		})
		It("should call pubSub.Publish w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedContainerId := "dummyContainerId"
			providedScgContainerCall := &model.ScgContainerCall{}
			providedOpRef := "dummyOpRef"
			providedOpGraphId := "dummyOpGraphId"

			expectedEvent := &model.Event{
				Timestamp: time.Now().UTC(),
				ContainerExited: &model.ContainerExitedEvent{
					ContainerId: providedContainerId,
					OpRef:       providedOpRef,
					OpGraphId:   providedOpGraphId,
				},
			}

			fakePubSub := new(pubsub.Fake)

			objectUnderTest := newContainerCaller(
				new(bundle.Fake),
				new(containerprovider.Fake),
				fakePubSub,
				new(fakeDcgNodeRepo),
			)

			/* act */
			objectUnderTest.Call(
				providedInboundScope,
				providedContainerId,
				providedScgContainerCall,
				providedOpRef,
				providedOpGraphId,
			)

			/* assert */
			actualEvent := fakePubSub.PublishArgsForCall(1)

			// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
			Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
			// set temporal fields to expected vals since they're already asserted
			actualEvent.Timestamp = expectedEvent.Timestamp

			Expect(actualEvent).To(Equal(expectedEvent))
		})
	})
})
