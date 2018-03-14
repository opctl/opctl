package data

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("coerceToBoolean", func() {
	Context("CoerceToBoolean", func() {
		Context("Value is nil", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _coerceToBoolean{}

				/* act */
				actualValue, actualErr := objectUnderTest.CoerceToBoolean(nil)

				/* assert */
				Expect(*actualValue).To(Equal(model.Value{Boolean: new(bool)}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Array isn't nil", func() {
			Context("Array empty", func() {
				It("should return expected result", func() {
					/* arrange */
					objectUnderTest := _coerceToBoolean{}

					expectedBoolean := false

					/* act */
					actualValue, actualErr := objectUnderTest.CoerceToBoolean(
						&model.Value{
							Array: []interface{}{},
						},
					)

					/* assert */
					Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
			Context("Array not empty", func() {
				It("should return expected result", func() {
					/* arrange */
					objectUnderTest := _coerceToBoolean{}

					expectedBoolean := true

					/* act */
					actualValue, actualErr := objectUnderTest.CoerceToBoolean(
						&model.Value{
							Array: []interface{}{
								// include item so len != 0
								"",
							},
						},
					)

					/* assert */
					Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
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

				objectUnderTest := _coerceToBoolean{}

				/* act */
				actualValue, actualErr := objectUnderTest.CoerceToBoolean(providedValue)

				/* assert */
				Expect(actualValue).To(Equal(providedValue))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Dir isn't nil", func() {
			It("should call ioutil.ReadDir w/ expected args", func() {
				/* arrange */
				providedDir := "dummyDir"

				providedValue := &model.Value{
					Dir: &providedDir,
				}

				fakeIOUtil := new(iioutil.Fake)
				// err to trigger immediate return
				fakeIOUtil.ReadDirReturns(nil, errors.New("dummyError"))

				objectUnderTest := _coerceToBoolean{
					ioUtil: fakeIOUtil,
				}

				/* act */
				objectUnderTest.CoerceToBoolean(providedValue)

				/* assert */
				Expect(fakeIOUtil.ReadDirArgsForCall(0)).To(Equal(providedDir))
			})
			Context("ioutil.ReadDir errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshalErr := errors.New("dummyError")
					fakeIOUtil.ReadDirReturns(nil, marshalErr)

					objectUnderTest := _coerceToBoolean{
						ioUtil: fakeIOUtil,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.CoerceToBoolean(
						&model.Value{Dir: new(string)},
					)

					/* assert */
					Expect(actualValue).To(BeNil())
					Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce dir to boolean; error was %v", marshalErr.Error())))
				})
			})
			Context("ioutil.ReadDir doesn't err", func() {
				Context("Directory empty", func() {
					It("should return expected result", func() {
						/* arrange */
						fakeIOUtil := new(iioutil.Fake)

						fileInfos := []os.FileInfo{}
						fakeIOUtil.ReadDirReturns(fileInfos, nil)

						expectedBoolean := false

						objectUnderTest := _coerceToBoolean{
							ioUtil: fakeIOUtil,
						}

						/* act */
						actualValue, actualErr := objectUnderTest.CoerceToBoolean(
							&model.Value{Dir: new(string)},
						)

						/* assert */
						Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
						Expect(actualErr).To(BeNil())
					})
				})
				Context("Directory not empty", func() {
					wd, err := os.Getwd()
					if nil != err {
						panic(err)
					}
					It("should return expected result", func() {
						/* arrange */
						fakeIOUtil := new(iioutil.Fake)

						// no good way to fake FileInfo's
						fileInfos, err := ioutil.ReadDir(wd)
						if nil != err {
							panic(err)
						}
						fakeIOUtil.ReadDirReturns(fileInfos, nil)

						expectedBoolean := true

						objectUnderTest := _coerceToBoolean{
							ioUtil: fakeIOUtil,
						}

						/* act */
						actualValue, actualErr := objectUnderTest.CoerceToBoolean(
							&model.Value{Dir: new(string)},
						)

						/* assert */
						Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
						Expect(actualErr).To(BeNil())
					})
				})
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

				objectUnderTest := _coerceToBoolean{
					ioUtil: fakeIOUtil,
				}

				/* act */
				objectUnderTest.CoerceToBoolean(providedValue)

				/* assert */
				Expect(fakeIOUtil.ReadFileArgsForCall(0)).To(Equal(providedFile))
			})
			Context("ioutil.ReadFile errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshalErr := errors.New("dummyError")
					fakeIOUtil.ReadFileReturns(nil, marshalErr)

					objectUnderTest := _coerceToBoolean{
						ioUtil: fakeIOUtil,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.CoerceToBoolean(
						&model.Value{File: new(string)},
					)

					/* assert */
					Expect(actualValue).To(BeNil())
					Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce file to boolean; error was %v", marshalErr.Error())))
				})
			})
			Context("ioutil.ReadFile doesn't err", func() {
				Context("File content truthy", func() {
					It("should return expected result", func() {
						/* arrange */
						fakeIOUtil := new(iioutil.Fake)

						marshaledBytes := []byte("t")
						fakeIOUtil.ReadFileReturns(marshaledBytes, nil)

						expectedBoolean := true

						objectUnderTest := _coerceToBoolean{
							ioUtil: fakeIOUtil,
						}

						/* act */
						actualValue, actualErr := objectUnderTest.CoerceToBoolean(
							&model.Value{File: new(string)},
						)

						/* assert */
						Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
						Expect(actualErr).To(BeNil())
					})
				})
				Context("File content falsy", func() {
					It("should return expected result", func() {
						/* arrange */
						fakeIOUtil := new(iioutil.Fake)

						marshaledBytes := []byte("")
						fakeIOUtil.ReadFileReturns(marshaledBytes, nil)

						expectedBoolean := false

						objectUnderTest := _coerceToBoolean{
							ioUtil: fakeIOUtil,
						}

						/* act */
						actualValue, actualErr := objectUnderTest.CoerceToBoolean(
							&model.Value{File: new(string)},
						)

						/* assert */
						Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
						Expect(actualErr).To(BeNil())
					})
				})
			})
		})
		Context("Value.Number isn't nil", func() {
			Context("Number == 0", func() {
				It("should return expected result", func() {
					/* arrange */
					providedNumber := 0.0
					providedValue := &model.Value{
						Number: &providedNumber,
					}

					expectedBoolean := false

					objectUnderTest := _coerceToBoolean{}

					/* act */
					actualValue, actualErr := objectUnderTest.CoerceToBoolean(providedValue)

					/* assert */
					Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
			Context("Number != 0", func() {
				It("should return expected result", func() {
					/* arrange */
					providedNumber := 1.0
					providedValue := &model.Value{
						Number: &providedNumber,
					}

					expectedBoolean := true

					objectUnderTest := _coerceToBoolean{}

					/* act */
					actualValue, actualErr := objectUnderTest.CoerceToBoolean(providedValue)

					/* assert */
					Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.Object isn't nil", func() {
			Context("Object has no properties", func() {
				It("should return expected result", func() {
					/* arrange */
					objectUnderTest := _coerceToBoolean{}

					expectedBoolean := false

					/* act */
					actualValue, actualErr := objectUnderTest.CoerceToBoolean(
						&model.Value{
							Object: map[string]interface{}{},
						},
					)

					/* assert */
					Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
			Context("Object has properties", func() {
				It("should return expected result", func() {
					/* arrange */
					objectUnderTest := _coerceToBoolean{}

					expectedBoolean := true

					/* act */
					actualValue, actualErr := objectUnderTest.CoerceToBoolean(
						&model.Value{
							Object: map[string]interface{}{
								// include item so len != 0
								"dummyProp": nil,
							},
						},
					)

					/* assert */
					Expect(*actualValue).To(Equal(model.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.Array,Boolean,Dir,File,Number,Object,Boolean nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValue := &model.Value{}

				objectUnderTest := _coerceToBoolean{}

				/* act */
				actualBoolean, actualErr := objectUnderTest.CoerceToBoolean(providedValue)

				/* assert */
				Expect(actualBoolean).To(BeNil())
				Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce '%+v' to boolean", providedValue)))
			})
		})
	})
})
