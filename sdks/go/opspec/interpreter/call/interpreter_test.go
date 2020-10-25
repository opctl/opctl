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
					"providedRootOpID",
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
						"providedRootOpID",
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
					Container: &model.CallContainerSpec{},
				}

				providedID := "providedID"
				providedOpPath := "providedOpPath"
				providedRootOpID := "providedRootOpID"

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
					providedRootOpID,
				)

				/* assert */
				actualScope,
					actualCallContainerSpec,
					actualContainerID,
					actualRootOpID,
					actualOpPath := fakeContainerCallInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualCallContainerSpec).To(Equal(providedCallSpec.Container))
				Expect(actualContainerID).To(Equal(providedID))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
				Expect(actualOpPath).To(Equal(providedOpPath))

			})
			It("should return expected result", func() {
				/* arrange */
				providedID := "providedID"

				providedCallSpec := &model.CallSpec{
					Container: &model.CallContainerSpec{},
				}

				providedParentIDValue := "providedParentID"
				providedParentID := &providedParentIDValue

				fakeContainerCallInterpreter := new(containerFakes.FakeInterpreter)
				expectedDCGContainerCall := &model.DCGContainerCall{}
				fakeContainerCallInterpreter.InterpretReturns(expectedDCGContainerCall, nil)

				expectedDCG := &model.DCG{
					Container: expectedDCGContainerCall,
					Id:        providedID,
					ParentID:  providedParentID,
				}

				objectUnderTest := _interpreter{
					containerCallInterpreter: fakeContainerCallInterpreter,
				}

				/* act */
				actualDCG,
					actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedCallSpec,
					providedID,
					"dummyOpPath",
					providedParentID,
					"providedRootOpID",
				)

				/* assert */
				Expect(actualDCG).To(Equal(expectedDCG))
				Expect(actualError).To(BeNil())

			})
		})
		Context("callSpec.Op not nil", func() {
			It("should call opCallInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{}

				providedCallSpec := &model.CallSpec{
					Op: &model.CallOpSpec{},
				}

				providedID := "providedID"
				providedOpPath := "providedOpPath"
				providedRootOpID := "providedRootOpID"

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
					providedRootOpID,
				)

				/* assert */
				actualScope,
					actualCallOpSpec,
					actualOpID,
					actualOpPath,
					actualRootOpID := fakeOpCallInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualCallOpSpec).To(Equal(providedCallSpec.Op))
				Expect(actualOpID).To(Equal(providedID))
				Expect(actualOpPath).To(Equal(providedOpPath))
				Expect(actualRootOpID).To(Equal(providedRootOpID))

			})
			It("should return expected result", func() {
				/* arrange */
				providedID := "providedID"

				providedCallSpec := &model.CallSpec{
					Op: &model.CallOpSpec{},
				}

				providedParentID := "providedParentID"

				fakeOpCallInterpreter := new(opFakes.FakeInterpreter)
				expectedDCGOpCall := &model.DCGOpCall{}
				fakeOpCallInterpreter.InterpretReturns(expectedDCGOpCall, nil)

				expectedDCG := &model.DCG{
					Id:       providedID,
					Op:       expectedDCGOpCall,
					ParentID: &providedParentID,
				}

				objectUnderTest := _interpreter{
					opCallInterpreter: fakeOpCallInterpreter,
				}

				/* act */
				actualDCG,
					actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedCallSpec,
					providedID,
					"dummyOpPath",
					&providedParentID,
					"providedRootOpID",
				)

				/* assert */
				Expect(actualDCG).To(Equal(expectedDCG))
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

				expectedDCG := &model.DCG{
					Id:       providedID,
					Parallel: *providedCallSpec.Parallel,
					ParentID: &providedParentID,
				}

				objectUnderTest := _interpreter{}

				/* act */
				actualDCG,
					actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedCallSpec,
					providedID,
					"dummyOpPath",
					&providedParentID,
					"providedRootOpID",
				)

				/* assert */
				Expect(*actualDCG).To(Equal(*expectedDCG))
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

				expectedDCG := &model.DCG{
					Id:       providedID,
					ParentID: &providedParentID,
					Serial:   *providedCallSpec.Serial,
				}

				objectUnderTest := _interpreter{}

				/* act */
				actualDCG,
					actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedCallSpec,
					providedID,
					"dummyOpPath",
					&providedParentID,
					"providedRootOpID",
				)

				/* assert */
				Expect(*actualDCG).To(Equal(*expectedDCG))
				Expect(actualError).To(BeNil())

			})
		})
	})
})
