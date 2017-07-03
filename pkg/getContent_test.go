package pkg

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"io/ioutil"
)

var _ = Context("pkg", func() {

	Context("GetContent", func() {

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
			objectUnderTest.GetContent(providedPkgRef, "dummyContentPath")

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
				_, actualErr := objectUnderTest.GetContent("pkgPath", "contentPath")

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("opener.Open doesn't err", func() {
			It("should call handle.GetContent w/ expected args & return result", func() {
				providedContentPath := "dummyContentPath"

				expectedReadSeekCloser, err := ioutil.TempFile("", "")
				if nil != err {
					panic(err)
				}
				defer expectedReadSeekCloser.Close()

				fakeHandle := new(fakeHandle)
				fakeHandle.GetContentReturns(expectedReadSeekCloser, nil)

				fakeOpener := new(fakeOpener)
				fakeOpener.OpenReturns(fakeHandle, nil)

				objectUnderTest := &_Pkg{
					opener: fakeOpener,
				}

				/* act */
				actualReadSeekCloser, _ := objectUnderTest.GetContent("pkgPath", providedContentPath)

				/* assert */
				Expect(fakeHandle.GetContentArgsForCall(0)).To(Equal(providedContentPath))
				Expect(actualReadSeekCloser).To(Equal(expectedReadSeekCloser))
			})
		})
	})
})
