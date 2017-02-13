package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/containerprovider"
	"github.com/opspec-io/opctl/util/pathnormalizer"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Context("core", func() {
	Context("KillOp", func() {
		It("should call dcgNodeRepo.DeleteIfExists w/ expected args", func() {
			/* arrange */
			providedReq := model.KillOpReq{OpGraphId: "dummyOpGraphId"}

			fakeDcgNodeRepo := new(fakeDcgNodeRepo)

			objectUnderTest := _core{
				containerProvider:   new(containerprovider.Fake),
				pubSub:              new(pubsub.Fake),
				opCaller:            new(fakeOpCaller),
				pathNormalizer:      pathnormalizer.NewPathNormalizer(),
				dcgNodeRepo:         fakeDcgNodeRepo,
				uniqueStringFactory: new(uniquestring.Fake),
			}

			/* act */
			objectUnderTest.KillOp(providedReq)

			/* assert */
			Expect(fakeDcgNodeRepo.DeleteIfExistsArgsForCall(0)).To(Equal(providedReq.OpGraphId))
		})
		It("should call dcgNodeRepo.ListWithOpGraphId w/ expected args", func() {
			/* arrange */
			providedReq := model.KillOpReq{OpGraphId: "dummyOpGraphId"}

			fakeDcgNodeRepo := new(fakeDcgNodeRepo)

			objectUnderTest := _core{
				containerProvider:   new(containerprovider.Fake),
				pubSub:              new(pubsub.Fake),
				opCaller:            new(fakeOpCaller),
				pathNormalizer:      pathnormalizer.NewPathNormalizer(),
				dcgNodeRepo:         fakeDcgNodeRepo,
				uniqueStringFactory: new(uniquestring.Fake),
			}

			/* act */
			objectUnderTest.KillOp(providedReq)

			/* assert */
			Expect(fakeDcgNodeRepo.ListWithOpGraphIdArgsForCall(0)).To(Equal(providedReq.OpGraphId))
		})
		Context("dcgNodeRepo.ListWithOpGraphId returns nodes", func() {
			It("should call dcgNodeRepo.DeleteIfExists w/ expected args for each", func() {
				/* arrange */
				providedReq := model.KillOpReq{OpGraphId: "dummyOpGraphId"}

				nodesReturnedFromDcgNodeRepo := []*dcgNodeDescriptor{
					{Id: "dummyNode1Id"},
					{Id: "dummyNode2Id"},
				}
				fakeDcgNodeRepo := new(fakeDcgNodeRepo)
				fakeDcgNodeRepo.ListWithOpGraphIdReturns(nodesReturnedFromDcgNodeRepo)

				objectUnderTest := _core{
					containerProvider:   new(containerprovider.Fake),
					pubSub:              new(pubsub.Fake),
					opCaller:            new(fakeOpCaller),
					pathNormalizer:      pathnormalizer.NewPathNormalizer(),
					dcgNodeRepo:         fakeDcgNodeRepo,
					uniqueStringFactory: new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.KillOp(providedReq)

				/* assert */
				for nodeIndex, node := range nodesReturnedFromDcgNodeRepo {
					Expect(fakeDcgNodeRepo.DeleteIfExistsArgsForCall(nodeIndex + 1)).To(Equal(node.Id))
				}
			})
			// @TODO: flickers?
			It("should call containerProvider.DeleteContainerIfExists w/ expected args for container nodes", func() {
				/* arrange */
				providedReq := model.KillOpReq{OpGraphId: "dummyOpGraphId"}

				containerNodeIds := []string{"dummyNode1Id", "dummyNode3Id"}

				nodesReturnedFromDcgNodeRepo := []*dcgNodeDescriptor{
					{Id: "dummyNode1Id", Container: &dcgContainerDescriptor{}},
					{Id: "dummyNode2Id"},
					{Id: "dummyNode3Id", Container: &dcgContainerDescriptor{}},
				}
				fakeDcgNodeRepo := new(fakeDcgNodeRepo)
				fakeDcgNodeRepo.ListWithOpGraphIdReturns(nodesReturnedFromDcgNodeRepo)

				fakeContainerProvider := new(containerprovider.Fake)

				objectUnderTest := &_core{
					containerProvider:   fakeContainerProvider,
					pubSub:              new(pubsub.Fake),
					opCaller:            new(fakeOpCaller),
					pathNormalizer:      pathnormalizer.NewPathNormalizer(),
					dcgNodeRepo:         fakeDcgNodeRepo,
					uniqueStringFactory: new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.KillOp(providedReq)

				/* assert */
				for nodeIndex, nodeId := range containerNodeIds {
					Expect(fakeContainerProvider.DeleteContainerIfExistsArgsForCall(nodeIndex)).To(Equal(nodeId))
				}
			})
		})
	})
})
