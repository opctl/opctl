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
			providedOpGraphId := "dummyOpGraphId"
			providedOpRef := "dummyOpRef"
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
				providedOpGraphId,
				providedOpRef,
				providedScgSerialCalls,
			)

			/* assert */
			for expectedScgIndex, expectedScg := range providedScgSerialCalls {
				actualNodeId,
					actualChildOutboundScope,
					actualScg,
					actualOpRef,
					actualOpGraphId := fakeCaller.CallArgsForCall(expectedScgIndex)
				Expect(actualNodeId).To(Equal(fmt.Sprintf("%v", expectedScgIndex)))
				Expect(actualChildOutboundScope).To(Equal(providedInboundScope))
				Expect(actualScg).To(Equal(expectedScg))
				Expect(actualOpRef).To(Equal(providedOpRef))
				Expect(actualOpGraphId).To(Equal(providedOpGraphId))
			}
		})
		Context("caller errors", func() {
			It("should return the expected error", func() {
				/* arrange */
				providedInboundScope := map[string]*model.Data{}
				providedOpGraphId := "dummyOpGraphId"
				providedOpRef := "dummyOpRef"
				providedScgSerialCalls := []*model.Scg{
					{
						Container: &model.ScgContainerCall{},
					},
				}

				expectedError := errors.New("dummyError")
				fakeCaller := new(fakeCaller)
				fakeCaller.CallReturns(map[string]*model.Data{}, expectedError)

				objectUnderTest := newSerialCaller(fakeCaller, new(uniquestring.Fake))

				/* act */
				actualErr := objectUnderTest.Call(
					providedInboundScope,
					providedOpGraphId,
					providedOpRef,
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
					expectedScopePassedToGrandchild := providedInboundScope
					providedOpGraphId := "dummyOpGraphId"
					providedOpRef := "dummyOpRef"
					providedScgSerialCalls := []*model.Scg{
						{
							Container: &model.ScgContainerCall{},
						},
						{
							Container: &model.ScgContainerCall{},
						},
					}

					fakeCaller := new(fakeCaller)

					objectUnderTest := newSerialCaller(fakeCaller, new(uniquestring.Fake))

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOpGraphId,
						providedOpRef,
						providedScgSerialCalls,
					)

					/* assert */
					_, actualScopePassedToGranchild, _, _, _ := fakeCaller.CallArgsForCall(1)
					Expect(actualScopePassedToGranchild).To(Equal(expectedScopePassedToGrandchild))
				})
			})
			Context("childOutboundScope not empty", func() {
				It("should call grandchild w/ childOutboundScope overlaying inboundScope", func() {
					/* arrange */
					providedInboundScope := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyParentVar1Data"},
						"dummyVar2Name": {Dir: "dummyParentVar2Data"},
						"dummyVar3Name": {File: "dummyParentVar3Data"},
					}
					childOutboundScope := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyChildVar1Data"},
						"dummyVar2Name": {Dir: "dummyChildVar2Data"},
					}
					expectedScopePassedToGrandchild := map[string]*model.Data{
						"dummyVar1Name": childOutboundScope["dummyVar1Name"],
						"dummyVar2Name": childOutboundScope["dummyVar2Name"],
						"dummyVar3Name": providedInboundScope["dummyVar3Name"],
					}
					providedOpGraphId := "dummyOpGraphId"
					providedOpRef := "dummyOpRef"
					providedScgSerialCalls := []*model.Scg{
						{
							Container: &model.ScgContainerCall{},
						},
						{
							Container: &model.ScgContainerCall{},
						},
					}

					fakeCaller := new(fakeCaller)
					fakeCaller.CallReturns(childOutboundScope, nil)

					objectUnderTest := newSerialCaller(fakeCaller, new(uniquestring.Fake))

					/* act */
					objectUnderTest.Call(
						providedInboundScope,
						providedOpGraphId,
						providedOpRef,
						providedScgSerialCalls,
					)

					/* assert */
					_, actualScopePassedToGranchild, _, _, _ := fakeCaller.CallArgsForCall(1)
					Expect(actualScopePassedToGranchild).To(Equal(expectedScopePassedToGrandchild))
				})
			})
		})
	})
})
