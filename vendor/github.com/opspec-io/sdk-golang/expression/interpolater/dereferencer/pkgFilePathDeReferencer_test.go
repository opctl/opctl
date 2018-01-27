package dereferencer

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/pkg/errors"
)

var _ = Context("pkgFilePathDeReferencer", func() {
	Context("ref is pkg file path ref", func() {
		It("should call pkgHandle.GetContent w/ expected args", func() {
			/* arrange */
			providedRef := "/dummyRef"
			fakePkgHandle := new(pkg.FakeHandle)
			// err to trigger immediate return
			fakePkgHandle.GetContentReturns(nil, errors.New("dummyError"))

			objectUnderTest := _pkgFilePathDeReferencer{}

			/* act */
			objectUnderTest.DeReferencePkgFilePath(
				providedRef,
				map[string]*model.Value{},
				fakePkgHandle,
			)

			/* assert */
			actualContext,
				actualContentPath := fakePkgHandle.GetContentArgsForCall(0)

			Expect(actualContext).To(Equal(context.TODO()))
			Expect(actualContentPath).To(Equal(providedRef))
		})
	})
})
