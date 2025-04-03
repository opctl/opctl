package oppath

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("tryFindGitDir", func() {
	Context("../../../.git exists", func() {
		It("should expected path", func() {
			/* arrange */
			cwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			/* act */
			actualPath, actualErr := tryFindGitDir(cwd)

			actualRelPath, err := filepath.Rel(cwd, actualPath)
			if err != nil {
				panic(err)
			}

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(actualRelPath).To(Equal("../../../.git"))
		})
	})
})
