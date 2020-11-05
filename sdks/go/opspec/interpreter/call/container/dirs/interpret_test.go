package dirs

import (
	"os"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("dir.Interpret errs", func() {
		It("should return expected error", func() {
			/* arrange */
			identifier:= "identifier"
			/* act */
			_, actualErr := Interpret(
				map[string]*model.Value{
					identifier: &model.Value{
						Socket: new(string),
					},
				},
				map[string]string{
					"/something": fmt.Sprintf("$(%s)", identifier),
				},
				os.TempDir(),
				"dataDirPath",
			)

			/* assert */
			Expect(actualErr).To(Equal(errors.New("unable to bind /something to $(identifier); error was unable to interpret $(identifier) to dir")))
		})
	})
	Context("dir.Interpret doesn't err", func() {
		Context("value.Dir not prefixed by dataDirPath", func() {
			It("should return expected results", func() {
				/* arrange */
				identifier := "identifier"

				dirPath, err := ioutil.TempDir("", "")
				if nil != err {
					panic(err)
				}

				expectedDirs := map[string]string{
					"/something": filepath.Join(os.TempDir(), "/something"),
				}

				/* act */
				actualContainerCallDirs, actualErr := Interpret(
					map[string]*model.Value{
						identifier: &model.Value{Dir: &dirPath},
					},
					map[string]string{
						// implicitly bound
						"/something": fmt.Sprintf("$(%s)", identifier),
					},
					os.TempDir(),
					filepath.Dir(dirPath),
				)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(actualContainerCallDirs).To(Equal(expectedDirs))

			})
		})
		Context("value.Dir prefixed by dataDirPath", func() {
			Context("dircopier.OS errs", func() {
				It("should return expected result", func() {
					/* arrange */
					identifier := "identifier"
					containerPath := "/something"
					dirValue, err := ioutil.TempDir("", "")
					if nil != err {
						panic(err)
					}

					scratchDirPath, err := ioutil.TempDir("", "")
					if nil != err {
						panic(err)
					}

					expectedResult := map[string]string{
						containerPath: filepath.Join(scratchDirPath, containerPath),
					}

					/* act */
					actualResult, actualErr := Interpret(
						map[string]*model.Value{
							identifier: &model.Value{Dir: &dirValue},
						},
						map[string]string{
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
