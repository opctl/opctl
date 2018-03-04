package pkg

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Context("resolver", func() {
	Context("Resolve", func() {
		It("should call providers[0].TryResolve w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedPkgRef := "dummyPkgRef"
			fakeProvider0 := new(FakeProvider)
			providedProviders := []Provider{fakeProvider0}

			objectUnderTest := _resolver{}

			/* act */
			objectUnderTest.Resolve(
				providedCtx,
				providedPkgRef,
				providedProviders...,
			)

			/* assert */
			actualCtx,
				actualPkgRef := fakeProvider0.TryResolveArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualPkgRef).To(Equal(providedPkgRef))
		})
		Context("providers[0].TryResolve errs", func() {
			It("should return error", func() {
				/* arrange */
				fakeProvider0 := new(FakeProvider)
				expectedErr := errors.New("dummyError")
				fakeProvider0.TryResolveReturns(nil, expectedErr)

				providedProviders := []Provider{fakeProvider0}

				objectUnderTest := _resolver{}

				/* act */
				_, actualErr := objectUnderTest.Resolve(
					context.Background(),
					"dummyPkgRef",
					providedProviders...,
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("providers[0].TryResolve doesn't err", func() {
			It("should return expected results", func() {
				/* arrange */
				fakeProvider0 := new(FakeProvider)
				expectedHandle := new(FakeHandle)
				fakeProvider0.TryResolveReturnsOnCall(0, expectedHandle, nil)

				providedProviders := []Provider{fakeProvider0}

				objectUnderTest := _resolver{}

				/* act */
				actualHandle, actualErr := objectUnderTest.Resolve(
					context.Background(),
					"dumyPkgRef",
					providedProviders...,
				)

				/* assert */
				Expect(actualHandle).To(Equal(expectedHandle))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
