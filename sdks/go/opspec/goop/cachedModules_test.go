package goop

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	containerInterpreter "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container"
)

var _ = Describe("fetched Go module bounds", func() {
	var data, repo string
	write := func(path, text string) {
		Expect(os.MkdirAll(filepath.Dir(path), 0700)).To(Succeed())
		Expect(os.WriteFile(path, []byte(text), 0600)).To(Succeed())
	}
	BeforeEach(func() {
		var err error
		data, err = os.MkdirTemp("", "goop-cache-")
		Expect(err).NotTo(HaveOccurred())
		repo = filepath.Join(data, "ops", "github.com", "example", "repo#1.0.0")
		Expect(os.MkdirAll(repo, 0700)).To(Succeed())
	})
	AfterEach(func() { os.RemoveAll(data) })
	It("copies cached parent modules and in-repo replacements while preserving mount destinations and workdir", func() {
		for _, nested := range []bool{false, true} {
			module := repo
			op := filepath.Join(repo, "ops", "hello")
			if nested {
				module = filepath.Join(repo, "app")
				op = filepath.Join(module, "ops", "hello")
			}
			contents := "module example.com/app\ngo 1.25.0\n"
			if nested {
				contents += "replace example.com/sdk => ../sdk\n"
				write(filepath.Join(repo, "sdk", "go.mod"), "module example.com/sdk\ngo 1.21\n")
				write(filepath.Join(repo, "sdk", "sdk.go"), "package sdk\n")
			}
			write(filepath.Join(module, "go.mod"), contents)
			write(filepath.Join(module, "helper", "helper.go"), "package helper\nfunc Text()string{return \"hello\"}\n")
			write(filepath.Join(op, "op.go"), "package main\nimport \"example.com/app/helper\"\nfunc main(){println(helper.Text())}\n")
			for _, phase := range []string{"describe", "plan"} {
				spec, mount, err := ContainerSpecInRepository(op, phase, repo)
				Expect(err).NotTo(HaveOccurred())
				Expect(mount).To(Equal(module))
				Expect(spec.Image.Ref).To(Equal("golang:1.25.0"))
				Expect(spec.WorkDir).To(Equal(op))
				request := filepath.Join(data, "request.json")
				write(request, "{}")
				call, err := containerInterpreter.Interpret(map[string]*model.Value{"module": {Dir: &mount}, "request": {File: &request}}, spec, "cached-"+phase, op, data)
				Expect(err).NotTo(HaveOccurred())
				copied := call.Dirs[module]
				Expect(copied).NotTo(Equal(module))
				Expect(call.WorkDir).To(Equal(op))
				Expect(copied).To(BeADirectory())
				Expect(filepath.Join(copied, "helper", "helper.go")).To(BeAnExistingFile())
				Expect(filepath.Join(copied, "go.mod")).To(BeAnExistingFile())
				write(filepath.Join(copied, "helper", "helper.go"), "mutated")
				original, err := os.ReadFile(filepath.Join(module, "helper", "helper.go"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(original)).To(ContainSubstring("package helper"))
				if nested {
					Expect(spec.Dirs).To(HaveLen(2))
					sdk := filepath.Join(repo, "sdk")
					Expect(call.Dirs).To(HaveKey(sdk))
					Expect(call.Dirs[sdk]).NotTo(Equal(sdk))
					Expect(filepath.Join(call.Dirs[sdk], "sdk.go")).To(BeAnExistingFile())
				}
			}
		}
	})
	It("rejects missing in-repo modules instead of searching above the fetched root", func() {
		write(filepath.Join(data, "go.mod"), "module outside\ngo 1.25\n")
		op := filepath.Join(repo, "ops", "hello")
		Expect(os.MkdirAll(op, 0700)).To(Succeed())
		_, _, err := ContainerSpecInRepository(op, "describe", repo)
		Expect(err).To(MatchError(ContainSubstring("fetched repository")))
		_, _, err = ContainerSpec(op, "describe")
		Expect(err).NotTo(HaveOccurred())
	})
	It("rejects a nearest go.mod symlink that resolves outside the fetched repository", func() {
		outside := filepath.Join(data, "outside.mod")
		write(outside, "module outside\ngo 1.25\n")
		Expect(os.Symlink(outside, filepath.Join(repo, "go.mod"))).To(Succeed())
		_, _, err := ContainerSpecInRepository(repo, "describe", repo)
		Expect(err).To(MatchError(ContainSubstring("fetched repository")))
	})

	It("rejects local replacements and symlinks outside the fetched root with clear errors", func() {
		outside := filepath.Join(data, "outside")
		write(filepath.Join(outside, "go.mod"), "module outside\ngo 1.25\n")
		for _, target := range []string{"../../../../outside", outside, "./linked", "./missing"} {
			if target == "./linked" {
				Expect(os.Symlink(outside, filepath.Join(repo, "linked"))).To(Succeed())
			}
			write(filepath.Join(repo, "go.mod"), "module fetched\ngo 1.25\nreplace example.com/sdk => "+target+"\n")
			_, _, err := ContainerSpecInRepository(repo, "plan", repo)
			Expect(err).To(HaveOccurred(), target)
			Expect(err.Error()).To(ContainSubstring("replace"))
			Expect(err.Error()).To(ContainSubstring("fetched repository"))
		}
		write(filepath.Join(repo, "go.mod"), "module fetched\ngo 1.25\n")
		_, _, err := ContainerSpecInRepository(outside, "describe", repo)
		Expect(err).To(MatchError(ContainSubstring("fetched repository")))
	})
})
