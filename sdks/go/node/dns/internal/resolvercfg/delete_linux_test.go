package resolvercfg

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("Delete", func() {
	Context("server exists", func() {
		It("should return expected result", func() {
			/* arrange */
			tmpFile, err := os.CreateTemp("", "")
			if err != nil {
				panic(err)
			}

			providedText := `
nameserver 127.0.0.1 # do not edit; managed by opctl
# other line
`
			etcResolvConfPath = tmpFile.Name()

			_, err = tmpFile.WriteString(providedText)
			if nil != err {
				panic(err)
			}

			expectedText := `
# other line
`

			/* act */
			err = Delete(
				context.Background(),
			)
			if err != nil {
				panic(err)
			}

			actualTextBytes, err := os.ReadFile(tmpFile.Name())
			if err != nil {
				panic(err)
			}

			/* assert */
			Expect(string(actualTextBytes)).To(Equal(expectedText))

		})
	})
})
