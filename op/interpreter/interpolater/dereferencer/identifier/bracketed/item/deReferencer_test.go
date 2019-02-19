package item

import (
	"errors"
	"fmt"

	"github.com/opctl/sdk-golang/op/interpreter/interpolater/dereferencer/identifier/value"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
)

var _ = Context("DeReferencer", func() {
	Context("NewDeReferencer", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewDeReferencer()).Should(Not(BeNil()))
		})
	})
	Context("DeReference", func() {
		It("should call parseIndexer.ParseIndex w/ expected args", func() {
			/* arrange */
			providedIndexString := "dummyIndexString"
			providedData := model.Value{Array: []interface{}{"dummyData"}}

			fakeParseIndexer := new(fakeParseIndexer)
			// err to trigger immediate return
			fakeParseIndexer.ParseIndexReturns(0, errors.New("dummyErr"))

			objectUnderTest := _deReferencer{
				parseIndexer: fakeParseIndexer,
			}

			/* act */
			objectUnderTest.DeReference(
				providedIndexString,
				providedData,
			)

			/* assert */
			actualIndexString,
				actualArray := fakeParseIndexer.ParseIndexArgsForCall(0)

			Expect(actualIndexString).To(Equal(providedIndexString))
			Expect(actualArray).To(Equal(providedData.Array))
		})
		Context("parseIndexer.ParseIndex errs", func() {
			It("should return expected result", func() {
				/* arrange */

				fakeParseIndexer := new(fakeParseIndexer)
				parseIndexErr := errors.New("dummyErr")
				fakeParseIndexer.ParseIndexReturns(0, parseIndexErr)

				expectedErr := fmt.Errorf("unable to deReference item; error was %v", parseIndexErr.Error())

				objectUnderTest := _deReferencer{
					parseIndexer: fakeParseIndexer,
				}

				/* act */
				_, actualErr := objectUnderTest.DeReference(
					"dummyIndexString",
					model.Value{},
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("parseIndexer.ParseIndex doesn't err", func() {

			It("should call valueConstructor.Construct w/ expected args", func() {
				/* arrange */
				providedData := model.Value{Array: []interface{}{"dummyData"}}

				fakeParseIndexer := new(fakeParseIndexer)
				parsedIndex := 0
				fakeParseIndexer.ParseIndexReturns(0, nil)

				fakeValueConstructor := new(value.FakeConstructor)
				// err to trigger immediate return
				fakeValueConstructor.ConstructReturns(nil, errors.New("dummyErr"))

				objectUnderTest := _deReferencer{
					parseIndexer:     fakeParseIndexer,
					valueConstructor: fakeValueConstructor,
				}

				/* act */
				objectUnderTest.DeReference(
					"dummyIndexString",
					providedData,
				)

				/* assert */
				actualValue := fakeValueConstructor.ConstructArgsForCall(0)

				Expect(actualValue).To(Equal(providedData.Array[parsedIndex]))
			})
			Context("valueConstructor.Construct errs", func() {

				It("should return expected result", func() {
					/* arrange */

					fakeParseIndexer := new(fakeParseIndexer)
					fakeParseIndexer.ParseIndexReturns(0, nil)

					fakeValueConstructor := new(value.FakeConstructor)
					constructValueErr := errors.New("constructValueErr")
					fakeValueConstructor.ConstructReturns(nil, errors.New("constructValueErr"))

					expectedErr := fmt.Errorf("unable to deReference item; error was %v", constructValueErr.Error())

					objectUnderTest := _deReferencer{
						parseIndexer:     fakeParseIndexer,
						valueConstructor: fakeValueConstructor,
					}

					/* act */
					_, actualErr := objectUnderTest.DeReference(
						"dummyIndexString",
						model.Value{Array: []interface{}{"dummyItem"}},
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("valueConstructor.Construct doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */

					fakeParseIndexer := new(fakeParseIndexer)
					fakeParseIndexer.ParseIndexReturns(0, nil)

					fakeValueConstructor := new(value.FakeConstructor)
					constructedValue := model.Value{String: new(string)}
					fakeValueConstructor.ConstructReturns(&constructedValue, nil)

					objectUnderTest := _deReferencer{
						parseIndexer:     fakeParseIndexer,
						valueConstructor: fakeValueConstructor,
					}

					/* act */
					actualItemValue, actualErr := objectUnderTest.DeReference(
						"dummyIndexString",
						model.Value{Array: []interface{}{"dummyItem"}},
					)

					/* assert */
					Expect(*actualItemValue).To(Equal(constructedValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
})
