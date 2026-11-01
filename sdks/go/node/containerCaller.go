package node

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
	"github.com/opctl/opctl/sdks/go/opspec"
)

//counterfeiter:generate -o internal/fakes/containerCaller.go . containerCaller
type containerCaller interface {
	// Executes a container call
	Call(
		ctx context.Context,
		containerCall *model.ContainerCall,
		inboundScope map[string]*model.Value,
		containerCallSpec *model.ContainerCallSpec,
		rootCallID string,
	) (
		map[string]*model.Value,
		error,
	)
}

func newContainerCaller(
	containerRuntime containerruntime.ContainerRuntime,
	pubSub pubsub.PubSub,
	stateStore stateStore,
) containerCaller {

	return _containerCaller{
		containerRuntime: containerRuntime,
		pubSub:           pubSub,
		stateStore:       stateStore,
	}

}

type _containerCaller struct {
	containerRuntime containerruntime.ContainerRuntime
	pubSub           pubsub.PubSub
	stateStore       stateStore
}

func (cc _containerCaller) Call(
	ctx context.Context,
	containerCall *model.ContainerCall,
	inboundScope map[string]*model.Value,
	containerCallSpec *model.ContainerCallSpec,
	rootCallID string,
) (
	map[string]*model.Value,
	error,
) {
	if containerCallSpec.Stdout != "" || containerCallSpec.Stderr != "" || len(containerCallSpec.Outputs) > 0 {
		return cc.callWithResults(ctx, containerCall, inboundScope, containerCallSpec, rootCallID, nil)
	}
	outputs := map[string]*model.Value{}
	var exitCode int64

	if !containerCall.Image.Prepared && containerCall.Image.Ref != nil && containerCall.Image.PullCreds == nil {
		if auth := cc.stateStore.TryGetAuth(*containerCall.Image.Ref); auth != nil {
			containerCall.Image.PullCreds = &auth.Creds
		}
	}

	logStdOutPR, logStdOutPW := io.Pipe()
	logStdErrPR, logStdErrPW := io.Pipe()

	// interpret logs
	logChan := make(chan error, 1)
	go func() {
		logChan <- cc.interpretLogs(
			logStdOutPR,
			logStdErrPR,
			containerCall,
			rootCallID,
		)
	}()

	outputs = cc.interpretOutputs(
		containerCallSpec,
		containerCall,
	)

	var stdout io.WriteCloser = logStdOutPW
	var stderr io.WriteCloser = logStdErrPW
	var eventPublisher pubsub.EventPublisher = cc.pubSub
	capture := stdoutCaptureFrom(ctx)
	if capture != nil {
		stdout = teeCloser{Writer: capture, Closer: logStdOutPW}
		stderr = teeCloser{Writer: io.MultiWriter(logStdErrPW, &capture.diagnostics), Closer: logStdErrPW}
		eventPublisher = bootstrapEventPublisher{cc.pubSub}
	}

	rawExitCode, err := cc.containerRuntime.RunContainer(
		ctx,
		containerCall,
		rootCallID,
		eventPublisher,
		stdout,
		stderr,
	)

	if err == nil && capture != nil {
		err = capture.err
	}

	// @TODO: handle no exit code
	if rawExitCode != nil {
		exitCode = *rawExitCode
	}

	if exitCode != 0 && (capture == nil || err == nil) {
		err = fmt.Errorf("nonzero container exit code: %d", exitCode)
		if capture != nil {
			err = &OpDeclarationError{Err: err}
		}
	}

	// wait on logChan
	if logChanErr := <-logChan; err == nil {
		// non-destructively set err
		err = logChanErr
	}

	return outputs, err
}

func (this _containerCaller) interpretLogs(
	stdOutReader io.Reader,
	stdErrReader io.Reader,
	containerCall *model.ContainerCall,
	rootCallID string,
) error {
	stdOutLogChan := make(chan error, 1)
	go func() {
		// interpret stdOut
		stdOutLogChan <- readChunks(
			stdOutReader,
			func(chunk []byte) {
				this.pubSub.Publish(
					model.Event{
						Timestamp: time.Now().UTC(),
						ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{
							Data:        chunk,
							ContainerID: containerCall.ContainerID,
							OpRef:       containerCall.OpPath,
							RootCallID:  rootCallID,
						},
					},
				)
			})
	}()

	stdErrLogChan := make(chan error, 1)
	go func() {
		// interpret stdErr
		stdErrLogChan <- readChunks(
			stdErrReader,
			func(chunk []byte) {
				this.pubSub.Publish(
					model.Event{
						Timestamp: time.Now().UTC(),
						ContainerStdErrWrittenTo: &model.ContainerStdErrWrittenTo{
							Data:        chunk,
							ContainerID: containerCall.ContainerID,
							OpRef:       containerCall.OpPath,
							RootCallID:  rootCallID,
						},
					},
				)
			})
	}()

	// wait on logs
	stdOutLogErr := <-stdOutLogChan
	stdErrLogErr := <-stdErrLogChan

	// return errs
	if stdOutLogErr != nil {
		return stdOutLogErr
	}
	if stdErrLogErr != nil {
		return stdErrLogErr
	}

	return nil
}

func (this _containerCaller) interpretOutputs(
	containerCallSpec *model.ContainerCallSpec,
	containerCall *model.ContainerCall,
) map[string]*model.Value {
	outputs := map[string]*model.Value{}

	for socketAddr, name := range containerCallSpec.Sockets {
		// add socket outputs
		if "0.0.0.0" == socketAddr {
			outputs[name] = &model.Value{Socket: &containerCall.ContainerID}
		}
	}
	for callSpecContainerFilePath, mountSrc := range containerCallSpec.Files {
		mountSrcStr, ok := mountSrc.(string)
		if !ok {
			continue
		}

		if mountSrcStr == "" {
			// skip embedded files
			continue
		}

		// add file outputs
		for callContainerFilePath, callHostFilePath := range containerCall.Files {
			if callSpecContainerFilePath == callContainerFilePath {
				// copy callHostFilePath before taking address; range vars have same address for every iteration
				value := callHostFilePath
				outputs[opspec.RefToName(mountSrcStr)] = &model.Value{File: &value}
			}
		}
	}
	for callSpecContainerDirPath, mountSrc := range containerCallSpec.Dirs {
		mountSrcStr, ok := mountSrc.(string)
		if !ok {
			continue
		}

		if mountSrcStr == "" {
			// skip embedded dirs
			continue
		}

		// add dir outputs
		for callContainerDirPath, callHostDirPath := range containerCall.Dirs {
			if callSpecContainerDirPath == callContainerDirPath {
				// copy callHostDirPath before taking address; range vars have same address for every iteration
				value := callHostDirPath
				outputs[opspec.RefToName(mountSrcStr)] = &model.Value{Dir: &value}
			}
		}
	}

	return outputs
}
