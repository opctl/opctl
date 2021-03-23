package core

import (
	"context"

	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
)

var _ = Context("opCaller", func() {
	Context("newOpCaller", func() {
		It("should return opCaller", func() {
			/* arrange/act/assert */
			Expect(newOpCaller(
				new(FakeCaller),
				"",
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call caller.Call w/ expected args", func() {
			/* arrange */
			providedOpPath := "providedOpPath"
			parentProvidedOpPath := filepath.Dir(providedOpPath)

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
				"./": {
					Dir: &providedOpPath,
				},
				"../": {
					Dir: &parentProvidedOpPath,
				},
				"/": {
					Dir: &providedOpPath,
				},
			}

			fakeCaller := new(FakeCaller)

			objectUnderTest := _opCaller{
				caller: fakeCaller,
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
			expectedOutputName := "expectedOutputName"

			providedOpPath := "testdata/opCaller"

			providedOpCall := &model.OpCall{
				BaseCall: model.BaseCall{
					OpPath: providedOpPath,
				},
				OpID: "providedOpId",
			}

			providedOpCallSpec := &model.OpCallSpec{
				Outputs: map[string]string{
					expectedOutputName: "",
				},
			}

			callOutputs := map[string]*model.Value{
				expectedOutputName: {
					String: new(string),
				},
				// include unbound output to ensure it's not added to scope
				"unexpectedOutputName": new(model.Value),
			}

			expectedOutputs := map[string]*model.Value{
				expectedOutputName: callOutputs[expectedOutputName],
			}

			fakeCaller := new(FakeCaller)
			fakeCaller.CallReturns(
				callOutputs,
				nil,
			)

			objectUnderTest := _opCaller{
				caller: fakeCaller,
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
			Expect(actualErr).To(BeNil())
			Expect(actualOutputs).To(Equal(expectedOutputs))
		})
	})
})
