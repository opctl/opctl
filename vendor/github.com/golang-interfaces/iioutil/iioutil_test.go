package iioutil

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

var _ = Context("_Iioutil", func() {
	wd, err := os.Getwd()
	if nil != err {
		panic(err)
	}
	Context("New", func() {
		It("should return Fs", func() {
			/* arrange/act/assert */
			Expect(New()).
				Should(Not(BeNil()))
		})
	})
	Context("ReadAll", func() {
		It("should return expected results", func() {
			/* arrange */
			expectedBytes := []byte{2, 3, 10}

			objectUnderTest := _IIOUtil{}

			/* act */
			actualBytes, actualErr := objectUnderTest.ReadAll(bytes.NewReader(expectedBytes))

			/* assert */
			Expect(actualBytes).To(Equal(expectedBytes))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("ReadDir", func() {
		It("should return expected fileinfos", func() {
			/* arrange */
			// use .opspec dir because it won't be modified during test
			providedDirName := path.Join(wd, ".opspec")

			expectedFileInfos, _ := ioutil.ReadDir(providedDirName)

			objectUnderTest := _IIOUtil{}

			/* act */
			actualFileinfos, actualErr := objectUnderTest.ReadDir(providedDirName)

			/* assert */
			Expect(actualFileinfos).To(Equal(expectedFileInfos))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("ReadFile", func() {
		It("should return expected fileinfo", func() {
			/* arrange */
			// use current file for test
			_, providedFileName, _, _ := runtime.Caller(1)

			expectedBytes, _ := ioutil.ReadFile(providedFileName)

			objectUnderTest := _IIOUtil{}

			/* act */
			actualBytes, actualErr := objectUnderTest.ReadFile(providedFileName)

			/* assert */
			Expect(actualBytes).To(Equal(expectedBytes))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("WriteFile", func() {
		It("should create expected file", func() {
			/* arrange */
			tempFile, err := ioutil.TempFile("", "dummyFile")
			if nil != err {
				panic(err)
			}
			providedFilename := tempFile.Name()

			providedData := bytes.NewBufferString("dummy file content").Bytes()
			providedPerm := os.FileMode(0777)

			objectUnderTest := _IIOUtil{}

			/* act */
			objectUnderTest.WriteFile(providedFilename, providedData, providedPerm)

			/* assert */
			actualData, err := ioutil.ReadFile(providedFilename)
			if nil != err {
				panic(err)
			}

			Expect(actualData).To(Equal(providedData))
		})
	})
})
