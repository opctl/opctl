package coerce

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"strconv"
)

var _ = Context("toString", func() {
	Context("Coerce", func() {
		Context("Value is nil", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _toString{}

				/* act */
				actualString, actualErr := objectUnderTest.ToString(nil)

				/* assert */
				Expect(actualString).To(Equal(""))
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

				objectUnderTest := _toString{}

				/* act */
				actualString, actualErr := objectUnderTest.ToString(providedValue)

				/* assert */
				Expect(actualString).To(Equal(""))
				Expect(actualErr).To(Equal(fmt.Errorf("Unable to coerce dir '%v' to string; incompatible types", providedDir)))
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

				fileUnderTest := _toString{
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

					fileUnderTest := _toString{
						ioUtil: fakeIOUtil,
					}

					/* act */
					_, actualErr := fileUnderTest.ToString(
						&model.Value{File: new(string)},
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("Unable to coerce file to string; error was %v", marshalErr.Error())))
				})
			})
			Context("ioutil.ReadFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshaledBytes := []byte{2, 3, 4}
					fakeIOUtil.ReadFileReturns(marshaledBytes, nil)

					fileUnderTest := _toString{
						ioUtil: fakeIOUtil,
					}

					/* act */
					actualString, actualErr := fileUnderTest.ToString(
						&model.Value{File: new(string)},
					)

					/* assert */
					Expect(actualString).To(Equal(string(marshaledBytes)))
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

				objectUnderTest := _toString{}

				/* act */
				actualString, actualErr := objectUnderTest.ToString(providedValue)

				/* assert */
				Expect(actualString).To(Equal(strconv.FormatFloat(providedNumber, 'f', -1, 64)))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Object isn't nil", func() {
			It("should call json.Marshal w/ expected args", func() {
				/* arrange */
				providedObject := map[string]interface{}{
					"dummyName": "dummyValue",
				}

				providedValue := &model.Value{
					Object: providedObject,
				}

				fakeJSON := new(ijson.Fake)
				// err to trigger immediate return
				fakeJSON.MarshalReturns(nil, errors.New("dummyError"))

				objectUnderTest := _toString{
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

					objectUnderTest := _toString{
						json: fakeJSON,
					}

					/* act */
					_, actualErr := objectUnderTest.ToString(
						&model.Value{Object: map[string]interface{}{"": ""}},
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("Unable to coerce object to string; error was %v", marshalErr.Error())))
				})
			})
			Context("json.Marshal doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeJSON := new(ijson.Fake)

					marshaledBytes := []byte{2, 3, 4}
					fakeJSON.MarshalReturns(marshaledBytes, nil)

					objectUnderTest := _toString{
						json: fakeJSON,
					}

					/* act */
					actualString, actualErr := objectUnderTest.ToString(
						&model.Value{Object: map[string]interface{}{"": ""}},
					)

					/* assert */
					Expect(actualString).To(Equal(string(marshaledBytes)))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.String isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedString := "dummyValue"
				providedValue := &model.Value{
					String: &providedString,
				}

				objectUnderTest := _toString{}

				/* act */
				actualString, actualErr := objectUnderTest.ToString(providedValue)

				/* assert */
				Expect(actualString).To(Equal(providedString))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Dir,File,Number,Object,String nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValue := &model.Value{}

				objectUnderTest := _toString{}

				/* act */
				actualString, actualErr := objectUnderTest.ToString(providedValue)

				/* assert */
				Expect(actualString).To(Equal(""))
				Expect(actualErr).To(Equal(fmt.Errorf("Unable to coerce '%v' to string", providedValue)))
			})
		})
	})
})
