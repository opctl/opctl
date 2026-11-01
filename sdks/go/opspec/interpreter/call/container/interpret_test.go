package container

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/internal/runcontext"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("cmd.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			dataDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Interpret(
				context.Background(),
				map[string]*model.Value{},
				&model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "ref",
					},
					Cmd: []interface{}{
						"$()",
					},
				},
				"dummyContainerID",
				"dummyOpPath",
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret [$()] to array: unable to interpret '$()' as array initializer item: unable to interpret '' as reference: '' not in scope"))
		})
	})

	Context("dirs.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			identifier := "identifier"
			dataDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Interpret(
				context.Background(),
				map[string]*model.Value{
					identifier: {
						Socket: new(string),
					},
				},
				&model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "ref",
					},
					Dirs: map[string]interface{}{
						"/something": fmt.Sprintf("$(%s)", identifier),
					},
				},
				"dummyContainerID",
				"dummyOpPath",
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to bind directory /something to $(identifier): unable to interpret $(identifier) to dir: unable to coerce socket to dir: incompatible types"))
		})
	})

	Context("caches", func() {
		for _, containerPath := range []string{"/caches", "/../caches"} {
			It(fmt.Sprintf("should isolate cold cache dirs from dirs and files scratch entries at %s", containerPath), func() {
				dataDirPath := GinkgoT().TempDir()
				sourceDirPath := filepath.Join(dataDirPath, "ops", "source")
				Expect(os.MkdirAll(sourceDirPath, 0700)).To(Succeed())
				sourceFilePath := filepath.Join(sourceDirPath, "file")
				Expect(os.WriteFile(sourceFilePath, []byte("source"), 0600)).To(Succeed())

				actualCall, err := Interpret(
					runcontext.WithNoCacheRoot(context.Background(), filepath.Join(dataDirPath, "nocache", "rootCallID")),
					map[string]*model.Value{
						"dir":  {Dir: &sourceDirPath},
						"file": {File: &sourceFilePath},
					},
					&model.ContainerCallSpec{
						Image:  &model.ContainerCallImageSpec{Ref: "node:22"},
						Dirs:   map[string]interface{}{containerPath: "$(dir)"},
						Files:  map[string]interface{}{containerPath + "/npm": "$(file)"},
						Caches: map[string]interface{}{"/root/.npm": "npm"},
					},
					"containerID",
					"opPath",
					dataDirPath,
				)

				Expect(err).To(BeNil())
				scratchDirPath := filepath.Join(dataDirPath, "dcg", "containerID", "fs")
				Expect(actualCall.Dirs[containerPath]).To(Equal(filepath.Join(scratchDirPath, "caches")))
				Expect(actualCall.Files[containerPath+"/npm"]).To(Equal(filepath.Join(scratchDirPath, "caches", "npm")))
				Expect(actualCall.Dirs["/root/.npm"]).To(Equal(filepath.Join(dataDirPath, "nocache", "rootCallID", "npm")))
				entries, err := os.ReadDir(actualCall.Dirs["/root/.npm"])
				Expect(err).To(BeNil())
				Expect(entries).To(BeEmpty())
				contents, err := os.ReadFile(actualCall.Files[containerPath+"/npm"])
				Expect(err).To(BeNil())
				Expect(string(contents)).To(Equal("source"))
				Expect(filepath.Join(dataDirPath, "caches")).NotTo(BeAnExistingFile())
			})
		}

		It("should derive cache names from the interpreted image ref", func() {
			/* arrange */
			dataDirPath := GinkgoT().TempDir()
			imageRef := "node:22"
			hostPath := filepath.Join(dataDirPath, "caches", "node-root-.npm")

			/* act */
			actualCall, err := Interpret(
				context.Background(),
				map[string]*model.Value{"imageRef": {String: &imageRef}},
				&model.ContainerCallSpec{
					Image:  &model.ContainerCallImageSpec{Ref: "$(imageRef)"},
					Caches: map[string]interface{}{"/root/.npm": nil},
				},
				"containerID",
				"opPath",
				dataDirPath,
			)

			/* assert */
			Expect(err).To(BeNil())
			Expect(actualCall.Dirs).To(Equal(map[string]string{"/root/.npm": hostPath}))
			Expect(hostPath).To(BeADirectory())
		})

		for _, name := range []interface{}{nil, "npm"} {
			It(fmt.Sprintf("should interpret cache name %v with an image src", name), func() {
				/* arrange */
				dataDirPath := GinkgoT().TempDir()
				imageDir := GinkgoT().TempDir()

				/* act */
				actualCall, err := Interpret(
					context.Background(),
					map[string]*model.Value{"imageDir": {Dir: &imageDir}},
					&model.ContainerCallSpec{
						Image:  &model.ContainerCallImageSpec{Ref: "$(imageDir)"},
						Caches: map[string]interface{}{"/root/.npm": name},
					},
					"containerID",
					"opPath",
					dataDirPath,
				)

				/* assert */
				if name == nil {
					Expect(err).To(MatchError(`unable to bind cache dir "/root/.npm": image has no ref; give an explicit cache name for this container path`))
					Expect(actualCall).To(BeNil())
					Expect(filepath.Join(dataDirPath, "caches")).NotTo(BeAnExistingFile())
				} else {
					Expect(err).To(BeNil())
					Expect(actualCall.Image.Ref).To(BeNil())
					Expect(actualCall.Image.Src.Dir).To(Equal(&imageDir))
					hostPath := filepath.Join(dataDirPath, "caches", "npm")
					Expect(actualCall.Dirs).To(Equal(map[string]string{"/root/.npm": hostPath}))
					Expect(hostPath).To(BeADirectory())
				}
			})
		}

		It("should merge cache dirs with dirs", func() {
			/* arrange */
			dataDirPath := GinkgoT().TempDir()
			dirPath := GinkgoT().TempDir()
			lockHash := "abc123"

			/* act */
			actualCall, err := Interpret(
				context.Background(),
				map[string]*model.Value{
					"dir":      {Dir: &dirPath},
					"lockHash": {String: &lockHash},
				},
				&model.ContainerCallSpec{
					Image:  &model.ContainerCallImageSpec{Ref: "ref"},
					Dirs:   map[string]interface{}{"/work": "$(dir)"},
					Caches: map[string]interface{}{"/go/pkg/mod": "gomod-$(lockHash)"},
				},
				"containerID",
				"opPath",
				dataDirPath,
			)

			/* assert */
			Expect(err).To(BeNil())
			Expect(actualCall.Dirs).To(Equal(map[string]string{
				"/work":       dirPath,
				"/go/pkg/mod": filepath.Join(dataDirPath, "caches", "gomod-abc123"),
			}))
		})

		It("should reject a container path declared in both dirs and caches", func() {
			/* arrange */
			dirPath := GinkgoT().TempDir()

			/* act */
			actualCall, err := Interpret(
				context.Background(),
				map[string]*model.Value{"dir": {Dir: &dirPath}},
				&model.ContainerCallSpec{
					Image:  &model.ContainerCallImageSpec{Ref: "ref"},
					Dirs:   map[string]interface{}{"/cache": "$(dir)"},
					Caches: map[string]interface{}{"/cache/": "gomod"},
				},
				"containerID",
				"opPath",
				GinkgoT().TempDir(),
			)

			/* assert */
			Expect(err).To(MatchError(`unable to bind cache dir "/cache": container path is also declared in dirs or files`))
			Expect(actualCall).To(BeNil())
		})

		It("should reject a container path declared in both files and caches", func() {
			/* arrange */
			filePath := filepath.Join(GinkgoT().TempDir(), "file")
			Expect(os.WriteFile(filePath, nil, 0600)).To(Succeed())

			/* act */
			actualCall, err := Interpret(
				context.Background(),
				map[string]*model.Value{"file": {File: &filePath}},
				&model.ContainerCallSpec{
					Image:  &model.ContainerCallImageSpec{Ref: "ref"},
					Files:  map[string]interface{}{"/cache": "$(file)"},
					Caches: map[string]interface{}{"/cache": "gomod"},
				},
				"containerID",
				"opPath",
				GinkgoT().TempDir(),
			)

			/* assert */
			Expect(err).To(MatchError(`unable to bind cache dir "/cache": container path is also declared in dirs or files`))
			Expect(actualCall).To(BeNil())
		})

		It("should return cache interpretation errors", func() {
			actualCall, err := Interpret(
				context.Background(),
				map[string]*model.Value{},
				&model.ContainerCallSpec{
					Image:  &model.ContainerCallImageSpec{Ref: "ref"},
					Caches: map[string]interface{}{"/cache": "../escape"},
				},
				"containerID",
				"opPath",
				GinkgoT().TempDir(),
			)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(`invalid cache name "../escape"`))
			Expect(actualCall).To(BeNil())
		})
	})

	Context("envVars.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			dataDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Interpret(
				context.Background(),
				map[string]*model.Value{},
				&model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "ref",
					},
					EnvVars: "$()",
				},
				"dummyContainerID",
				"dummyOpPath",
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret '$()' as envVars: unable to interpret $() to object: unable to interpret '' as reference: '' not in scope"))
		})
	})

	Context("files.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			dataDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Interpret(
				context.Background(),
				map[string]*model.Value{
					"not": {
						Socket: new(string),
					},
				},
				&model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "ref",
					},
					Files: map[string]interface{}{
						"/something": "$(not)",
					},
				},
				"dummyContainerID",
				"dummyOpPath",
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to bind file /something to $(not): unable to coerce '{\"socket\":\"\"}' to file"))
		})
	})

	Context("image.Interpret errors", func() {
		It("should return expected error", func() {
			/* arrange */
			dataDir, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			/* act */
			_, actualErr := Interpret(
				context.Background(),
				map[string]*model.Value{},
				&model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "$()",
					},
				},
				"dummyContainerID",
				"dummyOpPath",
				dataDir,
			)

			/* assert */
			Expect(actualErr).To(MatchError("unable to interpret $() to string: unable to interpret '' as reference: '' not in scope"))
		})
	})

	It("should return expected result", func() {
		/* arrange */
		providedContainerID := "providedContainerID"
		providedOpPath := "providedOpPath"

		expectedRef := "docker.io/library/ref"

		expectedResult := model.ContainerCall{
			BaseCall: model.BaseCall{
				OpPath: providedOpPath,
			},
			ContainerID: providedContainerID,
			Cmd:         []string{},
			Dirs:        map[string]string{},
			Files:       map[string]string{},
			Image: &model.ContainerCallImage{
				Ref: &expectedRef,
			},
			Sockets: map[string]string{},
			WorkDir: "",
		}

		dataDir, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}

		/* act */
		actualResult, actualErr := Interpret(
			context.Background(),
			map[string]*model.Value{},
			&model.ContainerCallSpec{
				Image: &model.ContainerCallImageSpec{
					Ref: "ref",
				},
			},
			providedContainerID,
			providedOpPath,
			dataDir,
		)

		/* assert */
		Expect(actualErr).To(BeNil())
		Expect(*actualResult).To(Equal(expectedResult))
	})
})
