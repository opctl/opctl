package bracketed

import (
	"errors"
	"fmt"

	"github.com/opctl/sdk-golang/op/interpreter/interpolater/dereferencer/identifier/bracketed/item"
	"github.com/opctl/sdk-golang/op/interpreter/interpolater/dereferencer/identifier/value"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
)

var _ = Context("DeReferencer", func() {
	Context("NewDeReferencer", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewDeReferencer()).Should(Not(BeNil()))
		})
	})
	Context("DeReference", func() {
		Context("ref doesn't start w/ '['", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"

				objectUnderTest := _deReferencer{}
				expectedErr := fmt.Errorf("unable to deReference '%v'; expected '['", providedRef)

				/* act */
				_, _, actualErr := objectUnderTest.DeReference(
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

				objectUnderTest := _deReferencer{}
				expectedErr := fmt.Errorf("unable to deReference '%v'; expected ']'", providedRef)

				/* act */
				_, _, actualErr := objectUnderTest.DeReference(
					providedRef,
					nil,
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		It("should call coerceToArrayOrObjecter.CoerceToArrayOrObject w/ expected args", func() {

			/* arrange */
			providedData := model.Value{String: new(string)}

			fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
			// err to trigger immediate return
			fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _deReferencer{
				coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
			}

			/* act */
			objectUnderTest.DeReference(
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
				providedData := model.Value{String: new(string)}

				fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
				coerceToArrayOrObjectErr := errors.New("coerceToArrayOrObjectErr")
				fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(nil, coerceToArrayOrObjectErr)

				expectedErr := fmt.Errorf("unable to deReference '%v'; error was %v", providedRef, coerceToArrayOrObjectErr.Error())

				objectUnderTest := _deReferencer{
					coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
				}

				/* act */
				_, _, actualErr := objectUnderTest.DeReference(
					providedRef,
					&providedData,
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("data is array", func() {
			It("should call itemDeReferencer.DeReference w/ expected args", func() {
				/* arrange */
				providedRefIdentifier := "dummyIdentifier"
				providedRef := fmt.Sprintf("[%v]", providedRefIdentifier)

				fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
				coercedArray := model.Value{Array: []interface{}{"dummyItem"}}
				fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(&coercedArray, nil)

				fakeItemDeReferencer := new(item.FakeDeReferencer)
				// err to trigger immediate return
				fakeItemDeReferencer.DeReferenceReturns(nil, errors.New("dummyErr"))

				objectUnderTest := _deReferencer{
					coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
					itemDeReferencer:        fakeItemDeReferencer,
				}

				/* act */
				objectUnderTest.DeReference(
					providedRef,
					new(model.Value),
				)

				/* assert */
				actualIdentifier,
					actualArray := fakeItemDeReferencer.DeReferenceArgsForCall(0)

				Expect(actualIdentifier).To(Equal(providedRefIdentifier))
				Expect(actualArray).To(Equal(coercedArray))
			})
			Context("itemDeReferencer.DeReference errs", func() {
				It("should return expected result", func() {
					/* arrange */

					fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
					coercedArray := model.Value{Array: []interface{}{"dummyItem"}}
					fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(&coercedArray, nil)

					fakeItemDeReferencer := new(item.FakeDeReferencer)
					itemDeReferencerErr := errors.New("itemDereferencerErr")
					fakeItemDeReferencer.DeReferenceReturns(nil, itemDeReferencerErr)

					objectUnderTest := _deReferencer{
						coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
						itemDeReferencer:        fakeItemDeReferencer,
					}

					/* act */
					_, _, actualErr := objectUnderTest.DeReference(
						"[]",
						new(model.Value),
					)

					/* assert */
					Expect(actualErr).To(Equal(itemDeReferencerErr))
				})
			})
			Context("itemDeReferencer.DeReferenceItem doesn't err", func() {

				It("should return expected result", func() {
					/* arrange */

					fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
					coercedArray := model.Value{Array: []interface{}{"dummyItem"}}
					fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(&coercedArray, nil)

					fakeItemDeReferencer := new(item.FakeDeReferencer)
					deReferencedItemValue := model.Value{}
					fakeItemDeReferencer.DeReferenceReturns(&deReferencedItemValue, nil)

					objectUnderTest := _deReferencer{
						coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
						itemDeReferencer:        fakeItemDeReferencer,
					}

					/* act */
					actualRefRemainder, actualData, actualErr := objectUnderTest.DeReference(
						"[]",
						new(model.Value),
					)

					/* assert */
					Expect(actualRefRemainder).To(BeEmpty())
					Expect(*actualData).To(Equal(deReferencedItemValue))
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
				coercedObject := model.Value{Object: map[string]interface{}{providedRefIdentifier: "dummyValue"}}
				fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(&coercedObject, nil)

				fakeValueConstructor := new(value.FakeConstructor)
				// err to trigger immediate return
				fakeValueConstructor.ConstructReturns(nil, errors.New("dummyErr"))

				objectUnderTest := _deReferencer{
					coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
					valueConstructor:        fakeValueConstructor,
				}

				/* act */
				objectUnderTest.DeReference(
					providedRef,
					new(model.Value),
				)

				/* assert */
				actualObject := fakeValueConstructor.ConstructArgsForCall(0)

				Expect(actualObject).To(Equal(coercedObject.Object[providedRefIdentifier]))
			})
			Context("valueConstructor.Construct errs", func() {
				It("should return expected result", func() {
					/* arrange */

					fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
					coercedObject := model.Value{Object: map[string]interface{}{"dummyItem": nil}}
					fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(&coercedObject, nil)

					fakeValueConstructor := new(value.FakeConstructor)
					constructErr := errors.New("constructErr")
					fakeValueConstructor.ConstructReturns(nil, constructErr)

					expectedErr := fmt.Errorf("unable to deReference property; error was %v", constructErr.Error())

					objectUnderTest := _deReferencer{
						coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
						valueConstructor:        fakeValueConstructor,
					}

					/* act */
					_, _, actualErr := objectUnderTest.DeReference(
						"[]",
						new(model.Value),
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("valueConstructor.Construct doesn't err", func() {

				It("should return expected result", func() {
					/* arrange */

					fakeCoerceToArrayOrObjecter := new(fakeCoerceToArrayOrObjecter)
					coercedObject := model.Value{Object: map[string]interface{}{"dummyItem": nil}}
					fakeCoerceToArrayOrObjecter.CoerceToArrayOrObjectReturns(&coercedObject, nil)

					fakeValueConstructor := new(value.FakeConstructor)
					deReferencedItemValue := model.Value{}
					fakeValueConstructor.ConstructReturns(&deReferencedItemValue, nil)

					objectUnderTest := _deReferencer{
						coerceToArrayOrObjecter: fakeCoerceToArrayOrObjecter,
						valueConstructor:        fakeValueConstructor,
					}

					/* act */
					actualRefRemainder, actualData, actualErr := objectUnderTest.DeReference(
						"[]",
						new(model.Value),
					)

					/* assert */
					Expect(actualRefRemainder).To(BeEmpty())
					Expect(*actualData).To(Equal(deReferencedItemValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
})
