package caches

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	It("should mount fresh scratch cache dirs without creating the cache root", func() {
		dataDirPath := GinkgoT().TempDir()
		for _, callID := range []string{"first", "second"} {
			noCacheRootPath := filepath.Join(dataDirPath, "nocache", callID)
			actualDirs, err := Interpret(
				nil,
				map[string]interface{}{"/cache": "shared"},
				nil,
				dataDirPath,
				noCacheRootPath,
			)

			Expect(err).To(BeNil())
			hostPath := filepath.Join(noCacheRootPath, "shared")
			Expect(actualDirs).To(Equal(map[string]string{"/cache": hostPath}))
			entries, err := os.ReadDir(hostPath)
			Expect(err).To(BeNil())
			Expect(entries).To(BeEmpty())
			Expect(os.WriteFile(filepath.Join(hostPath, "contents"), []byte("scratch"), 0600)).To(Succeed())
			Expect(filepath.Join(dataDirPath, "caches")).NotTo(BeAnExistingFile())
		}
	})

	It("should leave existing cache contents untouched during a cold run", func() {
		dataDirPath := GinkgoT().TempDir()
		warmDirs, err := Interpret(nil, map[string]interface{}{"/cache": "shared"}, nil, dataDirPath, "")
		Expect(err).To(BeNil())
		cacheFilePath := filepath.Join(warmDirs["/cache"], "contents")
		Expect(os.WriteFile(cacheFilePath, []byte("cached"), 0600)).To(Succeed())

		coldDirs, err := Interpret(nil, map[string]interface{}{"/cache": "shared"}, nil, dataDirPath, filepath.Join(dataDirPath, "nocache", "call"))

		Expect(err).To(BeNil())
		entries, err := os.ReadDir(coldDirs["/cache"])
		Expect(err).To(BeNil())
		Expect(entries).To(BeEmpty())
		Expect(os.WriteFile(filepath.Join(coldDirs["/cache"], "contents"), []byte("cold"), 0600)).To(Succeed())
		contents, err := os.ReadFile(cacheFilePath)
		Expect(err).To(BeNil())
		Expect(string(contents)).To(Equal("cached"))
	})

	It("should still validate interpreted cache names during a cold run", func() {
		dataDirPath := GinkgoT().TempDir()
		noCacheRootPath := filepath.Join(dataDirPath, "nocache", "call")
		invalidName := "../escape"

		actualDirs, err := Interpret(
			map[string]*model.Value{"name": {String: &invalidName}},
			map[string]interface{}{"/cache": "$(name)"},
			nil,
			dataDirPath,
			noCacheRootPath,
		)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(`invalid cache name "../escape"`))
		Expect(actualDirs).To(BeNil())
		Expect(filepath.Join(dataDirPath, "caches")).NotTo(BeAnExistingFile())
		Expect(noCacheRootPath).NotTo(BeAnExistingFile())
	})

	digest := "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	for _, test := range []struct {
		description   string
		imageRef      string
		containerPath string
		name          string
	}{
		{"plain image", "node", "/root/.npm", "node-root-.npm"},
		{"tag", "node:22", "/root/.npm", "node-root-.npm"},
		{"digest", "docker.io/library/node@" + digest, "/root/.npm", "node-root-.npm"},
		{"tag and digest", "node:22@" + digest, "/root/.npm", "node-root-.npm"},
		{"registry with port", "localhost:5000/x/y:tag", "/c", "localhost-5000-x-y-c"},
		{"nested repository", "ghcr.io/org/tool:1", "/cache", "ghcr.io-org-tool-cache"},
		{"uppercase ref", "GHCR.IO/Org/Tool:TAG", "/cache", "ghcr.io-org-tool-cache"},
		{"nested path", "golang:1.23", "/go/pkg/mod", "golang-go-pkg-mod"},
		{"cleaned path with trailing slash", "golang:1.23", "/go//pkg/./cache/../mod/", "golang-go-pkg-mod"},
		{"consecutive separators", "org/my--tool:1", "/Go_mod.v1---/cache", "org-my-tool-Go_mod.v1-cache"},
		{"trailing hyphens", "node:22", "/cache---/---", "node-cache"},
	} {
		It(fmt.Sprintf("should derive a cache name for %s", test.description), func() {
			/* arrange */
			dataDirPath := filepath.Join(GinkgoT().TempDir(), "data")
			hostPath := filepath.Join(dataDirPath, "caches", test.name)

			/* act */
			actualDirs, err := Interpret(
				nil,
				map[string]interface{}{test.containerPath: nil},
				&test.imageRef,
				dataDirPath,
				"",
			)

			/* assert */
			Expect(err).To(BeNil())
			Expect(actualDirs).To(Equal(map[string]string{path.Clean(test.containerPath): hostPath}))
			Expect(hostPath).To(BeADirectory())
		})
	}

	It("should require an explicit cache name when the image has no ref", func() {
		dataDirPath := GinkgoT().TempDir()
		actualDirs, err := Interpret(
			nil,
			map[string]interface{}{"/root/.npm": nil},
			nil,
			dataDirPath,
			"",
		)

		Expect(err).To(MatchError(`unable to bind cache dir "/root/.npm": image has no ref; give an explicit cache name for this container path`))
		Expect(actualDirs).To(BeNil())
		Expect(filepath.Join(dataDirPath, "caches")).NotTo(BeAnExistingFile())
	})

	for _, imageRef := range []string{digest, "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"} {
		It(fmt.Sprintf("should require an explicit cache name for image ID %q", imageRef), func() {
			dataDirPath := GinkgoT().TempDir()
			actualDirs, err := Interpret(
				nil,
				map[string]interface{}{"/cache": nil},
				&imageRef,
				dataDirPath,
				"",
			)

			Expect(err).To(MatchError(`unable to bind cache dir "/cache": image ref has no repository; give an explicit cache name for this container path`))
			Expect(actualDirs).To(BeNil())
			Expect(filepath.Join(dataDirPath, "caches")).NotTo(BeAnExistingFile())
		})
	}

	It("should return image ref parsing errors before creating a cache dir", func() {
		dataDirPath := GinkgoT().TempDir()
		imageRef := "not a ref"
		actualDirs, err := Interpret(
			nil,
			map[string]interface{}{"/cache": nil},
			&imageRef,
			dataDirPath,
			"",
		)

		Expect(err).To(MatchError(`unable to bind cache dir "/cache": unable to derive cache name from image ref: invalid reference format`))
		Expect(actualDirs).To(BeNil())
		Expect(filepath.Join(dataDirPath, "caches")).NotTo(BeAnExistingFile())
	})

	It("should create and mount a cache dir and its parents", func() {
		/* arrange */
		dataDirPath := filepath.Join(GinkgoT().TempDir(), "data")
		hostPath := filepath.Join(dataDirPath, "caches", "Go-mod_v1.2")

		/* act */
		actualDirs, err := Interpret(
			nil,
			map[string]interface{}{"/go/pkg/mod": "Go-mod_v1.2"},
			nil,
			dataDirPath,
			"",
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
			nil,
			dataDirPath,
			"",
		)
		Expect(err).To(BeNil())
		cacheFilePath := filepath.Join(firstDirs["/first"], "contents")
		Expect(os.WriteFile(cacheFilePath, []byte("cached"), 0600)).To(Succeed())

		/* act */
		actualDirs, err := Interpret(
			nil,
			map[string]interface{}{"/second": "shared"},
			nil,
			dataDirPath,
			"",
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
				nil,
				dataDirPath,
				"",
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
				nil,
				dataDirPath,
				"",
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
			nil,
			dataDirPath,
			"",
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
				nil,
				dataDirPath,
				"",
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
				nil,
				dataDirPath,
				"",
			)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("/cache"))
			Expect(err.Error()).To(ContainSubstring("cache name"))
			Expect(actualDirs).To(BeNil())
			Expect(filepath.Join(dataDirPath, "caches")).NotTo(BeAnExistingFile())
		})
	}

	for _, value := range []interface{}{42, true, []interface{}{}, map[string]interface{}{}} {
		It(fmt.Sprintf("should reject a cache name of type %T", value), func() {
			actualDirs, err := Interpret(
				nil,
				map[string]interface{}{"/cache": value},
				nil,
				GinkgoT().TempDir(),
				"",
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
			nil,
			dataDirPath,
			"",
		)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(`unable to bind cache dir "/cache" to "gomod"`))
		Expect(actualDirs).To(BeNil())
	})

	It("should reject the root container path", func() {
		actualDirs, err := Interpret(
			nil,
			map[string]interface{}{"/.": "gomod"},
			nil,
			GinkgoT().TempDir(),
			"",
		)

		Expect(err).To(MatchError(`unable to bind cache dir "/.": container path must not be /`))
		Expect(actualDirs).To(BeNil())
	})

	It("should reject equivalent container paths", func() {
		actualDirs, err := Interpret(
			nil,
			map[string]interface{}{"/cache": "a", "/cache/": "b"},
			nil,
			GinkgoT().TempDir(),
			"",
		)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("container path is declared more than once"))
		Expect(actualDirs).To(BeNil())
	})
})
