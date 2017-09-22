package interpolater

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Context("Interpolate", func() {
	Describe("called for every scenario", func() {
		It("should return result fulfilling scenario", func() {
			rootPath := "testdata/interpolater/Interpolate"
			pkg := pkg.New()
			pkgFsProvider := pkg.NewFSProvider()

			filepath.Walk(rootPath,
				func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						scenariosDotYmlFilePath := filepath.Join(path, "scenarios.yml")
						if _, err := os.Stat(scenariosDotYmlFilePath); nil == err {
							/* arrange */
							scenariosDotYmlBytes, err := ioutil.ReadFile(scenariosDotYmlFilePath)
							if nil != err {
								panic(err)
							}

							scenarioDotYml := []struct {
								Name     string
								Template string
								Scope    map[string]*model.Value
								Expected string
							}{}
							if err := yaml.Unmarshal(scenariosDotYmlBytes, &scenarioDotYml); nil != err {
								panic(fmt.Errorf("error unmarshalling scenario.yml for %v; error was %v", path, err))
							}

							absPath, err := filepath.Abs(path)
							if nil != err {
								panic(fmt.Errorf("error getting absPath for %v; error was %v", path, err))
							}

							pkgHandle, err := pkg.Resolve(absPath, pkgFsProvider)
							if nil != err {
								panic(fmt.Errorf("error getting pkgHandle for %v; error was %v", path, err))
							}

							for _, scenario := range scenarioDotYml {
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
									pkgHandle,
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
