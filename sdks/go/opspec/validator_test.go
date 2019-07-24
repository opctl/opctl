package op

import (
	"context"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"github.com/opctl/opctl/sdks/go/types"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("Validator", func() {
	Context("NewValidator", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewValidator()).Should(Not(BeNil()))
		})
	})
	Context("Validate", func() {
		Context("called w/ opspec ../../test-suite scenarios", func() {
			It("should return result fulfilling scenario.validate.expect", func() {
				rootPath := "../../../test-suite"

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
									Validate *struct {
										Expect string
									}
								}{}

								description := fmt.Sprintf("scenario '%v'", path)
								if err := yaml.Unmarshal(scenariosDotYmlBytes, &scenarioDotYml); nil != err {
									panic(fmt.Errorf("error unmarshalling %v; error was %v", description, err))
								}

								for _, scenario := range scenarioDotYml {
									if nil != scenario.Validate {
										/* act */
										fakeHandle := new(data.FakeHandle)
										fakeHandle.GetContentStub = func(ctx context.Context, contentPath string) (types.ReadSeekCloser, error) {
											return os.Open(filepath.Join(path, contentPath))
										}

										objectUnderTest := NewValidator()
										actualErrs := objectUnderTest.Validate(
											context.Background(),
											fakeHandle,
										)

										/* assert */
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
		It("should call dotYmlGetter.Get w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedOpHandle := new(data.FakeHandle)

			fakeDotYmlGetter := new(dotyml.FakeGetter)
			// error to trigger immediate return
			fakeDotYmlGetter.GetReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _validator{
				dotYmlGetter: fakeDotYmlGetter,
			}

			/* act */
			objectUnderTest.Validate(
				providedCtx,
				providedOpHandle,
			)

			/* assert */
			actualCtx,
				actualOpHandle := fakeDotYmlGetter.GetArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualOpHandle).To(Equal(providedOpHandle))
		})
		Context("dotYmlGetter.Get errs", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedErrors := []error{errors.New("dummyError")}

				fakeDotYmlGetter := new(dotyml.FakeGetter)
				fakeDotYmlGetter.GetReturns(nil, expectedErrors[0])

				objectUnderTest := _validator{
					dotYmlGetter: fakeDotYmlGetter,
				}

				/* act */
				actualErrors := objectUnderTest.Validate(
					context.Background(),
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))
			})
		})
		Context("dotYmlGetter.Get doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeDotYmlGetter := new(dotyml.FakeGetter)

				objectUnderTest := _validator{
					dotYmlGetter: fakeDotYmlGetter,
				}

				/* act */
				actualErrs := objectUnderTest.Validate(
					context.Background(),
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualErrs).To(BeEmpty())
			})
		})
	})
})
