package string

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/coerce"
	"github.com/pkg/errors"
)

var _ = Context("deReferencer", func() {
	Context("ref isn't in scope", func() {
		It("should return expected result", func() {
			/* arrange */
			providedRef := "dummyRef"

			objectUnderTest := _deReferencer{
				scope: map[string]*model.Value{},
			}

			/* act */
			actualString, actualOk, actualErr := objectUnderTest.DeReference(
				providedRef,
			)

			/* assert */
			Expect(actualString).To(Equal(providedRef))
			Expect(actualOk).To(Equal(false))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("ref is in scope", func() {
		It("should call coerce.Coerce w/ expected args", func() {
			/* arrange */
			providedRef := "dummyRef"

			providedScopeValue := &model.Value{}
			providedScope := map[string]*model.Value{
				providedRef: providedScopeValue,
			}

			fakeCoercer := new(coerce.Fake)
			// err to trigger immediate return
			fakeCoercer.ToStringReturns("", errors.New("dummyError"))

			objectUnderTest := _deReferencer{
				coerce: fakeCoercer,
				scope:  providedScope,
			}

			/* act */
			objectUnderTest.DeReference(providedRef)

			/* assert */
			Expect(fakeCoercer.ToStringArgsForCall(0)).To(Equal(providedScopeValue))
		})
		Context("coerce.Coerce errs", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"

				fakeCoercer := new(coerce.Fake)

				coerceError := errors.New("dummyError")
				fakeCoercer.ToStringReturns("", coerceError)

				objectUnderTest := _deReferencer{
					coerce: fakeCoercer,
					scope: map[string]*model.Value{
						providedRef: nil,
					},
				}

				/* act */
				_, actualOk, actualErr := objectUnderTest.DeReference(
					providedRef,
				)

				/* assert */
				Expect(actualOk).To(Equal(false))
				Expect(actualErr).To(Equal(fmt.Errorf("Unable to deReference '%v' as string; error was: %v", providedRef, coerceError.Error())))
			})
		})
		Context("coerce.Coerce doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"

				fakeCoercer := new(coerce.Fake)

				coercedString := "dummyString"
				fakeCoercer.ToStringReturns(coercedString, nil)

				objectUnderTest := _deReferencer{
					coerce: fakeCoercer,
					scope: map[string]*model.Value{
						providedRef: nil,
					},
				}

				/* act */
				actualString, actualOk, actualErr := objectUnderTest.DeReference(
					providedRef,
				)

				/* assert */
				Expect(actualString).To(Equal(coercedString))
				Expect(actualOk).To(Equal(true))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
