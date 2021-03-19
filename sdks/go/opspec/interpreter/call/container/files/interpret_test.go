package files

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("file.Interpret errs", func() {
		It("should return expected error", func() {
			/* arrange */
			identifier := "identifier"
			providedContainerCallSpecFiles := map[string]interface{}{
				// implicitly bound
				"/somewhere": fmt.Sprintf("$(%s)", identifier),
			}

			/* act */
			_, actualErr := Interpret(
				map[string]*model.Value{
					identifier: {Socket: new(string)},
				},
				providedContainerCallSpecFiles,
				"dummyScratchDirPath",
				"dataDirPath",
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to bind file /somewhere to $(identifier): unable to coerce '{\"socket\":\"\"}' to file"))
		})
	})
	Context("value.File not prefixed by dataDirPath", func() {
		It("should return expected results", func() {
			/* arrange */
			identifier := "identifier"
			providedScope := map[string]*model.Value{
				identifier: {File: new(string)},
			}

			containerPath := "/somewhere"

			providedContainerCallSpecFiles := map[string]interface{}{
				// implicitly bound
				containerPath: fmt.Sprintf("$(%s)", identifier),
			}

			expectedResult := map[string]string{
				containerPath: *providedScope[identifier].File,
			}

			/* act */
			actualResult, actualErr := Interpret(
				providedScope,
				providedContainerCallSpecFiles,
				"dummyScratchDirPath",
				"dataDirPath",
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(actualResult).To(Equal(expectedResult))

		})
	})
	Context("value.File prefixed by dataDirPath", func() {
		It("should return expected result", func() {
			/* arrange */
			identifier := "identifier"
			containerFilePath := "/somewhere"

			scratchDirPath, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

			referencedFile, err := ioutil.TempFile("", "")
			if err != nil {
				panic(err)
			}

			referencedFilePath := referencedFile.Name()

			expectedResult := map[string]string{
				containerFilePath: filepath.Join(scratchDirPath, containerFilePath),
			}

			/* act */
			actualResult, actualErr := Interpret(
				map[string]*model.Value{
					identifier: {
						File: &referencedFilePath,
					},
				},
				map[string]interface{}{
					// implicitly bound
					containerFilePath: fmt.Sprintf("$(%s)", identifier),
				},
				scratchDirPath,
				filepath.Dir(referencedFilePath),
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(actualResult).To(Equal(expectedResult))
		})
	})
})
