package call

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/container"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/predicates"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter(
				new(container.FakeInterpreter),
				"dummyDataDirPath",
			)).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		Context("scg.If not empty", func() {
			It("should call predicatesInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{}
				providedSCG := &model.SCG{
					If: []*model.SCGPredicate{
						&model.SCGPredicate{},
					},
				}

				providedOpHandle := new(data.FakeHandle)

				fakePredicatesInterpreter := new(predicates.FakeInterpreter)
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
					providedSCG,
					"providedID",
					providedOpHandle,
					nil,
					"providedRootOpID",
				)

				/* assert */
				actualOpHandle,
					actualSCG,
					actualScope := fakePredicatesInterpreter.InterpretArgsForCall(0)

				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualSCG).To(Equal(providedSCG))
				Expect(actualScope).To(Equal(providedScope))
			})
			Context("predicatesInterpreter returns err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakePredicatesInterpreter := new(predicates.FakeInterpreter)
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
						&model.SCG{
							If: []*model.SCGPredicate{
								&model.SCGPredicate{},
							},
						},
						"providedID",
						new(data.FakeHandle),
						nil,
						"providedRootOpID",
					)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
		})
		Context("scg.Container not nil", func() {
			It("should call containerCallInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{}

				providedSCG := &model.SCG{
					Container: &model.SCGContainerCall{},
				}

				providedID := "providedID"
				providedOpHandle := new(data.FakeHandle)
				providedRootOpID := "providedRootOpID"

				fakeContainerCallInterpreter := new(container.FakeInterpreter)

				objectUnderTest := _interpreter{
					containerCallInterpreter: fakeContainerCallInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedScope,
					providedSCG,
					providedID,
					providedOpHandle,
					nil,
					providedRootOpID,
				)

				/* assert */
				actualScope,
					actualSCGContainerCall,
					actualContainerID,
					actualRootOpID,
					actualOpHandle := fakeContainerCallInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualSCGContainerCall).To(Equal(providedSCG.Container))
				Expect(actualContainerID).To(Equal(providedID))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
				Expect(actualOpHandle).To(Equal(providedOpHandle))

			})
			It("should return expected result", func() {
				/* arrange */
				providedID := "providedID"

				providedSCG := &model.SCG{
					Container: &model.SCGContainerCall{},
					If: []*model.SCGPredicate{
						&model.SCGPredicate{},
					},
				}

				providedParentIDValue := "providedParentID"
				providedParentID := &providedParentIDValue

				fakePredicatesInterpreter := new(predicates.FakeInterpreter)
				expectedIf := true
				fakePredicatesInterpreter.InterpretReturns(
					true,
					nil,
				)

				fakeContainerCallInterpreter := new(container.FakeInterpreter)
				expectedDCGContainerCall := &model.DCGContainerCall{}

				expectedDCG := &model.DCG{
					Container: expectedDCGContainerCall,
					Id:        providedID,
					If:        &expectedIf,
					ParentID:  providedParentID,
				}

				objectUnderTest := _interpreter{
					containerCallInterpreter: fakeContainerCallInterpreter,
					predicatesInterpreter:    fakePredicatesInterpreter,
				}

				/* act */
				actualDCG,
					actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCG,
					providedID,
					new(data.FakeHandle),
					providedParentID,
					"providedRootOpID",
				)

				/* assert */
				Expect(actualDCG).To(Equal(expectedDCG))
				Expect(actualError).To(BeNil())

			})
		})
		Context("scg.Op not nil", func() {
			It("should call opCallInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{}

				providedSCG := &model.SCG{
					Op: &model.SCGOpCall{},
				}

				providedID := "providedID"
				providedOpHandle := new(data.FakeHandle)
				providedRootOpID := "providedRootOpID"

				fakeOpCallInterpreter := new(op.FakeInterpreter)

				objectUnderTest := _interpreter{
					opCallInterpreter: fakeOpCallInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedScope,
					providedSCG,
					providedID,
					providedOpHandle,
					nil,
					providedRootOpID,
				)

				/* assert */
				actualScope,
					actualSCGOpCall,
					actualOpID,
					actualParentOpHandle,
					actualRootOpID := fakeOpCallInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualSCGOpCall).To(Equal(providedSCG.Op))
				Expect(actualOpID).To(Equal(providedID))
				Expect(actualParentOpHandle).To(Equal(providedOpHandle))
				Expect(actualRootOpID).To(Equal(providedRootOpID))

			})
			It("should return expected result", func() {
				/* arrange */
				providedID := "providedID"

				providedSCG := &model.SCG{
					If: []*model.SCGPredicate{
						&model.SCGPredicate{},
					},
					Op: &model.SCGOpCall{},
				}

				providedParentID := "providedParentID"

				fakePredicatesInterpreter := new(predicates.FakeInterpreter)
				expectedIf := true
				fakePredicatesInterpreter.InterpretReturns(
					true,
					nil,
				)

				fakeOpCallInterpreter := new(op.FakeInterpreter)
				expectedDCGOpCall := &model.DCGOpCall{}

				expectedDCG := &model.DCG{
					Id:       providedID,
					If:       &expectedIf,
					Op:       expectedDCGOpCall,
					ParentID: &providedParentID,
				}

				objectUnderTest := _interpreter{
					opCallInterpreter:     fakeOpCallInterpreter,
					predicatesInterpreter: fakePredicatesInterpreter,
				}

				/* act */
				actualDCG,
					actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCG,
					providedID,
					new(data.FakeHandle),
					&providedParentID,
					"providedRootOpID",
				)

				/* assert */
				Expect(actualDCG).To(Equal(expectedDCG))
				Expect(actualError).To(BeNil())

			})
		})
		Context("scg.Parallel not empty", func() {
			It("should return expected result", func() {
				/* arrange */
				providedID := "providedID"

				providedSCG := &model.SCG{
					If: []*model.SCGPredicate{
						&model.SCGPredicate{},
					},
					Parallel: []*model.SCG{
						&model.SCG{},
					},
				}

				providedParentID := "providedParentID"

				fakePredicatesInterpreter := new(predicates.FakeInterpreter)
				expectedIf := true
				fakePredicatesInterpreter.InterpretReturns(
					true,
					nil,
				)

				expectedDCG := &model.DCG{
					If:       &expectedIf,
					Parallel: providedSCG.Parallel,
				}

				objectUnderTest := _interpreter{
					predicatesInterpreter: fakePredicatesInterpreter,
				}

				/* act */
				actualDCG,
					actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCG,
					providedID,
					new(data.FakeHandle),
					&providedParentID,
					"providedRootOpID",
				)

				/* assert */
				Expect(actualDCG).To(Equal(expectedDCG))
				Expect(actualError).To(BeNil())

			})
		})
		Context("scg.Serial not empty", func() {
			It("should return expected result", func() {
				/* arrange */
				providedID := "providedID"

				providedSCG := &model.SCG{
					If: []*model.SCGPredicate{
						&model.SCGPredicate{},
					},
					Serial: []*model.SCG{
						&model.SCG{},
					},
				}

				providedParentID := "providedParentID"

				fakePredicatesInterpreter := new(predicates.FakeInterpreter)
				expectedIf := true
				fakePredicatesInterpreter.InterpretReturns(
					true,
					nil,
				)

				expectedDCG := &model.DCG{
					If:     &expectedIf,
					Serial: providedSCG.Serial,
				}

				objectUnderTest := _interpreter{
					predicatesInterpreter: fakePredicatesInterpreter,
				}

				/* act */
				actualDCG,
					actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCG,
					providedID,
					new(data.FakeHandle),
					&providedParentID,
					"providedRootOpID",
				)

				/* assert */
				Expect(actualDCG).To(Equal(expectedDCG))
				Expect(actualError).To(BeNil())

			})
		})
	})
})
