package params

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
)

var _ = Context("ApplyDefaults", func() {
	Context("param is array", func() {
		Context("default exists", func() {
			It("should set output to default", func() {
				/* arrange */
				providedOutputName := "outputName"
				providedOutputDefault := []interface{}{}

				providedOutputParams := map[string]*model.ParamSpec{
					providedOutputName: {Array: &model.ArrayParamSpec{Default: providedOutputDefault}},
				}

				expectedOutputs := map[string]*ipld.Node{
					providedOutputName: {Array: &providedOutputDefault},
				}

				/* act */
				actualOutputs, actualErr := ApplyDefaults(
					map[string]*ipld.Node{},
					providedOutputParams,
					"dummyOpPath",
					"dummyOpScratchDir",
				)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(actualOutputs).To(Equal(expectedOutputs))
			})
			Context("array.Interpret errors", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputParams := map[string]*model.ParamSpec{
						"param0": {
							Array: &model.ArrayParamSpec{
								Default: "$(nonExistent)"},
						},
					}

					/* act */
					actualOutputs, actualErr := ApplyDefaults(
						map[string]*ipld.Node{},
						providedOutputParams,
						"dummyOpPath",
						"dummyOpScratchDir",
					)

					/* assert */
					Expect(actualErr).To(MatchError("unable to interpret $(nonExistent) to array: unable to interpret 'nonExistent' as reference: 'nonExistent' not in scope"))
					Expect(actualOutputs).To(BeNil())
				})
			})
		})
	})
	Context("param is boolean", func() {
		Context("default exists", func() {
			It("should set output to default", func() {
				/* arrange */
				providedOutputName := "outputName"
				providedOutputDefault := true

				providedOutputParams := map[string]*model.ParamSpec{
					providedOutputName: {Boolean: &model.BooleanParamSpec{Default: providedOutputDefault}},
				}

				expectedOutputs := map[string]*ipld.Node{
					providedOutputName: {Boolean: &providedOutputDefault},
				}

				/* act */
				actualOutputs, actualErr := ApplyDefaults(
					map[string]*ipld.Node{},
					providedOutputParams,
					"dummyOpPath",
					"dummyOpScratchDir",
				)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(actualOutputs).To(Equal(expectedOutputs))
			})
			Context("boolean.Interpret errors", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputParams := map[string]*model.ParamSpec{
						"param0": {
							Boolean: &model.BooleanParamSpec{
								Default: "$(nonExistent)"},
						},
					}

					/* act */
					actualOutputs, actualErr := ApplyDefaults(
						map[string]*ipld.Node{},
						providedOutputParams,
						"dummyOpPath",
						"dummyOpScratchDir",
					)

					/* assert */
					Expect(actualErr).To(MatchError("unable to interpret $(nonExistent) to boolean: unable to interpret 'nonExistent' as reference: 'nonExistent' not in scope"))
					Expect(actualOutputs).To(BeNil())
				})
			})
		})
	})
	Context("param is dir", func() {
		Context("default exists", func() {
			It("should set output to default", func() {
				/* arrange */
				wd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				providedOpPath := filepath.Join(wd, "testdata")

				providedOutputName := "outputName"
				defaultRelPath := "./"
				providedOutputDefault := opspec.NameToRef(defaultRelPath)

				providedOutputParams := map[string]*model.ParamSpec{
					providedOutputName: {Dir: &model.DirParamSpec{Default: providedOutputDefault}},
				}

				expectedOutputValue := filepath.Join(providedOpPath, defaultRelPath)
				expectedOutputs := map[string]*ipld.Node{
					providedOutputName: {Dir: &expectedOutputValue},
				}

				/* act */
				actualOutputs, actualErr := ApplyDefaults(
					map[string]*ipld.Node{},
					providedOutputParams,
					providedOpPath,
					"dummyOpScratchDir",
				)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(actualOutputs).To(Equal(expectedOutputs))
			})
			Context("dir.Interpret errors", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputParams := map[string]*model.ParamSpec{
						"param0": {
							Dir: &model.DirParamSpec{
								Default: "$(nonExistent)"},
						},
					}

					/* act */
					actualOutputs, actualErr := ApplyDefaults(
						map[string]*ipld.Node{},
						providedOutputParams,
						"dummyOpPath",
						"dummyOpScratchDir",
					)

					/* assert */
					Expect(actualErr).To(MatchError("unable to interpret $(nonExistent) to dir: unable to interpret 'nonExistent' as reference: 'nonExistent' not in scope"))
					Expect(actualOutputs).To(BeNil())
				})
			})
		})
	})
	Context("param is file", func() {
		Context("default exists", func() {
			It("should set output to default", func() {
				/* arrange */
				wd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				providedOpPath := filepath.Join(wd, "testdata")

				providedOutputName := "outputName"
				defaultRelPath := "./empty.txt"
				providedOutputDefault := opspec.NameToRef(defaultRelPath)

				providedOutputParams := map[string]*model.ParamSpec{
					providedOutputName: {File: &model.FileParamSpec{Default: providedOutputDefault}},
				}

				expectedOutputValue := filepath.Join(providedOpPath, defaultRelPath)
				expectedOutputs := map[string]*ipld.Node{
					providedOutputName: {File: &expectedOutputValue},
				}

				/* act */
				actualOutputs, actualErr := ApplyDefaults(
					map[string]*ipld.Node{},
					providedOutputParams,
					providedOpPath,
					"dummyOpScratchDir",
				)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(actualOutputs).To(Equal(expectedOutputs))
			})
			Context("file.Interpret errors", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputParams := map[string]*model.ParamSpec{
						"param0": {
							File: &model.FileParamSpec{
								Default: "$(nonExistent)"},
						},
					}

					/* act */
					actualOutputs, actualErr := ApplyDefaults(
						map[string]*ipld.Node{},
						providedOutputParams,
						"dummyOpPath",
						"dummyOpScratchDir",
					)

					/* assert */
					Expect(actualErr).To(MatchError("unable to interpret $(nonExistent) to file: unable to interpret 'nonExistent' as reference: 'nonExistent' not in scope"))
					Expect(actualOutputs).To(BeNil())
				})
			})
		})
	})
	Context("param is number", func() {
		Context("default exists", func() {
			It("should set output to default", func() {
				/* arrange */
				providedOutputName := "outputName"
				providedOutputDefault := 2.2

				providedOutputParams := map[string]*model.ParamSpec{
					providedOutputName: {Number: &model.NumberParamSpec{Default: providedOutputDefault}},
				}

				expectedOutputs := map[string]*ipld.Node{
					providedOutputName: {Number: &providedOutputDefault},
				}

				/* act */
				actualOutputs, actualErr := ApplyDefaults(
					map[string]*ipld.Node{},
					providedOutputParams,
					"dummyOpPath",
					"dummyOpScratchDir",
				)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(actualOutputs).To(Equal(expectedOutputs))
			})
			Context("number.Interpret errors", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputParams := map[string]*model.ParamSpec{
						"param0": {
							Number: &model.NumberParamSpec{
								Default: "$(nonExistent)"},
						},
					}

					/* act */
					actualOutputs, actualErr := ApplyDefaults(
						map[string]*ipld.Node{},
						providedOutputParams,
						"dummyOpPath",
						"dummyOpScratchDir",
					)

					/* assert */
					Expect(actualErr).To(MatchError("unable to interpret $(nonExistent) to number: unable to interpret 'nonExistent' as reference: 'nonExistent' not in scope"))
					Expect(actualOutputs).To(BeNil())
				})
			})
		})
	})
	Context("param is object", func() {
		Context("default exists", func() {
			It("should set output to default", func() {
				/* arrange */
				providedOutputName := "outputName"
				providedOutputDefault := map[string]interface{}{}

				providedOutputParams := map[string]*model.ParamSpec{
					providedOutputName: {Object: &model.ObjectParamSpec{Default: providedOutputDefault}},
				}

				expectedOutputs := map[string]*ipld.Node{
					providedOutputName: {Object: &providedOutputDefault},
				}

				/* act */
				actualOutputs, actualErr := ApplyDefaults(
					map[string]*ipld.Node{},
					providedOutputParams,
					"dummyOpPath",
					"dummyOpScratchDir",
				)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(actualOutputs).To(Equal(expectedOutputs))
			})
			Context("object.Interpret errors", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputParams := map[string]*model.ParamSpec{
						"param0": {
							Object: &model.ObjectParamSpec{
								Default: "$(nonExistent)"},
						},
					}

					/* act */
					actualOutputs, actualErr := ApplyDefaults(
						map[string]*ipld.Node{},
						providedOutputParams,
						"dummyOpPath",
						"dummyOpScratchDir",
					)

					/* assert */
					Expect(actualErr).To(MatchError("unable to interpret $(nonExistent) to object: unable to interpret 'nonExistent' as reference: 'nonExistent' not in scope"))
					Expect(actualOutputs).To(BeNil())
				})
			})
		})
	})
	Context("param is string", func() {
		Context("default exists", func() {
			It("should set output to default", func() {
				/* arrange */
				providedOutputName := "outputName"
				providedOutputDefault := "outputDefault"

				providedOutputParams := map[string]*model.ParamSpec{
					providedOutputName: {String: &model.StringParamSpec{Default: providedOutputDefault}},
				}

				expectedOutputs := map[string]*ipld.Node{
					providedOutputName: {String: &providedOutputDefault},
				}

				/* act */
				actualOutputs, actualErr := ApplyDefaults(
					map[string]*ipld.Node{},
					providedOutputParams,
					"dummyOpPath",
					"dummyOpScratchDir",
				)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(actualOutputs).To(Equal(expectedOutputs))
			})
			Context("string.Interpret errors", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputParams := map[string]*model.ParamSpec{
						"param0": {
							String: &model.StringParamSpec{
								Default: "$(nonExistent)"},
						},
					}

					/* act */
					actualOutputs, actualErr := ApplyDefaults(
						map[string]*ipld.Node{},
						providedOutputParams,
						"dummyOpPath",
						"dummyOpScratchDir",
					)

					/* assert */
					Expect(actualErr).To(MatchError("unable to interpret $(nonExistent) to string: unable to interpret 'nonExistent' as reference: 'nonExistent' not in scope"))
					Expect(actualOutputs).To(BeNil())
				})
			})
		})
	})
})
