package core

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	outputsFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/outputs/fakes"
	. "github.com/opctl/opctl/sdks/go/opspec/opfile/fakes"
)

var _ = Context("opCaller", func() {
	Context("newOpCaller", func() {
		It("should return opCaller", func() {
			/* arrange/act/assert */
			Expect(newOpCaller(
				new(FakeStateStore),
				new(FakeCaller),
				"",
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller.Call w/ expected args", func() {
			/* arrange */
			providedOpPath := "providedOpPath"

			dummyString := "dummyString"
			providedCtx := context.Background()
			providedOpCall := &model.OpCall{
				BaseCall: model.BaseCall{
					OpPath: providedOpPath,
				},
				ChildCallID: "dummyChildCallID",
				ChildCallCallSpec: &model.CallSpec{
					Parallel: &[]*model.CallSpec{
						{
							Container: &model.ContainerCallSpec{},
						},
					},
				},
				Inputs: map[string]*model.Value{
					"dummyScopeName": {String: &dummyString},
				},
				OpID: "providedOpID",
			}
			providedRootCallID := "providedRootCallID"

			expectedChildCallScope := map[string]*model.Value{
				"dummyScopeName": providedOpCall.Inputs["dummyScopeName"],
				"/": &model.Value{
					Dir: &providedOpPath,
				},
			}

			fakeCaller := new(FakeCaller)

			fakeOpFileGetter := new(FakeGetter)
			// err to trigger immediate return
			fakeOpFileGetter.GetReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _opCaller{
				caller:       fakeCaller,
				stateStore:   new(FakeStateStore),
				opFileGetter: fakeOpFileGetter,
			}

			/* act */
			objectUnderTest.Call(
				providedCtx,
				providedOpCall,
				map[string]*model.Value{},
				nil,
				providedRootCallID,
				&model.OpCallSpec{},
			)

			/* assert */
			actualCtx,
				actualChildCallID,
				actualChildCallScope,
				actualChildCallSpec,
				actualOpPath,
				actualParentCallID,
				actualRootCallID := fakeCaller.CallArgsForCall(0)

			Expect(actualCtx).To(Not(BeNil()))
			Expect(actualChildCallID).To(Equal(providedOpCall.ChildCallID))
			Expect(actualChildCallScope).To(Equal(expectedChildCallScope))
			Expect(actualChildCallSpec).To(Equal(providedOpCall.ChildCallCallSpec))
			Expect(actualOpPath).To(Equal(providedOpPath))
			Expect(actualParentCallID).To(Equal(&providedOpCall.OpID))
			Expect(actualRootCallID).To(Equal(providedRootCallID))
		})
		It("should return expected results", func() {
			/* arrange */
			providedOpPath := "providedOpPath"

			providedOpCall := &model.OpCall{
				BaseCall: model.BaseCall{
					OpPath: providedOpPath,
				},
				OpID: "providedOpId",
			}

			expectedOutputName := "expectedOutputName"

			providedOpCallSpec := &model.OpCallSpec{
				Outputs: map[string]string{
					expectedOutputName: "",
				},
			}

			fakeOutputsInterpreter := new(outputsFakes.FakeInterpreter)
			interpretedOutputs := map[string]*model.Value{
				expectedOutputName: new(model.Value),
				// include unbound output to ensure it's not added to scope
				"unexpectedOutputName": new(model.Value),
			}
			fakeOutputsInterpreter.InterpretReturns(interpretedOutputs, nil)

			expectedOutputs := map[string]*model.Value{
				expectedOutputName: interpretedOutputs[expectedOutputName],
			}

			fakeOpFileGetter := new(FakeGetter)
			fakeOpFileGetter.GetReturns(&model.OpSpec{}, nil)

			objectUnderTest := _opCaller{
				caller:             new(FakeCaller),
				stateStore:         new(FakeStateStore),
				opFileGetter:       fakeOpFileGetter,
				outputsInterpreter: fakeOutputsInterpreter,
			}

			/* act */
			actualOutputs, actualErr := objectUnderTest.Call(
				context.Background(),
				providedOpCall,
				map[string]*model.Value{},
				nil,
				"rootCallID",
				providedOpCallSpec,
			)

			/* assert */
			Expect(actualOutputs).To(Equal(expectedOutputs))
			Expect(actualErr).To(BeNil())
		})
	})
})
