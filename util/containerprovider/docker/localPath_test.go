package docker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/vruntime"
)

var _ = Context("localPath", func() {
	Context("when runtime.GOOS == windows", func() {

		fakeRuntime := new(vruntime.Fake)
		fakeRuntime.GOOSReturns("windows")
		objectUnderTest := _containerProvider{
			runtime: fakeRuntime,
		}

		Context("when path contains drive letter", func() {
			It("should trim preceeding slash and add colon proceeding drive letter", func() {
				/* arrange */
				expected := "c:/DummyPath"

				/* act */
				actual := objectUnderTest.localPath("/c/DummyPath")

				/* assert */
				Expect(actual).To(Equal(expected))
			})
		})
	})
})
