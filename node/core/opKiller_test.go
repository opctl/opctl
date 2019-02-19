package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/core/containerruntime"
)

var _ = Context("core", func() {
	Context("Kill", func() {
		It("should call dcgNodeRepo.DeleteIfExists w/ expected args", func() {
			/* arrange */
			providedReq := model.KillOpReq{OpID: "dummyOpID"}

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			objectUnderTest := _opKiller{
				containerRuntime: new(containerruntime.Fake),
				dcgNodeRepo:      fakeDCGNodeRepo,
			}

			/* act */
			objectUnderTest.Kill(providedReq)

			/* assert */
			Expect(fakeDCGNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedReq.OpID))
		})
		It("should call dcgNodeRepo.ListWithRootOpID w/ expected args", func() {
			/* arrange */
			providedReq := model.KillOpReq{OpID: "dummyOpID"}

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			objectUnderTest := _opKiller{
				containerRuntime: new(containerruntime.Fake),
				dcgNodeRepo:      fakeDCGNodeRepo,
			}

			/* act */
			objectUnderTest.Kill(providedReq)

			/* assert */
			Expect(fakeDCGNodeRepo.ListWithRootOpIDArgsForCall(0)).To(Equal(providedReq.OpID))
		})
		Context("dcgNodeRepo.ListWithRootOpID returns nodes", func() {
			It("should call dcgNodeRepo.DeleteIfExists w/ expected args for each", func() {
				/* arrange */
				providedReq := model.KillOpReq{OpID: "dummyOpID"}

				nodesReturnedFromDCGNodeRepo := []*dcgNodeDescriptor{
					{Id: "dummyNode1Id"},
					{Id: "dummyNode2Id"},
					{Id: "dummyNode3Id"},
				}
				fakeDCGNodeRepo := new(fakeDCGNodeRepo)
				fakeDCGNodeRepo.ListWithRootOpIDReturns(nodesReturnedFromDCGNodeRepo)

				// use map so order ignored; calls happen in parallel so all ordering bets are off
				expectedCalls := map[string]bool{
					providedReq.OpID:                   true,
					nodesReturnedFromDCGNodeRepo[0].Id: true,
					nodesReturnedFromDCGNodeRepo[1].Id: true,
					nodesReturnedFromDCGNodeRepo[2].Id: true,
				}

				objectUnderTest := _opKiller{
					containerRuntime: new(containerruntime.Fake),
					dcgNodeRepo:      fakeDCGNodeRepo,
				}

				/* act */
				objectUnderTest.Kill(providedReq)

				/* assert */
				actualCalls := map[string]bool{}
				callIndex := 0
				for callIndex < fakeDCGNodeRepo.DeleteIfExistsCallCount() {
					actualCalls[fakeDCGNodeRepo.DeleteIfExistsArgsForCall(callIndex)] = true
					callIndex++
				}

				Expect(actualCalls).To(Equal(expectedCalls))
			})
			It("should call containerRuntime.DeleteContainerIfExists w/ expected args for container nodes", func() {
				/* arrange */
				providedReq := model.KillOpReq{OpID: "dummyOpID"}

				nodesReturnedFromDCGNodeRepo := []*dcgNodeDescriptor{
					{Id: "dummyNode1Id", Container: &dcgContainerDescriptor{}},
					{Id: "dummyNode2Id", Container: &dcgContainerDescriptor{}},
					{Id: "dummyNode3Id", Container: &dcgContainerDescriptor{}},
				}
				fakeDCGNodeRepo := new(fakeDCGNodeRepo)
				fakeDCGNodeRepo.ListWithRootOpIDReturns(nodesReturnedFromDCGNodeRepo)

				// use map so order ignored; calls happen in parallel so all ordering bets are off
				expectedCalls := map[string]bool{
					nodesReturnedFromDCGNodeRepo[0].Id: true,
					nodesReturnedFromDCGNodeRepo[1].Id: true,
					nodesReturnedFromDCGNodeRepo[2].Id: true,
				}

				fakeContainerRuntime := new(containerruntime.Fake)

				objectUnderTest := &_opKiller{
					containerRuntime: fakeContainerRuntime,
					dcgNodeRepo:      fakeDCGNodeRepo,
				}

				/* act */
				objectUnderTest.Kill(providedReq)

				/* assert */
				actualCalls := map[string]bool{}
				callIndex := 0
				for callIndex < fakeContainerRuntime.DeleteContainerIfExistsCallCount() {
					actualCalls[fakeContainerRuntime.DeleteContainerIfExistsArgsForCall(callIndex)] = true
					callIndex++
				}

				Expect(actualCalls).To(Equal(expectedCalls))
			})
		})
	})
})
