package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Describe("parallelCaller", func() {
	Context("newParallelCaller", func() {
		It("should return parallelCaller", func() {
			/* arrange/act/assert */
			Expect(newParallelCaller(
				new(fakeCaller),
				new(uniquestring.FakeUniqueStringFactory),
			)).Should(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every parallelCall w/ expected args", func() {
			/* arrange */
			providedParentScope := map[string]*model.Data{}
			providedOpGraphId := "dummyOpGraphId"
			providedOpRef := "dummyOpRef"
			providedParallelCalls := []*model.Scg{
				{
					Container: &model.ScgContainerCall{},
				},
				{
					Op: &model.ScgOpCall{},
				},
				{
					Parallel: []*model.Scg{},
				},
				{
					Serial: []*model.Scg{},
				},
			}

			fakeCaller := new(fakeCaller)

			returnedUniqueString := "dummyUniqueString"
			fakeUniqueStringFactory := new(uniquestring.FakeUniqueStringFactory)
			fakeUniqueStringFactory.ConstructReturns(returnedUniqueString)

			objectUnderTest := newParallelCaller(fakeCaller, fakeUniqueStringFactory)

			/* act */
			objectUnderTest.Call(
				providedParentScope,
				providedOpGraphId,
				providedOpRef,
				providedParallelCalls,
			)

			/* assert */
			actualParallelCalls := []*model.Scg{}
			for callIndex := range providedParallelCalls {
				actualNodeId,
					actualChildScope,
					actualScg,
					actualOpRef,
					actualOpGraphId := fakeCaller.CallArgsForCall(callIndex)
				Expect(actualNodeId).To(Equal(returnedUniqueString))
				Expect(actualChildScope).To(Equal(providedParentScope))
				Expect(actualOpRef).To(Equal(providedOpRef))
				Expect(actualOpGraphId).To(Equal(providedOpGraphId))
				actualParallelCalls = append(actualParallelCalls, actualScg)
			}
			Expect(actualParallelCalls).To(ConsistOf(providedParallelCalls))
		})
		Context("caller errors", func() {
			It("shouldn't exit until all childCalls complete & return expected error", func() {
				/* arrange */
				providedParentScope := map[string]*model.Data{}
				providedOpGraphId := "dummyOpGraphId"
				providedOpRef := "dummyOpRef"
				providedParallelCalls := []*model.Scg{
					{
						Container: &model.ScgContainerCall{},
					},
					{
						Op: &model.ScgOpCall{},
					},
					{
						Parallel: []*model.Scg{},
					},
					{
						Serial: []*model.Scg{},
					},
				}

				expectedError := errors.New("One or more errors encountered in parallel run block")
				fakeCaller := new(fakeCaller)
				fakeCaller.CallReturns(map[string]*model.Data{}, errors.New("dummyError"))

				returnedUniqueString := "dummyUniqueString"
				fakeUniqueStringFactory := new(uniquestring.FakeUniqueStringFactory)
				fakeUniqueStringFactory.ConstructReturns(returnedUniqueString)

				objectUnderTest := newParallelCaller(fakeCaller, fakeUniqueStringFactory)

				/* act */
				actualError := objectUnderTest.Call(
					providedParentScope,
					providedOpGraphId,
					providedOpRef,
					providedParallelCalls,
				)

				/* assert */
				actualParallelCalls := []*model.Scg{}
				for callIndex := range providedParallelCalls {
					actualNodeId,
						actualChildScope,
						actualScg,
						actualOpRef,
						actualOpGraphId := fakeCaller.CallArgsForCall(callIndex)
					Expect(actualNodeId).To(Equal(returnedUniqueString))
					Expect(actualChildScope).To(Equal(providedParentScope))
					Expect(actualOpRef).To(Equal(providedOpRef))
					Expect(actualOpGraphId).To(Equal(providedOpGraphId))
					actualParallelCalls = append(actualParallelCalls, actualScg)
				}
				Expect(actualParallelCalls).To(ConsistOf(providedParallelCalls))
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("caller doesn't error", func() {
			It("shouldn't exit until all childCalls complete & not error", func() {
				/* arrange */
				providedParentScope := map[string]*model.Data{}
				providedOpGraphId := "dummyOpGraphId"
				providedOpRef := "dummyOpRef"
				providedParallelCalls := []*model.Scg{
					{
						Container: &model.ScgContainerCall{},
					},
					{
						Op: &model.ScgOpCall{},
					},
					{
						Parallel: []*model.Scg{},
					},
					{
						Serial: []*model.Scg{},
					},
				}

				fakeCaller := new(fakeCaller)

				returnedUniqueString := "dummyUniqueString"
				fakeUniqueStringFactory := new(uniquestring.FakeUniqueStringFactory)
				fakeUniqueStringFactory.ConstructReturns(returnedUniqueString)

				objectUnderTest := newParallelCaller(fakeCaller, fakeUniqueStringFactory)

				/* act */
				actualError := objectUnderTest.Call(
					providedParentScope,
					providedOpGraphId,
					providedOpRef,
					providedParallelCalls,
				)

				/* assert */
				actualParallelCalls := []*model.Scg{}
				for callIndex := range providedParallelCalls {
					actualNodeId,
						actualChildScope,
						actualScg,
						actualOpRef,
						actualOpGraphId := fakeCaller.CallArgsForCall(callIndex)
					Expect(actualNodeId).To(Equal(returnedUniqueString))
					Expect(actualChildScope).To(Equal(providedParentScope))
					Expect(actualOpRef).To(Equal(providedOpRef))
					Expect(actualOpGraphId).To(Equal(providedOpGraphId))
					actualParallelCalls = append(actualParallelCalls, actualScg)
				}
				Expect(actualParallelCalls).To(ConsistOf(providedParallelCalls))
				Expect(actualError).To(BeNil())
			})
		})
	})
})
