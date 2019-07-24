package params

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/types"
	"path/filepath"
)

var _ = Context("Defaulter", func() {
	Context("NewDefaulter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewDefaulter()).To(Not(BeNil()))
		})
	})
	Context("Default", func() {
		Context("param is array", func() {
			Context("default exists", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputName := "outputName"
					providedOutputDefault := new([]interface{})

					providedOutputParams := map[string]*types.Param{
						providedOutputName: {Array: &types.ArrayParam{Default: providedOutputDefault}},
					}

					objectUnderTest := _defaulter{}

					expectedOutputs := map[string]*types.Value{
						providedOutputName: {Array: providedOutputDefault},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*types.Value{},
						providedOutputParams,
						"dummyOpPath",
					)

					/* assert */
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
		Context("param is boolean", func() {
			Context("default exists", func() {
				It("should set output to default", func() {
					/* arrange */
					providedOutputName := "outputName"
					providedOutputDefault := true

					providedOutputParams := map[string]*types.Param{
						providedOutputName: {Boolean: &types.BooleanParam{Default: &providedOutputDefault}},
					}

					objectUnderTest := _defaulter{}

					expectedOutputs := map[string]*types.Value{
						providedOutputName: {Boolean: &providedOutputDefault},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*types.Value{},
						providedOutputParams,
						"dummyOpPath",
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

					providedOutputParams := map[string]*types.Param{
						providedOutputName: {Dir: &types.DirParam{Default: &providedOutputDefault}},
					}
					providedOpPath := "dummyOpPath"

					objectUnderTest := _defaulter{}

					expectedOutputValue := filepath.Join(providedOpPath, providedOutputDefault)
					expectedOutputs := map[string]*types.Value{
						providedOutputName: {Dir: &expectedOutputValue},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*types.Value{},
						providedOutputParams,
						providedOpPath,
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
					providedOpPath := "dummyOpPath"

					providedOutputParams := map[string]*types.Param{
						providedOutputName: {File: &types.FileParam{Default: &providedOutputDefault}},
					}

					objectUnderTest := _defaulter{}

					expectedOutputValue := filepath.Join(providedOpPath, providedOutputDefault)
					expectedOutputs := map[string]*types.Value{
						providedOutputName: {File: &expectedOutputValue},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*types.Value{},
						providedOutputParams,
						providedOpPath,
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

					providedOutputParams := map[string]*types.Param{
						providedOutputName: {Number: &types.NumberParam{Default: &providedOutputDefault}},
					}

					objectUnderTest := _defaulter{}

					expectedOutputs := map[string]*types.Value{
						providedOutputName: {Number: &providedOutputDefault},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*types.Value{},
						providedOutputParams,
						"dummyOpPath",
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
					providedOutputDefault := new(map[string]interface{})

					providedOutputParams := map[string]*types.Param{
						providedOutputName: {Object: &types.ObjectParam{Default: providedOutputDefault}},
					}

					objectUnderTest := _defaulter{}

					expectedOutputs := map[string]*types.Value{
						providedOutputName: {Object: providedOutputDefault},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*types.Value{},
						providedOutputParams,
						"dummyOpPath",
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

					providedOutputParams := map[string]*types.Param{
						providedOutputName: {String: &types.StringParam{Default: &providedOutputDefault}},
					}

					objectUnderTest := _defaulter{}

					expectedOutputs := map[string]*types.Value{
						providedOutputName: {String: &providedOutputDefault},
					}

					/* act */
					actualOutputs := objectUnderTest.Default(
						map[string]*types.Value{},
						providedOutputParams,
						"dummyOpPath",
					)

					/* assert */
					Expect(actualOutputs).To(Equal(expectedOutputs))
				})
			})
		})
	})
})
