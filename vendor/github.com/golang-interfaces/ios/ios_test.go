package ios

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Context("IOS", func() {
	wd, err := os.Getwd()
	if nil != err {
		panic(err)
	}
	Context("New", func() {
		It("should return IOS", func() {
			/* arrange/act/assert */
			Expect(New()).
				Should(Not(BeNil()))
		})
	})
	Context("FindProcess", func() {
		It("should return expected process", func() {
			/* arrange */
			providedPID := os.Getpid()
			expectedProcess, _ := os.FindProcess(providedPID)
			objectUnderTest := _IOS{}

			/* act */
			actualProcess, actualErr := objectUnderTest.FindProcess(providedPID)

			/* assert */
			Expect(actualProcess).To(Equal(expectedProcess))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Getpid", func() {
		It("should return expected PID", func() {
			/* arrange */
			expectedPID := os.Getpid()
			objectUnderTest := _IOS{}

			/* act */
			actualPID := objectUnderTest.Getpid()

			/* assert */
			Expect(actualPID).To(Equal(expectedPID))
		})
	})
	Context("Getenv proceeding Setenv", func() {
		It("should return value set by Setenv", func() {
			/* arrange */
			providedName := "dummyName"
			providedValue := "dummyValue"
			expectedValue := providedValue
			objectUnderTest := _IOS{}

			objectUnderTest.Setenv(providedName, providedValue)

			/* act */
			actualValue := objectUnderTest.Getenv(providedName)

			/* assert */
			Expect(actualValue).To(Equal(expectedValue))
		})
	})
	Context("Getwd", func() {
		It("should return expected process", func() {
			/* arrange */
			expectedWd, _ := os.Getwd()
			objectUnderTest := _IOS{}

			/* act */
			actualWd, actualErr := objectUnderTest.Getwd()

			/* assert */
			Expect(actualWd).To(Equal(expectedWd))
			Expect(actualErr).To(BeNil())
		})
	})

	Context("Open", func() {
		It("should return expected file", func() {
			/* arrange */
			providedName := wd

			expectedFile, _ := os.Open(providedName)
			expectedFileStat, _ := expectedFile.Stat()
			defer expectedFile.Close()

			objectUnderTest := New()

			/* act */
			actualFile, actualErr := objectUnderTest.Open(providedName)
			actualFileStat, _ := actualFile.Stat()
			defer actualFile.Close()

			/* assert */
			// compare FileInfo's; can't compare File's because FD will be unique
			Expect(actualFileStat).To(Equal(expectedFileStat))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("OpenFile", func() {
		It("should return expected file", func() {
			/* arrange */
			providedName := wd
			providedFlag := os.O_RDONLY
			providedPerm := os.FileMode(0)

			expectedFile, _ := os.OpenFile(providedName, providedFlag, providedPerm)
			expectedFileStat, _ := expectedFile.Stat()
			defer expectedFile.Close()

			objectUnderTest := New()

			/* act */
			actualFile, actualErr := objectUnderTest.OpenFile(providedName, providedFlag, providedPerm)
			actualFileStat, _ := actualFile.Stat()
			defer actualFile.Close()

			/* assert */
			// compare FileInfo's; can't compare File's because FD will be unique
			Expect(actualFileStat).To(Equal(expectedFileStat))
			Expect(actualErr).To(BeNil())
		})
	})
	Context("Stat", func() {
		It("should return expected fileinfo", func() {
			/* arrange */
			providedName := wd

			expectedFileInfo, _ := os.Stat(providedName)

			objectUnderTest := New()

			/* act */
			actualFileinfo, actualErr := objectUnderTest.Stat(providedName)

			/* assert */
			Expect(actualFileinfo).To(Equal(expectedFileInfo))
			Expect(actualErr).To(BeNil())
		})
	})
})
