package dereferencer

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/pkg/errors"
)

var _ = Context("deReferencer", func() {
	It("should call pkgFilePathDeReferencer.DeReferencePkgFilePath w/ expected args", func() {
		/* arrange */
		providedRef := "dummyRef"
		providedScope := map[string]*model.Value{"dummyName": {}}
		providedOpHandle := new(data.FakeHandle)

		fakePkgFilePathDeReferencer := new(fakePkgFilePathDeReferencer)
		// err to trigger immediate return
		fakePkgFilePathDeReferencer.DeReferencePkgFilePathReturns("", true, errors.New("dummyError"))

		objectUnderTest := _deReferencer{
			pkgFilePathDeReferencer: fakePkgFilePathDeReferencer,
		}

		/* act */
		objectUnderTest.DeReference(
			providedRef,
			providedScope,
			providedOpHandle,
		)

		/* assert */
		actualRef,
			actualScope,
			actualOpHandle := fakePkgFilePathDeReferencer.DeReferencePkgFilePathArgsForCall(0)

		Expect(actualRef).To(Equal(providedRef))
		Expect(actualScope).To(Equal(providedScope))
		Expect(actualOpHandle).To(Equal(providedOpHandle))
	})
	Context("ref is pkgFilePathRef", func() {
		It("should return expected result", func() {
			/* arrange */
			fakePkgFilePathDeReferencer := new(fakePkgFilePathDeReferencer)
			deReferencedValue := "dummyRef"
			err := errors.New("dummyError")
			fakePkgFilePathDeReferencer.DeReferencePkgFilePathReturns(deReferencedValue, true, err)

			objectUnderTest := _deReferencer{
				pkgFilePathDeReferencer: fakePkgFilePathDeReferencer,
			}

			/* act */
			actualValue,
				actualDidDeReference,
				actualErr := objectUnderTest.DeReference(
				"dummyRef",
				map[string]*model.Value{},
				new(data.FakeHandle),
			)

			/* assert */
			Expect(actualValue).To(Equal(deReferencedValue))
			Expect(actualDidDeReference).To(Equal(true))
			Expect(actualErr).To(Equal(err))
		})
	})
	Context("ref isn't pkgFilePathRef", func() {
		It("should call scopeDeReferencer.DeReferenceScope w/ expected args", func() {
			/* arrange */
			providedRef := "dummyRef"
			providedScope := map[string]*model.Value{"dummyName": {}}

			fakeScopeDeReferencer := new(fakeScopeDeReferencer)
			// err to trigger immediate return
			fakeScopeDeReferencer.DeReferenceScopeReturns("", true, errors.New("dummyError"))

			objectUnderTest := _deReferencer{
				pkgFilePathDeReferencer: new(fakePkgFilePathDeReferencer),
				scopeDeReferencer:       fakeScopeDeReferencer,
			}

			/* act */
			objectUnderTest.DeReference(
				providedRef,
				providedScope,
				new(data.FakeHandle),
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
					pkgFilePathDeReferencer: new(fakePkgFilePathDeReferencer),
					scopeDeReferencer:       fakeScopeDeReferencer,
				}

				/* act */
				actualValue,
					actualDidDeReference,
					actualErr := objectUnderTest.DeReference(
					"dummyRef",
					map[string]*model.Value{},
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualValue).To(Equal(deReferencedValue))
				Expect(actualDidDeReference).To(Equal(true))
				Expect(actualErr).To(Equal(err))
			})
		})
		Context("ref isn't scopeRef", func() {
			It("should call scopeFilePathDeReferencer.DeReferenceScopeFilePath w/ expected args", func() {
				/* arrange */
				providedRef := "dummyRef"
				providedScope := map[string]*model.Value{"dummyName": {}}

				fakeScopeFilePathDeReferencer := new(fakeScopeFilePathDeReferencer)
				// err to trigger immediate return
				fakeScopeFilePathDeReferencer.DeReferenceScopeFilePathReturns("", true, errors.New("dummyError"))

				objectUnderTest := _deReferencer{
					pkgFilePathDeReferencer:   new(fakePkgFilePathDeReferencer),
					scopeDeReferencer:         new(fakeScopeDeReferencer),
					scopeFilePathDeReferencer: fakeScopeFilePathDeReferencer,
				}

				/* act */
				objectUnderTest.DeReference(
					providedRef,
					providedScope,
					new(data.FakeHandle),
				)

				/* assert */
				actualRef,
					actualScope := fakeScopeFilePathDeReferencer.DeReferenceScopeFilePathArgsForCall(0)

				Expect(actualRef).To(Equal(providedRef))
				Expect(actualScope).To(Equal(providedScope))
			})
			Context("ref is scopeFilePathRef", func() {

				It("should return expected result", func() {
					/* arrange */
					fakeScopeFilePathDeReferencer := new(fakeScopeFilePathDeReferencer)
					deReferencedValue := "dummyRef"
					err := errors.New("dummyError")
					fakeScopeFilePathDeReferencer.DeReferenceScopeFilePathReturns(deReferencedValue, true, err)

					objectUnderTest := _deReferencer{
						pkgFilePathDeReferencer:   new(fakePkgFilePathDeReferencer),
						scopeDeReferencer:         new(fakeScopeDeReferencer),
						scopeFilePathDeReferencer: fakeScopeFilePathDeReferencer,
					}

					/* act */
					actualValue,
						actualDidDeReference,
						actualErr := objectUnderTest.DeReference(
						"dummyRef",
						map[string]*model.Value{},
						new(data.FakeHandle),
					)

					/* assert */
					Expect(actualValue).To(Equal(deReferencedValue))
					Expect(actualDidDeReference).To(Equal(true))
					Expect(actualErr).To(Equal(err))
				})
			})
			Context("ref isn't scopeFilePathRef", func() {
				It("should call scopeObjectPathDeReferencer.DeReferenceScopeObjectPath w/ expected args", func() {
					/* arrange */
					providedRef := "dummyRef"
					providedScope := map[string]*model.Value{"dummyName": {}}

					fakeScopeObjectPathDeReferencer := new(fakeScopeObjectPathDeReferencer)
					// err to trigger immediate return
					fakeScopeObjectPathDeReferencer.DeReferenceScopeObjectPathReturns("", true, errors.New("dummyError"))

					objectUnderTest := _deReferencer{
						pkgFilePathDeReferencer:     new(fakePkgFilePathDeReferencer),
						scopeDeReferencer:           new(fakeScopeDeReferencer),
						scopeFilePathDeReferencer:   new(fakeScopeFilePathDeReferencer),
						scopeObjectPathDeReferencer: fakeScopeObjectPathDeReferencer,
					}

					/* act */
					objectUnderTest.DeReference(
						providedRef,
						providedScope,
						new(data.FakeHandle),
					)

					/* assert */
					actualRef,
						actualScope := fakeScopeObjectPathDeReferencer.DeReferenceScopeObjectPathArgsForCall(0)

					Expect(actualRef).To(Equal(providedRef))
					Expect(actualScope).To(Equal(providedScope))
				})
				Context("ref is scopeObjectPathRef", func() {
					It("should return expected result", func() {
						/* arrange */
						fakeScopeObjectPathDeReferencer := new(fakeScopeObjectPathDeReferencer)
						deReferencedValue := "dummyRef"
						err := errors.New("dummyError")
						fakeScopeObjectPathDeReferencer.DeReferenceScopeObjectPathReturns(deReferencedValue, true, err)

						objectUnderTest := _deReferencer{
							pkgFilePathDeReferencer:     new(fakePkgFilePathDeReferencer),
							scopeDeReferencer:           new(fakeScopeDeReferencer),
							scopeFilePathDeReferencer:   new(fakeScopeFilePathDeReferencer),
							scopeObjectPathDeReferencer: fakeScopeObjectPathDeReferencer,
						}

						/* act */
						actualValue,
							actualDidDeReference,
							actualErr := objectUnderTest.DeReference(
							"dummyRef",
							map[string]*model.Value{},
							new(data.FakeHandle),
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
