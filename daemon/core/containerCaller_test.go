package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/pkg/containerengine/engines/fake"
	"github.com/opspec-io/opctl/util/eventbus"
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
				new(fake.ContainerEngine),
				new(eventbus.Fake),
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
				new(fake.ContainerEngine),
				new(eventbus.Fake),
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
					new(fake.ContainerEngine),
					new(eventbus.Fake),
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
					new(fake.ContainerEngine),
					new(eventbus.Fake),
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
			It("should call eventBus.Publish w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedContainerId := "dummyContainerId"
				providedScgContainerCall := &model.ScgContainerCall{}
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"

				expectedEvent := model.Event{
					Timestamp: time.Now().UTC(),
					ContainerStarted: &model.ContainerStartedEvent{
						ContainerId: providedContainerId,
						OpRef:       providedOpRef,
						OpGraphId:   providedOpGraphId,
					},
				}

				fakeEventBus := new(eventbus.Fake)

				objectUnderTest := newContainerCaller(
					new(bundle.Fake),
					new(fake.ContainerEngine),
					fakeEventBus,
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
				actualEvent := fakeEventBus.PublishArgsForCall(0)

				// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
				Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
				// set temporal fields to expected vals since they're already asserted
				actualEvent.Timestamp = expectedEvent.Timestamp

				Expect(actualEvent).To(Equal(expectedEvent))
			})
			It("should call containerEngine.StartContainer w/ expected args", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedContainerId := "dummyContainerId"
				providedScgContainerCall := &model.ScgContainerCall{}
				providedOpRef := "dummyOpRef"
				providedOpGraphId := "dummyOpGraphId"

				opViewReturnedFromBundle := model.OpView{
					Inputs: []*model.Param{
						{
							String: &model.StringParam{},
						},
					},
				}
				fakeBundle := new(bundle.Fake)
				fakeBundle.GetOpReturns(opViewReturnedFromBundle, nil)

				expectedReq := newContainerStartReq(
					providedInboundScope,
					providedScgContainerCall,
					providedContainerId,
					opViewReturnedFromBundle.Inputs,
					providedOpGraphId)

				fakeContainerEngine := new(fake.ContainerEngine)

				fakeEventBus := new(eventbus.Fake)

				objectUnderTest := newContainerCaller(
					fakeBundle,
					fakeContainerEngine,
					fakeEventBus,
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
				actualReq, actualEventPublisher := fakeContainerEngine.StartContainerArgsForCall(0)
				Expect(actualReq).To(Equal(expectedReq))
				Expect(actualEventPublisher).To(Equal(fakeEventBus))
			})
			Context("containerEngine.StartContainer errors", func() {
				It("should return expected error", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{}
					providedContainerId := "dummyContainerId"
					providedScgContainerCall := &model.ScgContainerCall{}
					providedOpRef := "dummyOpRef"
					providedOpGraphId := "dummyOpGraphId"

					expectedError := errors.New("dummyError")

					fakeContainerEngine := new(fake.ContainerEngine)
					fakeContainerEngine.StartContainerReturns(expectedError)

					objectUnderTest := newContainerCaller(
						new(bundle.Fake),
						fakeContainerEngine,
						new(eventbus.Fake),
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
			Context("containerEngine.StartContainer doesn't error", func() {
				It("should call containerEngine.InspectContainerIfExists w/ expected args", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{}
					providedContainerId := "dummyContainerId"
					providedScgContainerCall := &model.ScgContainerCall{}
					providedOpRef := "dummyOpRef"
					providedOpGraphId := "dummyOpGraphId"

					fakeContainerEngine := new(fake.ContainerEngine)

					objectUnderTest := newContainerCaller(
						new(bundle.Fake),
						fakeContainerEngine,
						new(eventbus.Fake),
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
					Expect(fakeContainerEngine.InspectContainerIfExistsArgsForCall(0)).Should(Equal(providedContainerId))
				})
				Context("containerEngine.InspectContainerIfExists errors", func() {
					It("should return expected error", func() {
						/* arrange */
						providedInboundScope := map[string]*model.Data{}
						providedContainerId := "dummyContainerId"
						providedScgContainerCall := &model.ScgContainerCall{}
						providedOpRef := "dummyOpRef"
						providedOpGraphId := "dummyOpGraphId"

						expectedError := errors.New("dummyError")

						fakeContainerEngine := new(fake.ContainerEngine)
						fakeContainerEngine.InspectContainerIfExistsReturns(&model.DcgContainerCall{}, expectedError)

						objectUnderTest := newContainerCaller(
							new(bundle.Fake),
							fakeContainerEngine,
							new(eventbus.Fake),
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
				Context("containerEngine.InspectContainerIfExists doesn't error", func() {
					It("should return expected outboundScope", func() {
						/* arrange */
						providedScgContainerCall := &model.ScgContainerCall{
							Dirs: map[string]*model.ScgContainerDir{
								"dir1ContainerPath": {Bind: "dir1VarName"},
								"dir2ContainerPath": {Bind: "dir2VarName"},
							},
							EnvVars: map[string]*model.ScgContainerEnvVar{
								"envVar1Name": {Bind: "string1VarName"},
								"envVar2Name": {Bind: "string2VarName"},
							},
							Files: map[string]*model.ScgContainerFile{
								"file1ContainerPath": {Bind: "file1VarName"},
								"file2ContainerPath": {Bind: "file2VarName"},
							},
							Sockets: map[string]*model.ScgContainerSocket{
								"socket1ContainerAddress": {Bind: "socket1VarName"},
								"socket2ContainerAddress": {Bind: "socket2VarName"},
							},
						}

						dcgContainerCallReturnedFromEngine := &model.DcgContainerCall{
							Dirs: map[string]string{
								"dir1ContainerPath": "dir1HostPath",
								"dir2ContainerPath": "dir2HostPath",
							},
							EnvVars: map[string]string{
								"envVar1Name": "envVar1Value",
								"envVar2Name": "envVar2Value",
							},
							Files: map[string]string{
								"file1ContainerPath": "file1HostPath",
								"file2ContainerPath": "file2HostPath",
							},
							Sockets: map[string]string{
								"socket1ContainerAddress": "socket1HostAddress",
								"socket2ContainerAddress": "socket2HostAddress",
							},
						}
						fakeContainerEngine := new(fake.ContainerEngine)
						fakeContainerEngine.InspectContainerIfExistsReturns(dcgContainerCallReturnedFromEngine, nil)

						expectedOutboundScope := map[string]*model.Data{
							// dirs
							"dir1VarName": {Dir: "dir1HostPath"},
							"dir2VarName": {Dir: "dir2HostPath"},

							// envVars
							"string1VarName": {String: "envVar1Value"},
							"string2VarName": {String: "envVar2Value"},

							// files
							"file1VarName": {File: "file1HostPath"},
							"file2VarName": {File: "file2HostPath"},

							// sockets
							"socket1VarName": {Socket: "socket1HostAddress"},
							"socket2VarName": {Socket: "socket2HostAddress"},
						}

						objectUnderTest := newContainerCaller(
							new(bundle.Fake),
							fakeContainerEngine,
							new(eventbus.Fake),
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
				new(fake.ContainerEngine),
				new(eventbus.Fake),
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
		It("should call eventBus.Publish w/ expected args", func() {
			/* arrange */
			providedInboundScope := map[string]*model.Data{}
			providedContainerId := "dummyContainerId"
			providedScgContainerCall := &model.ScgContainerCall{}
			providedOpRef := "dummyOpRef"
			providedOpGraphId := "dummyOpGraphId"

			expectedEvent := model.Event{
				Timestamp: time.Now().UTC(),
				ContainerExited: &model.ContainerExitedEvent{
					ContainerId: providedContainerId,
					OpRef:       providedOpRef,
					OpGraphId:   providedOpGraphId,
				},
			}

			fakeEventBus := new(eventbus.Fake)

			objectUnderTest := newContainerCaller(
				new(bundle.Fake),
				new(fake.ContainerEngine),
				fakeEventBus,
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
			actualEvent := fakeEventBus.PublishArgsForCall(1)

			// @TODO: implement/use VTime (similar to VOS & VFS) so we don't need custom assertions on temporal fields
			Expect(actualEvent.Timestamp).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
			// set temporal fields to expected vals since they're already asserted
			actualEvent.Timestamp = expectedEvent.Timestamp

			Expect(actualEvent).To(Equal(expectedEvent))
		})
	})
})
