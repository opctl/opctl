package k8s

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Describe("outputs file PVC mapping", func() {
	It("maps trusted shared scratch and rejects unshared or escaped paths", func() {
		image := "alpine"
		req := &model.ContainerCall{Image: &model.ContainerCallImage{Ref: &image}, OutputsFile: "/root/opctl/dcg/id/outputs/values.json", Files: map[string]string{"/opctl/outputs/values.json": "/root/opctl/dcg/id/outputs/values.json"}, EnvVars: map[string]string{"OPCTL_OUTPUTS": "/opctl/outputs/values.json"}}
		pod, err := constructPod(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(pod.Spec.Containers[0].VolumeMounts[1].SubPath).To(Equal("dcg/id/outputs/values.json"))
		Expect(pod.Spec.Containers[0].Env[0].Value).To(Equal("/opctl/outputs/values.json"))
		for _, path := range []string{"/tmp/outputs/values.json", "/root/opctl/../else/values.json", "/root/opctl-other/values.json"} {
			req.OutputsFile = path
			_, err = constructPod(req)
			Expect(err).To(MatchError(ContainSubstring("shared PVC")))
		}
	})
})
