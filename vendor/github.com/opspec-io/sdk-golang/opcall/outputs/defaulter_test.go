package outputs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

var _ = Context("defaulter", func() {
	Context("Default", func() {
		Context("param is array", func() {
			Context("default exists", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputName := "outputName"
					providedOutputDefault := []interface{}{}

					providedOutputParams := map[string]*model.Param{
						providedOutputName: {Array: &model.ArrayParam{Default: providedOutputDefault}},
					}

					objectUnderTest := _defaulter{}

					expectedOutputs := map[string]*model.Value{
						providedOutputName: {Array: providedOutputDefault},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*model.Value{},
						providedOutputParams,
						"dummyPkgPath",
					)

					/* assert */
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
		Context("param is dir", func() {
			Context("default exists", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputName := "outputName"
					providedOutputDefault := "/pkgDirDefault"

					providedOutputParams := map[string]*model.Param{
						providedOutputName: {Dir: &model.DirParam{Default: &providedOutputDefault}},
					}
					providedPkgPath := "dummyPkgPath"

					objectUnderTest := _defaulter{}

					expectedOutputValue := filepath.Join(providedPkgPath, providedOutputDefault)
					expectedOutputs := map[string]*model.Value{
						providedOutputName: {Dir: &expectedOutputValue},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*model.Value{},
						providedOutputParams,
						providedPkgPath,
					)

					/* assert */
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
		Context("param is file", func() {
			Context("default exists", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputName := "outputName"
					providedOutputDefault := "/pkgFileDefault"
					providedPkgPath := "dummyPkgPath"

					providedOutputParams := map[string]*model.Param{
						providedOutputName: {File: &model.FileParam{Default: &providedOutputDefault}},
					}

					objectUnderTest := _defaulter{}

					expectedOutputValue := filepath.Join(providedPkgPath, providedOutputDefault)
					expectedOutputs := map[string]*model.Value{
						providedOutputName: {File: &expectedOutputValue},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*model.Value{},
						providedOutputParams,
						providedPkgPath,
					)

					/* assert */
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
		Context("param is number", func() {
			Context("default exists", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputName := "outputName"
					providedOutputDefault := 2.2

					providedOutputParams := map[string]*model.Param{
						providedOutputName: {Number: &model.NumberParam{Default: &providedOutputDefault}},
					}

					objectUnderTest := _defaulter{}

					expectedOutputs := map[string]*model.Value{
						providedOutputName: {Number: &providedOutputDefault},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*model.Value{},
						providedOutputParams,
						"dummyPkgPath",
					)

					/* assert */
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
		Context("param is object", func() {
			Context("default exists", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputName := "outputName"
					providedOutputDefault := map[string]interface{}{}

					providedOutputParams := map[string]*model.Param{
						providedOutputName: {Object: &model.ObjectParam{Default: providedOutputDefault}},
					}

					objectUnderTest := _defaulter{}

					expectedOutputs := map[string]*model.Value{
						providedOutputName: {Object: providedOutputDefault},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*model.Value{},
						providedOutputParams,
						"dummyPkgPath",
					)

					/* assert */
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
		Context("param is string", func() {
			Context("default exists", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputName := "outputName"
					providedOutputDefault := "outputDefault"

					providedOutputParams := map[string]*model.Param{
						providedOutputName: {String: &model.StringParam{Default: &providedOutputDefault}},
					}

					objectUnderTest := _defaulter{}

					expectedOutputs := map[string]*model.Value{
						providedOutputName: {String: &providedOutputDefault},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*model.Value{},
						providedOutputParams,
						"dummyPkgPath",
					)

					/* assert */
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
	})
})
