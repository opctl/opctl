package local

import (
	"github.com/golang-utils/lockfile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/node"
)

var _ = Context("listNodes", func() {
	It("should call lockfile.PIdOfOwner w/ expected args", func() {
		/* arrange */
		fakeLockFile := new(lockfile.Fake)

		expectedFilePath := "dummyLockFilePath"

		objectUnderTest := nodeProvider{
			lockfile:     fakeLockFile,
			lockFilePath: expectedFilePath,
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
