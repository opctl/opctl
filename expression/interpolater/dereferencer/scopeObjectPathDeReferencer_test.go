package dereferencer

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/pkg/errors"
	"strings"
)

var _ = Context("scopeObjectPathDeReferencer", func() {
	Context("ref is scope object path ref", func() {
		Context("property doesn't exist", func() {
			It("should return expected result", func() {
				/* arrange */

				// build up object
				pathSegment1 := "pathSegment1"
				pathSegment1Value := map[string]interface{}{}

				objectRef := "dummyObjectRef"
				objectValue := map[string]interface{}{pathSegment1: pathSegment1Value}

				providedRef := strings.Join([]string{objectRef, pathSegment1, "doesNotExist"}, ".")

				fakeData := new(data.Fake)
				fakeData.CoerceToObjectReturns(&model.Value{Object: objectValue}, nil)

				objectUnderTest := _scopeObjectPathDeReferencer{
					data: fakeData,
				}

				/* act */
				actualString, actualOk, actualErr := objectUnderTest.DeReferenceScopeObjectPath(
					providedRef,
					map[string]*model.Value{
						objectRef: {Object: objectValue},
					},
				)

				/* assert */
				Expect(actualString).To(Equal(""))
				Expect(actualOk).To(Equal(false))
				Expect(actualErr).To(Equal(fmt.Errorf("unable to deReference '%v'; path doesn't exist", providedRef)))
			})
		})
		Context("property is float64", func() {
			It("should call data.Coerce w/ expected args", func() {
				/* arrange */

				// build up object
				pathSegment2 := "pathSegment2"
				pathSegment2Value := 2.2

				pathSegment1 := "pathSegment1"
				pathSegment1Value := map[string]interface{}{pathSegment2: pathSegment2Value}

				objectRef := "dummyObjectRef"
				objectValue := map[string]interface{}{pathSegment1: pathSegment1Value}

				providedRef := strings.Join([]string{objectRef, pathSegment1, pathSegment2}, ".")

				fakeData := new(data.Fake)
				fakeData.CoerceToObjectReturns(&model.Value{Object: objectValue}, nil)
				// err to trigger immediate return
				fakeData.CoerceToStringReturns(nil, errors.New("dummyError"))

				objectUnderTest := _scopeObjectPathDeReferencer{
					data: fakeData,
				}

				/* act */
				objectUnderTest.DeReferenceScopeObjectPath(
					providedRef,
					map[string]*model.Value{
						objectRef: {Object: objectValue},
					},
				)

				/* assert */
				actualValue := fakeData.CoerceToStringArgsForCall(0)
				Expect(*actualValue).To(Equal(model.Value{Number: &pathSegment2Value}))
			})
		})
		Context("property is map[string]interface{}", func() {
			It("should call data.Coerce w/ expected args", func() {
				/* arrange */

				// build up object
				pathSegment2 := "pathSegment2"
				pathSegment2Value := map[string]interface{}{"dummyKey": "dummyValue"}

				pathSegment1 := "pathSegment1"
				pathSegment1Value := map[string]interface{}{pathSegment2: pathSegment2Value}

				objectRef := "dummyObjectRef"
				objectValue := map[string]interface{}{pathSegment1: pathSegment1Value}

				providedRef := strings.Join([]string{objectRef, pathSegment1, pathSegment2}, ".")

				fakeData := new(data.Fake)
				fakeData.CoerceToObjectReturns(&model.Value{Object: objectValue}, nil)
				// err to trigger immediate return
				fakeData.CoerceToStringReturns(nil, errors.New("dummyError"))

				objectUnderTest := _scopeObjectPathDeReferencer{
					data: fakeData,
				}

				/* act */
				objectUnderTest.DeReferenceScopeObjectPath(
					providedRef,
					map[string]*model.Value{
						objectRef: {Object: objectValue},
					},
				)

				/* assert */
				actualValue := fakeData.CoerceToStringArgsForCall(0)
				Expect(*actualValue).To(Equal(model.Value{Object: pathSegment2Value}))
			})
		})
		Context("property is string", func() {
			It("should return expected result", func() {
				/* arrange */

				// build up object
				pathSegment2 := "pathSegment2"
				pathSegment2Value := "dummyString"

				pathSegment1 := "pathSegment1"
				pathSegment1Value := map[string]interface{}{pathSegment2: pathSegment2Value}

				objectRef := "dummyObjectRef"
				objectValue := map[string]interface{}{pathSegment1: pathSegment1Value}

				providedRef := strings.Join([]string{objectRef, pathSegment1, pathSegment2}, ".")

				fakeData := new(data.Fake)
				fakeData.CoerceToObjectReturns(&model.Value{Object: objectValue}, nil)

				objectUnderTest := _scopeObjectPathDeReferencer{
					data: fakeData,
				}

				/* act */
				actualString, actualOk, actualErr := objectUnderTest.DeReferenceScopeObjectPath(
					providedRef,
					map[string]*model.Value{
						objectRef: {Object: objectValue},
					},
				)

				/* assert */
				Expect(actualString).To(Equal(pathSegment2Value))
				Expect(actualOk).To(Equal(true))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("property is []interface{}", func() {
			It("should call data.Coerce w/ expected args", func() {
				/* arrange */

				// build up object
				pathSegment2 := "pathSegment2"
				pathSegment2Value := []interface{}{"string", 2.2}

				pathSegment1 := "pathSegment1"
				pathSegment1Value := map[string]interface{}{pathSegment2: pathSegment2Value}

				objectRef := "dummyObjectRef"
				objectValue := map[string]interface{}{pathSegment1: pathSegment1Value}

				providedRef := strings.Join([]string{objectRef, pathSegment1, pathSegment2}, ".")

				fakeData := new(data.Fake)
				fakeData.CoerceToObjectReturns(&model.Value{Object: objectValue}, nil)
				// err to trigger immediate return
				fakeData.CoerceToStringReturns(nil, errors.New("dummyError"))

				objectUnderTest := _scopeObjectPathDeReferencer{
					data: fakeData,
				}

				/* act */
				objectUnderTest.DeReferenceScopeObjectPath(
					providedRef,
					map[string]*model.Value{
						objectRef: {Object: objectValue},
					},
				)

				/* assert */
				actualValue := fakeData.CoerceToStringArgsForCall(0)
				Expect(*actualValue).To(Equal(model.Value{Array: pathSegment2Value}))
			})
		})
	})
})
