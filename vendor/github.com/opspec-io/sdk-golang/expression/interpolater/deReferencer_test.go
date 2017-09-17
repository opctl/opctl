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
