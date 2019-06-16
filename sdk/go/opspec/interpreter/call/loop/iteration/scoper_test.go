package iteration

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdk/go/data"
	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/loopable"
)

var _ = Context("scoper", func() {
	Context("NewScoper", func() {
		It("should return scoper", func() {
			/* arrange/act/assert */
			Expect(NewScoper()).To(Not(BeNil()))
		})
	})

	Context("Scope", func() {
		Context("nil != scg.Index", func() {
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
					&model.SCGLoop{
						Index: &indexName,
					},
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualScope).To(Equal(expectedScope))
			})
		})
		Context("nil != scg.For and nil != scg.For.Each", func() {
			It("should call loopableInterpreter w/ expected args", func() {
				/* arrange */
				providedForEach := "providedForEach"

				providedScope := map[string]*model.Value{
					"name1": &model.Value{String: new(string)},
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
					&model.SCGLoop{
						For: &model.SCGLoopFor{
							Each: providedForEach,
						},
					},
					providedOpHandle,
				)

				/* assert */
				actualExpression,
					actualOpHandle,
					actualScope := fakeLoopableInterpreter.InterpretArgsForCall(0)

				Expect(actualExpression).To(Equal(providedForEach))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
			})
			Context("loopableInterpreter errs", func() {

				It("should return expected result", func() {
					/* arrange */
					providedForEach := "providedForEach"

					providedScope := map[string]*model.Value{
						"name1": &model.Value{String: new(string)},
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
						&model.SCGLoop{
							For: &model.SCGLoopFor{
								Each: providedForEach,
							},
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
