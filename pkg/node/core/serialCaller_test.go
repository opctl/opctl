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
			providedInboundScope := map[string]*model.Data{}
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
			fakeCaller.CallStub = func(nodeId string, scope map[string]*model.Data, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
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
				providedInboundScope,
				providedOutputs,
				providedRootOpId,
				providedPkgRef,
				providedScgSerialCalls,
			)

			/* assert */
			for expectedScgIndex, expectedScg := range providedScgSerialCalls {
				actualNodeId,
					actualChildOutboundScope,
					_,
					actualScg,
					actualPkgRef,
					actualRootOpId := fakeCaller.CallArgsForCall(expectedScgIndex)
				Expect(actualNodeId).To(Equal(fmt.Sprintf("%v", expectedScgIndex)))
				Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
				Expect(actualScg).To(Equal(expectedScg))
				Expect(actualPkgRef).To(Equal(providedPkgRef))
				Expect(actualRootOpId).To(Equal(providedRootOpId))
			}
		})
		Context("caller errors", func() {
			It("should return the expected error", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
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
					providedInboundScope,
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
			Context("childOutboundScope empty", func() {
				It("should call grandchild w/ inboundScope", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyParentVar1Data"},
						"dummyVar2Name": {Dir: "dummyParentVar2Data"},
					}
					providedOutputs := make(chan *variable, 150)
					expectedScopePassedToGrandchild := providedInboundScope
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
					fakeCaller.CallStub = func(nodeId string, scope map[string]*model.Data, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
						close(outputs)
						return
					}

					objectUnderTest := newSerialCaller(fakeCaller, new(uniquestring.Fake))

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOutputs,
						providedRootOpId,
						providedPkgRef,
						providedScgSerialCalls,
					)

					/* assert */
					_, actualScopePassedToGranchild, _, _, _, _ := fakeCaller.CallArgsForCall(1)
					Expect(actualScopePassedToGranchild).To(Equal(expectedScopePassedToGrandchild))
				})
			})
			Context("childOutboundScope not empty", func() {
				It("should call secondChild w/ firstChildOutputs overlaying inboundScope", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyParentVar1Data"},
						"dummyVar2Name": {Dir: "dummyParentVar2Data"},
						"dummyVar3Name": {File: "dummyParentVar3Data"},
					}
					firstChildOutputs := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyFirstChildVar1Data"},
						"dummyVar2Name": {Dir: "dummyFirstChildVar2Data"},
					}
					expectedScopePassedToSecondChild := map[string]*model.Data{
						"dummyVar1Name": firstChildOutputs["dummyVar1Name"],
						"dummyVar2Name": firstChildOutputs["dummyVar2Name"],
						"dummyVar3Name": providedInboundScope["dummyVar3Name"],
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
					fakeCaller.CallStub = func(nodeId string, scope map[string]*model.Data, outputs chan *variable, scg *model.Scg, pkgRef string, rootOpId string) (err error) {
						// stub firstChildOutputs
						if scg == providedScgSerialCalls[0] {
							for varName, varValue := range firstChildOutputs {
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
						providedInboundScope,
						providedOutputs,
						providedRootOpId,
						providedPkgRef,
						providedScgSerialCalls,
					)

					/* assert */
					_, actualScopePassedToGranchild, _, _, _, _ := fakeCaller.CallArgsForCall(1)
					Expect(actualScopePassedToGranchild).To(Equal(expectedScopePassedToSecondChild))
				})
			})
		})
	})
})
