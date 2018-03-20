package dereferencer

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/pkg/errors"
)

var _ = Context("scopeDeReferencer", func() {
	Context("ref is scope ref", func() {
		Context("ref isn't in scope", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"

				objectUnderTest := _scopeDeReferencer{}

				/* act */
				actualString, actualOk, actualErr := objectUnderTest.DeReferenceScope(
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
			It("should call coerce.Coerce w/ expected args", func() {
				/* arrange */
				providedRef := "dummyRef"

				providedScopeValue := &model.Value{}

				fakeCoerce := new(coerce.Fake)
				// err to trigger immediate return
				fakeCoerce.ToStringReturns(nil, errors.New("dummyError"))

				objectUnderTest := _scopeDeReferencer{
					coerce: fakeCoerce,
				}

				/* act */
				objectUnderTest.DeReferenceScope(
					providedRef,
					map[string]*model.Value{
						providedRef: providedScopeValue,
					},
				)

				/* assert */
				Expect(fakeCoerce.ToStringArgsForCall(0)).To(Equal(providedScopeValue))
			})
			Context("coerce.Coerce errs", func() {
				It("should return expected result", func() {
					/* arrange */
					providedRef := "dummyRef"

					fakeCoerce := new(coerce.Fake)

					coerceError := errors.New("dummyError")
					fakeCoerce.ToStringReturns(nil, coerceError)

					objectUnderTest := _scopeDeReferencer{
						coerce: fakeCoerce,
					}

					/* act */
					_, actualOk, actualErr := objectUnderTest.DeReferenceScope(
						providedRef,
						map[string]*model.Value{
							providedRef: nil,
						},
					)

					/* assert */
					Expect(actualOk).To(Equal(false))
					Expect(actualErr).To(Equal(fmt.Errorf("unable to deReference '%v' as string; error was: %v", providedRef, coerceError.Error())))
				})
			})
			Context("coerce.Coerce doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedRef := "dummyRef"

					fakeCoerce := new(coerce.Fake)

					coercedString := "dummyString"
					fakeCoerce.ToStringReturns(&model.Value{String: &coercedString}, nil)

					objectUnderTest := _scopeDeReferencer{
						coerce: fakeCoerce,
					}

					/* act */
					actualString, actualOk, actualErr := objectUnderTest.DeReferenceScope(
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
})
