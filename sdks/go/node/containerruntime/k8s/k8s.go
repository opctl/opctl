package k8s

import (
	"context"
	"fmt"
	"io"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/pubsub"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func New() (
	containerRuntime containerruntime.ContainerRuntime,
	err error,
) {

	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}

	return _containerRuntime{
		k8sClient: k8sClient,
	}, nil
}

type _containerRuntime struct {
	k8sClient *kubernetes.Clientset
}

func (cr _containerRuntime) Delete(
	ctx context.Context,
) error {
	// for now this is a no-op
	return nil
}

func (cr _containerRuntime) DeleteContainerIfExists(
	ctx context.Context,
	containerID string,
) error {
	if err := cr.k8sClient.CoreV1().Pods("opctl").Delete(
		ctx,
		constructPodName(containerID),
		metaV1.DeleteOptions{},
	); err != nil {
		return fmt.Errorf("unable to delete k8s container: %w", err)
	}

	return nil
}

func (cr _containerRuntime) RunContainer(
	ctx context.Context,
	req *model.ContainerCall,
	rootCallID string,
	eventPublisher pubsub.EventPublisher,
	stdout io.WriteCloser,
	stderr io.WriteCloser,
) (*int64, error) {
	defer stdout.Close()
	defer stderr.Close()

	pod, err := constructPod(req)
	if err != nil {
		return nil, err
	}

	_, err = cr.k8sClient.CoreV1().Pods("opctl").Create(
		ctx,
		pod,
		metaV1.CreateOptions{},
	)
	if err != nil {
		return nil, err
	}

	watcher, err := cr.k8sClient.CoreV1().Pods("opctl").Watch(
		ctx,
		metaV1.ListOptions{
			FieldSelector: fmt.Sprintf("metadata.name=%s", pod.ObjectMeta.Name),
		},
	)
	if err != nil {
		return nil, err
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		var ok bool
		pod, ok = event.Object.(*coreV1.Pod)
		if !ok {
			continue
		}

		switch pod.Status.Phase {
		case coreV1.PodRunning:
			// https://stackoverflow.com/questions/53852530/how-to-get-logs-from-kubernetes-using-golang
			logsResult := cr.k8sClient.CoreV1().Pods("opctl").GetLogs(
				pod.ObjectMeta.Name,
				&coreV1.PodLogOptions{
					Follow: true,
				},
			)
			logSrc, err := logsResult.Stream(ctx)
			if err != nil {
				return nil, fmt.Errorf("unable to stream running pod logs from k8s: %w", err)
			}
			defer logSrc.Close()

			_, err = io.Copy(stdout, logSrc)
			if err != nil {
				return nil, err
			}

		// https://medium.com/programming-kubernetes/building-stuff-with-the-kubernetes-api-part-4-using-go-b1d0e3c1c899
		case coreV1.PodSucceeded:
			zero := int64(0)
			return &zero, nil
		case coreV1.PodFailed:
			// https://stackoverflow.com/questions/53852530/how-to-get-logs-from-kubernetes-using-golang
			logsResult := cr.k8sClient.CoreV1().Pods("opctl").GetLogs(
				pod.ObjectMeta.Name,
				&coreV1.PodLogOptions{
					Follow: true,
				},
			)
			logSrc, err := logsResult.Stream(ctx)
			if err != nil {
				return nil, fmt.Errorf("unable to stream failed pod logs from k8s: %w", err)
			}
			defer logSrc.Close()
			_, err = io.Copy(stdout, logSrc)
			if err != nil {
				return nil, err
			}

			exitCode := int64(pod.Status.ContainerStatuses[0].State.Terminated.ExitCode)
			return &exitCode, fmt.Errorf(
				"%s, %s",
				pod.Status.ContainerStatuses[0].State.Terminated.Reason,
				pod.Status.ContainerStatuses[0].State.Terminated.Message,
			)
		}
	}

	return nil, err
}
