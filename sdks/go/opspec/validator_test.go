package opspec

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	. "github.com/opctl/opctl/sdks/go/opspec/opfile/fakes"
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
									panic(fmt.Errorf("error unmarshalling %v; error was %v", description, err))
								}

								for _, scenario := range scenarioOpFile {
									if nil != scenario.Validate {
										/* act */
										fakeHandle := new(modelFakes.FakeDataHandle)
										fakeHandle.GetContentStub = func(ctx context.Context, contentPath string) (model.ReadSeekCloser, error) {
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
		It("should call opFileGetter.Get w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedOpHandle := new(modelFakes.FakeDataHandle)

			fakeOpFileGetter := new(FakeGetter)
			// error to trigger immediate return
			fakeOpFileGetter.GetReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _validator{
				opFileGetter: fakeOpFileGetter,
			}

			/* act */
			objectUnderTest.Validate(
				providedCtx,
				providedOpHandle,
			)

			/* assert */
			actualCtx,
				actualOpHandle := fakeOpFileGetter.GetArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualOpHandle).To(Equal(providedOpHandle))
		})
		Context("opFileGetter.Get errs", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedErrors := []error{errors.New("dummyError")}

				fakeOpFileGetter := new(FakeGetter)
				fakeOpFileGetter.GetReturns(nil, expectedErrors[0])

				objectUnderTest := _validator{
					opFileGetter: fakeOpFileGetter,
				}

				/* act */
				actualErrors := objectUnderTest.Validate(
					context.Background(),
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualErrors).To(Equal(expectedErrors))
			})
		})
		Context("opFileGetter.Get doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeOpFileGetter := new(FakeGetter)

				objectUnderTest := _validator{
					opFileGetter: fakeOpFileGetter,
				}

				/* act */
				actualErrs := objectUnderTest.Validate(
					context.Background(),
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualErrs).To(BeEmpty())
			})
		})
	})
})
