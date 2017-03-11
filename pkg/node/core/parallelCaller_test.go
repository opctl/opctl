package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Context("parallelCaller", func() {
	Context("newParallelCaller", func() {
		It("should return parallelCaller", func() {
			/* arrange/act/assert */
			Expect(newParallelCaller(
				new(fakeCaller),
				new(uniquestring.Fake),
			)).Should(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every parallelCall w/ expected args", func() {
			/* arrange */
			inputMap := map[string]*model.Data{}

			providedInputs := make(chan *variable)
			for varName, varValue := range inputMap {
				providedInputs <- &variable{Name: varName, Value: varValue}
			}
			// inputs chan must be closed for method under test to return
			close(providedInputs)

			providedRootOpId := "dummyRootOpId"
			providedPkgRef := "dummyPkgRef"
			providedScgParallelCalls := []*model.Scg{
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
			fakeUniqueStringFactory := new(uniquestring.Fake)
			fakeUniqueStringFactory.ConstructReturns(returnedUniqueString)

			objectUnderTest := newParallelCaller(fakeCaller, fakeUniqueStringFactory)

			/* act */
			objectUnderTest.Call(
				providedInputs,
				providedRootOpId,
				providedPkgRef,
				providedScgParallelCalls,
			)

			/* assert */
			actualScgParallelCalls := []*model.Scg{}
			for callIndex := range providedScgParallelCalls {
				actualNodeId,
					actualChildInputs,
					_,
					actualScg,
					actualPkgRef,
					actualRootOpId := fakeCaller.CallArgsForCall(callIndex)

				actualChildInputMap := map[string]*model.Data{}
				for input := range actualChildInputs {
					actualChildInputMap[input.Name] = input.Value
				}

				Expect(actualNodeId).To(Equal(returnedUniqueString))
				Expect(actualChildInputMap).To(Equal(inputMap))
				Expect(actualPkgRef).To(Equal(providedPkgRef))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
				actualScgParallelCalls = append(actualScgParallelCalls, actualScg)
			}
			Expect(actualScgParallelCalls).To(ConsistOf(providedScgParallelCalls))
		})
		Context("caller errors", func() {
			It("shouldn't exit until all childCalls complete & return expected error", func() {
				/* arrange */
				inputMap := map[string]*model.Data{}

				providedInputs := make(chan *variable)
				for varName, varValue := range inputMap {
					providedInputs <- &variable{Name: varName, Value: varValue}
				}
				// inputs chan must be closed for method under test to return
				close(providedInputs)

				providedRootOpId := "dummyRootOpId"
				providedPkgRef := "dummyPkgRef"
				providedScgParallelCalls := []*model.Scg{
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
				fakeCaller.CallReturns(errors.New("dummyError"))

				returnedUniqueString := "dummyUniqueString"
				fakeUniqueStringFactory := new(uniquestring.Fake)
				fakeUniqueStringFactory.ConstructReturns(returnedUniqueString)

				objectUnderTest := newParallelCaller(fakeCaller, fakeUniqueStringFactory)

				/* act */
				actualError := objectUnderTest.Call(
					providedInputs,
					providedRootOpId,
					providedPkgRef,
					providedScgParallelCalls,
				)

				/* assert */
				actualScgParallelCalls := []*model.Scg{}
				for callIndex := range providedScgParallelCalls {
					actualNodeId,
						actualChildInputs,
						_,
						actualScg,
						actualPkgRef,
						actualRootOpId := fakeCaller.CallArgsForCall(callIndex)

					actualChildInputMap := map[string]*model.Data{}
					for input := range actualChildInputs {
						actualChildInputMap[input.Name] = input.Value
					}

					Expect(actualNodeId).To(Equal(returnedUniqueString))
					Expect(actualChildInputMap).To(Equal(inputMap))
					Expect(actualPkgRef).To(Equal(providedPkgRef))
					Expect(actualRootOpId).To(Equal(providedRootOpId))
					actualScgParallelCalls = append(actualScgParallelCalls, actualScg)
				}
				Expect(actualScgParallelCalls).To(ConsistOf(providedScgParallelCalls))
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("caller doesn't error", func() {
			It("shouldn't exit until all childCalls complete & not error", func() {
				/* arrange */
				inputMap := map[string]*model.Data{}

				providedInputs := make(chan *variable)
				for varName, varValue := range inputMap {
					providedInputs <- &variable{Name: varName, Value: varValue}
				}
				// inputs chan must be closed for method under test to return
				close(providedInputs)

				providedRootOpId := "dummyRootOpId"
				providedPkgRef := "dummyPkgRef"
				providedScgParallelCalls := []*model.Scg{
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
				fakeUniqueStringFactory := new(uniquestring.Fake)
				fakeUniqueStringFactory.ConstructReturns(returnedUniqueString)

				objectUnderTest := newParallelCaller(fakeCaller, fakeUniqueStringFactory)

				/* act */
				actualError := objectUnderTest.Call(
					providedInputs,
					providedRootOpId,
					providedPkgRef,
					providedScgParallelCalls,
				)

				/* assert */
				actualScgParallelCalls := []*model.Scg{}
				for callIndex := range providedScgParallelCalls {
					actualNodeId,
						actualChildInputs,
						_,
						actualScg,
						actualPkgRef,
						actualRootOpId := fakeCaller.CallArgsForCall(callIndex)

					actualChildInputMap := map[string]*model.Data{}
					for input := range actualChildInputs {
						actualChildInputMap[input.Name] = input.Value
					}

					Expect(actualNodeId).To(Equal(returnedUniqueString))
					Expect(actualChildInputMap).To(Equal(inputMap))
					Expect(actualPkgRef).To(Equal(providedPkgRef))
					Expect(actualRootOpId).To(Equal(providedRootOpId))
					actualScgParallelCalls = append(actualScgParallelCalls, actualScg)
				}
				Expect(actualScgParallelCalls).To(ConsistOf(providedScgParallelCalls))
				Expect(actualError).To(BeNil())
			})
		})
	})
})
