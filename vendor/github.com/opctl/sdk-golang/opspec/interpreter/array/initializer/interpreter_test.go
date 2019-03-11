package initializer

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/gopkg.in-yaml.v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater"
	"reflect"
)

var _ = Context("Interpreter", func() {
	Context("Interpret", func() {
		It("should call yaml.Marshal w/ expected args", func() {
			/* arrange */
			providedExpression := []interface{}{
				"item1",
			}

			fakeYAML := new(iyaml.Fake)
			// err to trigger immediate return
			fakeYAML.MarshalReturns([]byte{}, errors.New("dummyError"))

			objectUnderTest := _interpreter{
				yaml: fakeYAML,
			}

			/* act */
			objectUnderTest.Interpret(
				providedExpression,
				map[string]*model.Value{},
				new(data.FakeHandle),
			)

			/* assert */
			actualExpression := fakeYAML.MarshalArgsForCall(0)

			Expect(actualExpression).To(Equal(providedExpression))
		})
		Context("yaml.Marshal errs", func() {
			It("should return expected result", func() {
				/* arrange */
				providedExpression := []interface{}{
					"item1",
				}

				fakeYAML := new(iyaml.Fake)
				marshalErr := errors.New("marshalErr")
				// err to trigger immediate return
				fakeYAML.MarshalReturns([]byte{}, marshalErr)

				expectedErr := fmt.Errorf(
					"unable to interpret %+v as array initializer; error was %v",
					providedExpression,
					marshalErr,
				)

				objectUnderTest := _interpreter{
					yaml: fakeYAML,
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					providedExpression,
					map[string]*model.Value{},
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))

			})
		})
		Context("yaml.Marshal doesn't err", func() {
			It("should call interpolater.Interpolate w/ expected args", func() {

				/* arrange */
				providedScope := map[string]*model.Value{"dummyName": {}}
				providedOpRef := new(data.FakeHandle)

				fakeYAML := new(iyaml.Fake)
				expectedExpression := []byte{2, 3, 4, 1}
				fakeYAML.MarshalReturns(expectedExpression, nil)

				fakeInterpolater := new(interpolater.Fake)
				// err to trigger immediate return
				fakeInterpolater.InterpolateReturns("", errors.New("dummyError"))

				objectUnderTest := _interpreter{
					yaml:         fakeYAML,
					interpolater: fakeInterpolater,
				}

				/* act */
				objectUnderTest.Interpret(
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
						"unable to interpret %+v as array initializer; error was %v",
						providedExpression,
						interpolateErr,
					)

					objectUnderTest := _interpreter{
						yaml:         new(iyaml.Fake),
						interpolater: fakeInterpolater,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						providedExpression,
						map[string]*model.Value{},
						new(data.FakeHandle),
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))

				})
			})
			Context("interpolater.Interpolate doesn't err", func() {
				It("should call yaml.Unmarshal w/ expected args", func() {

					/* arrange */
					providedScope := map[string]*model.Value{"dummyName": {}}
					providedOpRef := new(data.FakeHandle)

					fakeYAML := new(iyaml.Fake)

					fakeInterpolater := new(interpolater.Fake)
					expectedString := "dummyCoerce"
					// err to trigger immediate return
					fakeInterpolater.InterpolateReturns(expectedString, nil)

					objectUnderTest := _interpreter{
						yaml:         fakeYAML,
						interpolater: fakeInterpolater,
					}

					/* act */
					objectUnderTest.Interpret(
						[]interface{}{},
						providedScope,
						providedOpRef,
					)

					/* assert */
					actualBytes, _ := fakeYAML.UnmarshalArgsForCall(0)

					Expect(string(actualBytes)).To(Equal(string(expectedString)))

				})
				Context("yaml.Unmarshal errs", func() {
					It("should return expected result", func() {

						/* arrange */
						providedExpression := []interface{}{
							"item1",
						}

						fakeYAML := new(iyaml.Fake)
						unmarshalErr := errors.New("unmarshalErr")
						// err to trigger immediate return
						fakeYAML.UnmarshalReturns(unmarshalErr)

						expectedErr := fmt.Errorf(
							"unable to interpret %+v as array initializer; error was %v",
							providedExpression,
							unmarshalErr,
						)

						objectUnderTest := _interpreter{
							yaml:         fakeYAML,
							interpolater: new(interpolater.Fake),
						}

						/* act */
						_, actualErr := objectUnderTest.Interpret(
							providedExpression,
							map[string]*model.Value{},
							new(data.FakeHandle),
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
				Context("yaml.Unmarshal doesn't err", func() {
					It("should return expected result", func() {

						/* arrange */
						fakeYAML := new(iyaml.Fake)

						expectedValue := []interface{}{"item1"}

						fakeYAML.UnmarshalStub = func(data []byte, v interface{}) error {
							rs := reflect.ValueOf(expectedValue)
							rd := reflect.ValueOf(v).Elem()

							rt := reflect.MakeSlice(rs.Type(), 1, 1)
							reflect.Copy(rt, rs)
							rd.Set(rt)
							return nil
						}

						objectUnderTest := _interpreter{
							yaml:         fakeYAML,
							interpolater: new(interpolater.Fake),
						}

						/* act */
						actualValue, actualErr := objectUnderTest.Interpret(
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
