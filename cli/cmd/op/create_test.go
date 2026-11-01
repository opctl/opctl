package op

import (
	"bytes"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("create an op", func() {
	var parent string
	BeforeEach(func() {
		var err error
		parent, err = os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() { Expect(os.RemoveAll(parent)).To(Succeed()) })
	create := func(args ...string) string {
		cmd := newCreateCmd()
		var stderr bytes.Buffer
		cmd.SetErr(&stderr)
		cmd.SetArgs(append(args, "-p", parent, "hello"))
		Expect(cmd.Execute()).To(Succeed())
		return stderr.String()
	}

	It("writes an op file without --go", func() {
		Expect(create("-d", "Greets the caller.")).To(BeEmpty())
		data, err := os.ReadFile(filepath.Join(parent, "hello", "op.yml"))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal("description: Greets the caller.\nname: hello\n"))
		Expect(filepath.Join(parent, "hello", "op.go")).NotTo(BeAnExistingFile())
	})

	It("scaffolds a Go op and explains the replace directive when this build has no published SDK", func() {
		Expect(create("--go", "-d", "Greets the caller.")).To(Equal("This opctl build has no published Go op SDK. Add a replace directive for github.com/opctl/opctl/sdks/go/opspec/sdk to go.mod before running the op.\n"))
		data, err := os.ReadFile(filepath.Join(parent, "hello", "op.go"))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(HavePrefix("// Greets the caller.\npackage main\n"))
		Expect(filepath.Join(parent, "hello", "go.sum")).NotTo(BeAnExistingFile())
	})

	It("requires the published SDK linked into release builds", func() {
		sdkVersion, sdkSum, sdkGoModSum = "v0.1.77", "h1:zip=", "h1:mod="
		defer func() { sdkVersion, sdkSum, sdkGoModSum = "", "", "" }()
		Expect(create("--go")).To(BeEmpty())
		data, err := os.ReadFile(filepath.Join(parent, "hello", "go.mod"))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(HaveSuffix("require github.com/opctl/opctl/sdks/go/opspec/sdk v0.1.77\n"))
		Expect(filepath.Join(parent, "hello", "go.sum")).To(BeAnExistingFile())
	})

	It("links the SDK release and its sums in the compile op", func() {
		data, err := os.ReadFile("../../.opspec/compile/op.yml")
		Expect(err).NotTo(HaveOccurred())
		for _, link := range []string{"sdkVersion=v$(version)", "sdkSum=$(sdkSum)", "sdkGoModSum=$(sdkGoModSum)"} {
			Expect(string(data)).To(ContainSubstring("-X=github.com/opctl/opctl/cli/cmd/op." + link))
		}
	})
})
