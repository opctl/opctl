package interpolater

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpolate", func() {
	Describe("called for every scenario", func() {
		It("should return result fulfilling scenario", func() {
			rootPath := "testdata/interpolater/Interpolate"
			data := data.New()
			fsProvider := data.NewFSProvider()

			filepath.Walk(rootPath,
				func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						scenariosOpFilePath := filepath.Join(path, "scenarios.yml")
						if _, err := os.Stat(scenariosOpFilePath); nil == err {
							/* arrange */
							scenariosOpFileBytes, err := ioutil.ReadFile(scenariosOpFilePath)
							if nil != err {
								panic(err)
							}

							scenarioOpFile := []struct {
								Name     string
								Template string
								Scope    map[string]*model.Value
								Expected string
							}{}
							if err := yaml.Unmarshal(scenariosOpFileBytes, &scenarioOpFile); nil != err {
								panic(fmt.Errorf("error unmarshalling scenario.yml for %v; error was %v", path, err))
							}

							absPath, err := filepath.Abs(path)
							if nil != err {
								panic(fmt.Errorf("error getting absPath for %v; error was %v", path, err))
							}

							opHandle, err := data.Resolve(context.Background(), absPath, fsProvider)
							if nil != err {
								panic(fmt.Errorf("error getting opHandle for %v; error was %v", path, err))
							}

							for _, scenario := range scenarioOpFile {
								// add op dir to scope
								if 0 == len(scenario.Scope) {
									scenario.Scope = map[string]*model.Value{}
								}
								scenario.Scope["/"] = &model.Value{Dir: opHandle.Path()}

								for name, value := range scenario.Scope {
									// make file refs absolute
									if nil != value.File {
										absFilePath := filepath.Join(absPath, *value.File)
										scenario.Scope[name] = &model.Value{File: &absFilePath}
									}
								}

								/* act */
								objectUnderTest := New()
								actualResult, actualErr := objectUnderTest.Interpolate(
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
