package cliprompt

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Describe("cliPromptInputSrc", func() {
	Context("ReadString()", func() {
		Context("inputs contains entry for inputName", func() {
			Context("input already read", func() {
				It("should return expected results", func() {
					/* arrange */
					inputName := "dummyInputName"
					paramDefault := "dummyParamDefault"
					param := &model.Param{
						Description: "description",
						String: &model.StringParam{
							Default:     &paramDefault,
							Description: "deprecated description",
						},
					}

					objectUnderTest := New(
						map[string]*model.Param{inputName: param},
					)

					/* act */
					actualValue1, actualOk1 := objectUnderTest.ReadString(inputName)
					actualValue2, actualOk2 := objectUnderTest.ReadString(inputName)

					/* assert */
					Expect(actualValue1).To(BeNil())
					Expect(actualOk1).To(BeFalse())

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
						param := &model.Param{
							Array: &model.ArrayParam{Default: paramDefault},
						}

						objectUnderTest := New(
							map[string]*model.Param{inputName: param},
						)

						/* act */
						actualValue, actualOk := objectUnderTest.ReadString(inputName)

						/* assert */
						Expect(actualValue).To(BeNil())
						Expect(actualOk).To(BeFalse())
					})
				})
				Context("input is boolean", func() {
					It("should return expected result", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := true
						param := &model.Param{
							Boolean: &model.BooleanParam{Default: &paramDefault},
						}

						objectUnderTest := New(
							map[string]*model.Param{inputName: param},
						)

						/* act */
						actualValue, actualOk := objectUnderTest.ReadString(inputName)

						/* assert */
						Expect(actualValue).To(BeNil())
						Expect(actualOk).To(BeFalse())
					})
				})
				Context("input is dir", func() {
					Context("default is abs path", func() {
						It("should return expected result", func() {
							/* arrange */
							inputName := "dummyInputName"
							paramDefault := "/dummyParamDefault"
							param := &model.Param{
								Dir: &model.DirParam{Default: &paramDefault},
							}

							objectUnderTest := New(
								map[string]*model.Param{inputName: param},
							)

							/* act */
							actualValue, actualOk := objectUnderTest.ReadString(inputName)

							/* assert */
							Expect(actualValue).To(BeNil())
							Expect(actualOk).To(BeFalse())
						})
					})
				})
				Context("input is file", func() {
					Context("default is abs path", func() {
						It("should return expected result", func() {
							/* arrange */
							inputName := "dummyInputName"
							paramDefault := "/dummyParamDefault"
							param := &model.Param{
								File: &model.FileParam{Default: &paramDefault},
							}

							objectUnderTest := New(
								map[string]*model.Param{inputName: param},
							)

							/* act */
							actualValue, actualOk := objectUnderTest.ReadString(inputName)

							/* assert */
							Expect(actualValue).To(BeNil())
							Expect(actualOk).To(BeFalse())
						})
					})
				})
				Context("input is number", func() {
					It("should return expected result", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := 2.1
						param := &model.Param{
							Number: &model.NumberParam{Default: &paramDefault},
						}

						objectUnderTest := New(
							map[string]*model.Param{inputName: param},
						)

						/* act */
						actualValue, actualOk := objectUnderTest.ReadString(inputName)

						/* assert */
						Expect(actualValue).To(BeNil())
						Expect(actualOk).To(BeFalse())
					})
				})
				Context("input is object", func() {
					It("should return expected result", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := new(map[string]interface{})
						param := &model.Param{
							Object: &model.ObjectParam{Default: paramDefault},
						}

						objectUnderTest := New(
							map[string]*model.Param{inputName: param},
						)

						/* act */
						actualValue, actualOk := objectUnderTest.ReadString(inputName)

						/* assert */
						Expect(actualValue).To(BeNil())
						Expect(actualOk).To(BeFalse())
					})
				})
				Context("input is string", func() {
					It("should return expected result", func() {
						/* arrange */
						inputName := "dummyInputName"
						paramDefault := "dummyParamDefault"
						param := &model.Param{
							String: &model.StringParam{Default: &paramDefault},
						}

						objectUnderTest := New(
							map[string]*model.Param{inputName: param},
						)

						/* act */
						actualValue, actualOk := objectUnderTest.ReadString(inputName)

						/* assert */
						Expect(actualValue).To(BeNil())
						Expect(actualOk).To(BeFalse())
					})
				})
			})
		})
	})
	Context("inputs doesn't contain entry for inputName", func() {
		It("should return expected result", func() {
			/* arrange */
			objectUnderTest := New(
				map[string]*model.Param{},
			)

			/* act */
			actualValue, actualOk := objectUnderTest.ReadString("")

			/* assert */
			Expect(actualValue).To(BeNil())
			Expect(actualOk).To(BeFalse())
		})
	})
})
