package call

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdk/go/data"
	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/container"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/op"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/predicates"
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
		Context("scg.If not nil", func() {
			It("should call predicatesInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{}
				providedIf := new([]*model.SCGPredicate)
				providedSCG := &model.SCG{
					If: providedIf,
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
							If: new([]*model.SCGPredicate),
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
				}

				providedParentIDValue := "providedParentID"
				providedParentID := &providedParentIDValue

				fakeContainerCallInterpreter := new(container.FakeInterpreter)
				expectedDCGContainerCall := &model.DCGContainerCall{}

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
					Op: &model.SCGOpCall{},
				}

				providedParentID := "providedParentID"

				fakeOpCallInterpreter := new(op.FakeInterpreter)
				expectedDCGOpCall := &model.DCGOpCall{}

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
					Parallel: []*model.SCG{
						&model.SCG{},
					},
				}

				providedParentID := "providedParentID"

				expectedDCG := &model.DCG{
					Parallel: providedSCG.Parallel,
				}

				objectUnderTest := _interpreter{}

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
					Serial: []*model.SCG{
						&model.SCG{},
					},
				}

				providedParentID := "providedParentID"

				expectedDCG := &model.DCG{
					Serial: providedSCG.Serial,
				}

				objectUnderTest := _interpreter{}

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
