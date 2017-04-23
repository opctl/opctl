package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("Getter", func() {
	Context("newGetter()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(newGetter(nil, nil)).Should(Not(BeNil()))
		})
	})
	Context("Get", func() {
		It("should call resolver.Resolve w/ expected args", func() {
			/* arrange */
			providedBasePath := "dummyBasePath"
			providedPkgRef := "dummyPkgRef"

			resolvedPkgRef := "dummyPath"

			fakeResolver := new(fakeResolver)
			fakeResolver.ResolveReturns(resolvedPkgRef, true)

			objectUnderTest := _getter{
				manifestUnmarshaller: new(fakeManifestUnmarshaller),
				resolver:             fakeResolver,
			}

			/* act */
			objectUnderTest.Get(providedBasePath, providedPkgRef)

			/* assert */
			actualBasePath, actualPkgRef := fakeResolver.ResolveArgsForCall(0)

			Expect(actualBasePath).To(Equal(providedBasePath))
			Expect(actualPkgRef).To(Equal(providedPkgRef))
		})
		Context("resolver.Resolve returns true", func() {
			It("should call manifestUnmarshaller.Unmarshal w/ expected args", func() {
				/* arrange */
				providedBasePath := "dummyBasePath"
				providedPkgRef := "dummyPkgRef"

				resolvedPkgRef := "dummyPath"

				fakeResolver := new(fakeResolver)
				fakeResolver.ResolveReturns(resolvedPkgRef, true)

				fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)

				objectUnderTest := _getter{
					manifestUnmarshaller: fakeManifestUnmarshaller,
					resolver:             fakeResolver,
				}

				/* act */
				objectUnderTest.Get(providedBasePath, providedPkgRef)

				/* assert */
				Expect(fakeManifestUnmarshaller.UnmarshalArgsForCall(0)).To(Equal(resolvedPkgRef))
			})
			It("should return result of manifestUnmarshaller.Unmarshal", func() {
				/* arrange */
				providedBasePath := "dummyBasePath"
				providedPkgRef := "dummyPkgRef"

				resolvedPkgRef := "dummyPath"

				fakeResolver := new(fakeResolver)
				fakeResolver.ResolveReturns(resolvedPkgRef, true)

				expectedPkgManifest := &model.PkgManifest{Name: "dummyName"}
				expectedErr := errors.New("dummyError")

				fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)
				fakeManifestUnmarshaller.UnmarshalReturns(expectedPkgManifest, expectedErr)

				objectUnderTest := _getter{
					manifestUnmarshaller: fakeManifestUnmarshaller,
					resolver:             fakeResolver,
				}

				/* act */
				actualPkgManifest, actualErr := objectUnderTest.Get(providedBasePath, providedPkgRef)

				/* assert */
				Expect(actualPkgManifest).To(Equal(expectedPkgManifest))
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("resolver.Resolve returns false", func() {
			It("should return ErrPkgNotFound ", func() {
				/* arrange */
				expectedErr := ErrPkgNotFound{}

				fakeResolver := new(fakeResolver)
				fakeResolver.ResolveReturns("", false)

				objectUnderTest := _getter{
					manifestUnmarshaller: new(fakeManifestUnmarshaller),
					resolver:             fakeResolver,
				}

				/* act */
				actualPkgManifest, actualErr := objectUnderTest.Get("", "")

				/* assert */
				Expect(actualPkgManifest).To(BeNil())
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
	})
})
