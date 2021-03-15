package opspec

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Context("Validate", func() {
	Context("called w/ opspec ../../test-suite scenarios", func() {
		It("should return result fulfilling scenario.validate.expect", func() {
			rootPath := "../../../test-suite"

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
								Validate *struct {
									Expect string
								}
							}{}

							description := fmt.Sprintf("scenario '%v'", path)
							if err := yaml.Unmarshal(scenariosOpFileBytes, &scenarioOpFile); nil != err {
								panic(errors.Wrap(err, "error unmarshalling "+description))
							}

							for _, scenario := range scenarioOpFile {
								if nil != scenario.Validate {
									/* act */
									actualErr := Validate(
										context.Background(),
										path,
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
			/* act */
			actualErr := Validate(
				context.Background(),
				"dummyOpPath",
			)

			/* assert */
			Expect(actualErr.Error()).To(Equal("open dummyOpPath/op.yml: no such file or directory"))
		})
	})
	Context("opFileGetter.Get doesn't err", func() {
		It("should return expected result", func() {
			/* arrange */
			wd, err := os.Getwd()
			if nil != err {
				panic(err)
			}
			opRef := filepath.Join(wd, "testdata/testop")

			/* act */
			actualErr := Validate(
				context.Background(),
				opRef,
			)

			/* assert */
			Expect(actualErr).To(BeNil())
		})
	})
})
