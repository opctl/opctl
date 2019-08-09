package cliparamsatisfier

import (
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"reflect"
)

var _ = Describe("ymlFileInputSrc", func() {
	Context("ReadString()", func() {
		Context("yml doesn't contain entry w/ provided inputName", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, os.ErrNotExist)

				inputSrcFactory := _InputSrcFactory{
					os:   fakeOS,
					json: new(ijson.Fake),
				}

				objectUnderTest, err := inputSrcFactory.NewYMLFileInputSrc("")
				if err != nil {
					Fail(err.Error())
				}

				/* act */
				actualValue, actualOk := objectUnderTest.ReadString("nonExistentInputName")

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualOk).To(BeFalse())
			})
		})
		Context("yml contains entry w/ provided inputName", func() {
			It("should return expected result", func() {
				/* arrange */
				providedInputName := "dummyInputName"
				value := "dummyValue"

				fakeJSON := new(ijson.Fake)
				fakeJSON.UnmarshalStub = func(data []byte, v interface{}) error {
					reflect.ValueOf(v).Elem().SetMapIndex(
						reflect.ValueOf(providedInputName),
						reflect.ValueOf(value),
					)
					return nil
				}

				inputSrcFactory := _InputSrcFactory{
					os:     new(ios.Fake),
					json:   fakeJSON,
					ioutil: new(iioutil.Fake),
				}

				objectUnderTest, err := inputSrcFactory.NewYMLFileInputSrc("")
				if err != nil {
					Fail(err.Error())
				}

				/* act */
				actualValue, actualOk := objectUnderTest.ReadString(providedInputName)

				/* assert */
				Expect(*actualValue).To(Equal(string(value)))
				Expect(actualOk).To(BeTrue())
			})
			It("should return value only once", func() {
				/* arrange */
				providedInputName := "dummyInputName"
				value := "dummyValue"

				fakeJSON := new(ijson.Fake)
				fakeJSON.UnmarshalStub = func(data []byte, v interface{}) error {
					reflect.ValueOf(v).Elem().SetMapIndex(
						reflect.ValueOf(providedInputName),
						reflect.ValueOf(value),
					)
					return nil
				}

				inputSrcFactory := _InputSrcFactory{
					os:     new(ios.Fake),
					json:   fakeJSON,
					ioutil: new(iioutil.Fake),
				}

				objectUnderTest, err := inputSrcFactory.NewYMLFileInputSrc("")
				if err != nil {
					Fail(err.Error())
				}

				/* act */
				actualValue1, actualOk1 := objectUnderTest.ReadString(providedInputName)
				actualValue2, actualOk2 := objectUnderTest.ReadString(providedInputName)

				/* assert */
				Expect(*actualValue1).To(Equal(string(value)))
				Expect(actualOk1).To(BeTrue())

				Expect(actualValue2).To(BeNil())
				Expect(actualOk2).To(BeFalse())
			})
		})
	})
})
