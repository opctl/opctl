package data

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data/provider"
	. "github.com/opctl/opctl/sdks/go/data/provider/fakes"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
)

var _ = Context("resolver", func() {
	Context("Resolve", func() {
		It("should call providers[0].TryResolve w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedDataRef := "dummyDataRef"
			fakeProvider0 := new(FakeProvider)
			providedProviders := []provider.Provider{fakeProvider0}

			objectUnderTest := _resolver{}

			/* act */
			objectUnderTest.Resolve(
				providedCtx,
				providedDataRef,
				providedProviders...,
			)

			/* assert */
			actualCtx,
				actualDataRef := fakeProvider0.TryResolveArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualDataRef).To(Equal(providedDataRef))
		})
		Context("providers[0].TryResolve errs", func() {
			It("should return error", func() {
				/* arrange */
				fakeProvider0 := new(FakeProvider)
				expectedErr := errors.New("dummyError")
				fakeProvider0.TryResolveReturns(nil, expectedErr)

				providedProviders := []provider.Provider{fakeProvider0}

				objectUnderTest := _resolver{}

				/* act */
				_, actualErr := objectUnderTest.Resolve(
					context.Background(),
					"dummyDataRef",
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
				expectedHandle := new(modelFakes.FakeDataHandle)
				fakeProvider0.TryResolveReturnsOnCall(0, expectedHandle, nil)

				providedProviders := []provider.Provider{fakeProvider0}

				objectUnderTest := _resolver{}

				/* act */
				actualHandle, actualErr := objectUnderTest.Resolve(
					context.Background(),
					"dumyDataRef",
					providedProviders...,
				)

				/* assert */
				Expect(actualHandle).To(Equal(expectedHandle))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
