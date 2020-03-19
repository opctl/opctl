package local

import (
	"path/filepath"

	"github.com/golang-utils/lockfile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/model"
)

var _ = Context("listNodes", func() {
	dataDir, newDataDirErr := datadir.New(nil)
	if nil != newDataDirErr {
		panic(newDataDirErr)
	}

	It("should call lockfile.PIdOfOwner w/ expected args", func() {
		/* arrange */
		fakeLockFile := new(lockfile.Fake)

		objectUnderTest := nodeProvider{
			dataDir:  dataDir,
			lockfile: fakeLockFile,
		}

		/* act */
		objectUnderTest.ListNodes()

		/* assert */
		Expect(fakeLockFile.PIdOfOwnerArgsForCall(0)).To(Equal(filepath.Join(dataDir.Path(), "pid.lock")))
	})
	Context("lockfile.PIdOfOwner == 0", func() {
		It("should return expected results", func() {
			/* arrange */
			objectUnderTest := nodeProvider{
				dataDir:  dataDir,
				lockfile: new(lockfile.Fake),
			}

			/* act */
			actualNodes, actualError := objectUnderTest.ListNodes()

			/* assert */
			Expect(actualNodes).To(BeEmpty())
			Expect(actualError).To(BeNil())
		})
	})
	Context("lockfile.PIdOfOwner != 0", func() {
		It("should return expected results", func() {
			/* arrange */
			nodeHandle, _ := newNodeHandle()
			expectedNodes := []model.NodeHandle{
				nodeHandle,
			}

			fakeLockFile := new(lockfile.Fake)
			fakeLockFile.PIdOfOwnerReturns(333)

			objectUnderTest := nodeProvider{
				dataDir:  dataDir,
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
