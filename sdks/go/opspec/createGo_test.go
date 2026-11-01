package opspec

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("CreateGo", func() {
	var parent string
	BeforeEach(func() {
		var err error
		parent, err = os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() { Expect(os.RemoveAll(parent)).To(Succeed()) })
	read := func(path string) string {
		data, err := os.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())
		return string(data)
	}

	It("writes a Go op in its own module that needs a replace without a published SDK", func() {
		path := filepath.Join(parent, "hello")
		Expect(CreateGo(path, "hello", "", GoSDK{})).To(Succeed())

		Expect(read(filepath.Join(path, "op.go"))).To(Equal(`package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
	Who opctl.String ` + "`" + `opctl:"who,description='Who to greet',default=world"` + "`" + `
}

func plan(in Inputs) opctl.Call {
	return opctl.Container("alpine:3").
		Cmd("sh", "-ce", ` + "`" + `echo "hello $WHO"` + "`" + `).
		EnvFrom("WHO", in.Who)
}

func main() { opctl.Main(plan) }
`))
		Expect(read(filepath.Join(path, "go.mod"))).To(Equal("module ops.local/hello\n\ngo " + goDirective() + "\n\nrequire github.com/opctl/opctl/sdks/go/opspec/sdk v0.0.0\n"))
		Expect(read(filepath.Join(path, ".gitignore"))).To(Equal("/opctlgen/\n"))
		Expect(filepath.Join(path, "go.sum")).NotTo(BeAnExistingFile())
		Expect(filepath.Join(path, "op.yml")).NotTo(BeAnExistingFile())
	})

	It("requires a published SDK and writes its sums", func() {
		path := filepath.Join(parent, "hello")
		Expect(CreateGo(path, "hello", "", GoSDK{Version: "v0.1.77", Sum: "h1:zip=", GoModSum: "h1:mod="})).To(Succeed())
		Expect(read(filepath.Join(path, "go.mod"))).To(HaveSuffix("\nrequire github.com/opctl/opctl/sdks/go/opspec/sdk v0.1.77\n"))
		Expect(read(filepath.Join(path, "go.sum"))).To(Equal(
			"github.com/opctl/opctl/sdks/go/opspec/sdk v0.1.77 h1:zip=\n" +
				"github.com/opctl/opctl/sdks/go/opspec/sdk v0.1.77/go.mod h1:mod=\n",
		))
	})

	It("builds a scaffold whose name matches a standard library package", func() {
		path := filepath.Join(parent, "fmt")
		Expect(CreateGo(path, "fmt", "", GoSDK{})).To(Succeed())
		sdk, err := filepath.Abs("sdk")
		Expect(err).NotTo(HaveOccurred())
		goMod := read(filepath.Join(path, "go.mod"))
		Expect(os.WriteFile(filepath.Join(path, "go.mod"), []byte(goMod+"\nreplace "+goSDKModule+" => "+sdk+"\n"), 0600)).To(Succeed())

		cmd := exec.Command("go", "build", "-o", filepath.Join(parent, "op"), ".")
		cmd.Dir = path
		cmd.Env = append(os.Environ(), "GOFLAGS=-mod=readonly -buildvcs=false", "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
		output, err := cmd.CombinedOutput()
		Expect(err).NotTo(HaveOccurred(), string(output))
		Expect(goMod).To(HavePrefix("module ops.local/fmt\n"))
	})

	It("treats an SDK release without both sums as unpublished", func() {
		for _, sdk := range []GoSDK{{Version: "v0.1.77"}, {Version: "v0.1.77", Sum: "h1:zip="}, {Sum: "h1:zip=", GoModSum: "h1:mod="}} {
			Expect(sdk.Published()).To(BeFalse())
			path := filepath.Join(parent, "op"+sdk.Sum+sdk.Version)
			Expect(CreateGo(path, "hello", "", sdk)).To(Succeed())
			Expect(read(filepath.Join(path, "go.mod"))).To(HaveSuffix(" v0.0.0\n"))
			Expect(filepath.Join(path, "go.sum")).NotTo(BeAnExistingFile())
		}
	})

	It("writes each line of the description as the package doc comment", func() {
		path := filepath.Join(parent, "hello")
		Expect(CreateGo(path, "hello", "Greets the caller.\n\nUses alpine.", GoSDK{})).To(Succeed())
		Expect(read(filepath.Join(path, "op.go"))).To(HavePrefix("// Greets the caller.\n//\n// Uses alpine.\npackage main\n"))
	})

	It("never overwrites an existing op or module file", func() {
		for _, existing := range []string{"op.yml", "op.go", "go.mod", "go.sum", ".gitignore"} {
			path := filepath.Join(parent, existing)
			Expect(os.MkdirAll(path, 0700)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(path, existing), []byte("keep"), 0600)).To(Succeed())

			err := CreateGo(path, "hello", "", GoSDK{Version: "v0.1.77", Sum: "h1:zip=", GoModSum: "h1:mod="})

			Expect(err).To(MatchError(ContainSubstring(existing + " already exists")))
			entries, err := os.ReadDir(path)
			Expect(err).NotTo(HaveOccurred())
			Expect(entries).To(HaveLen(1))
			Expect(read(filepath.Join(path, existing))).To(Equal("keep"))
		}
	})

	It("rejects op names that aren't Go module paths", func() {
		path := filepath.Join(parent, "bad")
		Expect(CreateGo(path, "bad name", "", GoSDK{})).To(MatchError(`op name "bad name" is not a valid Go module path`))
		Expect(path).NotTo(BeAnExistingFile())
	})

	It("selects the go directive from the toolchain that built opctl", func() {
		for version, expected := range map[string]string{
			"go1.25.3": "1.25", "go1.26": "1.26", "go1.26rc1": "1.26", "devel go1.27-abc": "1.25", "": "1.25",
		} {
			Expect(goDirectiveFor(version)).To(Equal(expected), version)
		}
	})
})
