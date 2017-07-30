package pkg

import (
	"context"
	"errors"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
)

var _ = Context("Validate", func() {
	It("should call handle.GetContent w/ expected args", func() {
		/* arrange */
		providedFileHandle := new(FakeHandle)
		// error to trigger immediate return
		providedFileHandle.GetContentReturns(nil, errors.New("dummyError"))

		objectUnderTest := _Pkg{}

		/* act */
		objectUnderTest.Validate(providedFileHandle)

		/* assert */
		actualCtx,
			actualContentName := providedFileHandle.GetContentArgsForCall(0)

		Expect(actualCtx).To(Equal(context.TODO()))
		Expect(actualContentName).To(Equal(OpDotYmlFileName))
	})
	Context("handle.GetContent errs", func() {
		It("should return err", func() {
			/* arrange */
			expectedErrors := []error{errors.New("dummyError")}
			providedFileHandle := new(FakeHandle)
			// error to trigger immediate return
			providedFileHandle.GetContentReturns(nil, expectedErrors[0])

			objectUnderTest := _Pkg{}

			/* act */
			actualErrors := objectUnderTest.Validate(providedFileHandle)

			/* assert */
			Expect(actualErrors).To(Equal(expectedErrors))
		})
	})
	Context("handle.GetContent doesn't err", func() {
		It("should call manifestValidator.Validate w/ expected args & return result", func() {
			/* arrange */

			providedFileHandle := new(FakeHandle)

			expectedManifestBytes := []byte{2, 5, 61}
			fakeIOUtil := new(iioutil.Fake)
			fakeIOUtil.ReadAllReturns(expectedManifestBytes, nil)

			expectedErrs := []error{
				errors.New("dummyErr1"),
				errors.New("dummyErr2"),
			}
			fakeManifest := new(manifest.Fake)

			fakeManifest.ValidateReturns(expectedErrs)

			objectUnderTest := _Pkg{
				ioUtil:   fakeIOUtil,
				manifest: fakeManifest,
			}

			/* act */
			actualErrs := objectUnderTest.Validate(providedFileHandle)

			/* assert */
			Expect(fakeManifest.ValidateArgsForCall(0)).To(Equal(expectedManifestBytes))
			Expect(actualErrs).To(Equal(expectedErrs))
		})
	})
})
