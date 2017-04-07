package core

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opctl/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("parallelCaller", func() {
	Context("newParallelCaller", func() {
		It("should return parallelCaller", func() {
			/* arrange/act/assert */
			Expect(newParallelCaller(
				new(fakeCaller),
				new(pubsub.Fake),
				new(uniquestring.Fake),
			)).Should(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every parallelCall w/ expected args", func() {
			/* arrange */
			providedCallId := "dummyCallId"
			providedInboundScope := map[string]*model.Data{}
			providedRootOpId := "dummyRootOpId"
			providedPkgRef := "dummyPkgRef"
			providedSCGParallelCalls := []*model.SCG{
				{
					Container: &model.SCGContainerCall{},
				},
				{
					Op: &model.SCGOpCall{},
				},
				{
					Parallel: []*model.SCG{},
				},
				{
					Serial: []*model.SCG{},
				},
			}

			fakeCaller := new(fakeCaller)

			returnedUniqueString := "dummyUniqueString"
			fakeUniqueStringFactory := new(uniquestring.Fake)
			fakeUniqueStringFactory.ConstructReturns(returnedUniqueString)

			objectUnderTest := newParallelCaller(fakeCaller, new(pubsub.Fake), fakeUniqueStringFactory)

			/* act */
			objectUnderTest.Call(
				providedCallId,
				providedInboundScope,
				providedRootOpId,
				providedPkgRef,
				providedSCGParallelCalls,
			)

			/* assert */
			actualSCGParallelCalls := []*model.SCG{}
			for callIndex := range providedSCGParallelCalls {
				actualNodeId,
					actualChildOutboundScope,
					actualSCG,
					actualPkgRef,
					actualRootOpId := fakeCaller.CallArgsForCall(callIndex)
				Expect(actualNodeId).To(Equal(returnedUniqueString))
				Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
				Expect(actualPkgRef).To(Equal(providedPkgRef))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
				actualSCGParallelCalls = append(actualSCGParallelCalls, actualSCG)
			}
			Expect(actualSCGParallelCalls).To(ConsistOf(providedSCGParallelCalls))
		})
		Context("caller errors", func() {
			It("shouldn't exit until all childCalls complete & return expected error", func() {
				/* arrange */
				providedCallId := "dummyCallId"
				providedInboundScope := map[string]*model.Data{}
				providedRootOpId := "dummyRootOpId"
				providedPkgRef := "dummyPkgRef"
				providedSCGParallelCalls := []*model.SCG{
					{
						Container: &model.SCGContainerCall{},
					},
					{
						Op: &model.SCGOpCall{},
					},
					{
						Parallel: []*model.SCG{},
					},
					{
						Serial: []*model.SCG{},
					},
				}

				callErr := errors.New("dummyError")

				expectedError := fmt.Errorf(`
-
  Error during parallel call.
  Error(s):
    - %v
    - %v
    - %v
    - %v
-`,
					callErr,
					callErr,
					callErr,
					callErr,
				)

				fakeCaller := new(fakeCaller)
				fakeCaller.CallReturns(callErr)

				returnedUniqueString := "dummyUniqueString"
				fakeUniqueStringFactory := new(uniquestring.Fake)
				fakeUniqueStringFactory.ConstructReturns(returnedUniqueString)

				objectUnderTest := newParallelCaller(fakeCaller, new(pubsub.Fake), fakeUniqueStringFactory)

				/* act */
				actualError := objectUnderTest.Call(
					providedCallId,
					providedInboundScope,
					providedRootOpId,
					providedPkgRef,
					providedSCGParallelCalls,
				)

				/* assert */
				actualSCGParallelCalls := []*model.SCG{}
				for callIndex := range providedSCGParallelCalls {
					actualNodeId,
						actualChildOutboundScope,
						actualSCG,
						actualPkgRef,
						actualRootOpId := fakeCaller.CallArgsForCall(callIndex)
					Expect(actualNodeId).To(Equal(returnedUniqueString))
					Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
					Expect(actualPkgRef).To(Equal(providedPkgRef))
					Expect(actualRootOpId).To(Equal(providedRootOpId))
					actualSCGParallelCalls = append(actualSCGParallelCalls, actualSCG)
				}
				Expect(actualSCGParallelCalls).To(ConsistOf(providedSCGParallelCalls))
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("caller doesn't error", func() {
			It("shouldn't exit until all childCalls complete & not error", func() {
				/* arrange */
				providedCallId := "dummyCallId"
				providedInboundScope := map[string]*model.Data{}
				providedRootOpId := "dummyRootOpId"
				providedPkgRef := "dummyPkgRef"
				providedSCGParallelCalls := []*model.SCG{
					{
						Container: &model.SCGContainerCall{},
					},
					{
						Op: &model.SCGOpCall{},
					},
					{
						Parallel: []*model.SCG{},
					},
					{
						Serial: []*model.SCG{},
					},
				}

				fakeCaller := new(fakeCaller)

				returnedUniqueString := "dummyUniqueString"
				fakeUniqueStringFactory := new(uniquestring.Fake)
				fakeUniqueStringFactory.ConstructReturns(returnedUniqueString)

				objectUnderTest := newParallelCaller(fakeCaller, new(pubsub.Fake), fakeUniqueStringFactory)

				/* act */
				actualError := objectUnderTest.Call(
					providedCallId,
					providedInboundScope,
					providedRootOpId,
					providedPkgRef,
					providedSCGParallelCalls,
				)

				/* assert */
				actualSCGParallelCalls := []*model.SCG{}
				for callIndex := range providedSCGParallelCalls {
					actualNodeId,
						actualChildOutboundScope,
						actualSCG,
						actualPkgRef,
						actualRootOpId := fakeCaller.CallArgsForCall(callIndex)
					Expect(actualNodeId).To(Equal(returnedUniqueString))
					Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
					Expect(actualPkgRef).To(Equal(providedPkgRef))
					Expect(actualRootOpId).To(Equal(providedRootOpId))
					actualSCGParallelCalls = append(actualSCGParallelCalls, actualSCG)
				}
				Expect(actualSCGParallelCalls).To(ConsistOf(providedSCGParallelCalls))
				Expect(actualError).To(BeNil())
			})
		})
	})
})
