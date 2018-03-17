package core

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
)

var _ = Context("parallelCaller", func() {
	Context("newParallelCaller", func() {
		It("should return parallelCaller", func() {
			/* arrange/act/assert */
			Expect(newParallelCaller(
				new(fakeCaller),
				new(fakeOpKiller),
				new(pubsub.Fake),
				new(uniquestring.Fake),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every parallelCall w/ expected args", func() {
			/* arrange */
			providedCallId := "dummyCallId"
			providedInboundScope := map[string]*model.Value{}
			providedRootOpId := "dummyRootOpId"
			providedOpDirHandle := new(data.FakeHandle)
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
			fakeUniqueStringFactory.ConstructReturns(returnedUniqueString, nil)

			objectUnderTest := newParallelCaller(
				fakeCaller,
				new(fakeOpKiller),
				new(pubsub.Fake),
				fakeUniqueStringFactory,
			)

			/* act */
			objectUnderTest.Call(
				providedCallId,
				providedInboundScope,
				providedRootOpId,
				providedOpDirHandle,
				providedSCGParallelCalls,
			)

			/* assert */
			actualSCGParallelCalls := []*model.SCG{}
			for callIndex := range providedSCGParallelCalls {
				actualNodeId,
					actualChildOutboundScope,
					actualSCG,
					actualOpDirHandle,
					actualRootOpId := fakeCaller.CallArgsForCall(callIndex)
				Expect(actualNodeId).To(Equal(returnedUniqueString))
				Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
				Expect(actualOpDirHandle).To(Equal(providedOpDirHandle))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
				actualSCGParallelCalls = append(actualSCGParallelCalls, actualSCG)
			}
			Expect(actualSCGParallelCalls).To(ConsistOf(providedSCGParallelCalls))
		})
		Context("caller errors", func() {
			It("should fail fast on childCall error & return expected error", func() {
				/* arrange */
				providedCallId := "dummyCallId"
				providedInboundScope := map[string]*model.Value{}
				providedRootOpId := "dummyRootOpId"
				providedOpDirHandle := new(data.FakeHandle)
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
  Error:
    - %v
-`,
					callErr,
				)

				fakeCaller := new(fakeCaller)
				fakeCaller.CallReturnsOnCall(0, callErr)

				returnedUniqueString := "dummyUniqueString"
				fakeUniqueStringFactory := new(uniquestring.Fake)
				fakeUniqueStringFactory.ConstructReturns(returnedUniqueString, nil)

				objectUnderTest := newParallelCaller(
					fakeCaller,
					new(fakeOpKiller),
					new(pubsub.Fake),
					fakeUniqueStringFactory,
				)

				/* act */
				actualError := objectUnderTest.Call(
					providedCallId,
					providedInboundScope,
					providedRootOpId,
					providedOpDirHandle,
					providedSCGParallelCalls,
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("caller doesn't error", func() {
			It("shouldn't exit until all childCalls complete & not error", func() {
				/* arrange */
				providedCallId := "dummyCallId"
				providedInboundScope := map[string]*model.Value{}
				providedRootOpId := "dummyRootOpId"
				providedOpDirHandle := new(data.FakeHandle)
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
				fakeUniqueStringFactory.ConstructReturns(returnedUniqueString, nil)

				objectUnderTest := newParallelCaller(
					fakeCaller,
					new(fakeOpKiller),
					new(pubsub.Fake),
					fakeUniqueStringFactory,
				)

				/* act */
				actualError := objectUnderTest.Call(
					providedCallId,
					providedInboundScope,
					providedRootOpId,
					providedOpDirHandle,
					providedSCGParallelCalls,
				)

				/* assert */
				actualSCGParallelCalls := []*model.SCG{}
				for callIndex := range providedSCGParallelCalls {
					actualNodeId,
						actualChildOutboundScope,
						actualSCG,
						actualOpDirHandle,
						actualRootOpId := fakeCaller.CallArgsForCall(callIndex)
					Expect(actualNodeId).To(Equal(returnedUniqueString))
					Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
					Expect(actualOpDirHandle).To(Equal(providedOpDirHandle))
					Expect(actualRootOpId).To(Equal(providedRootOpId))
					actualSCGParallelCalls = append(actualSCGParallelCalls, actualSCG)
				}
				Expect(actualSCGParallelCalls).To(ConsistOf(providedSCGParallelCalls))
				Expect(actualError).To(BeNil())
			})
		})
	})
})
