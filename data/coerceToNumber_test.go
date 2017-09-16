package data

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"strconv"
)

var _ = Context("coerceToNumber", func() {
	Context("CoerceToNumber", func() {
		Context("Value is nil", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _coerceToNumber{}

				/* act */
				actualNumber, actualErr := objectUnderTest.CoerceToNumber(nil)

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

				objectUnderTest := _coerceToNumber{}

				/* act */
				actualNumber, actualErr := objectUnderTest.CoerceToNumber(providedValue)

				/* assert */
				Expect(actualNumber).To(Equal(float64(0)))
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

				fileUnderTest := _coerceToNumber{
					ioUtil: fakeIOUtil,
				}

				/* act */
				fileUnderTest.CoerceToNumber(providedValue)

				/* assert */
				Expect(fakeIOUtil.ReadFileArgsForCall(0)).To(Equal(providedFile))
			})
			Context("ioutil.ReadFile errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshalErr := errors.New("dummyError")
					fakeIOUtil.ReadFileReturns(nil, marshalErr)

					fileUnderTest := _coerceToNumber{
						ioUtil: fakeIOUtil,
					}

					/* act */
					_, actualErr := fileUnderTest.CoerceToNumber(
						&model.Value{File: new(string)},
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce file to number; error was %v", marshalErr.Error())))
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

					fileUnderTest := _coerceToNumber{
						ioUtil: fakeIOUtil,
					}

					/* act */
					actualNumber, actualErr := fileUnderTest.CoerceToNumber(
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

				objectUnderTest := _coerceToNumber{}

				/* act */
				actualNumber, actualErr := objectUnderTest.CoerceToNumber(providedValue)

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

				objectUnderTest := _coerceToNumber{}

				/* act */
				actualNumber, actualErr := objectUnderTest.CoerceToNumber(providedValue)

				/* assert */
				Expect(actualNumber).To(Equal(float64(0)))
				Expect(actualErr).To(Equal(errors.New("unable to coerce object to number; incompatible types")))
			})
		})
		Context("Value.Number isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedNumber := 2.2
				providedValue := &model.Value{
					Number: &providedNumber,
				}

				objectUnderTest := _coerceToNumber{}

				/* act */
				actualNumber, actualErr := objectUnderTest.CoerceToNumber(providedValue)

				/* assert */
				Expect(actualNumber).To(Equal(providedNumber))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Dir,File,Number,Object,Number nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValue := &model.Value{}

				objectUnderTest := _coerceToNumber{}

				/* act */
				actualNumber, actualErr := objectUnderTest.CoerceToNumber(providedValue)

				/* assert */
				Expect(actualNumber).To(Equal(float64(0)))
				Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce '%+v' to number", providedValue)))
			})
		})
	})
})
