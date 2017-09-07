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

var _ = Context("toNumber", func() {
	Context("ToNumber", func() {
		Context("Value is nil", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _number{}

				/* act */
				actualNumber, actualErr := objectUnderTest.ToNumber(nil)

				/* assert */
				Expect(actualNumber).To(Equal(float64(0)))
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

				objectUnderTest := _number{}

				/* act */
				actualNumber, actualErr := objectUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(actualNumber).To(Equal(float64(0)))
				Expect(actualErr).To(Equal(errors.New("Unable to coerce dir to number; incompatible types")))
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

				fileUnderTest := _number{
					ioUtil: fakeIOUtil,
				}

				/* act */
				fileUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(fakeIOUtil.ReadFileArgsForCall(0)).To(Equal(providedFile))
			})
			Context("ioutil.ReadFile errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshalErr := errors.New("dummyError")
					fakeIOUtil.ReadFileReturns(nil, marshalErr)

					fileUnderTest := _number{
						ioUtil: fakeIOUtil,
					}

					/* act */
					_, actualErr := fileUnderTest.ToNumber(
						&model.Value{File: new(string)},
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("Unable to coerce file to number; error was %v", marshalErr.Error())))
				})
			})
			Context("ioutil.ReadFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshaledBytes := []byte("2")
					fakeIOUtil.ReadFileReturns(marshaledBytes, nil)

					expectedNumber, err := strconv.ParseFloat(string(marshaledBytes), 64)
					if nil != err {
						panic(err)
					}

					fileUnderTest := _number{
						ioUtil: fakeIOUtil,
					}

					/* act */
					actualNumber, actualErr := fileUnderTest.ToNumber(
						&model.Value{File: new(string)},
					)

					/* assert */
					Expect(actualNumber).To(Equal(expectedNumber))
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

				objectUnderTest := _number{}

				/* act */
				actualNumber, actualErr := objectUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(actualNumber).To(Equal(providedNumber))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Object isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValue := &model.Value{
					Object: map[string]interface{}{},
				}

				objectUnderTest := _number{}

				/* act */
				actualNumber, actualErr := objectUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(actualNumber).To(Equal(float64(0)))
				Expect(actualErr).To(Equal(errors.New("Unable to coerce object to number; incompatible types")))
			})
		})
		Context("Value.Number isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedNumber := 2.2
				providedValue := &model.Value{
					Number: &providedNumber,
				}

				objectUnderTest := _number{}

				/* act */
				actualNumber, actualErr := objectUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(actualNumber).To(Equal(providedNumber))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Dir,File,Number,Object,Number nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValue := &model.Value{}

				objectUnderTest := _number{}

				/* act */
				actualNumber, actualErr := objectUnderTest.ToNumber(providedValue)

				/* assert */
				Expect(actualNumber).To(Equal(float64(0)))
				Expect(actualErr).To(Equal(fmt.Errorf("Unable to coerce '%#v' to number", providedValue)))
			})
		})
	})
})
