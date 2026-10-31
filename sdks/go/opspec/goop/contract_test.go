package goop

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"os"
	"path/filepath"
)

var _ = Describe("Go op bootstrap contract", func() {
	It("requests plan with the supported schema and empty inputs and known outputs", func() {
		data, err := json.Marshal(NewRequest())
		Expect(err).NotTo(HaveOccurred())
		Expect(data).To(MatchJSON(`{"contractVersion":1,"phase":"plan","opspecSchema":{"min":1,"max":1},"inputs":{},"knownOutputs":{}}`))
	})
	It("reads the SDK's golden plan through the existing op file schema", func() {
		data, err := os.ReadFile("../sdk/testdata/plan.hello.json")
		Expect(err).NotTo(HaveOccurred())
		spec, err := ParseResponse(data)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Name).To(Equal("hello"))
		Expect(spec.Run.Container.Image.Ref).To(Equal("alpine:3"))
		Expect(spec.Run.Container.Cmd).To(Equal([]interface{}{"echo", "hello"}))
	})
	It("rejects incompatible versions and malformed responses", func() {
		for _, data := range []string{
			`{`, `{} {}`, `{"contractVersion":2,"opspecSchema":1}`, `{"contractVersion":1,"opspecSchema":0}`, `{"contractVersion":1,"opspecSchema":2}`,
		} {
			_, err := ParseResponse([]byte(data))
			Expect(err).To(HaveOccurred(), data)
		}
		_, err := ParseResponse([]byte(`{"contractVersion":1,"opspecSchema":2}`))
		Expect(err.Error()).To(ContainSubstring("Go op requires opspec schema 2; upgrade opctl"))
	})
	It("preserves op file validation errors", func() {
		invalid := json.RawMessage(`{"name":"hello","run":{"container":{}}}`)
		_, expected := opfile.Unmarshal(invalid)
		data, err := json.Marshal(Response{ContractVersion: 1, OpspecSchema: 1, Opspec: invalid})
		Expect(err).NotTo(HaveOccurred())
		_, err = ParseResponse(data)
		Expect(err).To(MatchError(expected))
	})
	It("rejects plan responses without a run and unsupported inputs or outputs", func() {
		for _, spec := range []string{`{"name":"hello"}`, `{"name":"hello","inputs":{"a":{"string":{}}},"run":{"container":{"image":{"ref":"alpine:3"}}}}`, `{"name":"hello","outputs":{"a":{"string":{}}},"run":{"container":{"image":{"ref":"alpine:3"}}}}`} {
			data, err := json.Marshal(Response{ContractVersion: 1, OpspecSchema: 1, Opspec: json.RawMessage(spec)})
			Expect(err).NotTo(HaveOccurred())
			_, err = ParseResponse(data)
			Expect(err).To(HaveOccurred())
		}
	})
})

var _ = Describe("Go op module and container spec", func() {
	var root string
	BeforeEach(func() { var err error; root, err = os.MkdirTemp("", "goop-"); Expect(err).NotTo(HaveOccurred()) })
	AfterEach(func() { Expect(os.RemoveAll(root)).To(Succeed()) })
	write := func(path, data string) {
		Expect(os.MkdirAll(filepath.Dir(path), 0700)).To(Succeed())
		Expect(os.WriteFile(path, []byte(data), 0600)).To(Succeed())
	}
	It("walks up to the nearest module and uses its Go directive verbatim", func() {
		write(filepath.Join(root, "go.mod"), "module parent\ngo 1.21\n")
		module := filepath.Join(root, "app")
		write(filepath.Join(module, "go.mod"), "module app\ngo 1.25.0\n")
		opDir := filepath.Join(module, "ops", "hello")
		Expect(os.MkdirAll(opDir, 0700)).To(Succeed())
		actual, err := ModuleRoot(opDir)
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).To(Equal(module))
		version, err := GoVersion(module)
		Expect(err).NotTo(HaveOccurred())
		Expect(version).To(Equal("1.25.0"))
		spec, mount, err := ContainerSpec(opDir)
		Expect(err).NotTo(HaveOccurred())
		Expect(mount).To(Equal(module))
		Expect(spec.Image.Ref).To(Equal("golang:1.25.0"))
		Expect(spec.Cmd).To(Equal([]interface{}{"go", "run", ".", "plan"}))
		Expect(spec.Dirs["/src"]).To(Equal("$(module)"))
		Expect(spec.Files["/opctl/request.json"]).To(Equal("$(request)"))
		Expect(spec.WorkDir).To(Equal("/src/ops/hello"))
		Expect(spec.EnvVars).To(Equal(map[string]interface{}{"GOFLAGS": "-mod=readonly", "GOTOOLCHAIN": "local", "CGO_ENABLED": "0", "OPCTL_BOOTSTRAP_REQUEST": "/opctl/request.json"}))
	})
	It("mounts the common ancestor of all local replace targets", func() {
		module := filepath.Join(root, "app", "ops", "hello")
		write(filepath.Join(module, "go.mod"), "module hello\ngo 1.21\nreplace example.com/sdk => ../../../sdk\nreplace example.com/other => "+filepath.Join(root, "shared")+"\nreplace example.com/remote => example.com/remote v1.0.0\n")
		_, mount, err := ContainerSpec(module)
		Expect(err).NotTo(HaveOccurred())
		Expect(mount).To(Equal(root))
		spec, _, err := ContainerSpec(module)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.WorkDir).To(Equal("/src/app/ops/hello"))
	})
	It("reports a missing module or Go directive and malformed module files", func() {
		_, err := ModuleRoot(root)
		Expect(err).To(MatchError("Go op needs a go.mod in the op dir or a parent"))
		for _, text := range []string{"module hello\n", "not a module file"} {
			write(filepath.Join(root, "go.mod"), text)
			_, _, err := ContainerSpec(root)
			Expect(err).To(HaveOccurred())
		}
	})
	It("rejects paths that would interpolate in a container working directory", func() {
		path := filepath.Join(root, "$(secret)")
		write(filepath.Join(path, "go.mod"), "module hello\ngo 1.21\n")
		_, _, err := ContainerSpec(path)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("$("))
	})
})
