package container

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"os"
)

var _ = Describe("outputs file interpretation", func() {
	It("reserves only opted-in environments and overlapping mount destinations", func() {
		dir, err := os.MkdirTemp("", "outputs-interpret-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		scope := map[string]*model.Value{}
		spec := &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}, Outputs: map[string]string{"version": "$(v)"}}
		call, err := Interpret(scope, spec, "id", "op", dir)
		Expect(err).NotTo(HaveOccurred())
		Expect(call.OutputsFile).To(Equal(dir + "/dcg/id/outputs/values.json"))
		Expect(call.EnvVars).NotTo(HaveKey("OPCTL_OUTPUTS"))
		Expect(call.Dirs).NotTo(HaveKey("/opctl/outputs"))
		name := "OPCTL_OUTPUTS"
		scope["name"] = &model.Value{String: &name}
		spec.EnvVars = map[string]interface{}{"$(name)": "authored"}
		_, err = Interpret(scope, spec, "id", "op", dir)
		Expect(err).To(MatchError(ContainSubstring("OPCTL_OUTPUTS")))
		spec.EnvVars = nil
		for _, path := range []string{"/", "/opctl", "/opctl/outputs", "/opctl/outputs/child", "/tmp/../opctl/outputs"} {
			spec.Dirs = map[string]interface{}{path: map[string]interface{}{}}
			_, err = Interpret(scope, spec, "id", "op", dir)
			Expect(err).To(HaveOccurred(), path)
		}
		spec.Dirs = nil
		spec.Stdout = "$(v)"
		_, err = Interpret(scope, spec, "id", "op", dir)
		Expect(err).To(MatchError(ContainSubstring("duplicate")))
		spec.Stdout = ""
		spec.Outputs = nil
		spec.EnvVars = map[string]interface{}{"OPCTL_OUTPUTS": "old"}
		_, err = Interpret(scope, spec, "id", "op", dir)
		Expect(err).NotTo(HaveOccurred())
	})
})
