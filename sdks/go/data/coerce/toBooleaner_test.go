package coerce

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("toBooleaner", func() {
	Context("ToBoolean", func() {
		Context("Value is nil", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := _toBooleaner{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToBoolean(nil)

				/* assert */
				Expect(*actualValue).To(Equal(types.Value{Boolean: new(bool)}))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Array isn't nil", func() {
			Context("Array empty", func() {
				It("should return expected result", func() {
					/* arrange */
					objectUnderTest := _toBooleaner{}

					expectedBoolean := false

					/* act */
					actualValue, actualErr := objectUnderTest.ToBoolean(
						&types.Value{
							Array: new([]interface{}),
						},
					)

					/* assert */
					Expect(*actualValue).To(Equal(types.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
			Context("Array not empty", func() {
				It("should return expected result", func() {
					/* arrange */
					objectUnderTest := _toBooleaner{}
					array := &[]interface{}{
						"",
					}

					expectedBoolean := true

					/* act */
					actualValue, actualErr := objectUnderTest.ToBoolean(
						&types.Value{
							Array: array,
						},
					)

					/* assert */
					Expect(*actualValue).To(Equal(types.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.Boolean isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedBoolean := true
				providedValue := &types.Value{
					Boolean: &providedBoolean,
				}

				objectUnderTest := _toBooleaner{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToBoolean(providedValue)

				/* assert */
				Expect(actualValue).To(Equal(providedValue))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Dir isn't nil", func() {
			It("should call ioutil.ReadDir w/ expected args", func() {
				/* arrange */
				providedDir := "dummyDir"

				providedValue := &types.Value{
					Dir: &providedDir,
				}

				fakeIOUtil := new(iioutil.Fake)
				// err to trigger immediate return
				fakeIOUtil.ReadDirReturns(nil, errors.New("dummyError"))

				objectUnderTest := _toBooleaner{
					ioUtil: fakeIOUtil,
				}

				/* act */
				objectUnderTest.ToBoolean(providedValue)

				/* assert */
				Expect(fakeIOUtil.ReadDirArgsForCall(0)).To(Equal(providedDir))
			})
			Context("ioutil.ReadDir errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshalErr := errors.New("dummyError")
					fakeIOUtil.ReadDirReturns(nil, marshalErr)

					objectUnderTest := _toBooleaner{
						ioUtil: fakeIOUtil,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.ToBoolean(
						&types.Value{Dir: new(string)},
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

						objectUnderTest := _toBooleaner{
							ioUtil: fakeIOUtil,
						}

						/* act */
						actualValue, actualErr := objectUnderTest.ToBoolean(
							&types.Value{Dir: new(string)},
						)

						/* assert */
						Expect(*actualValue).To(Equal(types.Value{Boolean: &expectedBoolean}))
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

						objectUnderTest := _toBooleaner{
							ioUtil: fakeIOUtil,
						}

						/* act */
						actualValue, actualErr := objectUnderTest.ToBoolean(
							&types.Value{Dir: new(string)},
						)

						/* assert */
						Expect(*actualValue).To(Equal(types.Value{Boolean: &expectedBoolean}))
						Expect(actualErr).To(BeNil())
					})
				})
			})
		})
		Context("Value.File isn't nil", func() {
			It("should call ioutil.ReadFile w/ expected args", func() {
				/* arrange */
				providedFile := "dummyFile"

				providedValue := &types.Value{
					File: &providedFile,
				}

				fakeIOUtil := new(iioutil.Fake)
				// err to trigger immediate return
				fakeIOUtil.ReadFileReturns(nil, errors.New("dummyError"))

				objectUnderTest := _toBooleaner{
					ioUtil: fakeIOUtil,
				}

				/* act */
				objectUnderTest.ToBoolean(providedValue)

				/* assert */
				Expect(fakeIOUtil.ReadFileArgsForCall(0)).To(Equal(providedFile))
			})
			Context("ioutil.ReadFile errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)

					marshalErr := errors.New("dummyError")
					fakeIOUtil.ReadFileReturns(nil, marshalErr)

					objectUnderTest := _toBooleaner{
						ioUtil: fakeIOUtil,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.ToBoolean(
						&types.Value{File: new(string)},
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

						objectUnderTest := _toBooleaner{
							ioUtil: fakeIOUtil,
						}

						/* act */
						actualValue, actualErr := objectUnderTest.ToBoolean(
							&types.Value{File: new(string)},
						)

						/* assert */
						Expect(*actualValue).To(Equal(types.Value{Boolean: &expectedBoolean}))
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

						objectUnderTest := _toBooleaner{
							ioUtil: fakeIOUtil,
						}

						/* act */
						actualValue, actualErr := objectUnderTest.ToBoolean(
							&types.Value{File: new(string)},
						)

						/* assert */
						Expect(*actualValue).To(Equal(types.Value{Boolean: &expectedBoolean}))
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
					providedValue := &types.Value{
						Number: &providedNumber,
					}

					expectedBoolean := false

					objectUnderTest := _toBooleaner{}

					/* act */
					actualValue, actualErr := objectUnderTest.ToBoolean(providedValue)

					/* assert */
					Expect(*actualValue).To(Equal(types.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
			Context("Number != 0", func() {
				It("should return expected result", func() {
					/* arrange */
					providedNumber := 1.0
					providedValue := &types.Value{
						Number: &providedNumber,
					}

					expectedBoolean := true

					objectUnderTest := _toBooleaner{}

					/* act */
					actualValue, actualErr := objectUnderTest.ToBoolean(providedValue)

					/* assert */
					Expect(*actualValue).To(Equal(types.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.Object isn't nil", func() {
			Context("Object has no properties", func() {
				It("should return expected result", func() {
					/* arrange */
					objectUnderTest := _toBooleaner{}

					expectedBoolean := false

					/* act */
					actualValue, actualErr := objectUnderTest.ToBoolean(
						&types.Value{
							Object: new(map[string]interface{}),
						},
					)

					/* assert */
					Expect(*actualValue).To(Equal(types.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
			Context("Object has properties", func() {
				It("should return expected result", func() {
					/* arrange */
					objectUnderTest := _toBooleaner{}

					object := &map[string]interface{}{
						"dummyProp": nil,
					}

					expectedBoolean := true

					/* act */
					actualValue, actualErr := objectUnderTest.ToBoolean(
						&types.Value{
							Object: object,
						},
					)

					/* assert */
					Expect(*actualValue).To(Equal(types.Value{Boolean: &expectedBoolean}))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.Array,Boolean,Dir,File,Number,Object,Boolean nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedValue := &types.Value{}

				objectUnderTest := _toBooleaner{}

				/* act */
				actualBoolean, actualErr := objectUnderTest.ToBoolean(providedValue)

				/* assert */
				Expect(actualBoolean).To(BeNil())
				Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce '%+v' to boolean", providedValue)))
			})
		})
	})
})
