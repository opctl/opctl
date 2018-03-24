package dereferencer

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("pkgFilePathDeReferencer", func() {
	Context("ref is pkg file path ref", func() {
		It("should call opHandle.GetContent w/ expected args", func() {
			/* arrange */
			providedRef := "/dummyRef"
			fakeDataHandle := new(data.FakeHandle)
			// err to trigger immediate return
			fakeDataHandle.GetContentReturns(nil, errors.New("dummyError"))

			objectUnderTest := _pkgFilePathDeReferencer{}

			/* act */
			objectUnderTest.DeReferencePkgFilePath(
				providedRef,
				map[string]*model.Value{},
				fakeDataHandle,
			)

			/* assert */
			actualContext,
				actualContentPath := fakeDataHandle.GetContentArgsForCall(0)

			Expect(actualContext).To(Equal(context.TODO()))
			Expect(actualContentPath).To(Equal(providedRef))
		})
	})
})
