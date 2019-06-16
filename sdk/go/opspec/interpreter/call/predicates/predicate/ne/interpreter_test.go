package ne

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdk/go/data"
	"github.com/opctl/opctl/sdk/go/model"
	stringPkg "github.com/opctl/opctl/sdk/go/opspec/interpreter/string"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		It("should call stringInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedExpressions := []interface{}{
				"expression0",
				"expression1",
			}

			providedOpHandle := new(data.FakeHandle)
			providedScope := map[string]*model.Value{}

			fakeStringInterpreter := new(stringPkg.FakeInterpreter)
			fakeStringInterpreter.InterpretReturns(
				&model.Value{String: new(string)},
				nil,
			)

			objectUnderTest := _interpreter{
				stringInterpreter: fakeStringInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedExpressions,
				providedOpHandle,
				providedScope,
			)

			/* assert */
			actualScope0,
				actualExpression0,
				actualOpHandle0 := fakeStringInterpreter.InterpretArgsForCall(0)

			Expect(actualScope0).To(Equal(providedScope))
			Expect(actualExpression0).To(Equal(providedExpressions[0]))
			Expect(actualOpHandle0).To(Equal(providedOpHandle))

			actualScope1,
				actualExpression1,
				actualOpHandle1 := fakeStringInterpreter.InterpretArgsForCall(1)

			Expect(actualScope1).To(Equal(providedScope))
			Expect(actualExpression1).To(Equal(providedExpressions[1]))
			Expect(actualOpHandle1).To(Equal(providedOpHandle))
		})
		Context("stringInterpreter.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeStringInterpreter := new(stringPkg.FakeInterpreter)

				expectedError := errors.New("expectedError")
				fakeStringInterpreter.InterpretReturns(
					nil,
					expectedError,
				)

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStringInterpreter,
				}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					[]interface{}{"dummyExpression"},
					new(data.FakeHandle),
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("stringInterpreter.Interpret returns equal items", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeStringInterpreter := new(stringPkg.FakeInterpreter)

				str := "str"
				fakeStringInterpreter.InterpretReturns(
					&model.Value{String: &str},
					nil,
				)

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStringInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					[]interface{}{
						"expression0",
						"expression1",
					},
					new(data.FakeHandle),
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualResult).To(BeFalse())
			})
		})
		Context("stringInterpreter.Interpret returns unequal items", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeStringInterpreter := new(stringPkg.FakeInterpreter)

				zero := "zero"
				fakeStringInterpreter.InterpretReturnsOnCall(
					0,
					&model.Value{String: &zero},
					nil,
				)

				one := "one"
				fakeStringInterpreter.InterpretReturnsOnCall(
					1,
					&model.Value{String: &one},
					nil,
				)

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStringInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					[]interface{}{
						"expression0",
						"expression1",
					},
					new(data.FakeHandle),
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualResult).To(BeTrue())
			})
		})
	})
})
