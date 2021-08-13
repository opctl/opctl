package cliprompt

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Describe("cliPromptInputSrc", func() {
	cliOutput := clioutput.New(clicolorer.New(), ioutil.Discard, ioutil.Discard)

	Context("ReadString()", func() {
		Context("inputs contains entry for inputName", func() {
			Context("input already read", func() {
				It("should return expected results", func() {
					/* arrange */
					inputName := "dummyInputName"
					paramDefault := "dummyParamDefault"
					param := &model.ParamSpec{
						Description: "description",
						String: &model.StringParamSpec{
							Default:     &paramDefault,
							Description: "deprecated description",
						},
					}

					objectUnderTest := New(
						cliOutput,
						map[string]*model.ParamSpec{inputName: param},
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
						param := &model.ParamSpec{
							Array: &model.ArrayParamSpec{Default: paramDefault},
						}

						objectUnderTest := New(
							cliOutput,
							map[string]*model.ParamSpec{inputName: param},
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
						param := &model.ParamSpec{
							Boolean: &model.BooleanParamSpec{Default: &paramDefault},
						}

						objectUnderTest := New(
							cliOutput,
							map[string]*model.ParamSpec{inputName: param},
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
							param := &model.ParamSpec{
								Dir: &model.DirParamSpec{Default: &paramDefault},
							}

							objectUnderTest := New(
								cliOutput,
								map[string]*model.ParamSpec{inputName: param},
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
							param := &model.ParamSpec{
								File: &model.FileParamSpec{Default: &paramDefault},
							}

							objectUnderTest := New(
								cliOutput,
								map[string]*model.ParamSpec{inputName: param},
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
						param := &model.ParamSpec{
							Number: &model.NumberParamSpec{Default: &paramDefault},
						}

						objectUnderTest := New(
							cliOutput,
							map[string]*model.ParamSpec{inputName: param},
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
						param := &model.ParamSpec{
							Object: &model.ObjectParamSpec{Default: paramDefault},
						}

						objectUnderTest := New(
							cliOutput,
							map[string]*model.ParamSpec{inputName: param},
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
						param := &model.ParamSpec{
							String: &model.StringParamSpec{Default: &paramDefault},
						}

						objectUnderTest := New(
							cliOutput,
							map[string]*model.ParamSpec{inputName: param},
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
				cliOutput,
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
