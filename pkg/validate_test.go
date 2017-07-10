package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
	"io/ioutil"
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
		Expect(providedFileHandle.GetContentArgsForCall(0)).To(Equal(OpDotYmlFileName))
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
			expectedManifestReader, err := ioutil.TempFile("", "")
			if nil != err {
				Fail(err.Error())
			}

			providedFileHandle := new(FakeHandle)
			providedFileHandle.GetContentReturns(expectedManifestReader, nil)

			expectedErrs := []error{
				errors.New("dummyErr1"),
				errors.New("dummyErr2"),
			}
			fakeManifest := new(manifest.Fake)

			fakeManifest.ValidateReturns(expectedErrs)

			objectUnderTest := _Pkg{
				manifest: fakeManifest,
			}

			/* act */
			actualErrs := objectUnderTest.Validate(providedFileHandle)

			/* assert */
			Expect(fakeManifest.ValidateArgsForCall(0)).To(Equal(expectedManifestReader))
			Expect(actualErrs).To(Equal(expectedErrs))
		})
	})
})
