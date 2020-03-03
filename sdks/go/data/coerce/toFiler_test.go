package coerce

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uniquestringFakes "github.com/opctl/opctl/sdks/go/internal/uniquestring/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	"os"
	"path/filepath"
	"strconv"
)

var _ = Context("toFile", func() {
	Context("ToFile", func() {
		Context("Value is nil", func() {
			It("should call ioutil.WriteFile w/ expected args", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"
				uniqueString := "dummyUniqueString"

				fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
				fakeUniqueString.ConstructReturns(uniqueString, nil)

				fakeIOUtil := new(iioutil.Fake)

				objectUnderTest := _toFiler{
					uniqueString: fakeUniqueString,
					ioUtil:       fakeIOUtil,
				}

				/* act */
				objectUnderTest.ToFile(nil, providedScratchDir)

				/* assert */
				actualPath,
					actualData,
					actualPerms := fakeIOUtil.WriteFileArgsForCall(0)

				Expect(actualPath).To(Equal(filepath.Join(providedScratchDir, uniqueString)))
				Expect(actualData).To(BeEmpty())
				Expect(actualPerms).To(Equal(os.FileMode(0666)))
			})
			Context("ioutil.WriteFile errs", func() {
				It("should return expected result", func() {
					/* arrange */
					var providedValue *model.Value

					fakeIOUtil := new(iioutil.Fake)

					writeFileErr := errors.New("dummyError")
					fakeIOUtil.WriteFileReturns(writeFileErr)

					expectedErr := fmt.Errorf(
						"unable to coerce '%+v' to file; error was %v",
						providedValue,
						writeFileErr.Error(),
					)

					objectUnderTest := _toFiler{
						ioUtil:       fakeIOUtil,
						uniqueString: new(uniquestringFakes.FakeUniqueStringFactory),
					}

					/* act */
					_, actualErr := objectUnderTest.ToFile(
						providedValue,
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("ioutil.WriteFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedScratchDir := "dummyScratchDir"
					uniqueString := "dummyUniqueString"

					expectedValuePath := filepath.Join(providedScratchDir, uniqueString)
					expectedValue := model.Value{File: &expectedValuePath}

					fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
					fakeUniqueString.ConstructReturns(uniqueString, nil)

					objectUnderTest := _toFiler{
						uniqueString: fakeUniqueString,
						ioUtil:       new(iioutil.Fake),
					}

					/* act */
					actualValue, actualErr := objectUnderTest.ToFile(
						nil,
						providedScratchDir,
					)

					/* assert */
					Expect(*actualValue).To(Equal(expectedValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.Array isn't nil", func() {
			It("should call json.Marshal w/ expected args", func() {
				/* arrange */
				providedArray := &[]interface{}{"arrayItem"}

				providedValue := &model.Value{
					Array: providedArray,
				}

				expectedArray, _ := providedValue.Unbox()

				fakeJSON := new(ijson.Fake)
				// err to trigger immediate return
				fakeJSON.MarshalReturns(nil, errors.New("dummyError"))

				arrayUnderTest := _toFiler{
					json:         fakeJSON,
					uniqueString: new(uniquestringFakes.FakeUniqueStringFactory),
				}

				/* act */
				arrayUnderTest.ToFile(providedValue, "dummyScratchDir")

				/* assert */
				Expect(fakeJSON.MarshalArgsForCall(0)).To(Equal(expectedArray))
			})
			Context("json.Marshal errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeJSON := new(ijson.Fake)

					marshalErr := errors.New("dummyError")
					fakeJSON.MarshalReturns(nil, marshalErr)

					arrayUnderTest := _toFiler{
						json:         fakeJSON,
						uniqueString: new(uniquestringFakes.FakeUniqueStringFactory),
					}

					/* act */
					_, actualErr := arrayUnderTest.ToFile(
						&model.Value{Array: new([]interface{})},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce array to file; error was %v", marshalErr.Error())))
				})
			})
			Context("json.Marshal doesn't err", func() {
				It("should call ioutil.WriteFile w/ expected args", func() {
					/* arrange */
					providedScratchDir := "dummyScratchDir"

					fakeJSON := new(ijson.Fake)

					marshaledBytes := []byte{2, 3, 4}
					fakeJSON.MarshalReturns(marshaledBytes, nil)

					uniqueString := "dummyUniqueString"

					fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
					fakeUniqueString.ConstructReturns(uniqueString, nil)

					fakeIOUtil := new(iioutil.Fake)

					arrayUnderTest := _toFiler{
						json:         fakeJSON,
						uniqueString: fakeUniqueString,
						ioUtil:       fakeIOUtil,
					}

					/* act */
					arrayUnderTest.ToFile(
						&model.Value{Array: new([]interface{})},
						providedScratchDir,
					)

					/* assert */
					actualPath,
						actualData,
						actualPerms := fakeIOUtil.WriteFileArgsForCall(0)

					Expect(actualPath).To(Equal(filepath.Join(providedScratchDir, uniqueString)))
					Expect(actualData).To(Equal(marshaledBytes))
					Expect(actualPerms).To(Equal(os.FileMode(0666)))
				})
				Context("ioutil.WriteFile errs", func() {
					It("should return expected result", func() {
						/* arrange */
						providedValue := &model.Value{Array: new([]interface{})}

						fakeIOUtil := new(iioutil.Fake)

						writeFileErr := errors.New("dummyError")
						fakeIOUtil.WriteFileReturns(writeFileErr)

						expectedErr := fmt.Errorf(
							"unable to coerce '%+v' to file; error was %v",
							providedValue,
							writeFileErr.Error(),
						)

						arrayUnderTest := _toFiler{
							json:         new(ijson.Fake),
							ioUtil:       fakeIOUtil,
							uniqueString: new(uniquestringFakes.FakeUniqueStringFactory),
						}

						/* act */
						_, actualErr := arrayUnderTest.ToFile(
							providedValue,
							"dummyScratchDir",
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
				Context("ioutil.WriteFile doesn't err", func() {
					It("should return expected result", func() {
						/* arrange */
						providedScratchDir := "dummyScratchDir"
						uniqueString := "dummyUniqueString"

						fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
						fakeUniqueString.ConstructReturns(uniqueString, nil)

						expectedValuePath := filepath.Join(providedScratchDir, uniqueString)
						expectedValue := model.Value{File: &expectedValuePath}

						arrayUnderTest := _toFiler{
							json:         new(ijson.Fake),
							uniqueString: fakeUniqueString,
							ioUtil:       new(iioutil.Fake),
						}

						/* act */
						actualValue, actualErr := arrayUnderTest.ToFile(
							&model.Value{
								Array: new([]interface{}),
							},
							providedScratchDir,
						)

						/* assert */
						Expect(*actualValue).To(Equal(expectedValue))
						Expect(actualErr).To(BeNil())
					})
				})
			})
		})
		Context("Value.Boolean isn't nil", func() {
			It("should call ioutil.WriteFile w/ expected args", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"

				providedBoolean := true
				providedValue := &model.Value{
					Boolean: &providedBoolean,
				}

				uniqueString := "dummyUniqueString"

				fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
				fakeUniqueString.ConstructReturns(uniqueString, nil)

				fakeIOUtil := new(iioutil.Fake)

				objectUnderTest := _toFiler{
					uniqueString: fakeUniqueString,
					ioUtil:       fakeIOUtil,
				}

				/* act */
				objectUnderTest.ToFile(providedValue, providedScratchDir)

				/* assert */
				actualPath,
					actualData,
					actualPerms := fakeIOUtil.WriteFileArgsForCall(0)

				Expect(actualPath).To(Equal(filepath.Join(providedScratchDir, uniqueString)))
				Expect(actualData).To(Equal([]byte(strconv.FormatBool(providedBoolean))))
				Expect(actualPerms).To(Equal(os.FileMode(0666)))
			})
			Context("ioutil.WriteFile errs", func() {
				It("should return expected result", func() {
					/* arrange */
					providedValue := &model.Value{
						Boolean: new(bool),
					}

					fakeIOUtil := new(iioutil.Fake)

					writeFileErr := errors.New("dummyError")
					fakeIOUtil.WriteFileReturns(writeFileErr)

					expectedErr := fmt.Errorf(
						"unable to coerce '%+v' to file; error was %v",
						providedValue,
						writeFileErr.Error(),
					)

					objectUnderTest := _toFiler{
						ioUtil:       fakeIOUtil,
						uniqueString: new(uniquestringFakes.FakeUniqueStringFactory),
					}

					/* act */
					_, actualErr := objectUnderTest.ToFile(
						providedValue,
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("ioutil.WriteFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedScratchDir := "dummyScratchDir"
					uniqueString := "dummyUniqueString"

					fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
					fakeUniqueString.ConstructReturns(uniqueString, nil)

					expectedValuePath := filepath.Join(providedScratchDir, uniqueString)
					expectedValue := model.Value{File: &expectedValuePath}

					objectUnderTest := _toFiler{
						uniqueString: fakeUniqueString,
						ioUtil:       new(iioutil.Fake),
					}

					/* act */
					actualValue, actualErr := objectUnderTest.ToFile(
						&model.Value{
							Boolean: new(bool),
						},
						providedScratchDir,
					)

					/* assert */
					Expect(*actualValue).To(Equal(expectedValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.Dir isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"

				providedDir := "dummyValue"
				providedValue := &model.Value{
					Dir: &providedDir,
				}

				objectUnderTest := _toFiler{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToFile(providedValue, providedScratchDir)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce dir '%v' to file; incompatible types", providedDir)))
			})
		})
		Context("Value.File isn't nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"

				providedFile := "dummyFile"
				providedValue := &model.Value{
					File: &providedFile,
				}

				objectUnderTest := _toFiler{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToFile(providedValue, providedScratchDir)

				/* assert */
				Expect(actualValue).To(Equal(providedValue))
				Expect(actualErr).To(BeNil())
			})
		})
		Context("Value.Number isn't nil", func() {
			It("should call ioutil.WriteFile w/ expected args", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"

				providedNumber := 2.2
				providedValue := &model.Value{
					Number: &providedNumber,
				}

				uniqueString := "dummyUniqueString"

				fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
				fakeUniqueString.ConstructReturns(uniqueString, nil)

				fakeIOUtil := new(iioutil.Fake)

				objectUnderTest := _toFiler{
					uniqueString: fakeUniqueString,
					ioUtil:       fakeIOUtil,
				}

				/* act */
				objectUnderTest.ToFile(providedValue, providedScratchDir)

				/* assert */
				actualPath,
					actualData,
					actualPerms := fakeIOUtil.WriteFileArgsForCall(0)

				Expect(actualPath).To(Equal(filepath.Join(providedScratchDir, uniqueString)))
				Expect(actualData).To(Equal([]byte(strconv.FormatFloat(providedNumber, 'f', -1, 64))))
				Expect(actualPerms).To(Equal(os.FileMode(0666)))
			})
			Context("ioutil.WriteFile errs", func() {
				It("should return expected result", func() {
					/* arrange */
					providedValue := &model.Value{
						Number: new(float64),
					}

					fakeIOUtil := new(iioutil.Fake)

					writeFileErr := errors.New("dummyError")
					fakeIOUtil.WriteFileReturns(writeFileErr)

					expectedErr := fmt.Errorf(
						"unable to coerce '%+v' to file; error was %v",
						providedValue,
						writeFileErr.Error(),
					)

					objectUnderTest := _toFiler{
						ioUtil:       fakeIOUtil,
						uniqueString: new(uniquestringFakes.FakeUniqueStringFactory),
					}

					/* act */
					_, actualErr := objectUnderTest.ToFile(
						providedValue,
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("ioutil.WriteFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedScratchDir := "dummyScratchDir"
					uniqueString := "dummyUniqueString"

					fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
					fakeUniqueString.ConstructReturns(uniqueString, nil)

					expectedValuePath := filepath.Join(providedScratchDir, uniqueString)
					expectedValue := model.Value{File: &expectedValuePath}

					objectUnderTest := _toFiler{
						uniqueString: fakeUniqueString,
						ioUtil:       new(iioutil.Fake),
					}

					/* act */
					actualValue, actualErr := objectUnderTest.ToFile(
						&model.Value{
							Number: new(float64),
						},
						providedScratchDir,
					)

					/* assert */
					Expect(*actualValue).To(Equal(expectedValue))
					Expect(actualErr).To(BeNil())
				})
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

				expectedObject, _ := providedValue.Unbox()

				fakeJSON := new(ijson.Fake)
				// err to trigger immediate return
				fakeJSON.MarshalReturns(nil, errors.New("dummyError"))

				objectUnderTest := _toFiler{
					json:         fakeJSON,
					uniqueString: new(uniquestringFakes.FakeUniqueStringFactory),
				}

				/* act */
				objectUnderTest.ToFile(providedValue, "dummyScratchDir")

				/* assert */
				Expect(fakeJSON.MarshalArgsForCall(0)).To(Equal(expectedObject))
			})
			Context("json.Marshal errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeJSON := new(ijson.Fake)

					marshalErr := errors.New("dummyError")
					fakeJSON.MarshalReturns(nil, marshalErr)

					objectUnderTest := _toFiler{
						json:         fakeJSON,
						uniqueString: new(uniquestringFakes.FakeUniqueStringFactory),
					}

					/* act */
					_, actualErr := objectUnderTest.ToFile(
						&model.Value{Object: new(map[string]interface{})},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce object to file; error was %v", marshalErr.Error())))
				})
			})
			Context("json.Marshal doesn't err", func() {
				It("should call ioutil.WriteFile w/ expected args", func() {
					/* arrange */
					providedScratchDir := "dummyScratchDir"

					fakeJSON := new(ijson.Fake)

					marshaledBytes := []byte{2, 3, 4}
					fakeJSON.MarshalReturns(marshaledBytes, nil)

					uniqueString := "dummyUniqueString"

					fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
					fakeUniqueString.ConstructReturns(uniqueString, nil)

					fakeIOUtil := new(iioutil.Fake)

					objectUnderTest := _toFiler{
						json:         fakeJSON,
						uniqueString: fakeUniqueString,
						ioUtil:       fakeIOUtil,
					}

					/* act */
					objectUnderTest.ToFile(
						&model.Value{Object: new(map[string]interface{})},
						providedScratchDir,
					)

					/* assert */
					actualPath,
						actualData,
						actualPerms := fakeIOUtil.WriteFileArgsForCall(0)

					Expect(actualPath).To(Equal(filepath.Join(providedScratchDir, uniqueString)))
					Expect(actualData).To(Equal(marshaledBytes))
					Expect(actualPerms).To(Equal(os.FileMode(0666)))
				})
				Context("ioutil.WriteFile errs", func() {
					It("should return expected result", func() {
						/* arrange */
						providedValue := &model.Value{Object: new(map[string]interface{})}

						fakeIOUtil := new(iioutil.Fake)

						writeFileErr := errors.New("dummyError")
						fakeIOUtil.WriteFileReturns(writeFileErr)

						expectedErr := fmt.Errorf(
							"unable to coerce '%+v' to file; error was %v",
							providedValue,
							writeFileErr.Error(),
						)

						objectUnderTest := _toFiler{
							json:         new(ijson.Fake),
							ioUtil:       fakeIOUtil,
							uniqueString: new(uniquestringFakes.FakeUniqueStringFactory),
						}

						/* act */
						_, actualErr := objectUnderTest.ToFile(
							providedValue,
							"dummyScratchDir",
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
				Context("ioutil.WriteFile doesn't err", func() {
					It("should return expected result", func() {
						/* arrange */
						providedScratchDir := "dummyScratchDir"
						uniqueString := "dummyUniqueString"

						fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
						fakeUniqueString.ConstructReturns(uniqueString, nil)

						expectedValuePath := filepath.Join(providedScratchDir, uniqueString)
						expectedValue := model.Value{File: &expectedValuePath}

						objectUnderTest := _toFiler{
							json:         new(ijson.Fake),
							uniqueString: fakeUniqueString,
							ioUtil:       new(iioutil.Fake),
						}

						/* act */
						actualValue, actualErr := objectUnderTest.ToFile(
							&model.Value{
								Object: new(map[string]interface{}),
							},
							providedScratchDir,
						)

						/* assert */
						Expect(*actualValue).To(Equal(expectedValue))
						Expect(actualErr).To(BeNil())
					})
				})
			})
		})
		Context("Value.String isn't nil", func() {
			It("should call ioutil.WriteFile w/ expected args", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"

				providedString := "dummyString"
				providedValue := &model.Value{
					String: &providedString,
				}

				uniqueString := "dummyUniqueString"

				fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
				fakeUniqueString.ConstructReturns(uniqueString, nil)

				fakeIOUtil := new(iioutil.Fake)

				objectUnderTest := _toFiler{
					uniqueString: fakeUniqueString,
					ioUtil:       fakeIOUtil,
				}

				/* act */
				objectUnderTest.ToFile(providedValue, providedScratchDir)

				/* assert */
				actualPath,
					actualData,
					actualPerms := fakeIOUtil.WriteFileArgsForCall(0)

				Expect(actualPath).To(Equal(filepath.Join(providedScratchDir, uniqueString)))
				Expect(actualData).To(Equal([]byte(providedString)))
				Expect(actualPerms).To(Equal(os.FileMode(0666)))
			})
			Context("ioutil.WriteFile errs", func() {
				It("should return expected result", func() {
					/* arrange */
					providedValue := &model.Value{
						String: new(string),
					}

					fakeIOUtil := new(iioutil.Fake)

					writeFileErr := errors.New("dummyError")
					fakeIOUtil.WriteFileReturns(writeFileErr)

					expectedErr := fmt.Errorf(
						"unable to coerce '%+v' to file; error was %v",
						providedValue,
						writeFileErr.Error(),
					)

					objectUnderTest := _toFiler{
						ioUtil:       fakeIOUtil,
						uniqueString: new(uniquestringFakes.FakeUniqueStringFactory),
					}

					/* act */
					_, actualErr := objectUnderTest.ToFile(
						providedValue,
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("ioutil.WriteFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedScratchDir := "dummyScratchDir"
					uniqueString := "dummyUniqueString"

					fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
					fakeUniqueString.ConstructReturns(uniqueString, nil)

					expectedValuePath := filepath.Join(providedScratchDir, uniqueString)
					expectedValue := model.Value{File: &expectedValuePath}

					objectUnderTest := _toFiler{
						uniqueString: fakeUniqueString,
						ioUtil:       new(iioutil.Fake),
					}

					/* act */
					actualValue, actualErr := objectUnderTest.ToFile(
						&model.Value{
							String: new(string),
						},
						providedScratchDir,
					)

					/* assert */
					Expect(*actualValue).To(Equal(expectedValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
		Context("Value.Dir,File,Number,Object,String nil", func() {
			It("should return expected result", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"

				providedValue := &model.Value{}

				objectUnderTest := _toFiler{}

				/* act */
				actualValue, actualErr := objectUnderTest.ToFile(providedValue, providedScratchDir)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce '%+v' to file", providedValue)))
			})
		})
		Context("scratchDir doesn't exist", func() {
			It("should call os.MkdirAll w/ expected args", func() {
				/* arrange */
				providedValue := &model.Value{
					String: new(string),
				}
				providedScratchDir := "dummyScratchDir"

				fakeIOUtil := new(iioutil.Fake)
				fakeIOUtil.WriteFileReturns(os.ErrNotExist)

				fakeOS := new(ios.Fake)
				fakeOS.MkdirAllReturns(errors.New("dummyError"))

				fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
				fakeUniqueString.ConstructReturns("dummyUniqueString", nil)

				objectUnderTest := _toFiler{
					ioUtil:       fakeIOUtil,
					os:           fakeOS,
					uniqueString: fakeUniqueString,
				}

				/* act */
				objectUnderTest.ToFile(
					providedValue,
					"dummyScratchDir",
				)

				/* assert */
				actualPath,
					actualPerm := fakeOS.MkdirAllArgsForCall(0)

				Expect(actualPath).To(Equal(filepath.Join(providedScratchDir)))
				Expect(actualPerm).To(Equal(os.FileMode(0777)))
			})
			Context("os.MkdirAll errs", func() {
				It("should return expected result", func() {
					/* arrange */
					providedValue := &model.Value{
						String: new(string),
					}

					fakeIOUtil := new(iioutil.Fake)
					fakeIOUtil.WriteFileReturns(os.ErrNotExist)

					fakeOS := new(ios.Fake)
					mkdirAllError := errors.New("dummyError")
					fakeOS.MkdirAllReturns(mkdirAllError)

					expectedErr := fmt.Errorf(
						"unable to coerce '%+v' to file; error was %v",
						providedValue,
						mkdirAllError.Error(),
					)

					objectUnderTest := _toFiler{
						ioUtil:       fakeIOUtil,
						os:           fakeOS,
						uniqueString: new(uniquestringFakes.FakeUniqueStringFactory),
					}

					/* act */
					_, actualErr := objectUnderTest.ToFile(
						providedValue,
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("os.MkdirAll doesn't err", func() {
				It("should call ioutil.WriteFile w/ expected args", func() {
					/* arrange */
					providedScratchDir := "dummyScratchDir"

					providedString := "dummyString"
					providedValue := &model.Value{
						String: &providedString,
					}

					uniqueString := "dummyUniqueString"

					fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
					fakeUniqueString.ConstructReturns(uniqueString, nil)

					fakeIOUtil := new(iioutil.Fake)
					fakeIOUtil.WriteFileReturnsOnCall(0, os.ErrNotExist)

					objectUnderTest := _toFiler{
						uniqueString: fakeUniqueString,
						os:           new(ios.Fake),
						ioUtil:       fakeIOUtil,
					}

					/* act */
					objectUnderTest.ToFile(providedValue, providedScratchDir)

					/* assert */
					actualPath,
						actualData,
						actualPerms := fakeIOUtil.WriteFileArgsForCall(0)

					Expect(fakeIOUtil.WriteFileCallCount()).To(Equal(2))
					Expect(actualPath).To(Equal(filepath.Join(providedScratchDir, uniqueString)))
					Expect(actualData).To(Equal([]byte(providedString)))
					Expect(actualPerms).To(Equal(os.FileMode(0666)))
				})
				Context("ioutil.WriteFile errs", func() {
					It("should return expected result", func() {
						/* arrange */
						providedValue := &model.Value{
							String: new(string),
						}

						fakeIOUtil := new(iioutil.Fake)

						fakeIOUtil.WriteFileReturnsOnCall(0, os.ErrNotExist)
						writeFileErr := errors.New("dummyError")
						fakeIOUtil.WriteFileReturnsOnCall(1, writeFileErr)

						expectedErr := fmt.Errorf(
							"unable to coerce '%+v' to file; error was %v",
							providedValue,
							writeFileErr.Error(),
						)

						objectUnderTest := _toFiler{
							ioUtil:       fakeIOUtil,
							os:           new(ios.Fake),
							uniqueString: new(uniquestringFakes.FakeUniqueStringFactory),
						}

						/* act */
						_, actualErr := objectUnderTest.ToFile(
							providedValue,
							"dummyScratchDir",
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
				Context("ioutil.WriteFile doesn't err", func() {
					It("should return expected result", func() {
						/* arrange */
						providedScratchDir := "dummyScratchDir"
						uniqueString := "dummyUniqueString"

						fakeIOUtil := new(iioutil.Fake)
						fakeIOUtil.WriteFileReturnsOnCall(0, os.ErrNotExist)

						fakeUniqueString := new(uniquestringFakes.FakeUniqueStringFactory)
						fakeUniqueString.ConstructReturns(uniqueString, nil)

						expectedValuePath := filepath.Join(providedScratchDir, uniqueString)
						expectedValue := model.Value{File: &expectedValuePath}

						objectUnderTest := _toFiler{
							uniqueString: fakeUniqueString,
							os:           new(ios.Fake),
							ioUtil:       fakeIOUtil,
						}

						/* act */
						actualValue, actualErr := objectUnderTest.ToFile(
							&model.Value{
								String: new(string),
							},
							providedScratchDir,
						)

						/* assert */
						Expect(*actualValue).To(Equal(expectedValue))
						Expect(actualErr).To(BeNil())
					})
				})
			})
		})
	})
})
