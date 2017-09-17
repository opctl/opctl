package ijson

import (
	"bytes"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("IJSON", func() {
	Context("New", func() {
		It("should return IJSON", func() {
			/* arrange/act/assert */
			Expect(New()).
				Should(Not(BeNil()))
		})
	})
	Context("Marshal", func() {
		Context("errs", func() {
			It("should return expected err", func() {
				/* arrange */
				erroneousInput := make(chan int)
				objectUnderTest := New()

				_, expectedErr := json.Marshal(erroneousInput)

				/* act */
				_, actualErr := objectUnderTest.Marshal(erroneousInput)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("doesn't err", func() {
			It("should return expected bytes", func() {
				/* arrange */
				providedStruct := struct {
					Field1 string
					Field2 int
				}{
					Field1: "dummyString1",
					Field2: 1000,
				}

				expectedBytes, _ := json.Marshal(providedStruct)

				objectUnderTest := New()

				/* act */
				actualBytes, actualErr := objectUnderTest.Marshal(providedStruct)

				/* assert */
				Expect(actualBytes).To(Equal(expectedBytes))
				Expect(actualErr).To(BeNil())
			})
		})
	})
	Context("Marshal", func() {
		It("should return expected json.Encoder", func() {
			/* arrange */
			writer := bytes.NewBufferString("")

			expectedEncoder := json.NewEncoder(writer)

			objectUnderTest := New()

			/* act */
			actualEncoder := objectUnderTest.NewEncoder(writer)

			/* assert */
			Expect(actualEncoder).To(Equal(expectedEncoder))
		})
	})
	Context("Unmarshal", func() {
		Context("errs", func() {
			It("should return expected err", func() {
				/* arrange */
				erroneousInput := []byte("$$")
				objectUnderTest := New()

				expectedErr := json.Unmarshal(erroneousInput, &struct{}{})

				/* act */
				actualErr := objectUnderTest.Unmarshal(erroneousInput, &struct{}{})

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("doesn't err", func() {
			It("should unmarshal expected object", func() {
				/* arrange */
				expectedStruct := struct {
					Field1 string
					Field2 int
				}{
					Field1: "dummyString1",
					Field2: 1000,
				}

				providedBytes, err := json.Marshal(expectedStruct)
				if nil != err {
					Fail(err.Error())
				}

				actualStruct := struct {
					Field1 string
					Field2 int
				}{}

				objectUnderTest := New()

				/* act */
				actualErr := objectUnderTest.Unmarshal(providedBytes, &actualStruct)

				/* assert */
				Expect(actualStruct).To(Equal(expectedStruct))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
