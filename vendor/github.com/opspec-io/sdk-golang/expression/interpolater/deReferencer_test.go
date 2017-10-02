package interpolater

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/pkg/errors"
	"strings"
)

var _ = Context("deReferencer", func() {
	Context("pkg content ref", func() {
		It("should call pkgHandle.GetContent w/ expected args", func() {
			/* arrange */
			providedRef := "/dummyRef"
			fakePkgHandle := new(pkg.FakeHandle)
			// err to trigger immediate return
			fakePkgHandle.GetContentReturns(nil, errors.New("dummyError"))

			objectUnderTest := _deReferencer{}

			/* act */
			objectUnderTest.DeReference(
				providedRef,
				map[string]*model.Value{},
				fakePkgHandle,
			)

			/* assert */
			actualContext,
				actualContentPath := fakePkgHandle.GetContentArgsForCall(0)

			Expect(actualContext).To(Equal(context.TODO()))
			Expect(actualContentPath).To(Equal(providedRef))
		})
	})
	Context("scope object ref w/ path", func() {
		Context("float64 at path", func() {
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
				// err to trigger immediate return
				fakeData.CoerceToStringReturns(nil, errors.New("dummyError"))

				objectUnderTest := _deReferencer{
					data: fakeData,
				}

				/* act */
				objectUnderTest.DeReference(
					providedRef,
					map[string]*model.Value{
						objectRef: {Object: objectValue},
					},
					new(pkg.FakeHandle),
				)

				/* assert */
				actualValue := fakeData.CoerceToStringArgsForCall(0)
				Expect(*actualValue).To(Equal(model.Value{Number: &pathSegment2Value}))
			})
		})
		Context("map[string]interface{} at path", func() {
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
				// err to trigger immediate return
				fakeData.CoerceToStringReturns(nil, errors.New("dummyError"))

				objectUnderTest := _deReferencer{
					data: fakeData,
				}

				/* act */
				objectUnderTest.DeReference(
					providedRef,
					map[string]*model.Value{
						objectRef: {Object: objectValue},
					},
					new(pkg.FakeHandle),
				)

				/* assert */
				actualValue := fakeData.CoerceToStringArgsForCall(0)
				Expect(*actualValue).To(Equal(model.Value{Object: pathSegment2Value}))
			})
		})
		Context("string at path", func() {
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

				objectUnderTest := _deReferencer{}

				/* act */
				actualString, actualOk, actualErr := objectUnderTest.DeReference(
					providedRef,
					map[string]*model.Value{
						objectRef: {Object: objectValue},
					},
					new(pkg.FakeHandle),
				)

				/* assert */
				Expect(actualString).To(Equal(pathSegment2Value))
				Expect(actualOk).To(Equal(true))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("[]interface{} at path", func() {
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
				// err to trigger immediate return
				fakeData.CoerceToStringReturns(nil, errors.New("dummyError"))

				objectUnderTest := _deReferencer{
					data: fakeData,
				}

				/* act */
				objectUnderTest.DeReference(
					providedRef,
					map[string]*model.Value{
						objectRef: {Object: objectValue},
					},
					new(pkg.FakeHandle),
				)

				/* assert */
				actualValue := fakeData.CoerceToStringArgsForCall(0)
				Expect(*actualValue).To(Equal(model.Value{Array: pathSegment2Value}))
			})
		})
	})
	Context("scope ref", func() {
		Context("ref isn't in scope", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"

				objectUnderTest := _deReferencer{}

				/* act */
				actualString, actualOk, actualErr := objectUnderTest.DeReference(
					providedRef,
					map[string]*model.Value{},
					new(pkg.FakeHandle),
				)

				/* assert */
				Expect(actualString).To(Equal(providedRef))
				Expect(actualOk).To(Equal(false))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("ref is in scope", func() {
			It("should call data.Coerce w/ expected args", func() {
				/* arrange */
				providedRef := "dummyRef"

				providedScopeValue := &model.Value{}

				fakeData := new(data.Fake)
				// err to trigger immediate return
				fakeData.CoerceToStringReturns(nil, errors.New("dummyError"))

				objectUnderTest := _deReferencer{
					data: fakeData,
				}

				/* act */
				objectUnderTest.DeReference(
					providedRef,
					map[string]*model.Value{
						providedRef: providedScopeValue,
					},
					new(pkg.FakeHandle),
				)

				/* assert */
				Expect(fakeData.CoerceToStringArgsForCall(0)).To(Equal(providedScopeValue))
			})
			Context("data.Coerce errs", func() {
				It("should return expected result", func() {
					/* arrange */
					providedRef := "dummyRef"

					fakeData := new(data.Fake)

					coerceError := errors.New("dummyError")
					fakeData.CoerceToStringReturns(nil, coerceError)

					objectUnderTest := _deReferencer{
						data: fakeData,
					}

					/* act */
					_, actualOk, actualErr := objectUnderTest.DeReference(
						providedRef,
						map[string]*model.Value{
							providedRef: nil,
						},
						new(pkg.FakeHandle),
					)

					/* assert */
					Expect(actualOk).To(Equal(false))
					Expect(actualErr).To(Equal(fmt.Errorf("unable to deReference '%v' as string; error was: %v", providedRef, coerceError.Error())))
				})
			})
			Context("data.Coerce doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedRef := "dummyRef"

					fakeData := new(data.Fake)

					coercedString := "dummyString"
					fakeData.CoerceToStringReturns(&model.Value{String: &coercedString}, nil)

					objectUnderTest := _deReferencer{
						data: fakeData,
					}

					/* act */
					actualString, actualOk, actualErr := objectUnderTest.DeReference(
						providedRef,
						map[string]*model.Value{
							providedRef: nil,
						},
						new(pkg.FakeHandle),
					)

					/* assert */
					Expect(actualString).To(Equal(coercedString))
					Expect(actualOk).To(Equal(true))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
})
