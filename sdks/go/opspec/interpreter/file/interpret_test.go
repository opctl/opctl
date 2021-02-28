package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("expression is ref", func() {
		Context("reference.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				/* act */
				actualValue, actualErr := Interpret(
					map[string]*model.Value{},
					"$()",
					"providedScratchDir",
					false,
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(errors.New("unable to interpret $() to file; error was unable to interpret '' as reference; '' not in scope")))
			})
		})
		Context("reference.Interpret doesn't error", func() {
			It("should return expected result", func() {
				/* arrange */
				identifier := "identifier"
				tmpFile, err := ioutil.TempFile("", "")
				if nil != err {
					panic(err)
				}

				tmpFilePath := tmpFile.Name()

				expectedValue := model.Value{Link: &tmpFilePath}

				/* act */
				actualResultValue, actualErr := Interpret(
					map[string]*model.Value{
						identifier: &expectedValue,
					},
					fmt.Sprintf("$(%s)", identifier),
					os.TempDir(),
					false,
				)

				/* assert */
				Expect(*actualResultValue).To(Equal(expectedValue))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("value.Interpret errs", func() {
		It("should return expected result", func() {
			/* arrange */
			/* act */
			_, actualErr := Interpret(
				map[string]*model.Value{},
				nil,
				os.TempDir(),
				false,
			)

			/* assert */
			Expect(actualErr).To(Equal(errors.New("unable to interpret <nil> to file; error was unable to interpret <nil> as value; unsupported type")))
		})
	})
	Context("value.Interpret doesn't err", func() {
		It("should return expected result", func() {
			/* arrange */
			tmpFile, err := ioutil.TempFile("", "")
			if nil != err {
				panic(err)
			}

			tmpFilePath := tmpFile.Name()
			providedExpression := model.Value{Link: &tmpFilePath}

			/* act */
			actualResultValue, actualErr := Interpret(
				map[string]*model.Value{},
				providedExpression,
				os.TempDir(),
				false,
			)

			/* assert */
			Expect(*actualResultValue).To(Equal(providedExpression))
			Expect(actualErr).To(BeNil())
		})
	})
})
