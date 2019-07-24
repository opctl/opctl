package iteration

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("scoper", func() {
	Context("NewScoper", func() {
		It("should return scoper", func() {
			/* arrange/act/assert */
			Expect(NewScoper()).To(Not(BeNil()))
		})
	})

	Context("Scope", func() {
		Context("nil != scg.Vars", func() {
			Context("nil == scg.Range", func() {
				Context("nil != scg.Vars.Index", func() {
					It("should return expected result", func() {
						/* arrange */
						indexValue := 2
						indexValueAsFloat64 := float64(indexValue)

						indexName := "indexName"

						objectUnderTest := _scoper{}

						expectedScope := map[string]*types.Value{
							indexName: &types.Value{Number: &indexValueAsFloat64},
						}

						/* act */
						actualScope, _ := objectUnderTest.Scope(
							indexValue,
							map[string]*types.Value{},
							nil,
							&types.SCGLoopVars{
								Index: &indexName,
							},
							new(data.FakeHandle),
						)

						/* assert */
						Expect(actualScope).To(Equal(expectedScope))
					})
				})
			})
			Context("nil != scg.Range", func() {
				It("should call loopableInterpreter w/ expected args", func() {
					/* arrange */
					providedLoopRange := "providedLoopRange"

					providedScope := map[string]*types.Value{
						"name1": &types.Value{String: new(string)},
					}
					providedOpHandle := new(data.FakeHandle)

					fakeLoopableInterpreter := new(loopable.FakeInterpreter)
					// err to trigger immediate return
					fakeLoopableInterpreter.InterpretReturns(nil, errors.New("dummyErr"))

					objectUnderTest := _scoper{
						loopableInterpreter: fakeLoopableInterpreter,
					}

					/* act */
					objectUnderTest.Scope(
						0,
						providedScope,
						providedLoopRange,
						&types.SCGLoopVars{},
						providedOpHandle,
					)

					/* assert */
					actualExpression,
						actualOpHandle,
						actualScope := fakeLoopableInterpreter.InterpretArgsForCall(0)

					Expect(actualExpression).To(Equal(providedLoopRange))
					Expect(actualScope).To(Equal(providedScope))
					Expect(actualOpHandle).To(Equal(providedOpHandle))
				})
				Context("loopableInterpreter errs", func() {

					It("should return expected result", func() {
						/* arrange */
						providedLoopRange := "providedLoopRange"

						providedScope := map[string]*types.Value{
							"name1": &types.Value{String: new(string)},
						}
						providedOpHandle := new(data.FakeHandle)

						fakeLoopableInterpreter := new(loopable.FakeInterpreter)
						expectedErr := errors.New("expectedErr")
						fakeLoopableInterpreter.InterpretReturns(nil, expectedErr)

						objectUnderTest := _scoper{
							loopableInterpreter: fakeLoopableInterpreter,
						}

						/* act */
						_, actualErr := objectUnderTest.Scope(
							0,
							providedScope,
							providedLoopRange,
							&types.SCGLoopVars{
								Index: new(string),
							},
							providedOpHandle,
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
			})
		})
	})
})
