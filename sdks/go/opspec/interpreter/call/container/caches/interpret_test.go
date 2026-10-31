package caches

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("Interpret", func() {
	It("should create and mount a cache dir and its parents", func() {
		/* arrange */
		dataDirPath := filepath.Join(GinkgoT().TempDir(), "data")
		hostPath := filepath.Join(dataDirPath, "caches", "Go-mod_v1.2")

		/* act */
		actualDirs, err := Interpret(
			map[string]interface{}{"/go/pkg/mod": "Go-mod_v1.2"},
			dataDirPath,
		)

		/* assert */
		Expect(err).To(BeNil())
		Expect(actualDirs).To(Equal(map[string]string{"/go/pkg/mod": hostPath}))
		Expect(hostPath).To(BeADirectory())
	})

	It("should share existing cache contents by name", func() {
		/* arrange */
		dataDirPath := GinkgoT().TempDir()
		firstDirs, err := Interpret(
			map[string]interface{}{"/first": "shared"},
			dataDirPath,
		)
		Expect(err).To(BeNil())
		cacheFilePath := filepath.Join(firstDirs["/first"], "contents")
		Expect(os.WriteFile(cacheFilePath, []byte("cached"), 0600)).To(Succeed())

		/* act */
		actualDirs, err := Interpret(
			map[string]interface{}{"/second": "shared"},
			dataDirPath,
		)

		/* assert */
		Expect(err).To(BeNil())
		Expect(actualDirs["/second"]).To(Equal(firstDirs["/first"]))
		contents, err := os.ReadFile(cacheFilePath)
		Expect(err).To(BeNil())
		Expect(string(contents)).To(Equal("cached"))
	})

	for _, name := range []string{"../x", "a/b", ".", "..", "", "a b", `/absolute`, `a\b`, "$(name)", "a\n"} {
		It(fmt.Sprintf("should reject cache name %q", name), func() {
			dataDirPath := GinkgoT().TempDir()
			actualDirs, err := Interpret(
				map[string]interface{}{"/cache": name},
				dataDirPath,
			)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("/cache"))
			Expect(err.Error()).To(ContainSubstring("cache name"))
			Expect(actualDirs).To(BeNil())
			Expect(filepath.Join(dataDirPath, "caches")).NotTo(BeAnExistingFile())
		})
	}

	for _, value := range []interface{}{nil, 42, true, []interface{}{}, map[string]interface{}{}} {
		It(fmt.Sprintf("should reject a cache name of type %T", value), func() {
			actualDirs, err := Interpret(
				map[string]interface{}{"/cache": value},
				GinkgoT().TempDir(),
			)

			Expect(err).To(MatchError(`unable to bind cache dir "/cache": cache name must be a non-empty string literal`))
			Expect(actualDirs).To(BeNil())
		})
	}

	It("should return directory creation errors", func() {
		dataDirPath := GinkgoT().TempDir()
		Expect(os.WriteFile(filepath.Join(dataDirPath, "caches"), nil, 0600)).To(Succeed())

		actualDirs, err := Interpret(
			map[string]interface{}{"/cache": "gomod"},
			dataDirPath,
		)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(`unable to bind cache dir "/cache" to "gomod"`))
		Expect(actualDirs).To(BeNil())
	})

	It("should reject the root container path", func() {
		actualDirs, err := Interpret(
			map[string]interface{}{"/.": "gomod"},
			GinkgoT().TempDir(),
		)

		Expect(err).To(MatchError(`unable to bind cache dir "/.": container path must not be /`))
		Expect(actualDirs).To(BeNil())
	})

	It("should reject equivalent container paths", func() {
		actualDirs, err := Interpret(
			map[string]interface{}{"/cache": "a", "/cache/": "b"},
			GinkgoT().TempDir(),
		)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("container path is declared more than once"))
		Expect(actualDirs).To(BeNil())
	})
})
