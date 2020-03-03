package iteration

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	loopableFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable/fakes"
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

						expectedScope := map[string]*model.Value{
							indexName: &model.Value{Number: &indexValueAsFloat64},
						}

						/* act */
						actualScope, _ := objectUnderTest.Scope(
							indexValue,
							map[string]*model.Value{},
							nil,
							&model.SCGLoopVars{
								Index: &indexName,
							},
							new(modelFakes.FakeDataHandle),
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

					providedScope := map[string]*model.Value{
						"name1": &model.Value{String: new(string)},
					}
					providedOpHandle := new(modelFakes.FakeDataHandle)

					fakeLoopableInterpreter := new(loopableFakes.FakeInterpreter)
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
						&model.SCGLoopVars{},
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

						providedScope := map[string]*model.Value{
							"name1": &model.Value{String: new(string)},
						}
						providedOpHandle := new(modelFakes.FakeDataHandle)

						fakeLoopableInterpreter := new(loopableFakes.FakeInterpreter)
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
							&model.SCGLoopVars{
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
