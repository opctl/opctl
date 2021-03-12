package local

import (
	"io/ioutil"
	"path/filepath"

	"github.com/golang-utils/lockfile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/sdks/go/node"
)

var _ = Context("listNodes", func() {
	tmpDir, err := ioutil.TempDir("", "")
	if nil != err {
		panic(err)
	}

	dataDir, newDataDirErr := datadir.New(tmpDir)
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
			listenAddress := "127.0.0.1:42224"
			expectedNode, err := newAPIClientNode(listenAddress)
			if nil != err {
				panic(err)
			}

			expectedNodes := []node.Node{
				expectedNode,
			}

			fakeLockFile := new(lockfile.Fake)
			fakeLockFile.PIdOfOwnerReturns(333)

			objectUnderTest := nodeProvider{
				dataDir:       dataDir,
				lockfile:      fakeLockFile,
				listenAddress: listenAddress,
			}

			/* act */
			actualNodes, actualError := objectUnderTest.ListNodes()

			/* assert */
			Expect(actualNodes).To(HaveLen(len(expectedNodes)))
			Expect(actualError).To(BeNil())
		})
	})
})
