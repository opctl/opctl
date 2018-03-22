package expression

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression/interpolater"
	"reflect"
)

var _ = Context("evalArrayInitializerer", func() {
	Context("Eval", func() {
		It("should call json.Marshal w/ expected args", func() {
			/* arrange */
			providedExpression := []interface{}{
				"item1",
			}

			fakeJSON := new(ijson.Fake)
			// err to trigger immediate return
			fakeJSON.MarshalReturns([]byte{}, errors.New("dummyError"))

			arrayUnderTest := _evalArrayInitializerer{
				json: fakeJSON,
			}

			/* act */
			arrayUnderTest.Eval(
				providedExpression,
				map[string]*model.Value{},
				new(data.FakeHandle),
			)

			/* assert */
			actualExpression := fakeJSON.MarshalArgsForCall(0)

			Expect(actualExpression).To(Equal(providedExpression))
		})
		Context("json.Marshal errs", func() {
			It("should return expected result", func() {
				/* arrange */
				providedExpression := []interface{}{
					"item1",
				}

				fakeJSON := new(ijson.Fake)
				marshalErr := errors.New("marshalErr")
				// err to trigger immediate return
				fakeJSON.MarshalReturns([]byte{}, marshalErr)

				expectedErr := fmt.Errorf(
					"unable to eval %+v as arrayInitializer; error was %v",
					providedExpression,
					marshalErr,
				)

				arrayUnderTest := _evalArrayInitializerer{
					json: fakeJSON,
				}

				/* act */
				_, actualErr := arrayUnderTest.Eval(
					providedExpression,
					map[string]*model.Value{},
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))

			})
		})
		Context("json.Marshal doesn't err", func() {
			It("should call interpolater.Interpolate w/ expected args", func() {

				/* arrange */
				providedScope := map[string]*model.Value{"dummyName": {}}
				providedOpRef := new(data.FakeHandle)

				fakeJSON := new(ijson.Fake)
				expectedExpression := []byte{2, 3, 4, 1}
				fakeJSON.MarshalReturns(expectedExpression, nil)

				fakeInterpolater := new(interpolater.Fake)
				// err to trigger immediate return
				fakeInterpolater.InterpolateReturns("", errors.New("dummyError"))

				arrayUnderTest := _evalArrayInitializerer{
					json:         fakeJSON,
					interpolater: fakeInterpolater,
				}

				/* act */
				arrayUnderTest.Eval(
					[]interface{}{},
					providedScope,
					providedOpRef,
				)

				/* assert */
				actualExpression,
					actualScope,
					actualOpRef := fakeInterpolater.InterpolateArgsForCall(0)

				Expect(actualExpression).To(Equal(string(expectedExpression)))
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualOpRef).To(Equal(providedOpRef))

			})
			Context("interpolater.Interpolate errs", func() {
				It("should return expected result", func() {

					/* arrange */
					providedExpression := []interface{}{
						"item1",
					}

					fakeInterpolater := new(interpolater.Fake)
					interpolateErr := errors.New("interpolateErr")
					// err to trigger immediate return
					fakeInterpolater.InterpolateReturns("", interpolateErr)

					expectedErr := fmt.Errorf(
						"unable to eval %+v as arrayInitializer; error was %v",
						providedExpression,
						interpolateErr,
					)

					arrayUnderTest := _evalArrayInitializerer{
						json:         new(ijson.Fake),
						interpolater: fakeInterpolater,
					}

					/* act */
					_, actualErr := arrayUnderTest.Eval(
						providedExpression,
						map[string]*model.Value{},
						new(data.FakeHandle),
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))

				})
			})
			Context("interpolater.Interpolate doesn't err", func() {
				It("should call json.Unmarshal w/ expected args", func() {

					/* arrange */
					providedScope := map[string]*model.Value{"dummyName": {}}
					providedOpRef := new(data.FakeHandle)

					fakeJSON := new(ijson.Fake)

					fakeInterpolater := new(interpolater.Fake)
					expectedString := "dummyCoerce"
					// err to trigger immediate return
					fakeInterpolater.InterpolateReturns(expectedString, nil)

					arrayUnderTest := _evalArrayInitializerer{
						json:         fakeJSON,
						interpolater: fakeInterpolater,
					}

					/* act */
					arrayUnderTest.Eval(
						[]interface{}{},
						providedScope,
						providedOpRef,
					)

					/* assert */
					actualBytes, _ := fakeJSON.UnmarshalArgsForCall(0)

					Expect(string(actualBytes)).To(Equal(string(expectedString)))

				})
				Context("json.Unmarshal errs", func() {
					It("should return expected result", func() {

						/* arrange */
						providedExpression := []interface{}{
							"item1",
						}

						fakeJSON := new(ijson.Fake)
						unmarshalErr := errors.New("unmarshalErr")
						// err to trigger immediate return
						fakeJSON.UnmarshalReturns(unmarshalErr)

						expectedErr := fmt.Errorf(
							"unable to eval %+v as arrayInitializer; error was %v",
							providedExpression,
							unmarshalErr,
						)

						arrayUnderTest := _evalArrayInitializerer{
							json:         fakeJSON,
							interpolater: new(interpolater.Fake),
						}

						/* act */
						_, actualErr := arrayUnderTest.Eval(
							providedExpression,
							map[string]*model.Value{},
							new(data.FakeHandle),
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
				Context("json.Unmarshal doesn't err", func() {
					It("should return expected result", func() {

						/* arrange */
						fakeJSON := new(ijson.Fake)

						expectedValue := []interface{}{"item1"}

						fakeJSON.UnmarshalStub = func(data []byte, v interface{}) error {
							rs := reflect.ValueOf(expectedValue)
							rd := reflect.ValueOf(v).Elem()

							rt := reflect.MakeSlice(rs.Type(), 1, 1)
							reflect.Copy(rt, rs)
							rd.Set(rt)
							return nil
						}

						arrayUnderTest := _evalArrayInitializerer{
							json:         fakeJSON,
							interpolater: new(interpolater.Fake),
						}

						/* act */
						actualValue, actualErr := arrayUnderTest.Eval(
							[]interface{}{},
							map[string]*model.Value{},
							new(data.FakeHandle),
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
