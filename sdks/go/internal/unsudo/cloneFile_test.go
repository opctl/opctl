package unsudo

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("cloneFile", func() {
	It("should maintain mode, gid, uid, and content & no err", func() {
		/* arrange */
		srcFile, err := os.CreateTemp("", "")
		if err != nil {
			panic(err)
		}

		srcContent := []byte(`
			line1
			line2`,
		)

		if _, err := srcFile.Write(
			srcContent,
		); err != nil {
			panic(err)
		}

		srcInfo, err := srcFile.Stat()
		if err != nil {
			panic(err)
		}

		// make executable
		if err := os.Chmod(
			srcFile.Name(),
			srcInfo.Mode()|0111,
		); err != nil {
			panic(err)
		}

		srcInfo, err = srcFile.Stat()
		if err != nil {
			panic(err)
		}

		uid := 4000
		gid := 5000

		if err := os.Chown(
			srcFile.Name(),
			uid,
			gid,
		); err != nil {
			panic(err)
		}

		dstPath := srcFile.Name() + "copy"

		/* act */
		err = CloneFile(
			srcFile.Name(),
			dstPath,
		)

		/* assert */
		dstFile, openErr := os.Open(dstPath)
		if err != nil {
			panic(openErr)
		}

		dstInfo, statErr := dstFile.Stat()
		if err != nil {
			panic(statErr)
		}

		dstContent, err := os.ReadFile(dstFile.Name())

		Expect(err).To(BeNil())
		Expect(dstPath).Should(BeAnExistingFile())
		Expect(dstInfo.Mode()).To(Equal(srcInfo.Mode()))
		Expect(dstContent).To(Equal(srcContent))
	})
})
