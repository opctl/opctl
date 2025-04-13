package opfile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data/fs"
)

var _ = Context("Validate", func() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	Context("called w/ opspec ../../test-suite scenarios", func() {
		It("should return result fulfilling scenario.validate.expect", func() {
			rootPath := "../../../../test-suite"

			filepath.Walk(rootPath,
				func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						scenariosOpFilePath := filepath.Join(path, "scenarios.json")
						if _, err := os.Stat(scenariosOpFilePath); err == nil {
							/* arrange */
							scenariosOpFileBytes, err := os.ReadFile(scenariosOpFilePath)
							if err != nil {
								panic(err)
							}

							scenarioOpFile := []struct {
								Validate *struct {
									Expect string
								}
							}{}

							description := fmt.Sprintf("scenario '%v'", path)
							if err := yaml.Unmarshal(scenariosOpFileBytes, &scenarioOpFile); err != nil {
								panic(fmt.Errorf("error unmarshalling %s: %w", description, err))
							}

							for _, scenario := range scenarioOpFile {
								if scenario.Validate != nil {
									/* arrange */

									providedCtx := context.Background()

									opDir, err := fs.New(wd).TryResolve(providedCtx, path)
									if err != nil {
										panic(fmt.Errorf("error resolving %s", path))
									}

									/* act */
									_, actualErr := Get(
										providedCtx,
										opDir,
									)

									/* assert */
									switch expect := scenario.Validate.Expect; expect {
									case "success":
										Expect(actualErr).To(BeNil(), description)
									case "failure":
										Expect(actualErr).To(Not(BeNil()), description)
									}
								}
							}
						}
					}
					return nil
				})
		})
	})
	Context("opFileGetter.Get errs", func() {
		It("should return expected result", func() {
			/* arrange */
			providedCtx := context.Background()
			opDir, err := fs.New(wd).TryResolve(providedCtx, "testdata")
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Get(
				providedCtx,
				opDir,
			)

			/* assert */
			Expect(actualErr).To(Not(BeNil()))
		})
	})
	Context("opFileGetter.Get doesn't err", func() {
		It("should return expected result", func() {
			/* arrange */
			providedCtx := context.Background()
			opDir, err := fs.New().TryResolve(providedCtx, filepath.Join(wd, "testdata/testop"))
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Get(
				providedCtx,
				opDir,
			)

			/* assert */
			Expect(actualErr).To(BeNil())
		})
	})
})
