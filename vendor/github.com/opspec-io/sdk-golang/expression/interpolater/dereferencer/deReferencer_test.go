package dereferencer

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/pkg/errors"
)

var _ = Context("deReferencer", func() {
	It("should call pkgFileDeReferencer.DeReferencePkgFile w/ expected args", func() {
		/* arrange */
		providedRef := "dummyRef"
		providedScope := map[string]*model.Value{"dummyName": {}}
		providedPkgHandle := new(pkg.FakeHandle)

		fakePkgFileDeReferencer := new(fakePkgFileDeReferencer)
		// err to trigger immediate return
		fakePkgFileDeReferencer.DeReferencePkgFileReturns("", true, errors.New("dummyError"))

		objectUnderTest := _deReferencer{
			pkgFileDeReferencer: fakePkgFileDeReferencer,
		}

		/* act */
		objectUnderTest.DeReference(
			providedRef,
			providedScope,
			providedPkgHandle,
		)

		/* assert */
		actualRef,
			actualScope,
			actualPkgHandle := fakePkgFileDeReferencer.DeReferencePkgFileArgsForCall(0)

		Expect(actualRef).To(Equal(providedRef))
		Expect(actualScope).To(Equal(providedScope))
		Expect(actualPkgHandle).To(Equal(providedPkgHandle))
	})
	Context("ref is pkgFileRef", func() {
		It("should return expected result", func() {
			/* arrange */
			fakePkgFileDeReferencer := new(fakePkgFileDeReferencer)
			deReferencedValue := "dummyRef"
			err := errors.New("dummyError")
			fakePkgFileDeReferencer.DeReferencePkgFileReturns(deReferencedValue, true, err)

			objectUnderTest := _deReferencer{
				pkgFileDeReferencer: fakePkgFileDeReferencer,
			}

			/* act */
			actualValue,
				actualDidDeReference,
				actualErr := objectUnderTest.DeReference(
				"dummyRef",
				map[string]*model.Value{},
				new(pkg.FakeHandle),
			)

			/* assert */
			Expect(actualValue).To(Equal(deReferencedValue))
			Expect(actualDidDeReference).To(Equal(true))
			Expect(actualErr).To(Equal(err))
		})
	})
	Context("ref isn't pkgFileRef", func() {
		It("should call scopeDeReferencer.DeReferenceScope w/ expected args", func() {
			/* arrange */
			providedRef := "dummyRef"
			providedScope := map[string]*model.Value{"dummyName": {}}

			fakeScopeDeReferencer := new(fakeScopeDeReferencer)
			// err to trigger immediate return
			fakeScopeDeReferencer.DeReferenceScopeReturns("", true, errors.New("dummyError"))

			objectUnderTest := _deReferencer{
				pkgFileDeReferencer: new(fakePkgFileDeReferencer),
				scopeDeReferencer:   fakeScopeDeReferencer,
			}

			/* act */
			objectUnderTest.DeReference(
				providedRef,
				providedScope,
				new(pkg.FakeHandle),
			)

			/* assert */
			actualRef,
				actualScope := fakeScopeDeReferencer.DeReferenceScopeArgsForCall(0)

			Expect(actualRef).To(Equal(providedRef))
			Expect(actualScope).To(Equal(providedScope))
		})
		Context("ref is scopeRef", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeScopeDeReferencer := new(fakeScopeDeReferencer)
				deReferencedValue := "dummyRef"
				err := errors.New("dummyError")
				fakeScopeDeReferencer.DeReferenceScopeReturns(deReferencedValue, true, err)

				objectUnderTest := _deReferencer{
					pkgFileDeReferencer: new(fakePkgFileDeReferencer),
					scopeDeReferencer:   fakeScopeDeReferencer,
				}

				/* act */
				actualValue,
					actualDidDeReference,
					actualErr := objectUnderTest.DeReference(
					"dummyRef",
					map[string]*model.Value{},
					new(pkg.FakeHandle),
				)

				/* assert */
				Expect(actualValue).To(Equal(deReferencedValue))
				Expect(actualDidDeReference).To(Equal(true))
				Expect(actualErr).To(Equal(err))
			})
		})
		Context("ref isn't scopeRef", func() {
			It("should call scopeFileDeReferencer.DeReferenceScopeFile w/ expected args", func() {
				/* arrange */
				providedRef := "dummyRef"
				providedScope := map[string]*model.Value{"dummyName": {}}

				fakeScopeFileDeReferencer := new(fakeScopeFileDeReferencer)
				// err to trigger immediate return
				fakeScopeFileDeReferencer.DeReferenceScopeFileReturns("", true, errors.New("dummyError"))

				objectUnderTest := _deReferencer{
					pkgFileDeReferencer:   new(fakePkgFileDeReferencer),
					scopeDeReferencer:     new(fakeScopeDeReferencer),
					scopeFileDeReferencer: fakeScopeFileDeReferencer,
				}

				/* act */
				objectUnderTest.DeReference(
					providedRef,
					providedScope,
					new(pkg.FakeHandle),
				)

				/* assert */
				actualRef,
					actualScope := fakeScopeFileDeReferencer.DeReferenceScopeFileArgsForCall(0)

				Expect(actualRef).To(Equal(providedRef))
				Expect(actualScope).To(Equal(providedScope))
			})
			Context("ref is scopeFileRef", func() {

				It("should return expected result", func() {
					/* arrange */
					fakeScopeFileDeReferencer := new(fakeScopeFileDeReferencer)
					deReferencedValue := "dummyRef"
					err := errors.New("dummyError")
					fakeScopeFileDeReferencer.DeReferenceScopeFileReturns(deReferencedValue, true, err)

					objectUnderTest := _deReferencer{
						pkgFileDeReferencer:   new(fakePkgFileDeReferencer),
						scopeDeReferencer:     new(fakeScopeDeReferencer),
						scopeFileDeReferencer: fakeScopeFileDeReferencer,
					}

					/* act */
					actualValue,
						actualDidDeReference,
						actualErr := objectUnderTest.DeReference(
						"dummyRef",
						map[string]*model.Value{},
						new(pkg.FakeHandle),
					)

					/* assert */
					Expect(actualValue).To(Equal(deReferencedValue))
					Expect(actualDidDeReference).To(Equal(true))
					Expect(actualErr).To(Equal(err))
				})
			})
			Context("ref isn't scopeFileRef", func() {
				It("should call scopePropertyDeReferencer.DeReferenceScopeProperty w/ expected args", func() {
					/* arrange */
					providedRef := "dummyRef"
					providedScope := map[string]*model.Value{"dummyName": {}}

					fakeScopePropertyDeReferencer := new(fakeScopePropertyDeReferencer)
					// err to trigger immediate return
					fakeScopePropertyDeReferencer.DeReferenceScopePropertyReturns("", true, errors.New("dummyError"))

					objectUnderTest := _deReferencer{
						pkgFileDeReferencer:       new(fakePkgFileDeReferencer),
						scopeDeReferencer:         new(fakeScopeDeReferencer),
						scopeFileDeReferencer:     new(fakeScopeFileDeReferencer),
						scopePropertyDeReferencer: fakeScopePropertyDeReferencer,
					}

					/* act */
					objectUnderTest.DeReference(
						providedRef,
						providedScope,
						new(pkg.FakeHandle),
					)

					/* assert */
					actualRef,
						actualScope := fakeScopePropertyDeReferencer.DeReferenceScopePropertyArgsForCall(0)

					Expect(actualRef).To(Equal(providedRef))
					Expect(actualScope).To(Equal(providedScope))
				})
				Context("ref is scopePropertyRef", func() {
					It("should return expected result", func() {
						/* arrange */
						fakeScopePropertyDeReferencer := new(fakeScopePropertyDeReferencer)
						deReferencedValue := "dummyRef"
						err := errors.New("dummyError")
						fakeScopePropertyDeReferencer.DeReferenceScopePropertyReturns(deReferencedValue, true, err)

						objectUnderTest := _deReferencer{
							pkgFileDeReferencer:       new(fakePkgFileDeReferencer),
							scopeDeReferencer:         new(fakeScopeDeReferencer),
							scopeFileDeReferencer:     new(fakeScopeFileDeReferencer),
							scopePropertyDeReferencer: fakeScopePropertyDeReferencer,
						}

						/* act */
						actualValue,
							actualDidDeReference,
							actualErr := objectUnderTest.DeReference(
							"dummyRef",
							map[string]*model.Value{},
							new(pkg.FakeHandle),
						)

						/* assert */
						Expect(actualValue).To(Equal(deReferencedValue))
						Expect(actualDidDeReference).To(Equal(true))
						Expect(actualErr).To(Equal(err))
					})
				})
			})
		})
	})
})
