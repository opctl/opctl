package coerce

import (
	"errors"
	"fmt"
	"strconv"

	ijson "github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdk/go/model"
)

var _ = Context("toStringer", func() {
	Context("ToString", func() {
		Context("Value is nil", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _toStringer{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToString(nil)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{String: new(string)}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Array isn't nil", func() {
			It("should call json.Marshal w/ expected args", func() {
				/* arrange */
				providedArray := &[]interface{}{"dummyItem"}

				providedValue := &model.Value{
					Array: providedArray,
				}

				fakeJSON := new(ijson.Fake)
				// err to trigger immediate return
				fakeJSON.MarshalReturns(nil, errors.New("dummyError"))

				arrayUnderTest := _toStringer{
					json: fakeJSON,
				}

				/* act */
				arrayUnderTest.ToString(providedValue)

				/* assert */
				Expect(fakeJSON.MarshalArgsForCall(0)).To(Equal(providedArray))
			})
			Context("json.Marshal errs", func() {
				It("should return expected result", func() {
					/* arrange */

					fakeJSON := new(ijson.Fake)

					marshalErr := errors.New("dummyError")
					fakeJSON.MarshalReturns(nil, marshalErr)

					arrayUnderTest := _toStringer{
						json: fakeJSON,
					}

					/* act */
					actualValue, actualErr := arrayUnderTest.ToString(
						&model.Value{Array: new([]interface{})},
					)

					/* assert */
					Expect(actualValue).To(BeNil())
					Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce array to string; error was %v", marshalErr.Error())))
				})
			})
			Context("json.Marshal doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeJSON := new(ijson.Fake)

					marshaledBytes := []byte{2, 3, 4}
					fakeJSON.MarshalReturns(marshaledBytes, nil)

					marshaledString := string(marshaledBytes)
					expectedValue := model.Value{String: &marshaledString}

					arrayUnderTest := _toStringer{
						json: fakeJSON,
					}

					/* act */
					actualValue, actualErr := arrayUnderTest.ToString(
						&model.Value{Array: new([]interface{})},
					)

					/* assert */
					Expect(*actualValue).To(Equal(expectedValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.Boolean isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedBoolean := true
				providedValue := &model.Value{
					Boolean: &providedBoolean,
				}

				booleanString := strconv.FormatBool(providedBoolean)
				expectedValue := model.Value{String: &booleanString}

				objectUnderTest := _toStringer{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToString(providedValue)

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
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

				objectUnderTest := _toStringer{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToString(providedValue)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce dir '%v' to string; incompatible types", providedDir)))
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

				fileUnderTest := _toStringer{
					ioUtil: fakeIOUtil,
				}

				/* act */
				fileUnderTest.ToString(providedValue)

				/* assert */
				Expect(fakeIOUtil.ReadFileArgsForCall(0)).To(Equal(providedFile))
			})
			Context("ioutil.ReadFile errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshalErr := errors.New("dummyError")
					fakeIOUtil.ReadFileReturns(nil, marshalErr)

					fileUnderTest := _toStringer{
						ioUtil: fakeIOUtil,
					}

					/* act */
					actualValue, actualErr := fileUnderTest.ToString(
						&model.Value{File: new(string)},
					)

					/* assert */
					Expect(actualValue).To(BeNil())
					Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce file to string; error was %v", marshalErr.Error())))
				})
			})
			Context("ioutil.ReadFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshaledBytes := []byte{2, 3, 4}
					fakeIOUtil.ReadFileReturns(marshaledBytes, nil)

					marshaledString := string(marshaledBytes)

					expectedValue := model.Value{String: &marshaledString}

					fileUnderTest := _toStringer{
						ioUtil: fakeIOUtil,
					}

					/* act */
					actualValue, actualErr := fileUnderTest.ToString(
						&model.Value{File: new(string)},
					)

					/* assert */
					Expect(*actualValue).To(Equal(expectedValue))
					Expect(actualErr).To(BeNil())
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

				numberString := strconv.FormatFloat(providedNumber, 'f', -1, 64)
				expectedValue := model.Value{String: &numberString}

				objectUnderTest := _toStringer{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToString(providedValue)

				/* assert */
				Expect(*actualValue).To(Equal(expectedValue))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Object isn't nil", func() {
			It("should call json.Marshal w/ expected args", func() {
				/* arrange */
				providedObject := &map[string]interface{}{
					"dummyName": "dummyValue",
				}

				providedValue := &model.Value{
					Object: providedObject,
				}

				fakeJSON := new(ijson.Fake)
				// err to trigger immediate return
				fakeJSON.MarshalReturns(nil, errors.New("dummyError"))

				objectUnderTest := _toStringer{
					json: fakeJSON,
				}

				/* act */
				objectUnderTest.ToString(providedValue)

				/* assert */
				Expect(fakeJSON.MarshalArgsForCall(0)).To(Equal(providedObject))
			})
			Context("json.Marshal errs", func() {
				It("should return expected result", func() {
					/* arrange */

					fakeJSON := new(ijson.Fake)

					marshalErr := errors.New("dummyError")
					fakeJSON.MarshalReturns(nil, marshalErr)

					objectUnderTest := _toStringer{
						json: fakeJSON,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.ToString(
						&model.Value{Object: new(map[string]interface{})},
					)

					/* assert */
					Expect(actualValue).To(BeNil())
					Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce object to string; error was %v", marshalErr.Error())))
				})
			})
			Context("json.Marshal doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeJSON := new(ijson.Fake)

					marshaledBytes := []byte{2, 3, 4}
					fakeJSON.MarshalReturns(marshaledBytes, nil)

					marshaledString := string(marshaledBytes)
					expectedValue := model.Value{String: &marshaledString}

					objectUnderTest := _toStringer{
						json: fakeJSON,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.ToString(
						&model.Value{Object: new(map[string]interface{})},
					)

					/* assert */
					Expect(*actualValue).To(Equal(expectedValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.String isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedString := "dummyValue"
				providedValue := model.Value{
					String: &providedString,
				}

				objectUnderTest := _toStringer{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToString(&providedValue)

				/* assert */
				Expect(*actualValue).To(Equal(providedValue))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Dir,File,Number,Object,String nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValue := &model.Value{}

				objectUnderTest := _toStringer{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToString(providedValue)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce '%+v' to string", providedValue)))
			})
		})
	})
})
