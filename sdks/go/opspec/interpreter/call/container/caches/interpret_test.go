package caches

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	It("should create and mount a cache dir and its parents", func() {
		/* arrange */
		dataDirPath := filepath.Join(GinkgoT().TempDir(), "data")
		hostPath := filepath.Join(dataDirPath, "caches", "Go-mod_v1.2")

		/* act */
		actualDirs, err := Interpret(
			nil,
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
			nil,
			map[string]interface{}{"/first": "shared"},
			dataDirPath,
		)
		Expect(err).To(BeNil())
		cacheFilePath := filepath.Join(firstDirs["/first"], "contents")
		Expect(os.WriteFile(cacheFilePath, []byte("cached"), 0600)).To(Succeed())

		/* act */
		actualDirs, err := Interpret(
			nil,
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

	for _, expression := range []string{"$(lockHash)", "npm-$(lockHash)"} {
		It(fmt.Sprintf("should resolve cache name %q from scope", expression), func() {
			/* arrange */
			dataDirPath := GinkgoT().TempDir()
			lockHash := "abc123"
			name := lockHash
			if expression != "$(lockHash)" {
				name = "npm-" + lockHash
			}
			hostPath := filepath.Join(dataDirPath, "caches", name)

			/* act */
			actualDirs, err := Interpret(
				map[string]*model.Value{"lockHash": {String: &lockHash}},
				map[string]interface{}{"/root/.npm": expression},
				dataDirPath,
			)

			/* assert */
			Expect(err).To(BeNil())
			Expect(actualDirs).To(Equal(map[string]string{"/root/.npm": hostPath}))
			Expect(hostPath).To(BeADirectory())
		})
	}

	for _, name := range []string{"a/b", "../../etc", ".", "..", "a b", ""} {
		It(fmt.Sprintf("should reject interpolated cache name %q before creating any cache dir", name), func() {
			dataDirPath := GinkgoT().TempDir()
			actualDirs, err := Interpret(
				map[string]*model.Value{"name": {String: &name}},
				map[string]interface{}{"/cache": "$(name)"},
				dataDirPath,
			)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("/cache"))
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("invalid cache name %q", name)))
			Expect(actualDirs).To(BeNil())
			Expect(filepath.Join(dataDirPath, "caches")).NotTo(BeAnExistingFile())
		})
	}

	It("should use a prior call's file contents in a cache name", func() {
		/* arrange */
		dataDirPath := GinkgoT().TempDir()
		lockHashPath := filepath.Join(GinkgoT().TempDir(), "lockHash")
		Expect(os.WriteFile(lockHashPath, []byte("abc123"), 0600)).To(Succeed())
		hostPath := filepath.Join(dataDirPath, "caches", "npm-abc123")

		/* act */
		actualDirs, err := Interpret(
			map[string]*model.Value{"lockHash": {File: &lockHashPath}},
			map[string]interface{}{"/root/.npm": "npm-$(lockHash)"},
			dataDirPath,
		)

		/* assert */
		Expect(err).To(BeNil())
		Expect(actualDirs).To(Equal(map[string]string{"/root/.npm": hostPath}))
		Expect(hostPath).To(BeADirectory())
	})

	for _, expression := range []string{"$(name)", "npm-$(lockHash)"} {
		It(fmt.Sprintf("should reject unresolvable reference in cache name %q", expression), func() {
			dataDirPath := GinkgoT().TempDir()
			actualDirs, err := Interpret(
				nil,
				map[string]interface{}{"/cache": expression},
				dataDirPath,
			)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(`unable to bind cache dir "/cache"`))
			Expect(err.Error()).To(ContainSubstring("not in scope"))
			Expect(actualDirs).To(BeNil())
			Expect(filepath.Join(dataDirPath, "caches")).NotTo(BeAnExistingFile())
		})
	}

	for _, name := range []string{"../x", "a/b", ".", "..", "", "a b", `/absolute`, `a\b`, "a\n"} {
		It(fmt.Sprintf("should reject cache name %q", name), func() {
			dataDirPath := GinkgoT().TempDir()
			actualDirs, err := Interpret(
				nil,
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
				nil,
				map[string]interface{}{"/cache": value},
				GinkgoT().TempDir(),
			)

			Expect(err).To(MatchError(`unable to bind cache dir "/cache": cache name must be a string expression`))
			Expect(actualDirs).To(BeNil())
		})
	}

	It("should return directory creation errors", func() {
		dataDirPath := GinkgoT().TempDir()
		Expect(os.WriteFile(filepath.Join(dataDirPath, "caches"), nil, 0600)).To(Succeed())

		actualDirs, err := Interpret(
			nil,
			map[string]interface{}{"/cache": "gomod"},
			dataDirPath,
		)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(`unable to bind cache dir "/cache" to "gomod"`))
		Expect(actualDirs).To(BeNil())
	})

	It("should reject the root container path", func() {
		actualDirs, err := Interpret(
			nil,
			map[string]interface{}{"/.": "gomod"},
			GinkgoT().TempDir(),
		)

		Expect(err).To(MatchError(`unable to bind cache dir "/.": container path must not be /`))
		Expect(actualDirs).To(BeNil())
	})

	It("should reject equivalent container paths", func() {
		actualDirs, err := Interpret(
			nil,
			map[string]interface{}{"/cache": "a", "/cache/": "b"},
			GinkgoT().TempDir(),
		)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("container path is declared more than once"))
		Expect(actualDirs).To(BeNil())
	})
})
