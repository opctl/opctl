package interpolater

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpolate", func() {
	Describe("called for every scenario", func() {
		It("should return result fulfilling scenario", func() {
			rootPath := "testdata/Interpolate"
			fsProvider := fs.New()

			filepath.Walk(rootPath,
				func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						scenariosOpFilePath := filepath.Join(path, "scenarios.yml")
						if _, err := os.Stat(scenariosOpFilePath); err == nil {
							/* arrange */
							scenariosOpFileBytes, err := os.ReadFile(scenariosOpFilePath)
							if err != nil {
								panic(err)
							}

							scenarioOpFile := []struct {
								Name     string
								Template string
								Scope    map[string]*ipld.Node
								Expected string
							}{}
							if err := yaml.Unmarshal(scenariosOpFileBytes, &scenarioOpFile); err != nil {
								panic(fmt.Errorf("error unmarshalling scenario.yml for %s: %w", path, err))
							}

							absPath, err := filepath.Abs(path)
							if err != nil {
								panic(fmt.Errorf("error getting absPath for %s: %w", path, err))
							}

							opHandle, err := data.Resolve(context.Background(), absPath, fsProvider)
							if err != nil {
								panic(fmt.Errorf("error getting opHandle for %s: %w", path, err))
							}

							for _, scenario := range scenarioOpFile {
								// add op dir to scope
								if len(scenario.Scope) == 0 {
									scenario.Scope = map[string]*ipld.Node{}
								}
								scenario.Scope["/"] = &ipld.Node{Dir: opHandle.Path()}

								for name, value := range scenario.Scope {
									// make file refs absolute
									if value.File != nil {
										absFilePath := filepath.Join(absPath, *value.File)
										scenario.Scope[name] = &ipld.Node{File: &absFilePath}
									}
								}

								/* act */
								actualResult, actualErr := Interpolate(
									scenario.Template,
									scenario.Scope,
								)

								/* assert */
								description := fmt.Sprintf("scenario:\n  path: '%v'\n  name: '%v'", path, scenario.Name)
								Expect(actualErr).To(BeNil(), description)
								Expect(actualResult).To(Equal(scenario.Expected), description)
							}
						}
					}
					return nil
				})
		})
	})
})
