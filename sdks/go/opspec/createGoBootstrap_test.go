package opspec_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

// bootstrapLocally runs a Go op's describe and plan phases with `go run`, the
// way bootstrap runs them, but without a node.
func bootstrapLocally(opDir string, inputs map[string]interface{}) (*model.OpSpec, *goop.Plan) {
	phase := func(name string, inputs map[string]interface{}) []byte {
		request, err := json.Marshal(goop.NewRequest(name, inputs))
		Expect(err).NotTo(HaveOccurred())
		dir, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		requestPath := filepath.Join(dir, "request.json")
		Expect(os.WriteFile(requestPath, request, 0600)).To(Succeed())
		cmd := exec.Command("go", "run", ".", name)
		cmd.Dir = opDir
		cmd.Env = append(os.Environ(), "GOFLAGS=-mod=readonly", "GOTOOLCHAIN=local", "CGO_ENABLED=0", "GOWORK=off", "OPCTL_BOOTSTRAP_REQUEST="+requestPath)
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		Expect(cmd.Run()).To(Succeed(), stderr.String())
		return stdout.Bytes()
	}
	described, err := goop.ParseDescribe(phase("describe", nil))
	Expect(err).NotTo(HaveOccurred())
	planned, err := goop.ParsePlan(phase("plan", inputs))
	Expect(err).NotTo(HaveOccurred())
	return described, planned
}

var _ = Context("CreateGo bootstrap", func() {
	It("scaffolds an op that describes and plans with the in-repo SDK", func() {
		sdk, err := filepath.Abs("sdk")
		Expect(err).NotTo(HaveOccurred())
		parent, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(parent)
		path := filepath.Join(parent, "hello")
		Expect(opspec.CreateGo(path, "hello", "Greets the caller.", opspec.GoSDK{})).To(Succeed())
		goMod, err := os.OpenFile(filepath.Join(path, "go.mod"), os.O_APPEND|os.O_WRONLY, 0)
		Expect(err).NotTo(HaveOccurred())
		_, err = goMod.WriteString("\nreplace github.com/opctl/opctl/sdks/go/opspec/sdk => " + sdk + "\n")
		Expect(err).NotTo(HaveOccurred())
		Expect(goMod.Close()).To(Succeed())

		described, planned := bootstrapLocally(path, map[string]interface{}{"who": "world"})

		Expect(described.Name).To(Equal("hello"))
		Expect(described.Description).To(Equal("Greets the caller."))
		Expect(described.Inputs).To(HaveKey("who"))
		Expect(described.Inputs["who"].String.Default).To(Equal("world"))
		Expect(planned.Run.Container.Image.Ref).To(Equal("alpine:3"))
		Expect(planned.Run.Container.EnvVars).To(Equal(map[string]interface{}{"WHO": "$(who)"}))
	})
})
