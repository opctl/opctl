package container

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"os"
)

var _ = Describe("container stream interpretation", func() {
	It("preserves old spec and interpreted JSON byte for byte", func() {
		dir, err := os.MkdirTemp("", "old-container-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		spec := &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}}
		encoded, err := json.Marshal(spec)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(encoded)).To(Equal(`{"image":{"ref":"alpine"}}`))
		call, err := Interpret(nil, spec, "id", "op", dir)
		Expect(err).NotTo(HaveOccurred())
		encoded, err = json.Marshal(call)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(encoded)).To(Equal(`{"opPath":"op","containerId":"id","cmd":[],"dirs":{},"envVars":null,"files":{},"image":{"ref":"docker.io/library/alpine"},"sockets":{},"workDir":""}`))
	})
	It("validates targets without resolving them and rejects new collisions", func() {
		dir, err := os.MkdirTemp("", "stream-bindings-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		spec := &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}, Stdout: "$(fresh)", Stderr: "$(diagnostic)"}
		_, err = Interpret(nil, spec, "id", "op", dir)
		Expect(err).NotTo(HaveOccurred())
		for _, bad := range []string{"bare", "$(v.key)", "$(v/path)", "$($(v))"} {
			spec.Stdout = bad
			_, err = Interpret(nil, spec, "id", "op", dir)
			Expect(err).To(HaveOccurred())
		}
		spec.Stdout = spec.Stderr
		_, err = Interpret(nil, spec, "id", "op", dir)
		Expect(err).To(MatchError(ContainSubstring("duplicate")))
		spec.Stdout = "$(fresh)"
		spec.Files = map[string]interface{}{"/out": "$(fresh)"}
		_, err = Interpret(nil, spec, "id", "op", dir)
		Expect(err).To(MatchError(ContainSubstring("duplicate")))
		spec.Files = nil
		spec.Sockets = map[string]string{"0.0.0.0": "fresh"}
		_, err = Interpret(nil, spec, "id", "op", dir)
		Expect(err).To(MatchError(ContainSubstring("duplicate")))

	})
})
