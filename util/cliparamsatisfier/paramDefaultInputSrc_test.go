package cliparamsatisfier

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"strconv"
)

var _ = Describe("paramDefaultInputSrc", func() {
	Context("NewParamDefaultInputSrc()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewParamDefaultInputSrc(map[string]*model.Param{})).To(Not(BeNil()))
		})
	})
	Context("Read()", func() {
		Context("inputs contains entry for inputName", func() {
			Context("input already read", func() {
				It("should return nil", func() {
					/* arrange */
					inputName := "dummyInputName"
					paramDefault := "dummyParamDefault"
					param := &model.Param{
						String: &model.StringParam{Default: &paramDefault},
					}
					objectUnderTest := NewParamDefaultInputSrc(
						map[string]*model.Param{inputName: param},
					)

					/* act */
					actualValue1 := objectUnderTest.Read(inputName)
					actualValue2 := objectUnderTest.Read(inputName)

					/* assert */
					Expect(actualValue1).To(Equal(param.String.Default))
					Expect(actualValue2).To(BeNil())
				})
			})
			Context("input not yet read", func() {
				Context("input is dir", func() {
					It("should return param default", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := "dummyParamDefault"
						param := &model.Param{
							Dir: &model.DirParam{Default: &paramDefault},
						}

						expectedValue := param.Dir.Default

						objectUnderTest := NewParamDefaultInputSrc(
							map[string]*model.Param{inputName: param},
						)

						/* act */
						actualValue := objectUnderTest.Read(inputName)

						/* assert */
						Expect(actualValue).To(Equal(expectedValue))
					})
				})
				Context("input is file", func() {
					It("should return param default", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := "dummyParamDefault"
						param := &model.Param{
							File: &model.FileParam{Default: &paramDefault},
						}

						expectedValue := param.File.Default

						objectUnderTest := NewParamDefaultInputSrc(
							map[string]*model.Param{inputName: param},
						)

						/* act */
						actualValue := objectUnderTest.Read(inputName)

						/* assert */
						Expect(actualValue).To(Equal(expectedValue))
					})
				})
				Context("input is number", func() {
					It("should return param default", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := 2.1
						param := &model.Param{
							Number: &model.NumberParam{Default: &paramDefault},
						}

						expectedValue := strconv.FormatFloat(*param.Number.Default, 'E', -1, 64)

						objectUnderTest := NewParamDefaultInputSrc(
							map[string]*model.Param{inputName: param},
						)

						/* act */
						actualValue := objectUnderTest.Read(inputName)

						/* assert */
						Expect(actualValue).To(Equal(&expectedValue))
					})
				})
			})
			Context("input is string", func() {
				It("should return param default", func() {
					/* arrange */
					inputName := "dummyInputName"
					paramDefault := "dummyParamDefault"
					param := &model.Param{
						String: &model.StringParam{Default: &paramDefault},
					}

					expectedValue := param.String.Default

					objectUnderTest := NewParamDefaultInputSrc(
						map[string]*model.Param{inputName: param},
					)

					/* act */
					actualValue := objectUnderTest.Read(inputName)

					/* assert */
					Expect(actualValue).To(Equal(expectedValue))
				})
			})
		})
	})
	Context("inputs doesn't contain entry for inputName", func() {
		It("should return nil", func() {
			/* arrange */
			objectUnderTest := NewParamDefaultInputSrc(
				map[string]*model.Param{},
			)

			/* act */
			actualValue := objectUnderTest.Read("")

			/* assert */
			Expect(actualValue).To(BeNil())
		})
	})
})
