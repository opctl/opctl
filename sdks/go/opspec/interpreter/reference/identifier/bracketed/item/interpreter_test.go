package item

import (
	"errors"
	"fmt"

	. "github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/bracketed/item/internal/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).Should(Not(BeNil()))
		})
	})
	Context("Interpret", func() {

		It("should call parseIndexer.ParseIndex w/ expected args", func() {
			/* arrange */
			providedIndexString := "dummyIndexString"
			providedData := model.Value{Array: new([]interface{})}

			fakeParseIndexer := new(FakeParseIndexer)
			// err to trigger immediate return
			fakeParseIndexer.ParseIndexReturns(0, errors.New("dummyErr"))

			objectUnderTest := _interpreter{
				parseIndexer: fakeParseIndexer,
			}

			/* act */
			objectUnderTest.Interpret(
				providedIndexString,
				providedData,
			)

			/* assert */
			actualIndexString,
				actualArray := fakeParseIndexer.ParseIndexArgsForCall(0)

			Expect(actualIndexString).To(Equal(providedIndexString))
			Expect(actualArray).To(Equal(*providedData.Array))
		})
		Context("parseIndexer.ParseIndex errs", func() {
			It("should return expected result", func() {
				/* arrange */

				fakeParseIndexer := new(FakeParseIndexer)
				parseIndexErr := errors.New("dummyErr")
				fakeParseIndexer.ParseIndexReturns(0, parseIndexErr)

				expectedErr := fmt.Errorf("unable to interpret item; error was %v", parseIndexErr.Error())

				objectUnderTest := _interpreter{
					parseIndexer: fakeParseIndexer,
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					"dummyIndexString",
					model.Value{Array: new([]interface{})},
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("parseIndexer.ParseIndex doesn't err", func() {
			array := &[]interface{}{"dummyItem"}
			nonEmptyArrayValue := model.Value{Array: array}

			It("should call valueConstructor.Construct w/ expected args", func() {
				/* arrange */
				providedData := nonEmptyArrayValue

				fakeParseIndexer := new(FakeParseIndexer)
				parsedIndex := 0
				fakeParseIndexer.ParseIndexReturns(0, nil)

				fakeValueConstructor := new(value.FakeConstructor)
				// err to trigger immediate return
				fakeValueConstructor.ConstructReturns(nil, errors.New("dummyErr"))

				objectUnderTest := _interpreter{
					parseIndexer:     fakeParseIndexer,
					valueConstructor: fakeValueConstructor,
				}

				/* act */
				objectUnderTest.Interpret(
					"dummyIndexString",
					providedData,
				)

				/* assert */
				actualValue := fakeValueConstructor.ConstructArgsForCall(0)

				Expect(actualValue).To(Equal((*providedData.Array)[parsedIndex]))
			})
			Context("valueConstructor.Construct errs", func() {

				It("should return expected result", func() {
					/* arrange */
					fakeParseIndexer := new(FakeParseIndexer)
					fakeParseIndexer.ParseIndexReturns(0, nil)

					fakeValueConstructor := new(value.FakeConstructor)
					constructValueErr := errors.New("constructValueErr")
					fakeValueConstructor.ConstructReturns(nil, errors.New("constructValueErr"))

					expectedErr := fmt.Errorf("unable to interpret item; error was %v", constructValueErr.Error())

					objectUnderTest := _interpreter{
						parseIndexer:     fakeParseIndexer,
						valueConstructor: fakeValueConstructor,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						"dummyIndexString",
						nonEmptyArrayValue,
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("valueConstructor.Construct doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */

					fakeParseIndexer := new(FakeParseIndexer)
					fakeParseIndexer.ParseIndexReturns(0, nil)

					fakeValueConstructor := new(value.FakeConstructor)
					constructedValue := model.Value{String: new(string)}
					fakeValueConstructor.ConstructReturns(&constructedValue, nil)

					objectUnderTest := _interpreter{
						parseIndexer:     fakeParseIndexer,
						valueConstructor: fakeValueConstructor,
					}

					/* act */
					actualItemValue, actualErr := objectUnderTest.Interpret(
						"dummyIndexString",
						nonEmptyArrayValue,
					)

					/* assert */
					Expect(*actualItemValue).To(Equal(constructedValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
})
