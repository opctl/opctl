package pkg

import (
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Context("pkg", func() {

	Context("Create", func() {

		It("should call os.MkdirAll with expected args", func() {
			/* arrange */
			providedPkgPath := "dummyPkgPath"
			providedContentPath := "dummyContentPath"

			fakeOS := new(ios.Fake)

			objectUnderTest := &_Pkg{
				os: fakeOS,
			}

			/* act */
			objectUnderTest.GetContent(providedPkgPath, providedContentPath)

			/* assert */
			Expect(fakeOS.OpenArgsForCall(0)).To(Equal(filepath.Join(providedPkgPath, providedContentPath)))
		})
	})
})
