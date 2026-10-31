package k8s

import (
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	coreV1 "k8s.io/api/core/v1"
)

func TestConstructPodCacheDirs(t *testing.T) {
	for _, subPath := range []string{
		"caches/gomod",
		"nocache/run-id/gomod",
	} {
		t.Run(subPath, func(t *testing.T) {
			imageRef := "alpine:3"
			pod, err := constructPod(&model.ContainerCall{
				ContainerID: "cache-test",
				Image:       &model.ContainerCallImage{Ref: &imageRef},
				Dirs:        map[string]string{"/go/pkg/mod": "/root/opctl/" + subPath},
			})
			if err != nil {
				t.Fatal(err)
			}

			mounts := pod.Spec.Containers[0].VolumeMounts
			if len(mounts) != 2 {
				t.Fatalf("expected data dir and cache dir mounts, got %v", mounts)
			}
			expectedMount := coreV1.VolumeMount{
				Name:      "opctl",
				MountPath: "/go/pkg/mod",
				SubPath:   subPath,
			}
			if mounts[1] != expectedMount {
				t.Errorf("expected writable cache mount %+v, got %+v", expectedMount, mounts[1])
			}
			claim := pod.Spec.Volumes[0].PersistentVolumeClaim
			if claim == nil || claim.ClaimName != "opctl" || claim.ReadOnly {
				t.Errorf("expected writable opctl PVC, got %+v", claim)
			}
		})
	}
}
