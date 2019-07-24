package cliparamsatisfier

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Describe("paramDefaultInputSrc", func() {
	Context("ReadString()", func() {
		Context("inputs contains entry for inputName", func() {
			Context("input already read", func() {
				It("should return expected results", func() {
					/* arrange */
					inputName := "dummyInputName"
					paramDefault := "dummyParamDefault"
					param := &types.Param{
						String: &types.StringParam{Default: &paramDefault},
					}

					inputSrcFactory := _InputSrcFactory{}
					objectUnderTest := inputSrcFactory.NewParamDefaultInputSrc(
						map[string]*types.Param{inputName: param},
					)

					/* act */
					actualValue1, actualOk1 := objectUnderTest.ReadString(inputName)
					actualValue2, actualOk2 := objectUnderTest.ReadString(inputName)

					/* assert */
					Expect(actualValue1).To(BeNil())
					Expect(actualOk1).To(BeTrue())

					Expect(actualValue2).To(BeNil())
					Expect(actualOk2).To(BeFalse())
				})
			})
			Context("input not yet read", func() {
				Context("input is array", func() {
					It("should return expected result", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := new([]interface{})
						param := &types.Param{
							Array: &types.ArrayParam{Default: paramDefault},
						}

						inputSrcFactory := _InputSrcFactory{}
						objectUnderTest := inputSrcFactory.NewParamDefaultInputSrc(
							map[string]*types.Param{inputName: param},
						)

						/* act */
						actualValue, actualOk := objectUnderTest.ReadString(inputName)

						/* assert */
						Expect(actualValue).To(BeNil())
						Expect(actualOk).To(BeTrue())
					})
				})
				Context("input is boolean", func() {
					It("should return expected result", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := true
						param := &types.Param{
							Boolean: &types.BooleanParam{Default: &paramDefault},
						}

						inputSrcFactory := _InputSrcFactory{}
						objectUnderTest := inputSrcFactory.NewParamDefaultInputSrc(
							map[string]*types.Param{inputName: param},
						)

						/* act */
						actualValue, actualOk := objectUnderTest.ReadString(inputName)

						/* assert */
						Expect(actualValue).To(BeNil())
						Expect(actualOk).To(BeTrue())
					})
				})
				Context("input is dir", func() {
					Context("default is abs path", func() {
						It("should return expected result", func() {
							/* arrange */
							inputName := "dummyInputName"
							paramDefault := "/dummyParamDefault"
							param := &types.Param{
								Dir: &types.DirParam{Default: &paramDefault},
							}

							inputSrcFactory := _InputSrcFactory{}
							objectUnderTest := inputSrcFactory.NewParamDefaultInputSrc(
								map[string]*types.Param{inputName: param},
							)

							/* act */
							actualValue, actualOk := objectUnderTest.ReadString(inputName)

							/* assert */
							Expect(actualValue).To(BeNil())
							Expect(actualOk).To(BeTrue())
						})
					})
					Context("default is relative path", func() {
						It("should return expected result", func() {
							/* arrange */
							inputName := "dummyInputName"
							paramDefault := "dummyParamDefault"
							param := &types.Param{
								Dir: &types.DirParam{Default: &paramDefault},
							}

							expectedValue := param.Dir.Default

							inputSrcFactory := _InputSrcFactory{}
							objectUnderTest := inputSrcFactory.NewParamDefaultInputSrc(
								map[string]*types.Param{inputName: param},
							)

							/* act */
							actualValue, actualOk := objectUnderTest.ReadString(inputName)

							/* assert */
							Expect(actualValue).To(Equal(expectedValue))
							Expect(actualOk).To(BeTrue())
						})
					})
				})
				Context("input is file", func() {
					Context("default is abs path", func() {
						It("should return expected result", func() {
							/* arrange */
							inputName := "dummyInputName"
							paramDefault := "/dummyParamDefault"
							param := &types.Param{
								File: &types.FileParam{Default: &paramDefault},
							}

							inputSrcFactory := _InputSrcFactory{}
							objectUnderTest := inputSrcFactory.NewParamDefaultInputSrc(
								map[string]*types.Param{inputName: param},
							)

							/* act */
							actualValue, actualOk := objectUnderTest.ReadString(inputName)

							/* assert */
							Expect(actualValue).To(BeNil())
							Expect(actualOk).To(BeTrue())
						})
					})
					Context("default is relative path", func() {
						It("should return expected result", func() {
							/* arrange */
							inputName := "dummyInputName"
							paramDefault := "dummyParamDefault"
							param := &types.Param{
								File: &types.FileParam{Default: &paramDefault},
							}

							expectedValue := param.File.Default

							inputSrcFactory := _InputSrcFactory{}
							objectUnderTest := inputSrcFactory.NewParamDefaultInputSrc(
								map[string]*types.Param{inputName: param},
							)

							/* act */
							actualValue, actualOk := objectUnderTest.ReadString(inputName)

							/* assert */
							Expect(actualValue).To(Equal(expectedValue))
							Expect(actualOk).To(BeTrue())
						})
					})
				})
				Context("input is number", func() {
					It("should return expected result", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := 2.1
						param := &types.Param{
							Number: &types.NumberParam{Default: &paramDefault},
						}

						inputSrcFactory := _InputSrcFactory{}
						objectUnderTest := inputSrcFactory.NewParamDefaultInputSrc(
							map[string]*types.Param{inputName: param},
						)

						/* act */
						actualValue, actualOk := objectUnderTest.ReadString(inputName)

						/* assert */
						Expect(actualValue).To(BeNil())
						Expect(actualOk).To(BeTrue())
					})
				})
				Context("input is object", func() {
					It("should return expected result", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := new(map[string]interface{})
						param := &types.Param{
							Object: &types.ObjectParam{Default: paramDefault},
						}

						inputSrcFactory := _InputSrcFactory{}
						objectUnderTest := inputSrcFactory.NewParamDefaultInputSrc(
							map[string]*types.Param{inputName: param},
						)

						/* act */
						actualValue, actualOk := objectUnderTest.ReadString(inputName)

						/* assert */
						Expect(actualValue).To(BeNil())
						Expect(actualOk).To(BeTrue())
					})
				})
				Context("input is string", func() {
					It("should return expected result", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := "dummyParamDefault"
						param := &types.Param{
							String: &types.StringParam{Default: &paramDefault},
						}

						inputSrcFactory := _InputSrcFactory{}
						objectUnderTest := inputSrcFactory.NewParamDefaultInputSrc(
							map[string]*types.Param{inputName: param},
						)

						/* act */
						actualValue, actualOk := objectUnderTest.ReadString(inputName)

						/* assert */
						Expect(actualValue).To(BeNil())
						Expect(actualOk).To(BeTrue())
					})
				})
			})
		})
	})
	Context("inputs doesn't contain entry for inputName", func() {
		It("should return expected result", func() {
			/* arrange */
			inputSrcFactory := _InputSrcFactory{}
			objectUnderTest := inputSrcFactory.NewParamDefaultInputSrc(
				map[string]*types.Param{},
			)

			/* act */
			actualValue, actualOk := objectUnderTest.ReadString("")

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualOk).To(BeFalse())
		})
	})
})
