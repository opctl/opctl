package dirs

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("dir.Interpret errs", func() {
		It("should return expected error", func() {
			/* arrange */
			identifier := "identifier"
			dataDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Interpret(
				map[string]*ipld.Node{
					identifier: {
						Socket: new(string),
					},
				},
				map[string]interface{}{
					"/something": fmt.Sprintf("$(%s)", identifier),
				},
				dataDir,
				"dataDirPath",
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to bind directory /something to $(identifier): unable to interpret $(identifier) to dir: unable to coerce socket to dir: incompatible types"))
		})
	})
	Context("dir.Interpret doesn't err", func() {
		Context("value.Dir not prefixed by dataCachePath", func() {
			It("should return expected results", func() {
				/* arrange */
				identifier := "identifier"

				dataDir, err := os.MkdirTemp("", "")
				if err != nil {
					panic(err)
				}
				dirPath, err := os.MkdirTemp("", "")
				if err != nil {
					panic(err)
				}

				expectedDirs := model.NewStringMap(
					map[string]string{
						"/something": filepath.Join(dataDir, "/something"),
					},
				)

				/* act */
				actualContainerCallDirs, actualErr := Interpret(
					map[string]*ipld.Node{
						identifier: {Dir: &dirPath},
					},
					map[string]interface{}{
						// implicitly bound
						"/something": fmt.Sprintf("$(%s)", identifier),
					},
					dataDir,
					filepath.Dir(dirPath),
				)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(actualContainerCallDirs).To(Equal(expectedDirs))

			})
		})
		Context("value.Dir prefixed by dataCachePath", func() {
			Context("dircopier.OS errs", func() {
				It("should return expected result", func() {
					/* arrange */
					identifier := "identifier"
					containerPath := "/something"
					dirValue, err := os.MkdirTemp("", "")
					if err != nil {
						panic(err)
					}

					scratchDirPath, err := os.MkdirTemp("", "")
					if err != nil {
						panic(err)
					}

					expectedResult := model.NewStringMap(
						map[string]string{
							containerPath: filepath.Join(scratchDirPath, containerPath),
						},
					)

					/* act */
					actualResult, actualErr := Interpret(
						map[string]*ipld.Node{
							identifier: {Dir: &dirValue},
						},
						map[string]interface{}{
							// implicitly bound
							containerPath: fmt.Sprintf("$(%s)", identifier),
						},
						scratchDirPath,
						filepath.Dir(dirValue),
					)

					/* assert */
					Expect(actualErr).To(BeNil())
					Expect(actualResult).To(Equal(expectedResult))
				})
			})
		})
	})
})
