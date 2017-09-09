package interpolater

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/pkg/errors"
)

var _ = Context("deReferencer", func() {
	Context("ref isn't in scope", func() {
		It("should return expected result", func() {
			/* arrange */
			providedRef := "dummyRef"

			objectUnderTest := _deReferencer{}

			/* act */
			actualString, actualOk, actualErr := objectUnderTest.DeReference(
				providedRef,
				map[string]*model.Value{},
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
			fakeData.CoerceToStringReturns("", errors.New("dummyError"))

			objectUnderTest := _deReferencer{
				data: fakeData,
			}

			/* act */
			objectUnderTest.DeReference(
				providedRef,
				map[string]*model.Value{
					providedRef: providedScopeValue,
				},
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
				fakeData.CoerceToStringReturns("", coerceError)

				objectUnderTest := _deReferencer{
					data: fakeData,
				}

				/* act */
				_, actualOk, actualErr := objectUnderTest.DeReference(
					providedRef,
					map[string]*model.Value{
						providedRef: nil,
					},
				)

				/* assert */
				Expect(actualOk).To(Equal(false))
				Expect(actualErr).To(Equal(fmt.Errorf("Unable to deReference '%v' as string; error was: %v", providedRef, coerceError.Error())))
			})
		})
		Context("data.Coerce doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"

				fakeData := new(data.Fake)

				coercedString := "dummyString"
				fakeData.CoerceToStringReturns(coercedString, nil)

				objectUnderTest := _deReferencer{
					data: fakeData,
				}

				/* act */
				actualString, actualOk, actualErr := objectUnderTest.DeReference(
					providedRef,
					map[string]*model.Value{
						providedRef: nil,
					},
				)

				/* assert */
				Expect(actualString).To(Equal(coercedString))
				Expect(actualOk).To(Equal(true))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
