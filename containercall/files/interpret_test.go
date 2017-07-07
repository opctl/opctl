package files

import (
	"fmt"
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/filecopier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

var _ = Context("Files", func() {
	Context("Interpret", func() {
		Context("bound to socket", func() {
			It("should return expected error", func() {
				/* arrange */
				containerFileBind := "dummyContainerFileBind"

				providedScope := map[string]*model.Value{
					containerFileBind: {Socket: new(string)},
				}

				containerPath := "dummyContainerPath"
				providedSCGContainerCallFiles := map[string]string{
					// explicitly bound
					containerPath: containerFileBind,
				}

				expectedErr := fmt.Errorf("Unable to bind file '%v' to '%v'. '%v' not a file, number, or string", containerPath, containerFileBind, containerFileBind)

				objectUnderTest := _Files{
					os: new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					"dummyPkgPath",
					providedScope,
					providedSCGContainerCallFiles,
					"dummyScratchDirPath",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("bound to dir", func() {
			It("should return expected error", func() {
				/* arrange */
				containerFileBind := "dummyContainerFileBind"

				providedScope := map[string]*model.Value{
					containerFileBind: {Dir: new(string)},
				}

				containerPath := "dummyContainerPath"
				providedSCGContainerCallFiles := map[string]string{
					// explicitly bound
					containerPath: containerFileBind,
				}

				expectedErr := fmt.Errorf("Unable to bind file '%v' to '%v'. '%v' not a file, number, or string", containerPath, containerFileBind, containerFileBind)

				objectUnderTest := _Files{
					os: new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					"dummyPkgPath",
					providedScope,
					providedSCGContainerCallFiles,
					"dummyScratchDirPath",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		It("should return expected dcg.Files", func() {

			/* arrange */
			providedRootFSPath := "/dummyRootFSPath"

			providedScratchDir := "dummyScratchDir"

			expectedFile1Path := "/dummyFile1Path.txt"
			expectedDCGContainerCallFiles := map[string]string{
				expectedFile1Path: filepath.Join(providedScratchDir, expectedFile1Path),
			}

			providedSCGContainerCallFiles := map[string]string{
				// implicitly bound
				expectedFile1Path: "",
			}

			objectUnderTest := _Files{
				fileCopier: new(filecopier.Fake),
				os:         new(ios.Fake),
				rootFSPath: providedRootFSPath,
			}

			/* act */
			actualDCGContainerCallFiles, _ := objectUnderTest.Interpret(
				"dummyPkgPath",
				map[string]*model.Value{},
				providedSCGContainerCallFiles,
				providedScratchDir,
			)

			/* assert */
			Expect(actualDCGContainerCallFiles).To(Equal(expectedDCGContainerCallFiles))
		})
	})
})
