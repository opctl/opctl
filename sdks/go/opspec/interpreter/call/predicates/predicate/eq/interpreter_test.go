package eq

import (
	"errors"

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
		It("should call stringInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedExpressions := []interface{}{
				"expression0",
				"expression1",
			}

			providedOpHandle := new(modelFakes.FakeDataHandle)
			providedScope := map[string]*model.Value{}

			fakeStrInterpreter := new(strFakes.FakeInterpreter)
			fakeStrInterpreter.InterpretReturns(
				&model.Value{String: new(string)},
				nil,
			)

			objectUnderTest := _interpreter{
				stringInterpreter: fakeStrInterpreter,
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
				actualOpHandle0 := fakeStrInterpreter.InterpretArgsForCall(0)

			Expect(actualScope0).To(Equal(providedScope))
			Expect(actualExpression0).To(Equal(providedExpressions[0]))
			Expect(actualOpHandle0).To(Equal(providedOpHandle))

			actualScope1,
				actualExpression1,
				actualOpHandle1 := fakeStrInterpreter.InterpretArgsForCall(1)

			Expect(actualScope1).To(Equal(providedScope))
			Expect(actualExpression1).To(Equal(providedExpressions[1]))
			Expect(actualOpHandle1).To(Equal(providedOpHandle))
		})
		Context("stringInterpreter.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeStrInterpreter := new(strFakes.FakeInterpreter)

				expectedError := errors.New("expectedError")
				fakeStrInterpreter.InterpretReturns(
					nil,
					expectedError,
				)

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStrInterpreter,
				}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					[]interface{}{"dummyExpression"},
					new(modelFakes.FakeDataHandle),
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("stringInterpreter.Interpret returns equal items", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeStrInterpreter := new(strFakes.FakeInterpreter)

				str := "str"
				fakeStrInterpreter.InterpretReturns(
					&model.Value{String: &str},
					nil,
				)

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStrInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					[]interface{}{
						"expression0",
						"expression1",
					},
					new(modelFakes.FakeDataHandle),
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualResult).To(BeTrue())
			})
		})
		Context("stringInterpreter.Interpret returns unequal items", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeStrInterpreter := new(strFakes.FakeInterpreter)

				zero := "zero"
				fakeStrInterpreter.InterpretReturnsOnCall(
					0,
					&model.Value{String: &zero},
					nil,
				)

				one := "one"
				fakeStrInterpreter.InterpretReturnsOnCall(
					1,
					&model.Value{String: &one},
					nil,
				)

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStrInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					[]interface{}{
						"expression0",
						"expression1",
					},
					new(modelFakes.FakeDataHandle),
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualResult).To(BeFalse())
			})
		})
	})
})
