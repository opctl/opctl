package expression

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/expression/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"reflect"
)

var _ = Context("evalObjectInitializerer", func() {
	Context("Eval", func() {
		Context("expression contains shorthand property", func() {
			It("should call json.Marshal w/ expected args", func() {
				/* arrange */
				shortHandPropName := "prop1Name"
				providedExpression := map[string]interface{}{
					shortHandPropName: nil,
				}

				expectedExpression := map[string]interface{}{
					shortHandPropName: fmt.Sprintf(
						"%v%v%v",
						interpolater.RefStart,
						shortHandPropName,
						interpolater.RefEnd),
				}

				fakeJSON := new(ijson.Fake)
				// err to trigger immediate return
				fakeJSON.MarshalReturns([]byte{}, errors.New("dummyError"))

				objectUnderTest := _evalObjectInitializerer{
					json: fakeJSON,
				}

				/* act */
				objectUnderTest.Eval(
					providedExpression,
					map[string]*model.Value{},
					new(pkg.FakeHandle),
				)

				/* assert */
				actualExpression := fakeJSON.MarshalArgsForCall(0)

				Expect(actualExpression).To(Equal(expectedExpression))
			})
		})
		It("should call json.Marshal w/ expected args", func() {
			/* arrange */
			providedExpression := map[string]interface{}{
				"prop1Name": "prop1Value",
			}

			fakeJSON := new(ijson.Fake)
			// err to trigger immediate return
			fakeJSON.MarshalReturns([]byte{}, errors.New("dummyError"))

			objectUnderTest := _evalObjectInitializerer{
				json: fakeJSON,
			}

			/* act */
			objectUnderTest.Eval(
				providedExpression,
				map[string]*model.Value{},
				new(pkg.FakeHandle),
			)

			/* assert */
			actualExpression := fakeJSON.MarshalArgsForCall(0)

			Expect(actualExpression).To(Equal(providedExpression))
		})
		Context("json.Marshal errs", func() {
			It("should return expected result", func() {
				/* arrange */
				providedExpression := map[string]interface{}{
					"prop1Name": "prop1Value",
				}

				fakeJSON := new(ijson.Fake)
				marshalErr := errors.New("marshalErr")
				// err to trigger immediate return
				fakeJSON.MarshalReturns([]byte{}, marshalErr)

				expectedErr := fmt.Errorf(
					"unable to eval %+v as objectInitializer; error was %v",
					providedExpression,
					marshalErr,
				)

				objectUnderTest := _evalObjectInitializerer{
					json: fakeJSON,
				}

				/* act */
				_, actualErr := objectUnderTest.Eval(
					providedExpression,
					map[string]*model.Value{},
					new(pkg.FakeHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))

			})
		})
		Context("json.Marshal doesn't err", func() {
			It("should call interpolater.Interpolate w/ expected args", func() {

				/* arrange */
				providedScope := map[string]*model.Value{"dummyName": {}}
				providedPkgRef := new(pkg.FakeHandle)

				fakeJSON := new(ijson.Fake)
				expectedExpression := []byte{2, 3, 4, 1}
				fakeJSON.MarshalReturns(expectedExpression, nil)

				fakeInterpolater := new(interpolater.Fake)
				// err to trigger immediate return
				fakeInterpolater.InterpolateReturns("", errors.New("dummyError"))

				objectUnderTest := _evalObjectInitializerer{
					json:         fakeJSON,
					interpolater: fakeInterpolater,
				}

				/* act */
				objectUnderTest.Eval(
					map[string]interface{}{},
					providedScope,
					providedPkgRef,
				)

				/* assert */
				actualExpression,
					actualScope,
					actualPkgRef := fakeInterpolater.InterpolateArgsForCall(0)

				Expect(actualExpression).To(Equal(string(expectedExpression)))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualPkgRef).To(Equal(providedPkgRef))

			})
			Context("interpolater.Interpolate errs", func() {
				It("should return expected result", func() {

					/* arrange */
					providedExpression := map[string]interface{}{
						"prop1Name": "prop1Value",
					}

					fakeInterpolater := new(interpolater.Fake)
					interpolateErr := errors.New("interpolateErr")
					// err to trigger immediate return
					fakeInterpolater.InterpolateReturns("", interpolateErr)

					expectedErr := fmt.Errorf(
						"unable to eval %+v as objectInitializer; error was %v",
						providedExpression,
						interpolateErr,
					)

					objectUnderTest := _evalObjectInitializerer{
						json:         new(ijson.Fake),
						interpolater: fakeInterpolater,
					}

					/* act */
					_, actualErr := objectUnderTest.Eval(
						providedExpression,
						map[string]*model.Value{},
						new(pkg.FakeHandle),
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))

				})
			})
			Context("interpolater.Interpolate doesn't err", func() {
				It("should call json.Unmarshal w/ expected args", func() {

					/* arrange */
					providedScope := map[string]*model.Value{"dummyName": {}}
					providedPkgRef := new(pkg.FakeHandle)

					fakeJSON := new(ijson.Fake)

					fakeInterpolater := new(interpolater.Fake)
					expectedString := "dummyString"
					// err to trigger immediate return
					fakeInterpolater.InterpolateReturns(expectedString, nil)

					objectUnderTest := _evalObjectInitializerer{
						json:         fakeJSON,
						interpolater: fakeInterpolater,
					}

					/* act */
					objectUnderTest.Eval(
						map[string]interface{}{},
						providedScope,
						providedPkgRef,
					)

					/* assert */
					actualBytes, _ := fakeJSON.UnmarshalArgsForCall(0)

					Expect(string(actualBytes)).To(Equal(string(expectedString)))

				})
				Context("json.Unmarshal errs", func() {
					It("should return expected result", func() {

						/* arrange */
						providedExpression := map[string]interface{}{
							"prop1Name": "prop1Value",
						}

						fakeJSON := new(ijson.Fake)
						unmarshalErr := errors.New("unmarshalErr")
						// err to trigger immediate return
						fakeJSON.UnmarshalReturns(unmarshalErr)

						expectedErr := fmt.Errorf(
							"unable to eval %+v as objectInitializer; error was %v",
							providedExpression,
							unmarshalErr,
						)

						objectUnderTest := _evalObjectInitializerer{
							json:         fakeJSON,
							interpolater: new(interpolater.Fake),
						}

						/* act */
						_, actualErr := objectUnderTest.Eval(
							providedExpression,
							map[string]*model.Value{},
							new(pkg.FakeHandle),
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
				Context("json.Unmarshal doesn't err", func() {
					It("should return expected result", func() {

						/* arrange */
						providedExpression := map[string]interface{}{
							"prop1Name": "prop1Value",
						}

						fakeJSON := new(ijson.Fake)

						mapKey := "dummyMapKey"
						mapValue := "dummyMapValue"
						expectedValue := map[string]interface{}{mapKey: mapValue}

						fakeJSON.UnmarshalStub = func(data []byte, v interface{}) error {
							reflect.ValueOf(v).Elem().SetMapIndex(
								reflect.ValueOf(mapKey),
								reflect.ValueOf(mapValue),
							)
							return nil
						}

						objectUnderTest := _evalObjectInitializerer{
							json:         fakeJSON,
							interpolater: new(interpolater.Fake),
						}

						/* act */
						actualValue, actualErr := objectUnderTest.Eval(
							providedExpression,
							map[string]*model.Value{},
							new(pkg.FakeHandle),
						)

						/* assert */
						Expect(actualValue).To(Equal(expectedValue))
						Expect(actualErr).To(BeNil())

					})
				})
			})
		})
	})
})
