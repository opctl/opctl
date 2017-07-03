package pkg

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/pkg/errors"
)

var _ = Context("pkg", func() {

	Context("ListContents", func() {

		It("should call opener.Open w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"

			fakeOpener := new(fakeOpener)
			// err to trigger immediate return
			fakeOpener.OpenReturns(nil, errors.New("dummyError"))

			objectUnderTest := &_Pkg{
				opener: fakeOpener,
			}

			/* act */
			objectUnderTest.ListContents(providedPkgRef)

			/* assert */
			Expect(fakeOpener.OpenArgsForCall(0)).To(Equal(providedPkgRef))
		})
		Context("opener.Open errs", func() {
			It("should return expected error", func() {
				expectedErr := errors.New("dummyError")

				fakeOpener := new(fakeOpener)
				fakeOpener.OpenReturns(nil, expectedErr)

				objectUnderTest := &_Pkg{
					opener: fakeOpener,
				}

				/* act */
				_, actualErr := objectUnderTest.ListContents("pkgPath")

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("opener.Open doesn't err", func() {
			It("should call handle.ListContents w/ expected args & return result", func() {

				expectedPkgContents := []*model.PkgContent{
					{
						Path: "dummyPath1",
					},
					{
						Path: "dummyPath2",
					},
				}

				fakeHandle := new(fakeHandle)
				fakeHandle.ListContentsReturns(expectedPkgContents, nil)

				fakeOpener := new(fakeOpener)
				fakeOpener.OpenReturns(fakeHandle, nil)

				objectUnderTest := &_Pkg{
					opener: fakeOpener,
				}

				/* act */
				actualReadSeekCloser, _ := objectUnderTest.ListContents("pkgPath")

				/* assert */
				Expect(actualReadSeekCloser).To(Equal(expectedPkgContents))
			})
		})
	})
})
