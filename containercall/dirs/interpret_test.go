package dirs

import (
	"fmt"
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/dircopier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

var _ = Context("Files", func() {
	Context("Interpret", func() {
		Context("bound to non dir", func() {
			It("should return expected error", func() {
				/* arrange */
				containerDirBind := "dummyContainerDirBind"

				providedScope := map[string]*model.Value{
					containerDirBind: {String: new(string)},
				}

				containerPath := "dummyContainerPath"
				providedSCGContainerCallDirs := map[string]string{
					// explicitly bound
					containerPath: containerDirBind,
				}

				expectedErr := fmt.Errorf("Unable to bind dir '%v' to '%v'. '%v' not a dir", containerPath, containerDirBind, containerDirBind)

				objectUnderTest := _Dirs{
					os: new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					"dummyPkgRef",
					providedScope,
					providedSCGContainerCallDirs,
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		It("should return expected dcg.Dirs", func() {

			/* arrange */
			rootFSPath := "/dummyRootFSPath"
			providedContainerId := "dummyContainerId"
			providedRootOpId := "dummyRootOpId"
			providedPkgPath := "pkgPath"

			providedScratchDirPath := filepath.Join(
				rootFSPath,
				"dcg",
				providedRootOpId,
				"containers",
				providedContainerId,
				"fs",
			)

			expectedDir1Path := "/dummyFile1Path.txt"
			expectedDirs := map[string]string{
				expectedDir1Path: filepath.Join(providedScratchDirPath, expectedDir1Path),
			}

			providedSCGContainerCallDirs := map[string]string{
				// implicitly bound
				expectedDir1Path: "",
			}

			objectUnderTest := _Dirs{
				dirCopier:  new(dircopier.Fake),
				os:         new(ios.Fake),
				rootFSPath: rootFSPath,
			}

			/* act */
			actualDCGContainerCallDirs, _ := objectUnderTest.Interpret(
				providedPkgPath,
				map[string]*model.Value{},
				providedSCGContainerCallDirs,
				providedScratchDirPath,
			)

			/* assert */
			Expect(actualDCGContainerCallDirs).To(Equal(expectedDirs))
		})
	})
})
