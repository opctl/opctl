package opfile

import (
	"context"
	"errors"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	"io"
	"os"
	"path/filepath"
	"strings"

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

var _ = Describe("Go op detection", func() {
	It("detects a Go op without reading it as YAML", func() {
		path, err := filepath.Abs("testdata/goop")
		Expect(err).NotTo(HaveOccurred())
		dir, err := fs.New().TryResolve(context.Background(), path)
		Expect(err).NotTo(HaveOccurred())
		_, err = Get(context.Background(), dir)
		Expect(errors.Is(err, ErrGoOp)).To(BeTrue())
		isGo, err := IsGoOp(context.Background(), dir)
		Expect(err).NotTo(HaveOccurred())
		Expect(isGo).To(BeTrue())
	})
	It("rejects an op dir with both declarations", func() {
		path, err := filepath.Abs("testdata/both")
		Expect(err).NotTo(HaveOccurred())
		dir, err := fs.New().TryResolve(context.Background(), path)
		Expect(err).NotTo(HaveOccurred())
		_, err = Get(context.Background(), dir)
		Expect(err).To(MatchError("op dir contains both op.yml and op.go; an op must be a YAML op or a Go op, not both"))
		_, err = IsGoOp(context.Background(), dir)
		Expect(err).To(HaveOccurred())
	})
	It("preserves the YAML read error when the Go probe also fails", func() {
		dir := new(modelFakes.FakeDataHandle)
		original := errors.New("remote YAML read failed")
		dir.GetContentReturns(nil, original)
		_, err := Get(context.Background(), dir)
		Expect(err).To(Equal(original))
	})
	It("ignores arbitrary probe errors for a YAML op", func() {
		dir := new(modelFakes.FakeDataHandle)
		dir.GetContentStub = func(_ context.Context, path string) (model.ReadSeekCloser, error) {
			if path == "op.yml" {
				return struct {
					io.Reader
					io.Seeker
					io.Closer
				}{strings.NewReader("name: yaml\nrun:\n  container:\n    image: {ref: alpine:3}\n"), nil, io.NopCloser(strings.NewReader(""))}, nil
			}
			return nil, errors.New("HTTP 404")
		}
		spec, err := Get(context.Background(), dir)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Name).To(Equal("yaml"))
	})
})
