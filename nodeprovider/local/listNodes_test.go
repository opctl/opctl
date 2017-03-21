package local

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/node"
	"github.com/opctl/opctl/util/lockfile"
)

var _ = Context("listNodes", func() {
	It("should call lockfile.PIdOfOwner w/ expected args", func() {
		/* arrange */
		expectedFilePath := lockFilePath()

		fakeLockFile := new(lockfile.Fake)

		objectUnderTest := nodeProvider{
			lockfile: fakeLockFile,
		}

		/* act */
		objectUnderTest.ListNodes()

		/* assert */
		Expect(fakeLockFile.PIdOfOwnerArgsForCall(0)).To(Equal(expectedFilePath))
	})
	Context("lockfile.PIdOfOwner == 0", func() {
		It("should return expected results", func() {
			/* arrange */
			objectUnderTest := nodeProvider{
				lockfile: new(lockfile.Fake),
			}

			/* act */
			actualNodes, actualError := objectUnderTest.ListNodes()

			/* assert */
			Expect(actualNodes).To(BeNil())
			Expect(actualError).To(BeNil())
		})
	})
	Context("lockfile.PIdOfOwner != 0", func() {
		It("should return expected results", func() {
			/* arrange */
			expectedNodes := []*node.InfoView{{}}

			fakeLockFile := new(lockfile.Fake)
			fakeLockFile.PIdOfOwnerReturns(333)

			objectUnderTest := nodeProvider{
				lockfile: fakeLockFile,
			}

			/* act */
			actualNodes, actualError := objectUnderTest.ListNodes()

			/* assert */
			Expect(actualNodes).To(Equal(expectedNodes))
			Expect(actualError).To(BeNil())
		})
	})
})
