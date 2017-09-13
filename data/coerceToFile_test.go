package data

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
	"os"
	"path/filepath"
	"strconv"
)

var _ = Context("coerceToFile", func() {
	Context("Coerce", func() {
		Context("Value is nil", func() {
			It("should call ioutil.WriteFile w/ expected args", func() {
				/* arrange */
				providedScratchDir := "dummyScratchDir"
				uniqueString := "dummyUniqueString"

				fakeUniqueString := new(uniquestring.Fake)
				fakeUniqueString.ConstructReturns(uniqueString)

				fakeIOUtil := new(iioutil.Fake)

				objectUnderTest := _coerceToFile{
					uniqueString: fakeUniqueString,
					ioUtil:       fakeIOUtil,
				}

				/* act */
				objectUnderTest.CoerceToFile(nil, providedScratchDir)

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
					fakeIOUtil := new(iioutil.Fake)

					writeFileErr := errors.New("dummyError")
					fakeIOUtil.WriteFileReturns(writeFileErr)

					objectUnderTest := _coerceToFile{
						ioUtil:       fakeIOUtil,
						uniqueString: new(uniquestring.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.CoerceToFile(
						nil,
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce nil to file; error was %v", writeFileErr.Error())))
				})
			})
			Context("ioutil.WriteFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedScratchDir := "dummyScratchDir"
					uniqueString := "dummyUniqueString"

					expectedValuePath := filepath.Join(providedScratchDir, uniqueString)
					expectedValue := model.Value{File: &expectedValuePath}

					fakeUniqueString := new(uniquestring.Fake)
					fakeUniqueString.ConstructReturns(uniqueString)

					objectUnderTest := _coerceToFile{
						uniqueString: fakeUniqueString,
						ioUtil:       new(iioutil.Fake),
					}

					/* act */
					actualValue, actualErr := objectUnderTest.CoerceToFile(
						nil,
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

				objectUnderTest := _coerceToFile{}

				/* act */
				actualValue, actualErr := objectUnderTest.CoerceToFile(providedValue, providedScratchDir)

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

				objectUnderTest := _coerceToFile{}

				/* act */
				actualValue, actualErr := objectUnderTest.CoerceToFile(providedValue, providedScratchDir)

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

				fakeUniqueString := new(uniquestring.Fake)
				fakeUniqueString.ConstructReturns(uniqueString)

				fakeIOUtil := new(iioutil.Fake)

				objectUnderTest := _coerceToFile{
					uniqueString: fakeUniqueString,
					ioUtil:       fakeIOUtil,
				}

				/* act */
				objectUnderTest.CoerceToFile(providedValue, providedScratchDir)

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
					fakeIOUtil := new(iioutil.Fake)

					writeFileErr := errors.New("dummyError")
					fakeIOUtil.WriteFileReturns(writeFileErr)

					objectUnderTest := _coerceToFile{
						ioUtil:       fakeIOUtil,
						uniqueString: new(uniquestring.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.CoerceToFile(
						&model.Value{
							Number: new(float64),
						},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce number to file; error was %v", writeFileErr.Error())))
				})
			})
			Context("ioutil.WriteFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedScratchDir := "dummyScratchDir"
					uniqueString := "dummyUniqueString"

					fakeUniqueString := new(uniquestring.Fake)
					fakeUniqueString.ConstructReturns(uniqueString)

					expectedValuePath := filepath.Join(providedScratchDir, uniqueString)
					expectedValue := model.Value{File: &expectedValuePath}

					objectUnderTest := _coerceToFile{
						uniqueString: fakeUniqueString,
						ioUtil:       new(iioutil.Fake),
					}

					/* act */
					actualValue, actualErr := objectUnderTest.CoerceToFile(
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
				providedObject := map[string]interface{}{
					"dummyName": "dummyValue",
				}

				providedValue := &model.Value{
					Object: providedObject,
				}

				fakeJSON := new(ijson.Fake)
				// err to trigger immediate return
				fakeJSON.MarshalReturns(nil, errors.New("dummyError"))

				objectUnderTest := _coerceToFile{
					json:         fakeJSON,
					uniqueString: new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.CoerceToFile(providedValue, "dummyScratchDir")

				/* assert */
				Expect(fakeJSON.MarshalArgsForCall(0)).To(Equal(providedObject))
			})
			Context("json.Marshal errs", func() {
				It("should return expected result", func() {
					/* arrange */
					fakeJSON := new(ijson.Fake)

					marshalErr := errors.New("dummyError")
					fakeJSON.MarshalReturns(nil, marshalErr)

					objectUnderTest := _coerceToFile{
						json:         fakeJSON,
						uniqueString: new(uniquestring.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.CoerceToFile(
						&model.Value{Object: map[string]interface{}{"": ""}},
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

					fakeUniqueString := new(uniquestring.Fake)
					fakeUniqueString.ConstructReturns(uniqueString)

					fakeIOUtil := new(iioutil.Fake)

					objectUnderTest := _coerceToFile{
						json:         fakeJSON,
						uniqueString: fakeUniqueString,
						ioUtil:       fakeIOUtil,
					}

					/* act */
					objectUnderTest.CoerceToFile(
						&model.Value{Object: map[string]interface{}{"": ""}},
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
						fakeIOUtil := new(iioutil.Fake)

						writeFileErr := errors.New("dummyError")
						fakeIOUtil.WriteFileReturns(writeFileErr)

						objectUnderTest := _coerceToFile{
							json:         new(ijson.Fake),
							ioUtil:       fakeIOUtil,
							uniqueString: new(uniquestring.Fake),
						}

						/* act */
						_, actualErr := objectUnderTest.CoerceToFile(
							&model.Value{Object: map[string]interface{}{"": ""}},
							"dummyScratchDir",
						)

						/* assert */
						Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce object to file; error was %v", writeFileErr.Error())))
					})
				})
				Context("ioutil.WriteFile doesn't err", func() {
					It("should return expected result", func() {
						/* arrange */
						providedScratchDir := "dummyScratchDir"
						uniqueString := "dummyUniqueString"

						fakeUniqueString := new(uniquestring.Fake)
						fakeUniqueString.ConstructReturns(uniqueString)

						expectedValuePath := filepath.Join(providedScratchDir, uniqueString)
						expectedValue := model.Value{File: &expectedValuePath}

						objectUnderTest := _coerceToFile{
							json:         new(ijson.Fake),
							uniqueString: fakeUniqueString,
							ioUtil:       new(iioutil.Fake),
						}

						/* act */
						actualValue, actualErr := objectUnderTest.CoerceToFile(
							&model.Value{
								Object: map[string]interface{}{"": ""},
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

				fakeUniqueString := new(uniquestring.Fake)
				fakeUniqueString.ConstructReturns(uniqueString)

				fakeIOUtil := new(iioutil.Fake)

				objectUnderTest := _coerceToFile{
					uniqueString: fakeUniqueString,
					ioUtil:       fakeIOUtil,
				}

				/* act */
				objectUnderTest.CoerceToFile(providedValue, providedScratchDir)

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
					fakeIOUtil := new(iioutil.Fake)

					writeFileErr := errors.New("dummyError")
					fakeIOUtil.WriteFileReturns(writeFileErr)

					objectUnderTest := _coerceToFile{
						ioUtil:       fakeIOUtil,
						uniqueString: new(uniquestring.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.CoerceToFile(
						&model.Value{
							String: new(string),
						},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce string to file; error was %v", writeFileErr.Error())))
				})
			})
			Context("ioutil.WriteFile doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedScratchDir := "dummyScratchDir"
					uniqueString := "dummyUniqueString"

					fakeUniqueString := new(uniquestring.Fake)
					fakeUniqueString.ConstructReturns(uniqueString)

					expectedValuePath := filepath.Join(providedScratchDir, uniqueString)
					expectedValue := model.Value{File: &expectedValuePath}

					objectUnderTest := _coerceToFile{
						uniqueString: fakeUniqueString,
						ioUtil:       new(iioutil.Fake),
					}

					/* act */
					actualValue, actualErr := objectUnderTest.CoerceToFile(
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

				objectUnderTest := _coerceToFile{}

				/* act */
				actualValue, actualErr := objectUnderTest.CoerceToFile(providedValue, providedScratchDir)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(fmt.Errorf("unable to coerce '%#v' to file", providedValue)))
			})
		})
	})
})
