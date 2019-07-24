package bracketed

import (
	"errors"
	"fmt"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/bracketed/item"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).Should(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		Context("ref doesn't start w/ '['", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"

				objectUnderTest := _interpreter{}
				expectedErr := fmt.Errorf("unable to interpret '%v'; expected '['", providedRef)

				/* act */
				_, _, actualErr := objectUnderTest.Interpret(
					providedRef,
					nil,
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("ref doesn't contain ']", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "[dummyRef"

				objectUnderTest := _interpreter{}
				expectedErr := fmt.Errorf("unable to interpret '%v'; expected ']'", providedRef)

				/* act */
				_, _, actualErr := objectUnderTest.Interpret(
					providedRef,
					nil,
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		It("should call coerceToArrayOrObjecter.CoerceToArrayOrObject w/ expected args", func() {

			/* arrange */
			providedData := types.Value{String: new(string)}

			fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
			// err to trigger immediate return
			fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _interpreter{
				coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
			}

			/* act */
			objectUnderTest.Interpret(
				"[]",
				&providedData,
			)

			/* assert */
			actualData := fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectArgsForCall(0)

			Expect(*actualData).To(Equal(providedData))
		})
		Context("coerceToArrayOrObjecter.CoerceToArrayOrObject errs", func() {

			It("should return expected results", func() {

				/* arrange */
				providedRef := "[]"
				providedData := types.Value{String: new(string)}

				fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
				coerceToArrayOrObjectErr := errors.New("coerceToArrayOrObjectErr")
				fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(nil, coerceToArrayOrObjectErr)

				expectedErr := fmt.Errorf("unable to interpret '%v'; error was %v", providedRef, coerceToArrayOrObjectErr.Error())

				objectUnderTest := _interpreter{
					coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
				}

				/* act */
				_, _, actualErr := objectUnderTest.Interpret(
					providedRef,
					&providedData,
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("data is array", func() {
			It("should call itemInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedRefIdentifier := "dummyIdentifier"
				providedRef := fmt.Sprintf("[%v]", providedRefIdentifier)

				fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
				coercedArray := types.Value{Array: new([]interface{})}
				fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(&coercedArray, nil)

				fakeItemInterpreter := new(item.FakeInterpreter)
				// err to trigger immediate return
				fakeItemInterpreter.InterpretReturns(nil, errors.New("dummyErr"))

				objectUnderTest := _interpreter{
					coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
					itemInterpreter:         fakeItemInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedRef,
					new(types.Value),
				)

				/* assert */
				actualIdentifier,
					actualArray := fakeItemInterpreter.InterpretArgsForCall(0)

				Expect(actualIdentifier).To(Equal(providedRefIdentifier))
				Expect(actualArray).To(Equal(coercedArray))
			})
			Context("itemInterpreter.Interpret errs", func() {
				It("should return expected result", func() {
					/* arrange */

					fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
					coercedArray := types.Value{Array: new([]interface{})}
					fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(&coercedArray, nil)

					fakeItemInterpreter := new(item.FakeInterpreter)
					itemInterpreterErr := errors.New("itemDereferencerErr")
					fakeItemInterpreter.InterpretReturns(nil, itemInterpreterErr)

					objectUnderTest := _interpreter{
						coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
						itemInterpreter:         fakeItemInterpreter,
					}

					/* act */
					_, _, actualErr := objectUnderTest.Interpret(
						"[]",
						new(types.Value),
					)

					/* assert */
					Expect(actualErr).To(Equal(itemInterpreterErr))
				})
			})
			Context("itemInterpreter.InterpretItem doesn't err", func() {

				It("should return expected result", func() {
					/* arrange */

					fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
					coercedArray := types.Value{Array: new([]interface{})}
					fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(&coercedArray, nil)

					fakeItemInterpreter := new(item.FakeInterpreter)
					interpretdItemValue := types.Value{}
					fakeItemInterpreter.InterpretReturns(&interpretdItemValue, nil)

					objectUnderTest := _interpreter{
						coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
						itemInterpreter:         fakeItemInterpreter,
					}

					/* act */
					actualRefRemainder, actualData, actualErr := objectUnderTest.Interpret(
						"[]",
						new(types.Value),
					)

					/* assert */
					Expect(actualRefRemainder).To(BeEmpty())
					Expect(*actualData).To(Equal(interpretdItemValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("data is Object", func() {
			It("should call valueConstructor.Construct w/ expected args", func() {
				/* arrange */
				providedRefIdentifier := "dummyIdentifier"
				providedRef := fmt.Sprintf("[%v]", providedRefIdentifier)

				fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
				object := &map[string]interface{}{providedRefIdentifier: "dummyValue"}
				coercedObject := types.Value{Object: object}
				fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(&coercedObject, nil)

				fakeValueConstructor := new(value.FakeConstructor)
				// err to trigger immediate return
				fakeValueConstructor.ConstructReturns(nil, errors.New("dummyErr"))

				objectUnderTest := _interpreter{
					coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
					valueConstructor:        fakeValueConstructor,
				}

				/* act */
				objectUnderTest.Interpret(
					providedRef,
					new(types.Value),
				)

				/* assert */
				actualObject := fakeValueConstructor.ConstructArgsForCall(0)

				Expect(actualObject).To(Equal((*coercedObject.Object)[providedRefIdentifier]))
			})
			Context("valueConstructor.Construct errs", func() {
				It("should return expected result", func() {
					/* arrange */

					fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
					coercedObject := types.Value{Object: new(map[string]interface{})}
					fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(&coercedObject, nil)

					fakeValueConstructor := new(value.FakeConstructor)
					constructErr := errors.New("constructErr")
					fakeValueConstructor.ConstructReturns(nil, constructErr)

					expectedErr := fmt.Errorf("unable to interpret property; error was %v", constructErr.Error())

					objectUnderTest := _interpreter{
						coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
						valueConstructor:        fakeValueConstructor,
					}

					/* act */
					_, _, actualErr := objectUnderTest.Interpret(
						"[]",
						new(types.Value),
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("valueConstructor.Construct doesn't err", func() {

				It("should return expected result", func() {
					/* arrange */

					fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
					coercedObject := types.Value{Object: new(map[string]interface{})}
					fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(&coercedObject, nil)

					fakeValueConstructor := new(value.FakeConstructor)
					interpretdItemValue := types.Value{}
					fakeValueConstructor.ConstructReturns(&interpretdItemValue, nil)

					objectUnderTest := _interpreter{
						coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
						valueConstructor:        fakeValueConstructor,
					}

					/* act */
					actualRefRemainder, actualData, actualErr := objectUnderTest.Interpret(
						"[]",
						new(types.Value),
					)

					/* assert */
					Expect(actualRefRemainder).To(BeEmpty())
					Expect(*actualData).To(Equal(interpretdItemValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
})
