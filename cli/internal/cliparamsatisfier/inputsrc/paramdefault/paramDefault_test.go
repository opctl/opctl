package paramdefault

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Describe("paramDefaultInputSrc", func() {
	Context("ReadString()", func() {
		Context("inputs contains entry for inputName", func() {
			Context("input already read", func() {
				It("should return expected results", func() {
					/* arrange */
					inputName := "dummyInputName"
					paramDefault := "dummyParamDefault"
					param := &model.ParamSpec{
						String: &model.StringParamSpec{Default: &paramDefault},
					}

					objectUnderTest := New(
						map[string]*model.ParamSpec{inputName: param},
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
						param := &model.ParamSpec{
							Array: &model.ArrayParamSpec{Default: paramDefault},
						}

						objectUnderTest := New(
							map[string]*model.ParamSpec{inputName: param},
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
						param := &model.ParamSpec{
							Boolean: &model.BooleanParamSpec{Default: &paramDefault},
						}

						objectUnderTest := New(
							map[string]*model.ParamSpec{inputName: param},
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
							param := &model.ParamSpec{
								Dir: &model.DirParamSpec{Default: &paramDefault},
							}

							objectUnderTest := New(
								map[string]*model.ParamSpec{inputName: param},
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
							param := &model.ParamSpec{
								Dir: &model.DirParamSpec{Default: "./dummyParamDefault"},
							}

							expectedValue := param.Dir.Default

							objectUnderTest := New(
								map[string]*model.ParamSpec{inputName: param},
							)

							/* act */
							actualValue, actualOk := objectUnderTest.ReadString(inputName)

							/* assert */
							Expect(*actualValue).To(Equal(expectedValue))
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
							param := &model.ParamSpec{
								File: &model.FileParamSpec{Default: &paramDefault},
							}

							objectUnderTest := New(
								map[string]*model.ParamSpec{inputName: param},
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
							param := &model.ParamSpec{
								File: &model.FileParamSpec{Default: "./dummyParamDefault"},
							}

							expectedValue := param.File.Default

							objectUnderTest := New(
								map[string]*model.ParamSpec{inputName: param},
							)

							/* act */
							actualValue, actualOk := objectUnderTest.ReadString(inputName)

							/* assert */
							Expect(*actualValue).To(Equal(expectedValue))
							Expect(actualOk).To(BeTrue())
						})
					})
				})
				Context("input is number", func() {
					It("should return expected result", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := 2.1
						param := &model.ParamSpec{
							Number: &model.NumberParamSpec{Default: &paramDefault},
						}

						objectUnderTest := New(
							map[string]*model.ParamSpec{inputName: param},
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
						param := &model.ParamSpec{
							Object: &model.ObjectParamSpec{Default: paramDefault},
						}

						objectUnderTest := New(
							map[string]*model.ParamSpec{inputName: param},
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
						param := &model.ParamSpec{
							String: &model.StringParamSpec{Default: &paramDefault},
						}

						objectUnderTest := New(
							map[string]*model.ParamSpec{inputName: param},
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
			objectUnderTest := New(
				map[string]*model.ParamSpec{},
			)

			/* act */
			actualValue, actualOk := objectUnderTest.ReadString("")

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualOk).To(BeFalse())
		})
	})
})
