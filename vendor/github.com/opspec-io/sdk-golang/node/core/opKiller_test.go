package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/core/containerruntime"
)

var _ = Context("core", func() {
	Context("Kill", func() {
		It("should call dcgNodeRepo.DeleteIfExists w/ expected args", func() {
			/* arrange */
			providedReq := model.KillOpReq{OpId: "dummyOpId"}

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			objectUnderTest := _opKiller{
				containerRuntime: new(containerruntime.Fake),
				dcgNodeRepo:      fakeDCGNodeRepo,
			}

			/* act */
			objectUnderTest.Kill(providedReq)

			/* assert */
			Expect(fakeDCGNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedReq.OpId))
		})
		It("should call dcgNodeRepo.ListWithRootOpId w/ expected args", func() {
			/* arrange */
			providedReq := model.KillOpReq{OpId: "dummyOpId"}

			fakeDCGNodeRepo := new(fakeDCGNodeRepo)

			objectUnderTest := _opKiller{
				containerRuntime: new(containerruntime.Fake),
				dcgNodeRepo:      fakeDCGNodeRepo,
			}

			/* act */
			objectUnderTest.Kill(providedReq)

			/* assert */
			Expect(fakeDCGNodeRepo.ListWithRootOpIdArgsForCall(0)).To(Equal(providedReq.OpId))
		})
		Context("dcgNodeRepo.ListWithRootOpId returns nodes", func() {
			It("should call dcgNodeRepo.DeleteIfExists w/ expected args for each", func() {
				/* arrange */
				providedReq := model.KillOpReq{OpId: "dummyOpId"}

				nodesReturnedFromDCGNodeRepo := []*dcgNodeDescriptor{
					{Id: "dummyNode1Id"},
					{Id: "dummyNode2Id"},
					{Id: "dummyNode3Id"},
				}
				fakeDCGNodeRepo := new(fakeDCGNodeRepo)
				fakeDCGNodeRepo.ListWithRootOpIdReturns(nodesReturnedFromDCGNodeRepo)

				// use map so order ignored; calls happen in parallel so all ordering bets are off
				expectedCalls := map[string]bool{
					providedReq.OpId:                   true,
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
				providedReq := model.KillOpReq{OpId: "dummyOpId"}

				nodesReturnedFromDCGNodeRepo := []*dcgNodeDescriptor{
					{Id: "dummyNode1Id", Container: &dcgContainerDescriptor{}},
					{Id: "dummyNode2Id", Container: &dcgContainerDescriptor{}},
					{Id: "dummyNode3Id", Container: &dcgContainerDescriptor{}},
				}
				fakeDCGNodeRepo := new(fakeDCGNodeRepo)
				fakeDCGNodeRepo.ListWithRootOpIdReturns(nodesReturnedFromDCGNodeRepo)

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
