package coerce

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"reflect"
)

var _ = Context("toObject", func() {
	Context("Coerce", func() {
		Context("Value is nil", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _toObject{}

				/* act */
				actualObject, actualErr := objectUnderTest.ToObject(nil)

				/* assert */
				Expect(actualObject).To(BeNil())
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Dir isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedDir := "dummyValue"
				providedValue := &model.Value{
					Dir: &providedDir,
				}

				objectUnderTest := _toObject{}

				/* act */
				actualObject, actualErr := objectUnderTest.ToObject(providedValue)

				/* assert */
				Expect(actualObject).To(BeNil())
				Expect(actualErr).To(Equal(fmt.Errorf("Unable to coerce dir '%v' to object; incompatible types", providedDir)))
			})
		})
		Context("Value.File isn't nil", func() {
			It("should call ioutil.ReadFile w/ expected args", func() {
				/* arrange */
				providedFile := "dummyFile"

				providedValue := &model.Value{
					File: &providedFile,
				}

				fakeIOUtil := new(iioutil.Fake)
				// err to trigger immediate return
				fakeIOUtil.ReadFileReturns(nil, errors.New("dummyError"))

				fileUnderTest := _toObject{
					ioUtil: fakeIOUtil,
				}

				/* act */
				fileUnderTest.ToObject(providedValue)

				/* assert */
				Expect(fakeIOUtil.ReadFileArgsForCall(0)).To(Equal(providedFile))
			})
			Context("ioutil.ReadFile errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshalErr := errors.New("dummyError")
					fakeIOUtil.ReadFileReturns(nil, marshalErr)

					fileUnderTest := _toObject{
						ioUtil: fakeIOUtil,
					}

					/* act */
					_, actualErr := fileUnderTest.ToObject(
						&model.Value{File: new(string)},
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("Unable to coerce file to object; error was %v", marshalErr.Error())))
				})
			})
			Context("ioutil.ReadFile doesn't err", func() {
				It("should call json.Unmarshal w/ expected args", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshaledBytes := []byte{2, 3, 5}
					fakeIOUtil.ReadFileReturns(marshaledBytes, nil)

					fakeJSON := new(ijson.Fake)
					// err to trigger immediate return
					fakeJSON.UnmarshalReturns(errors.New("dummyError"))

					objectUnderTest := _toObject{
						ioUtil: fakeIOUtil,
						json:   fakeJSON,
					}

					/* act */
					objectUnderTest.ToObject(
						&model.Value{File: new(string)},
					)

					/* assert */
					actualBytes,
						_ := fakeJSON.UnmarshalArgsForCall(0)
					Expect(actualBytes).To(Equal(marshaledBytes))
				})
				Context("json.Unmarshal errs", func() {
					It("should return expected result", func() {
						/* arrange */

						fakeJSON := new(ijson.Fake)

						unmarshalError := errors.New("dummyError")
						fakeJSON.UnmarshalReturns(unmarshalError)

						objectUnderTest := _toObject{
							ioUtil: new(iioutil.Fake),
							json:   fakeJSON,
						}

						/* act */
						_, actualErr := objectUnderTest.ToObject(
							&model.Value{File: new(string)},
						)

						/* assert */
						Expect(actualErr).To(Equal(fmt.Errorf("Unable to coerce file to object; error was %v", unmarshalError.Error())))
					})
				})
				Context("json.Unmarshal doesn't err", func() {
					It("should return expected result", func() {
						/* arrange */
						fakeJSON := new(ijson.Fake)

						mapKey := "dummyMapKey"
						mapValue := "dummyMapValue"
						expectedObject := map[string]interface{}{
							mapKey: mapValue,
						}

						fakeJSON.UnmarshalStub = func(data []byte, v interface{}) error {
							reflect.ValueOf(v).Elem().SetMapIndex(
								reflect.ValueOf(mapKey),
								reflect.ValueOf(mapValue),
							)
							return nil
						}

						objectUnderTest := _toObject{
							ioUtil: new(iioutil.Fake),
							json:   fakeJSON,
						}

						/* act */
						actualObject, actualErr := objectUnderTest.ToObject(
							&model.Value{File: new(string)},
						)

						/* assert */
						Expect(actualObject).To(Equal(expectedObject))
						Expect(actualErr).To(BeNil())
					})
				})
			})
		})
		Context("Value.Number isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedNumber := 2.2
				providedValue := &model.Value{
					Number: &providedNumber,
				}

				objectUnderTest := _toObject{}

				/* act */
				actualObject, actualErr := objectUnderTest.ToObject(providedValue)

				/* assert */
				Expect(actualObject).To(BeNil())
				Expect(actualErr).To(Equal(fmt.Errorf("Unable to coerce number '%v' to object; incompatible types", providedNumber)))
			})
		})
		Context("Value.Object isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedObject := map[string]interface{}{
					"dummyName": "dummyValue",
				}
				providedValue := &model.Value{
					Object: providedObject,
				}

				objectUnderTest := _toObject{}

				/* act */
				actualObject, actualErr := objectUnderTest.ToObject(providedValue)

				/* assert */
				Expect(actualObject).To(Equal(providedObject))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.String isn't nil", func() {
			It("should call json.Unmarshal w/ expected args", func() {
				/* arrange */
				providedString := "{}"

				providedValue := &model.Value{
					String: &providedString,
				}

				fakeJSON := new(ijson.Fake)
				// err to trigger immediate return
				fakeJSON.UnmarshalReturns(errors.New("dummyError"))

				objectUnderTest := _toObject{
					json: fakeJSON,
				}

				/* act */
				objectUnderTest.ToObject(providedValue)

				/* assert */
				actualBytes,
					_ := fakeJSON.UnmarshalArgsForCall(0)
				Expect(actualBytes).To(Equal([]byte(providedString)))
			})
			Context("json.Unmarshal errs", func() {
				It("should return expected result", func() {
					/* arrange */

					fakeJSON := new(ijson.Fake)

					unmarshalError := errors.New("dummyError")
					fakeJSON.UnmarshalReturns(unmarshalError)

					objectUnderTest := _toObject{
						json: fakeJSON,
					}

					/* act */
					_, actualErr := objectUnderTest.ToObject(
						&model.Value{String: new(string)},
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("Unable to coerce string to object; error was %v", unmarshalError.Error())))
				})
			})
			Context("json.Unmarshal doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeJSON := new(ijson.Fake)

					mapKey := "dummyMapKey"
					mapValue := "dummyMapValue"
					expectedObject := map[string]interface{}{
						mapKey: mapValue,
					}

					fakeJSON.UnmarshalStub = func(data []byte, v interface{}) error {
						reflect.ValueOf(v).Elem().SetMapIndex(
							reflect.ValueOf(mapKey),
							reflect.ValueOf(mapValue),
						)
						return nil
					}

					objectUnderTest := _toObject{
						json: fakeJSON,
					}

					/* act */
					actualObject, actualErr := objectUnderTest.ToObject(
						&model.Value{String: new(string)},
					)

					/* assert */
					Expect(actualObject).To(Equal(expectedObject))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.Dir,File,Number,Object,String nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValue := &model.Value{}
				var expectedObject map[string]interface{}

				objectUnderTest := _toObject{}

				/* act */
				actualObject, actualErr := objectUnderTest.ToObject(providedValue)

				/* assert */
				Expect(actualObject).To(Equal(expectedObject))
				Expect(actualErr).To(Equal(fmt.Errorf("Unable to coerce '%#v' to object", providedValue)))
			})
		})
	})
})
