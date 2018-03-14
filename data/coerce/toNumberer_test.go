package coerce

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"strconv"
)

var _ = Context("toNumberer", func() {
	Context("ToNumber", func() {
		Context("Value is nil", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _toNumberer{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToNumber(nil)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Number: new(float64)}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Array isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValue := &model.Value{
					Array: []interface{}{},
				}

				arrayUnderTest := _toNumberer{}

				/* act */
				actualValue, actualErr := arrayUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(errors.New("unable to coerce array to number; incompatible types")))
			})
		})
		Context("Value.Dir isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedDir := "dummyValue"
				providedValue := &model.Value{
					Dir: &providedDir,
				}

				objectUnderTest := _toNumberer{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(errors.New("unable to coerce dir to number; incompatible types")))
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

				objectUnderTest := _toNumberer{
					ioUtil: fakeIOUtil,
				}

				/* act */
				objectUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(fakeIOUtil.ReadFileArgsForCall(0)).To(Equal(providedFile))
			})
			Context("ioutil.ReadFile errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshalErr := errors.New("dummyError")
					fakeIOUtil.ReadFileReturns(nil, marshalErr)

					objectUnderTest := _toNumberer{
						ioUtil: fakeIOUtil,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.ToNumber(
						&model.Value{File: new(string)},
					)

					/* assert */
					Expect(actualValue).To(BeNil())
					Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce file to number; error was %v", marshalErr.Error())))
				})
			})
			Context("ioutil.ReadFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshaledBytes := []byte("2")
					fakeIOUtil.ReadFileReturns(marshaledBytes, nil)

					parsedNumber, err := strconv.ParseFloat(string(marshaledBytes), 64)
					if nil != err {
						panic(err)
					}

					objectUnderTest := _toNumberer{
						ioUtil: fakeIOUtil,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.ToNumber(
						&model.Value{File: new(string)},
					)

					/* assert */
					Expect(*actualValue).To(Equal(model.Value{Number: &parsedNumber}))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.Number isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedNumber := float64(2.2)
				providedValue := &model.Value{
					Number: &providedNumber,
				}

				objectUnderTest := _toNumberer{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(actualValue).To(Equal(providedValue))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Object isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValue := &model.Value{
					Object: map[string]interface{}{},
				}

				objectUnderTest := _toNumberer{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(errors.New("unable to coerce object to number; incompatible types")))
			})
		})
		Context("Value.String isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedString := "2.2"
				providedValue := &model.Value{
					String: &providedString,
				}

				parsedNumber, err := strconv.ParseFloat(providedString, 64)
				if nil != err {
					panic(err.Error)
				}

				objectUnderTest := _toNumberer{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Number: &parsedNumber}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Array,Value.Dir,File,Number,Object,Number nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValue := &model.Value{}

				objectUnderTest := _toNumberer{}

				/* act */
				actualNumber, actualErr := objectUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(actualNumber).To(BeNil())
				Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce '%+v' to number", providedValue)))
			})
		})
	})
})
