package core

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/pkg/errors"
)

var _ = Describe("serialCaller", func() {
	Context("newSerialCaller", func() {
		It("should return serialCaller", func() {
			/* arrange/act/assert */
			Expect(newSerialCaller(
				new(fakeCaller),
				new(uniquestring.FakeUniqueStringFactory),
			)).Should(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller for every serialCall w/ expected args", func() {
			/* arrange */
			providedParentScope := map[string]*model.Data{}
			providedOpGraphId := "dummyOpGraphId"
			providedOpRef := "dummyOpRef"
			providedSerialCall := []*model.Scg{
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

			fakeUniqueStringFactory := new(uniquestring.FakeUniqueStringFactory)
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
				providedParentScope,
				providedOpGraphId,
				providedOpRef,
				providedSerialCall,
			)

			/* assert */
			for expectedScgIndex, expectedScg := range providedSerialCall {
				actualNodeId,
					_, // actualArgs,
					actualScg,
					actualOpRef,
					actualOpGraphId := fakeCaller.CallArgsForCall(expectedScgIndex)
				Expect(actualNodeId).To(Equal(fmt.Sprintf("%v", expectedScgIndex)))
				Expect(actualScg).To(Equal(expectedScg))
				Expect(actualOpRef).To(Equal(providedOpRef))
				Expect(actualOpGraphId).To(Equal(providedOpGraphId))
			}
		})
		Describe("caller errors", func() {
			It("should return the expected error", func() {
				/* arrange */
				providedParentScope := map[string]*model.Data{}
				providedOpGraphId := "dummyOpGraphId"
				providedOpRef := "dummyOpRef"
				providedSerialCall := []*model.Scg{
					{
						Container: &model.ScgContainerCall{},
					},
				}

				expectedError := errors.New("dummyError")
				fakeCaller := new(fakeCaller)
				fakeCaller.CallReturns(map[string]*model.Data{}, expectedError)

				objectUnderTest := newSerialCaller(fakeCaller, new(uniquestring.FakeUniqueStringFactory))

				/* act */
				actualErr := objectUnderTest.Call(
					providedParentScope,
					providedOpGraphId,
					providedOpRef,
					providedSerialCall,
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedError))
			})
		})
		Describe("caller doesn't error", func() {
			Describe("childScope empty", func() {
				It("should call grandchild w/ parentScope", func() {
					/* arrange */
					providedParentScope := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyParentVar1Data"},
						"dummyVar2Name": {Dir: "dummyParentVar2Data"},
					}
					expectedScopePassedToGrandchild := providedParentScope
					providedOpGraphId := "dummyOpGraphId"
					providedOpRef := "dummyOpRef"
					providedSerialCall := []*model.Scg{
						{
							Container: &model.ScgContainerCall{},
						},
						{
							Container: &model.ScgContainerCall{},
						},
					}

					fakeCaller := new(fakeCaller)

					objectUnderTest := newSerialCaller(fakeCaller, new(uniquestring.FakeUniqueStringFactory))

					/* act */
					objectUnderTest.Call(
						providedParentScope,
						providedOpGraphId,
						providedOpRef,
						providedSerialCall,
					)

					/* assert */
					_, actualScopePassedToGranchild, _, _, _ := fakeCaller.CallArgsForCall(1)
					Expect(actualScopePassedToGranchild).To(Equal(expectedScopePassedToGrandchild))
				})
			})
			Describe("childScope not empty", func() {
				It("should call grandchild w/ childScope overlaying parentScope", func() {
					/* arrange */
					providedParentScope := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyParentVar1Data"},
						"dummyVar2Name": {Dir: "dummyParentVar2Data"},
						"dummyVar3Name": {File: "dummyParentVar3Data"},
					}
					childScope := map[string]*model.Data{
						"dummyVar1Name": {String: "dummyChildVar1Data"},
						"dummyVar2Name": {Dir: "dummyChildVar2Data"},
					}
					expectedScopePassedToGrandchild := map[string]*model.Data{
						"dummyVar1Name": childScope["dummyVar1Name"],
						"dummyVar2Name": childScope["dummyVar2Name"],
						"dummyVar3Name": providedParentScope["dummyVar3Name"],
					}
					providedOpGraphId := "dummyOpGraphId"
					providedOpRef := "dummyOpRef"
					providedSerialCall := []*model.Scg{
						{
							Container: &model.ScgContainerCall{},
						},
						{
							Container: &model.ScgContainerCall{},
						},
					}

					fakeCaller := new(fakeCaller)
					fakeCaller.CallReturns(childScope, nil)

					objectUnderTest := newSerialCaller(fakeCaller, new(uniquestring.FakeUniqueStringFactory))

					/* act */
					objectUnderTest.Call(
						providedParentScope,
						providedOpGraphId,
						providedOpRef,
						providedSerialCall,
					)

					/* assert */
					_, actualScopePassedToGranchild, _, _, _ := fakeCaller.CallArgsForCall(1)
					Expect(actualScopePassedToGranchild).To(Equal(expectedScopePassedToGrandchild))
				})
			})
		})
	})
})
