package loopable

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	arrayFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/array/fakes"
	objectFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/object/fakes"
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

			providedScope := map[string]*model.Value{}

			fakeArrayInterpreter := new(arrayFakes.FakeInterpreter)

			objectUnderTest := _interpreter{
				arrayInterpreter: fakeArrayInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedExpression,
				providedScope,
			)

			/* assert */
			actualScope,
				actualExpression := fakeArrayInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualExpression).To(Equal(providedExpression))
		})
		Context("arrayInterpreter.Interpret doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeArrayInterpreter := new(arrayFakes.FakeInterpreter)

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

				providedScope := map[string]*model.Value{}

				fakeArrayInterpreter := new(arrayFakes.FakeInterpreter)
				fakeArrayInterpreter.InterpretReturns(
					nil,
					errors.New(""),
				)

				fakeObjectInterpreter := new(objectFakes.FakeInterpreter)

				objectUnderTest := _interpreter{
					arrayInterpreter:  fakeArrayInterpreter,
					objectInterpreter: fakeObjectInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedExpression,
					providedScope,
				)

				/* assert */
				actualScope,
					actualExpression := fakeObjectInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualExpression).To(Equal(providedExpression))
			})
			It("should return expected result", func() {
				/* arrange */
				fakeArrayInterpreter := new(arrayFakes.FakeInterpreter)
				fakeArrayInterpreter.InterpretReturns(
					nil,
					errors.New(""),
				)

				fakeObjectInterpreter := new(objectFakes.FakeInterpreter)

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
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
	})
})
