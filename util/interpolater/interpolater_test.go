package interpolater

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Context("Interpolate", func() {
	Describe("called for every scenario", func() {
		It("should return result fulfilling scenario", func() {
			rootPath := "testdata/interpolater/Interpolate"

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
								Name       string
								Expression string
								Scope      map[string]*model.Value
								Expected   string
							}{}
							if err := yaml.Unmarshal(scenariosDotYmlBytes, &scenarioDotYml); nil != err {
								panic(fmt.Errorf("Error unmarshalling scenario.yml for %v; error was %v", path, err))
							}

							for _, scenario := range scenarioDotYml {
								/* act */
								objectUnderTest := New()
								actualResult, actualErr := objectUnderTest.Interpolate(scenario.Expression, scenario.Scope)

								/* assert */
								description := fmt.Sprintf("scenario:\n  path: '%v'\n  name: '%v'", path, scenario.Name)
								Expect(actualResult).To(Equal(scenario.Expected), description)
								Expect(actualErr).To(BeNil(), description)
							}
						}
					}
					return nil
				})
		})
	})
})
