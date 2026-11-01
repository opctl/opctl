package goop

import (
	"encoding/json"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	containerInterpreter "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

var _ = Describe("Go op bootstrap contract", func() {
	It("requests plan with the supported schema and empty inputs and known outputs", func() {
		data, err := json.Marshal(NewRequest("plan", nil))
		Expect(err).NotTo(HaveOccurred())
		Expect(data).To(MatchJSON(`{"contractVersion":1,"phase":"plan","opspecSchema":{"min":1,"max":4},"inputs":{},"knownOutputs":{}}`))
	})
	It("reads the SDK's golden plan through the existing op file schema", func() {
		data, err := os.ReadFile("../sdk/testdata/plan.hello.json")
		Expect(err).NotTo(HaveOccurred())
		spec, err := ParsePlan(data)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Name).To(Equal("hello"))
		Expect(spec.Run.Container.Image.Ref).To(Equal("alpine:3"))
		Expect(spec.Run.Container.Cmd).To(Equal([]interface{}{"echo", "hello"}))
	})
	It("reads serial and parallel calls from the SDK's composed plan", func() {
		data, err := os.ReadFile("../sdk/testdata/plan.compose.json")
		Expect(err).NotTo(HaveOccurred())
		spec, err := ParsePlan(data)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Name).To(Equal("compose"))
		Expect(spec.Run.Serial).NotTo(BeNil())
		Expect(*spec.Run.Serial).To(HaveLen(2))
		Expect((*spec.Run.Serial)[0].Parallel).NotTo(BeNil())
		Expect(*(*spec.Run.Serial)[0].Parallel).To(HaveLen(2))
		Expect((*(*spec.Run.Serial)[0].Parallel)[1].Container.Cmd).To(Equal([]interface{}{"echo", "b"}))
		Expect((*spec.Run.Serial)[1].Container.Cmd).To(Equal([]interface{}{"echo", "done"}))
	})
	It("rejects incompatible versions and malformed responses", func() {
		for _, data := range []string{
			`{`, `{} {}`, `{"contractVersion":2,"opspecSchema":1}`, `{"contractVersion":1,"opspecSchema":0}`, `{"contractVersion":1,"opspecSchema":999}`,
		} {
			_, err := ParsePlan([]byte(data))
			Expect(err).To(HaveOccurred(), data)
		}
		_, err := ParsePlan([]byte(`{"contractVersion":1,"opspecSchema":999}`))
		Expect(err.Error()).To(ContainSubstring("Go op requires opspec schema 999; upgrade opctl"))
	})
	It("preserves op file validation errors", func() {
		invalid := json.RawMessage(`{"name":"hello","run":{"container":{}}}`)
		_, expected := opfile.Unmarshal(invalid)
		data, err := json.Marshal(Response{ContractVersion: 1, OpspecSchema: 1, Opspec: invalid})
		Expect(err).NotTo(HaveOccurred())
		_, err = ParsePlan(data)
		Expect(err).To(MatchError(expected))
	})
	It("rejects plan responses without a run and unsupported inputs or outputs", func() {
		for _, spec := range []string{`{"name":"hello"}`, `{"name":"hello","inputs":{"a":{"string":{}}},"run":{"container":{"image":{"ref":"alpine:3"}}}}`, `{"name":"hello","outputs":{"a":{"string":{}}},"run":{"container":{"image":{"ref":"alpine:3"}}}}`} {
			data, err := json.Marshal(Response{ContractVersion: 1, OpspecSchema: 1, Opspec: json.RawMessage(spec)})
			Expect(err).NotTo(HaveOccurred())
			_, err = ParsePlan(data)
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
		spec, mount, err := sourceSpec(opDir, "plan")
		Expect(err).NotTo(HaveOccurred())
		Expect(mount).To(Equal(module))
		Expect(spec.Image.Ref).To(Equal("golang:1.25.0"))
		Expect(spec.Dirs[module]).To(Equal("$(module)"))
		Expect(spec.Files["/opctl/request.json"]).To(Equal("$(request)"))
		Expect(spec.WorkDir).To(Equal(opDir))
		Expect(spec.EnvVars).To(Equal(map[string]interface{}{"GOFLAGS": "-mod=readonly", "GOTOOLCHAIN": "local", "CGO_ENABLED": "0", "OPCTL_BOOTSTRAP_REQUEST": "/opctl/request.json"}))
	})
	It("mounts the module and each local replacement at their own host paths", func() {
		module := filepath.Join(root, "app", "ops", "hello")
		sdk := filepath.Join(root, "sdk")
		shared := filepath.Join(root, "absolute-replacement")
		Expect(os.MkdirAll(sdk, 0700)).To(Succeed())
		Expect(os.MkdirAll(shared, 0700)).To(Succeed())
		write(filepath.Join(module, "go.mod"), "module hello\ngo 1.21\nreplace example.com/sdk => ../../../sdk\nreplace example.com/other => "+shared+"\nreplace example.com/remote => example.com/remote v1.0.0\n")
		spec, mount, err := sourceSpec(module, "plan")
		Expect(err).NotTo(HaveOccurred())
		Expect(mount).To(Equal(module))
		Expect(spec.WorkDir).To(Equal(module))
		Expect(spec.Dirs).To(HaveLen(3))
		requestPath := filepath.Join(root, "request.json")
		write(requestPath, "{}")
		actual, err := containerInterpreter.Interpret(map[string]*model.Value{"module": {Dir: &mount}, "request": {File: &requestPath}}, spec, "bootstrap", module, filepath.Join(root, "cache"))
		Expect(err).NotTo(HaveOccurred())
		Expect(actual.Dirs).To(Equal(map[string]string{module: module, sdk: sdk, shared: shared}))
		Expect(actual.WorkDir).To(Equal(module))
	})
	It("rejects replacement mounts that hide required bootstrap image directories", func() {
		for _, target := range []string{"/", "/usr", "/usr/local", "/usr/local/..", "/go", "/opctl"} {
			write(filepath.Join(root, "go.mod"), "module hello\ngo 1.21\nreplace example.com/sdk => "+target+"\n")
			_, _, err := sourceSpec(root, "plan")
			Expect(err).To(HaveOccurred(), target)
			Expect(err.Error()).To(ContainSubstring("mount destination"))
			Expect(err.Error()).To(ContainSubstring(filepath.Clean(target)))
		}
	})
	It("allows replacements below image directories and duplicate module replacements", func() {
		write(filepath.Join(root, "go.mod"), "module hello\ngo 1.21\nreplace example.com/sdk => ./sdk\nreplace example.com/self => .\n")
		spec, _, err := sourceSpec(root, "plan")
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Dirs).To(HaveLen(1))
		Expect(spec.Dirs[root]).To(Equal("$(module)"))
	})
	It("uses the root mount for replacements within the module, including local links", func() {
		write(filepath.Join(root, "go.mod"), "module app\ngo 1.21\nreplace example.com/lib => ./lib\n")
		outside := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-outside")
		Expect(os.MkdirAll(outside, 0700)).To(Succeed())
		defer os.RemoveAll(outside)
		Expect(os.Symlink(outside, filepath.Join(root, "lib"))).To(Succeed())
		spec, _, err := sourceSpec(root, "plan")
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Dirs).To(HaveLen(1))
	})
	It("rejects symlinked separate replacement mounts and module roots", func() {
		module := filepath.Join(root, "app")
		write(filepath.Join(module, "go.mod"), "module app\ngo 1.21\nreplace example.com/lib => ../link\n")
		Expect(os.MkdirAll(filepath.Join(root, "private"), 0700)).To(Succeed())
		Expect(os.Symlink(filepath.Join(root, "private"), filepath.Join(root, "link"))).To(Succeed())
		_, _, err := sourceSpec(module, "plan")
		Expect(err).To(MatchError(ContainSubstring("symlink")))
		Expect(os.Symlink(module, filepath.Join(root, "alias"))).To(Succeed())
		_, _, err = sourceSpec(filepath.Join(root, "alias"), "plan")
		Expect(err).To(MatchError(ContainSubstring("symlink")))
	})

	It("reports a missing module or Go directive and malformed module files", func() {
		_, err := ModuleRoot(root)
		Expect(err).To(MatchError("Go op needs a go.mod in the op dir or a parent"))
		for _, text := range []string{"module hello\n", "not a module file"} {
			write(filepath.Join(root, "go.mod"), text)
			_, _, err := sourceSpec(root, "plan")
			Expect(err).To(HaveOccurred())
		}
	})
	It("rejects paths that would interpolate in a container working directory", func() {
		path := filepath.Join(root, "$(secret)")
		write(filepath.Join(path, "go.mod"), "module hello\ngo 1.21\n")
		_, _, err := sourceSpec(path, "plan")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("$("))
	})
})

var _ = Describe("Go op interfaces and public plan inputs", func() {
	It("parses the SDK's typed input interface and plan independently", func() {
		data, err := os.ReadFile("../sdk/testdata/describe.inputs.json")
		Expect(err).NotTo(HaveOccurred())
		spec, err := ParseDescribe(data)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Inputs["token"].String.IsSecret).To(BeTrue())
		Expect(spec.Inputs["count"].Number.Default).To(BeNumerically("==", 2))
		_, err = ParsePlan(data)
		Expect(err).To(HaveOccurred())
		data, err = os.ReadFile("../sdk/testdata/plan.inputs.json")
		Expect(err).NotTo(HaveOccurred())
		_, err = ParsePlan(data)
		Expect(err).NotTo(HaveOccurred())
		_, err = ParseDescribe(data)
		Expect(err).To(HaveOccurred())
	})
	It("parses the description the SDK derives from the package doc comment", func() {
		data, err := os.ReadFile("../sdk/testdata/describe.metadata.json")
		Expect(err).NotTo(HaveOccurred())
		spec, err := ParseDescribe(data)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Name).To(Equal("greeter"))
		Expect(spec.Description).To(Equal("Greets the caller.\n\nRuns echo in alpine:\n\n\tsh -ce 'echo hello'"))
		Expect(spec.Inputs).To(HaveKey("who"))
	})
	It("parses op-directory handles and dir subpaths from the SDK's plan", func() {
		data, err := os.ReadFile("../sdk/testdata/plan.op-dir.json")
		Expect(err).NotTo(HaveOccurred())
		spec, err := ParsePlan(data)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Run.Container.Dirs).To(Equal(map[string]interface{}{
			"/cache": "$(HOME/.cache/go-build/v1)", "/go/pkg": "$(HOME/go/pkg)", "/nested": "$(../../sub/dir)", "/scripts": "$(./scripts)", "/src": "$(../..)",
		}))
		Expect(spec.Run.Container.Files).To(Equal(map[string]interface{}{"/config.yml": "$(src/config/app.yml)", "/expected.txt": "$(./expected.txt)"}))
	})
	It("keeps secret and filesystem values out of plan requests", func() {
		text := "private"
		number := float64(3)
		object := map[string]interface{}{"nested": model.Value{Number: &number}}
		params := map[string]*model.ParamSpec{
			"token": {String: &model.StringParamSpec{IsSecret: true}}, "dir": {Dir: &model.DirParamSpec{}}, "file": {File: &model.FileParamSpec{}}, "socket": {Socket: &model.SocketParamSpec{}},
			"number": {Number: &model.NumberParamSpec{}}, "object": {Object: &model.ObjectParamSpec{}}, "secretNumber": {Number: &model.NumberParamSpec{IsSecret: true}}, "secretArray": {Array: &model.ArrayParamSpec{IsSecret: true}}, "secretObject": {Object: &model.ObjectParamSpec{IsSecret: true}},
		}
		values := map[string]*model.Value{"token": {String: &text}, "dir": {Dir: &text}, "file": {File: &text}, "socket": {Socket: &text}, "number": {Number: &number}, "object": {Object: &object}}
		public, err := PlanInputs(params, values)
		Expect(err).NotTo(HaveOccurred())
		Expect(public).To(Equal(map[string]interface{}{"number": number, "object": map[string]interface{}{"nested": number}}))
		_, err = PlanInputs(map[string]*model.ParamSpec{"missing": {String: &model.StringParamSpec{}}}, nil)
		Expect(err).To(HaveOccurred())
		_, err = PlanInputs(map[string]*model.ParamSpec{"bad": {String: &model.StringParamSpec{}}}, map[string]*model.Value{"bad": {}})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).NotTo(ContainSubstring("private"))
	})
	It("rejects run in describe", func() {
		for _, spec := range []string{`{"run":{"container":{"image":{"ref":"alpine"}}}}`} {
			data, err := json.Marshal(Response{ContractVersion: 1, OpspecSchema: 1, Opspec: json.RawMessage(spec)})
			Expect(err).NotTo(HaveOccurred())
			_, err = ParseDescribe(data)
			Expect(err).To(HaveOccurred())
		}
	})
})

var _ = Describe("Go op call contract", func() {
	It("parses SDK op calls with the existing opspec schema", func() {
		data, err := os.ReadFile("../sdk/testdata/plan.op-call.json")
		Expect(err).NotTo(HaveOccurred())
		spec, err := ParsePlan(data)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Run.Op.Ref).To(Equal("./yaml"))
	})
})

var _ = Describe("Go op output call contract", func() {
	It("uses callee output names as keys and caller references as values", func() {
		data, err := os.ReadFile("../sdk/testdata/plan.op-output.json")
		Expect(err).NotTo(HaveOccurred())
		spec, err := ParsePlan(data)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Run.Op.Outputs).To(Equal(map[string]string{"artifact": "$(opctl_dir_1)", "label": "$(opctl_string_2)"}))
	})
})

var _ = Describe("Go output binding contract", func() {
	It("accepts output declarations and plan bindings separately", func() {
		desc, err := ParseDescribe([]byte(`{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"hello","outputs":{"out":{"string":{}}}}}`))
		Expect(err).NotTo(HaveOccurred())
		Expect(desc.Outputs).To(HaveKey("out"))
		plan, err := ParsePlan([]byte(`{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"hello","run":{"op":{"ref":"./yaml"}}},"outputBindings":{"out":"$(opctl_string_1)"}}`))
		Expect(err).NotTo(HaveOccurred())
		Expect(plan.OutputBindings).To(HaveKeyWithValue("out", "$(opctl_string_1)"))
		Expect(plan.ValidateOutputs(desc.Outputs)).To(Succeed())
	})
	It("rejects malformed bindings and declarations in the wrong phase", func() {
		for _, bindings := range []string{`{"bad/name":"$(x)"}`, `{"opctl_private":"$(x)"}`, `{"out":"literal"}`, `{"out":"$(x) trailing"}`, `{"out":null}`} {
			_, err := ParsePlan([]byte(`{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"hello","run":{"op":{"ref":"x"}}},"outputBindings":` + bindings + `}`))
			Expect(err).To(HaveOccurred(), bindings)
		}
		_, err := ParseDescribe([]byte(`{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"hello"},"outputBindings":{"out":"$(x)"}}`))
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("Go output fixtures", func() {
	It("accepts the shared output declaration and binding fixtures", func() {
		data, err := os.ReadFile("../sdk/testdata/describe.outputs.json")
		Expect(err).NotTo(HaveOccurred())
		desc, err := ParseDescribe(data)
		Expect(err).NotTo(HaveOccurred())
		data, err = os.ReadFile("../sdk/testdata/plan.outputs.json")
		Expect(err).NotTo(HaveOccurred())
		plan, err := ParsePlan(data)
		Expect(err).NotTo(HaveOccurred())
		Expect(plan.ValidateOutputs(desc.Outputs)).To(Succeed())
	})
})
