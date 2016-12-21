package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("workDirPathGetter", func() {
	Context("Get", func() {
		It("should return current slash separated work dir path", func() {
			/* arrange */
			currentWorkDir, err := os.Getwd()
			if nil != err {
				Fail(err.Error())
			}

			expectedWorkDirPath := filepath.ToSlash(currentWorkDir)

			objectUnderTest := newWorkDirPathGetter()

			/* act */
			actualWorkDirPath := objectUnderTest.Get()

			/* assert */
			Expect(actualWorkDirPath).Should(Equal(expectedWorkDirPath))

		})
	})
})
