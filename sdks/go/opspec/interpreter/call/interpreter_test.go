package call

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	containerFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/fakes"
	opFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/fakes"
	predicatesFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates/fakes"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter(
				new(containerFakes.FakeInterpreter),
				"dummyDataDirPath",
			)).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		Context("callSpec.If not nil", func() {
			It("should call predicatesInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{}
				providedCallSpec := &model.CallSpec{
					If: new([]*model.PredicateSpec),
				}

				fakePredicatesInterpreter := new(predicatesFakes.FakeInterpreter)
				fakePredicatesInterpreter.InterpretReturns(
					true,
					nil,
				)

				objectUnderTest := _interpreter{
					predicatesInterpreter: fakePredicatesInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedScope,
					providedCallSpec,
					"providedID",
					"dummyOpPath",
					nil,
					"providedRootCallID",
				)

				/* assert */
				actualCallSpecIf,
					actualScope := fakePredicatesInterpreter.InterpretArgsForCall(0)

				Expect(actualCallSpecIf).To(Equal(*providedCallSpec.If))
				Expect(actualScope).To(Equal(providedScope))
			})
			Context("predicatesInterpreter returns err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakePredicatesInterpreter := new(predicatesFakes.FakeInterpreter)
					expectedError := errors.New("expectedError")
					fakePredicatesInterpreter.InterpretReturns(
						true,
						expectedError,
					)

					objectUnderTest := _interpreter{
						predicatesInterpreter: fakePredicatesInterpreter,
					}

					/* act */
					_, actualError := objectUnderTest.Interpret(
						map[string]*model.Value{},
						&model.CallSpec{
							If: new([]*model.PredicateSpec),
						},
						"providedID",
						"dummyOpPath",
						nil,
						"providedRootCallID",
					)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
		})
		Context("callSpec.Container not nil", func() {
			It("should call containerCallInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{}

				providedCallSpec := &model.CallSpec{
					Container: &model.ContainerCallSpec{},
				}

				providedID := "providedID"
				providedOpPath := "providedOpPath"
				providedRootCallID := "providedRootCallID"

				fakeContainerCallInterpreter := new(containerFakes.FakeInterpreter)

				objectUnderTest := _interpreter{
					containerCallInterpreter: fakeContainerCallInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedScope,
					providedCallSpec,
					providedID,
					providedOpPath,
					nil,
					providedRootCallID,
				)

				/* assert */
				actualScope,
					actualContainerCallSpec,
					actualContainerID,
					actualOpPath := fakeContainerCallInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualContainerCallSpec).To(Equal(providedCallSpec.Container))
				Expect(actualContainerID).To(Equal(providedID))
				Expect(actualOpPath).To(Equal(providedOpPath))

			})
			It("should return expected result", func() {
				/* arrange */
				providedID := "providedID"

				providedCallSpec := &model.CallSpec{
					Container: &model.ContainerCallSpec{},
				}

				providedParentIDValue := "providedParentID"
				providedParentID := &providedParentIDValue
				providedRootCallID := "providedRootCallID"

				fakeContainerCallInterpreter := new(containerFakes.FakeInterpreter)
				expectedContainerCall := &model.ContainerCall{}
				fakeContainerCallInterpreter.InterpretReturns(expectedContainerCall, nil)

				expectedCall := &model.Call{
					Container: expectedContainerCall,
					ID:        providedID,
					ParentID:  providedParentID,
					RootID:    providedRootCallID,
				}

				objectUnderTest := _interpreter{
					containerCallInterpreter: fakeContainerCallInterpreter,
				}

				/* act */
				actualCall,
					actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedCallSpec,
					providedID,
					"dummyOpPath",
					providedParentID,
					"providedRootCallID",
				)

				/* assert */
				Expect(actualCall).To(Equal(expectedCall))
				Expect(actualError).To(BeNil())

			})
		})
		Context("callSpec.Op not nil", func() {
			It("should call opCallInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{}

				providedCallSpec := &model.CallSpec{
					Op: &model.OpCallSpec{},
				}

				providedID := "providedID"
				providedOpPath := "providedOpPath"
				providedRootCallID := "providedRootCallID"

				fakeOpCallInterpreter := new(opFakes.FakeInterpreter)

				objectUnderTest := _interpreter{
					opCallInterpreter: fakeOpCallInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedScope,
					providedCallSpec,
					providedID,
					providedOpPath,
					nil,
					providedRootCallID,
				)

				/* assert */
				actualScope,
					actualOpCallSpec,
					actualOpID,
					actualOpPath := fakeOpCallInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualOpCallSpec).To(Equal(providedCallSpec.Op))
				Expect(actualOpID).To(Equal(providedID))
				Expect(actualOpPath).To(Equal(providedOpPath))

			})
			It("should return expected result", func() {
				/* arrange */
				providedID := "providedID"

				providedCallSpec := &model.CallSpec{
					Op: &model.OpCallSpec{},
				}

				providedParentID := "providedParentID"
				providedRootCallID := "providedRootCallID"

				fakeOpCallInterpreter := new(opFakes.FakeInterpreter)
				expectedOpCall := &model.OpCall{}
				fakeOpCallInterpreter.InterpretReturns(expectedOpCall, nil)

				expectedCall := &model.Call{
					ID:       providedID,
					Op:       expectedOpCall,
					ParentID: &providedParentID,
					RootID:   providedRootCallID,
				}

				objectUnderTest := _interpreter{
					opCallInterpreter: fakeOpCallInterpreter,
				}

				/* act */
				actualCall,
					actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedCallSpec,
					providedID,
					"dummyOpPath",
					&providedParentID,
					"providedRootCallID",
				)

				/* assert */
				Expect(actualCall).To(Equal(expectedCall))
				Expect(actualError).To(BeNil())

			})
		})
		Context("callSpec.Parallel not empty", func() {
			It("should return expected result", func() {
				/* arrange */
				providedID := "providedID"

				providedCallSpec := &model.CallSpec{
					Parallel: &[]*model.CallSpec{
						&model.CallSpec{},
					},
				}

				providedParentID := "providedParentID"
				providedRootCallID := "providedRootCallID"

				expectedCall := &model.Call{
					ID:       providedID,
					Parallel: *providedCallSpec.Parallel,
					ParentID: &providedParentID,
					RootID:   providedRootCallID,
				}

				objectUnderTest := _interpreter{}

				/* act */
				actualCall,
					actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedCallSpec,
					providedID,
					"dummyOpPath",
					&providedParentID,
					"providedRootCallID",
				)

				/* assert */
				Expect(*actualCall).To(Equal(*expectedCall))
				Expect(actualError).To(BeNil())

			})
		})
		Context("callSpec.Serial not empty", func() {
			It("should return expected result", func() {
				/* arrange */
				providedID := "providedID"

				providedCallSpec := &model.CallSpec{
					Serial: &[]*model.CallSpec{
						&model.CallSpec{},
					},
				}

				providedParentID := "providedParentID"
				providedRootCallID := "providedRootCallID"

				expectedCall := &model.Call{
					ID:       providedID,
					ParentID: &providedParentID,
					RootID:   providedRootCallID,
					Serial:   *providedCallSpec.Serial,
				}

				objectUnderTest := _interpreter{}

				/* act */
				actualCall,
					actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedCallSpec,
					providedID,
					"dummyOpPath",
					&providedParentID,
					providedRootCallID,
				)

				/* assert */
				Expect(*actualCall).To(Equal(*expectedCall))
				Expect(actualError).To(BeNil())

			})
		})
	})
})
