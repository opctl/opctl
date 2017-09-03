package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/manifest"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("Validate", func() {
	Context("called w/ opspec test-suite scenarios", func() {
		It("should return result fulfilling scenario.expect.validate", func() {
			rootPath := "../github.com/opspec-io/test-suite/scenarios/pkg"

			pendingScenarios := map[string]interface{}{
				// these scenarios are currently pending;
				filepath.Join(rootPath, "inputs/dir/default/is-file"):                 nil,
				filepath.Join(rootPath, "inputs/file/default/is-dir"):                 nil,
				filepath.Join(rootPath, "run/op/inputs/file-to-number/isnt-number"):   nil,
				filepath.Join(rootPath, "run/op/inputs/file-to-object/isnt-object"):   nil,
				filepath.Join(rootPath, "run/op/inputs/string-to-number/isnt-number"): nil,
				filepath.Join(rootPath, "run/op/inputs/string-to-object/isnt-object"): nil,
			}

			filepath.Walk(rootPath,
				func(path string, info os.FileInfo, err error) error {
					_, isPending := pendingScenarios[path]
					if !isPending && info.IsDir() {
						scenariosDotYmlFilePath := filepath.Join(path, "scenarios.yml")
						if _, err := os.Stat(scenariosDotYmlFilePath); nil == err {
							/* arrange */
							scenariosDotYmlBytes, err := ioutil.ReadFile(scenariosDotYmlFilePath)
							if nil != err {
								panic(err)
							}

							scenarioDotYml := []struct {
								Call *struct {
									Expect string
								}
								Validate *struct {
									Expect string
								}
							}{}

							yaml.Unmarshal(scenariosDotYmlBytes, &scenarioDotYml)

							for _, scenario := range scenarioDotYml {
								if nil != scenario.Validate {
									/* act */
									objectUnderTest := New()
									actualErrs := objectUnderTest.Validate(newFSHandle(path))

									/* assert */
									description := fmt.Sprintf("scenario path: '%v'", path)
									switch expect := scenario.Validate.Expect; expect {
									case "success":
										Expect(actualErrs).To(BeEmpty(), description)
									case "failure":
										Expect(actualErrs).To(Not(BeEmpty()), description)
									}
								}
							}
						}
					}
					return nil
				})
		})
	})
	It("should call handle.GetContent w/ expected args", func() {
		/* arrange */
		providedFileHandle := new(FakeHandle)
		// error to trigger immediate return
		providedFileHandle.GetContentReturns(nil, errors.New("dummyError"))

		objectUnderTest := _Pkg{}

		/* act */
		objectUnderTest.Validate(providedFileHandle)

		/* assert */
		actualCtx,
			actualContentName := providedFileHandle.GetContentArgsForCall(0)

		Expect(actualCtx).To(Equal(context.TODO()))
		Expect(actualContentName).To(Equal(OpDotYmlFileName))
	})
	Context("handle.GetContent errs", func() {
		It("should return err", func() {
			/* arrange */
			expectedErrors := []error{errors.New("dummyError")}
			providedFileHandle := new(FakeHandle)
			// error to trigger immediate return
			providedFileHandle.GetContentReturns(nil, expectedErrors[0])

			objectUnderTest := _Pkg{}

			/* act */
			actualErrors := objectUnderTest.Validate(providedFileHandle)

			/* assert */
			Expect(actualErrors).To(Equal(expectedErrors))
		})
	})
	Context("handle.GetContent doesn't err", func() {
		It("should call manifestValidator.Validate w/ expected args & return result", func() {
			/* arrange */

			providedFileHandle := new(FakeHandle)

			expectedManifestBytes := []byte{2, 5, 61}
			fakeIOUtil := new(iioutil.Fake)
			fakeIOUtil.ReadAllReturns(expectedManifestBytes, nil)

			expectedErrs := []error{
				errors.New("dummyErr1"),
				errors.New("dummyErr2"),
			}
			fakeManifest := new(manifest.Fake)

			fakeManifest.ValidateReturns(expectedErrs)

			objectUnderTest := _Pkg{
				ioUtil:   fakeIOUtil,
				manifest: fakeManifest,
			}

			/* act */
			actualErrs := objectUnderTest.Validate(providedFileHandle)

			/* assert */
			Expect(fakeManifest.ValidateArgsForCall(0)).To(Equal(expectedManifestBytes))
			Expect(actualErrs).To(Equal(expectedErrs))
		})
	})
})
