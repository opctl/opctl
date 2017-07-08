package cliparamsatisfier

import (
	"encoding/json"
	"errors"
	"github.com/ghodss/yaml"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ijson"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"reflect"
)

var _ = Describe("inputSrcFactory", func() {
	Context("NewCLIPromptInputSrc()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(_InputSrcFactory{}.NewCliPromptInputSrc(nil)).To(Not(BeNil()))
		})
	})
	Context("NewEnvVarInputSrc()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(_InputSrcFactory{}.NewEnvVarInputSrc()).To(Not(BeNil()))
		})
	})
	Context("NewParamDefaultInputSrc()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(_InputSrcFactory{}.NewParamDefaultInputSrc(
				map[string]*model.Param{},
			)).To(Not(BeNil()))
		})
	})
	Context("NewSliceInputSrc()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(_InputSrcFactory{}.NewSliceInputSrc([]string{}, "")).To(Not(BeNil()))
		})
	})
	Context("NewYMLFileInputSrc", func() {
		It("should call os.Stat w/ expected args", func() {
			/* arrange */
			providedFilePath := "dummyFilePath"
			fakeOS := new(ios.Fake)
			// error to trigger immediate return
			fakeOS.StatReturns(nil, errors.New("dummyError"))

			objectUnderTest := _InputSrcFactory{
				os:   fakeOS,
				json: new(ijson.Fake),
			}

			/* act */
			objectUnderTest.NewYMLFileInputSrc(providedFilePath)

			/* assert */
			Expect(fakeOS.StatArgsForCall(0)).To(Equal(providedFilePath))
		})
		Context("os.Stat errs", func() {
			Context("IsNotExist", func() {

			})
			It("should return expected error", func() {
				/* arrange */
				fakeOS := new(ios.Fake)
				expectedErr := errors.New("dummyError")
				fakeOS.StatReturns(nil, expectedErr)

				objectUnderTest := _InputSrcFactory{
					os:   fakeOS,
					json: new(ijson.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.NewYMLFileInputSrc("dummyFilePath")

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("ioutil.ReadFile doesn't err", func() {
			It("should call os.Open w/ expected args", func() {
				/* arrange */
				providedFilePath := "dummyFilePath"

				fakeIOUtil := new(iioutil.Fake)
				// error to trigger immediate return
				fakeIOUtil.ReadFileReturns(nil, errors.New("dummyError"))

				objectUnderTest := _InputSrcFactory{
					os:     new(ios.Fake),
					ioutil: fakeIOUtil,
					json:   new(ijson.Fake),
				}

				/* act */
				objectUnderTest.NewYMLFileInputSrc(providedFilePath)

				/* assert */
				Expect(fakeIOUtil.ReadFileArgsForCall(0)).To(Equal(providedFilePath))
			})
			Context("ioutil.ReadFile errs", func() {
				It("should return expected error", func() {
					/* arrange */
					fakeIOUtil := new(iioutil.Fake)
					expectedErr := errors.New("dummyError")
					fakeIOUtil.ReadFileReturns(nil, expectedErr)

					objectUnderTest := _InputSrcFactory{
						os:     new(ios.Fake),
						ioutil: fakeIOUtil,
						json:   new(ijson.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.NewYMLFileInputSrc("dummyFilePath")

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("ioutil.ReadFile doesn't err", func() {
				It("should call json.Unmarshal w/ expected args", func() {
					/* arrange */
					ymlString := "someKey: someValue"
					ymlBytes := []byte(ymlString)
					jsonString, err := yaml.YAMLToJSON(ymlBytes)
					if nil != err {
						Fail(err.Error())
					}

					fakeIOUtil := new(iioutil.Fake)
					fakeIOUtil.ReadFileReturns(ymlBytes, nil)

					fakeJSON := new(ijson.Fake)
					// error to trigger immediate return
					fakeJSON.UnmarshalReturns(errors.New("dummyError"))

					objectUnderTest := _InputSrcFactory{
						os:     new(ios.Fake),
						ioutil: fakeIOUtil,
						json:   fakeJSON,
					}

					/* act */
					objectUnderTest.NewYMLFileInputSrc("dummyFilePath")

					/* assert */
					actualBytes, _ := fakeJSON.UnmarshalArgsForCall(0)
					Expect(actualBytes).To(Equal(jsonString))
				})
				Context("json.Unmarshal errs", func() {
					It("should return expected error", func() {
						/* arrange */
						fakeJSON := new(ijson.Fake)
						expectedErr := errors.New("dummyError")
						fakeJSON.UnmarshalReturns(expectedErr)

						objectUnderTest := _InputSrcFactory{
							os:     new(ios.Fake),
							ioutil: new(iioutil.Fake),
							json:   fakeJSON,
						}

						/* act */
						_, actualErr := objectUnderTest.NewYMLFileInputSrc("dummyFilePath")

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
				Context("json.Unmarshal doesn't err", func() {
					It("should return expected InputSrc", func() {
						/* arrange */

						input1Name := "input1Name"
						input1Value := json.RawMessage("input1Value")

						expectedMap := map[string]string{
							input1Name: string(input1Value),
						}

						fakeJSON := new(ijson.Fake)
						fakeJSON.UnmarshalStub = func(data []byte, v interface{}) error {
							reflect.ValueOf(v).Elem().SetMapIndex(
								reflect.ValueOf(input1Name),
								reflect.ValueOf(&input1Value),
							)
							return nil
						}

						objectUnderTest := _InputSrcFactory{
							os:     new(ios.Fake),
							ioutil: new(iioutil.Fake),
							json:   fakeJSON,
						}

						/* act */
						actualInputSrc, actualErr := objectUnderTest.NewYMLFileInputSrc("dummyFilePath")

						/* assert */
						Expect(actualInputSrc).To(Equal(ymlFileInputSrc{
							argMap: expectedMap,
						}))
						Expect(actualErr).To(BeNil())
					})
				})
			})
		})
	})
})
