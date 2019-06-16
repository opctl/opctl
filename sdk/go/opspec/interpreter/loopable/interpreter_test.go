package loopable

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdk/go/data"
	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/array"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/object"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		It("should call arrayInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedExpression := "providedExpression"

			providedOpHandle := new(data.FakeHandle)
			providedScope := map[string]*model.Value{}

			fakeArrayInterpreter := new(array.FakeInterpreter)

			objectUnderTest := _interpreter{
				arrayInterpreter: fakeArrayInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedExpression,
				providedOpHandle,
				providedScope,
			)

			/* assert */
			actualScope,
				actualExpression,
				actualOpHandle := fakeArrayInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualExpression).To(Equal(providedExpression))
			Expect(actualOpHandle).To(Equal(providedOpHandle))
		})
		Context("arrayInterpreter.Interpret doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeArrayInterpreter := new(array.FakeInterpreter)

				expectedValue := &model.Value{}
				fakeArrayInterpreter.InterpretReturns(
					expectedValue,
					nil,
				)

				expectedResult := &model.Value{}

				objectUnderTest := _interpreter{
					arrayInterpreter: fakeArrayInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					"dummyExpression",
					new(data.FakeHandle),
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
		Context("arrayInterpreter.Interpret errs", func() {
			It("should call objectInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedExpression := "providedExpression"

				providedOpHandle := new(data.FakeHandle)
				providedScope := map[string]*model.Value{}

				fakeArrayInterpreter := new(array.FakeInterpreter)
				fakeArrayInterpreter.InterpretReturns(
					nil,
					errors.New(""),
				)

				fakeObjectInterpreter := new(object.FakeInterpreter)

				objectUnderTest := _interpreter{
					arrayInterpreter:  fakeArrayInterpreter,
					objectInterpreter: fakeObjectInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedExpression,
					providedOpHandle,
					providedScope,
				)

				/* assert */
				actualScope,
					actualExpression,
					actualOpHandle := fakeObjectInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualExpression).To(Equal(providedExpression))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
			})
			It("should return expected result", func() {
				/* arrange */
				fakeArrayInterpreter := new(array.FakeInterpreter)
				fakeArrayInterpreter.InterpretReturns(
					nil,
					errors.New(""),
				)

				fakeObjectInterpreter := new(object.FakeInterpreter)

				expectedValue := &model.Value{}
				fakeObjectInterpreter.InterpretReturns(
					expectedValue,
					nil,
				)

				expectedResult := &model.Value{}

				objectUnderTest := _interpreter{
					arrayInterpreter:  fakeArrayInterpreter,
					objectInterpreter: fakeObjectInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					"dummyExpression",
					new(data.FakeHandle),
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
	})
})
