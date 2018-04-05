package dereferencer

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("opFilePathDeReferencer", func() {
	Context("ref is op file path ref", func() {
		It("should call opHandle.GetContent w/ expected args", func() {
			/* arrange */
			providedRef := "/dummyRef"
			fakeDataHandle := new(data.FakeHandle)
			// err to trigger immediate return
			fakeDataHandle.GetContentReturns(nil, errors.New("dummyError"))

			objectUnderTest := _opFilePathDeReferencer{}

			/* act */
			objectUnderTest.DeReferenceOpFilePath(
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
