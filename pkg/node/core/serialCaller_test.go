package core

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/pkg/errors"
)

var _ = Context("serialCaller", func() {
	Context("newSerialCaller", func() {
		It("should return serialCaller", func() {
			/* arrange/act/assert */
			Expect(newSerialCaller(
				new(fakeCaller),
				new(uniquestring.Fake),
			)).Should(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every serialCall w/ expected args", func() {
			/* arrange */
			inputMap := map[string]*model.Data{}
			providedInputs := make(chan *variable, 150)
			for varName, varValue := range inputMap {
				providedInputs <- &variable{Name: varName, Value: varValue}
			}
			// inputs chan must be closed for method under test to return
			close(providedInputs)

			providedOutputs := make(chan *variable, 150)
			providedRootOpId := "dummyRootOpId"
			providedPkgRef := "dummyPkgRef"
			providedScgSerialCalls := []*model.Scg{
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
			// outputs chan must be closed for method under test to return
			fakeCaller.CallStub = func(nodeId string, inputs chan *variable, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
				close(outputs)
				return
			}

			fakeUniqueStringFactory := new(uniquestring.Fake)
			uniqueStringCallIndex := 0
			fakeUniqueStringFactory.ConstructStub = func() (uniqueString string) {
				defer func() {
					uniqueStringCallIndex++
				}()
				return fmt.Sprintf("%v", uniqueStringCallIndex)
			}

			objectUnderTest := newSerialCaller(fakeCaller, fakeUniqueStringFactory)

			/* act */
			objectUnderTest.Call(
				providedInputs,
				providedOutputs,
				providedRootOpId,
				providedPkgRef,
				providedScgSerialCalls,
			)

			/* assert */
			for expectedScgIndex, expectedScg := range providedScgSerialCalls {
				actualNodeId,
					actualChildInputs,
					_,
					actualScg,
					actualPkgRef,
					actualRootOpId := fakeCaller.CallArgsForCall(expectedScgIndex)

				actualChildInputMap := map[string]*model.Data{}
				for input := range actualChildInputs {
					actualChildInputMap[input.Name] = input.Value
				}

				Expect(actualNodeId).To(Equal(fmt.Sprintf("%v", expectedScgIndex)))
				Expect(actualChildInputMap).To(Equal(inputMap))
				Expect(actualScg).To(Equal(expectedScg))
				Expect(actualPkgRef).To(Equal(providedPkgRef))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
			}
		})
		Context("caller errors", func() {
			It("should return the expected error", func() {
				/* arrange */
				providedInputs := make(chan *variable, 150)
				// inputs chan must be closed for method under test to return
				close(providedInputs)

				providedOutputs := make(chan *variable, 150)
				providedRootOpId := "dummyRootOpId"
				providedPkgRef := "dummyPkgRef"
				providedScgSerialCalls := []*model.Scg{
					{
						Container: &model.ScgContainerCall{},
					},
				}

				expectedError := errors.New("dummyError")
				fakeCaller := new(fakeCaller)
				fakeCaller.CallReturns(expectedError)

				objectUnderTest := newSerialCaller(fakeCaller, new(uniquestring.Fake))

				/* act */
				actualErr := objectUnderTest.Call(
					providedInputs,
					providedOutputs,
					providedRootOpId,
					providedPkgRef,
					providedScgSerialCalls,
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedError))
			})
		})
		Context("caller doesn't error", func() {
			Context("childOutputs empty", func() {
				It("should call secondChild w/ inboundScope", func() {
					/* arrange */
					inputMap := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyParentVar1Data"},
						"dummyVar2Name": {Dir: "dummyParentVar2Data"},
					}
					providedInputs := make(chan *variable, 150)
					for varName, varValue := range inputMap {
						providedInputs <- &variable{Name: varName, Value: varValue}
					}
					// inputs chan must be closed for method under test to return
					close(providedInputs)

					providedOutputs := make(chan *variable, 150)
					expectedSecondChildInputMap := inputMap
					providedRootOpId := "dummyRootOpId"
					providedPkgRef := "dummyPkgRef"
					providedScgSerialCalls := []*model.Scg{
						{
							Container: &model.ScgContainerCall{},
						},
						{
							Container: &model.ScgContainerCall{},
						},
					}

					fakeCaller := new(fakeCaller)
					// outputs chan must be closed for method under test to return
					fakeCaller.CallStub = func(nodeId string, inputs chan *variable, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
						close(outputs)
						return
					}

					objectUnderTest := newSerialCaller(fakeCaller, new(uniquestring.Fake))

					/* act */
					objectUnderTest.Call(
						providedInputs,
						providedOutputs,
						providedRootOpId,
						providedPkgRef,
						providedScgSerialCalls,
					)

					/* assert */
					_, actualSecondChildInputs, _, _, _, _ := fakeCaller.CallArgsForCall(1)
					actualSecondChildInputMap := map[string]*model.Data{}
					for input := range actualSecondChildInputs {
						actualSecondChildInputMap[input.Name] = input.Value
					}

					Expect(actualSecondChildInputMap).To(Equal(expectedSecondChildInputMap))
				})
			})
			Context("childOutputs not empty", func() {
				It("should call secondChild w/ firstChildOutputs overlaying inboundScope", func() {
					/* arrange */
					inputMap := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyParentVar1Data"},
						"dummyVar2Name": {Dir: "dummyParentVar2Data"},
						"dummyVar3Name": {File: "dummyParentVar3Data"},
					}
					providedInputs := make(chan *variable, 150)
					for varName, varValue := range inputMap {
						providedInputs <- &variable{Name: varName, Value: varValue}
					}
					// inputs chan must be closed for method under test to return
					close(providedInputs)

					firstChildOutputMap := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyFirstChildVar1Data"},
						"dummyVar2Name": {Dir: "dummyFirstChildVar2Data"},
					}
					expectedSecondChildInputMap := map[string]*model.Data{
						"dummyVar1Name": firstChildOutputMap["dummyVar1Name"],
						"dummyVar2Name": firstChildOutputMap["dummyVar2Name"],
						"dummyVar3Name": inputMap["dummyVar3Name"],
					}

					providedOutputs := make(chan *variable, 150)
					providedRootOpId := "dummyRootOpId"
					providedPkgRef := "dummyPkgRef"
					providedScgSerialCalls := []*model.Scg{
						{
							Container: &model.ScgContainerCall{},
						},
						{
							Container: &model.ScgContainerCall{},
						},
					}

					fakeCaller := new(fakeCaller)
					fakeCaller.CallStub = func(nodeId string, inputs chan *variable, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
						// stub firstChildOutputs
						if scg == providedScgSerialCalls[0] {
							for varName, varValue := range firstChildOutputMap {
								outputs <- &variable{
									Name:  varName,
									Value: varValue,
								}
							}
						}
						// outputs chan must be closed for method under test to return
						close(outputs)
						return
					}

					objectUnderTest := newSerialCaller(fakeCaller, new(uniquestring.Fake))

					/* act */
					objectUnderTest.Call(
						providedInputs,
						providedOutputs,
						providedRootOpId,
						providedPkgRef,
						providedScgSerialCalls,
					)

					/* assert */
					_, actualSecondChildInputs, _, _, _, _ := fakeCaller.CallArgsForCall(1)
					actualSecondChildInputMap := map[string]*model.Data{}
					for input := range actualSecondChildInputs {
						actualSecondChildInputMap[input.Name] = input.Value
					}

					Expect(actualSecondChildInputMap).To(Equal(expectedSecondChildInputMap))
				})
			})
		})
	})
})
