package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	stringPkg "github.com/opctl/opctl/sdks/go/opspec/interpreter/string"
	"github.com/opctl/opctl/sdks/go/types"
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
				providedCurrentScope := map[string]*types.Value{
					"name1": {String: &providedString1},
				}
				providedOpHandle := new(data.FakeHandle)

				providedSCGContainerCallCmd := []interface{}{
					"dummy1",
					"dummy2",
				}

				fakeStringInterpreter := new(stringPkg.FakeInterpreter)
				fakeStringInterpreter.InterpretReturns(&types.Value{String: new(string)}, nil)

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStringInterpreter,
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
						actualOpHandle := fakeStringInterpreter.InterpretArgsForCall(expectedCmdIndex)

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

				fakeStringInterpreter := new(stringPkg.FakeInterpreter)
				fakeStringInterpreter.InterpretReturnsOnCall(0, &types.Value{String: &expectedResult[0]}, nil)
				fakeStringInterpreter.InterpretReturnsOnCall(1, &types.Value{String: &expectedResult[1]}, nil)

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStringInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*types.Value{},
					providedSCGContainerCallCmd,
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
	})
})
