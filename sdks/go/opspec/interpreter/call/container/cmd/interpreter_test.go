package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	strFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/str/fakes"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		Context("scgContainerCallCmd not empty", func() {
			It("should call stringInterpreter.Interpret w/ expected args for each container.Cmd entry", func() {
				/* arrange */
				providedString1 := "dummyString1"
				providedCurrentScope := map[string]*model.Value{
					"name1": {String: &providedString1},
				}
				providedOpHandle := new(modelFakes.FakeDataHandle)

				providedSCGContainerCallCmd := []interface{}{
					"dummy1",
					"dummy2",
				}

				fakeStrInterpreter := new(strFakes.FakeInterpreter)
				fakeStrInterpreter.InterpretReturns(&model.Value{String: new(string)}, nil)

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStrInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedCurrentScope,
					providedSCGContainerCallCmd,
					providedOpHandle,
				)

				/* assert */
				for expectedCmdIndex, expectedCmdEntry := range providedSCGContainerCallCmd {
					actualScope,
						actualCmdEntry,
						actualOpHandle := fakeStrInterpreter.InterpretArgsForCall(expectedCmdIndex)

					Expect(actualScope).To(Equal(providedCurrentScope))
					Expect(actualCmdEntry).To(Equal(expectedCmdEntry))
					Expect(actualOpHandle).To(Equal(providedOpHandle))
				}
			})
			It("should return expected dcg.Cmd", func() {
				/* arrange */
				expectedResult := []string{
					"dummyCmdEntry1",
					"dummyCmdEntry2",
				}

				providedSCGContainerCallCmd := []interface{}{
					"dummy1",
					"dummy2",
				}

				fakeStrInterpreter := new(strFakes.FakeInterpreter)
				fakeStrInterpreter.InterpretReturnsOnCall(0, &model.Value{String: &expectedResult[0]}, nil)
				fakeStrInterpreter.InterpretReturnsOnCall(1, &model.Value{String: &expectedResult[1]}, nil)

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStrInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCGContainerCallCmd,
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
	})
})
