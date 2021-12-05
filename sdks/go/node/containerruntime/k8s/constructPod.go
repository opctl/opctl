package k8s

import (
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func constructPod(
	req *model.ContainerCall,
) (*coreV1.Pod, error) {

	podName := constructPodName(req.ContainerID)

	container := coreV1.Container{
		Name:            podName,
		Image:           *req.Image.Ref,
		Command:         req.Cmd,
		WorkingDir:      req.WorkDir,
		ImagePullPolicy: coreV1.PullAlways,
		VolumeMounts: []coreV1.VolumeMount{
			{
				Name:      "opctl",
				MountPath: "/root/opctl",
			},
		},
	}

	for _, cmd := range req.Cmd {
		container.Command = append(
			container.Command,
			cmd,
		)
	}

	for envVarName, envVarValue := range req.EnvVars {
		container.Env = append(
			container.Env,
			coreV1.EnvVar{
				Name:  envVarName,
				Value: envVarValue,
			},
		)
	}

	pathPrefix := "/root/opctl/"
	for containerPath, hostPath := range req.Dirs {
		container.VolumeMounts = append(
			container.VolumeMounts,
			coreV1.VolumeMount{
				Name:      "opctl",
				MountPath: containerPath,
				SubPath:   strings.TrimPrefix(hostPath, pathPrefix),
			},
		)
	}

	for containerPath, hostPath := range req.Files {
		container.VolumeMounts = append(
			container.VolumeMounts,
			coreV1.VolumeMount{
				Name:      "opctl",
				MountPath: containerPath,
				SubPath:   strings.TrimPrefix(hostPath, pathPrefix),
			},
		)
	}

	return &coreV1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name: podName,
		},
		Spec: coreV1.PodSpec{
			Containers: []coreV1.Container{
				container,
			},
			RestartPolicy: coreV1.RestartPolicyNever,
			Volumes: []coreV1.Volume{
				{
					Name: "opctl",
					VolumeSource: coreV1.VolumeSource{
						PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
							ClaimName: "opctl",
						},
					},
				},
			},
		},
	}, nil
}
