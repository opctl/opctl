package core

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/golang-interfaces/iio"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o internal/fakes/containerCaller.go . containerCaller
type containerCaller interface {
	// Executes a container call
	Call(
		ctx context.Context,
		containerCall *model.ContainerCall,
		inboundScope map[string]*model.Value,
		containerCallSpec *model.ContainerCallSpec,
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
		io:               iio.New(),
	}

}

type _containerCaller struct {
	containerRuntime containerruntime.ContainerRuntime
	pubSub           pubsub.PubSub
	stateStore       stateStore
	io               iio.IIO
}

func (cc _containerCaller) Call(
	ctx context.Context,
	containerCall *model.ContainerCall,
	inboundScope map[string]*model.Value,
	containerCallSpec *model.ContainerCallSpec,
) {
	var err error
	outputs := map[string]*model.Value{}
	var exitCode int64

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		event := model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExited{
				ContainerID: containerCall.ContainerID,
				OpRef:       containerCall.OpPath,
				RootCallID:  containerCall.RootCallID,
				ExitCode:    exitCode,
				Outputs:     outputs,
			},
		}

		if nil != err {
			event.ContainerExited.Error = &model.CallEndedError{
				Message: err.Error(),
			}
		}

		cc.pubSub.Publish(event)
	}()

	if nil != containerCall.Image.Ref && nil == containerCall.Image.PullCreds {
		if auth := cc.stateStore.TryGetAuth(*containerCall.Image.Ref); nil != auth {
			containerCall.Image.PullCreds = &auth.Creds
		}
	}

	logStdOutPR, logStdOutPW := cc.io.Pipe()
	logStdErrPR, logStdErrPW := cc.io.Pipe()

	// interpret logs
	logChan := make(chan error, 1)
	go func() {
		logChan <- cc.interpretLogs(
			logStdOutPR,
			logStdErrPR,
			containerCall,
		)
	}()

	outputs = cc.interpretOutputs(
		containerCallSpec,
		containerCall,
	)

	var rawExitCode *int64
	rawExitCode, err = cc.containerRuntime.RunContainer(
		ctx,
		containerCall,
		cc.pubSub,
		logStdOutPW,
		logStdErrPW,
	)

	// @TODO: handle no exit code
	if nil != rawExitCode {
		exitCode = *rawExitCode
	}

	if exitCode != 0 {
		err = fmt.Errorf("nonzero container exit code. Exit code was: %v", exitCode)
	}

	// wait on logChan
	if logChanErr := <-logChan; nil == err {
		// non-destructively set err
		err = logChanErr
	}
}

func (this _containerCaller) interpretLogs(
	stdOutReader io.Reader,
	stdErrReader io.Reader,
	containerCall *model.ContainerCall,
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
							RootCallID:  containerCall.RootCallID,
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
							RootCallID:  containerCall.RootCallID,
						},
					},
				)
			})
	}()

	// wait on logs
	stdOutLogErr := <-stdOutLogChan
	stdErrLogErr := <-stdErrLogChan

	// return errs
	if nil != stdOutLogErr {
		return stdOutLogErr
	}
	if nil != stdErrLogErr {
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
	for callSpecContainerFilePath, name := range containerCallSpec.Files {
		if "" == name {
			// skip embedded files
			continue
		}

		// add file outputs
		for callContainerFilePath, callHostFilePath := range containerCall.Files {
			if callSpecContainerFilePath == callContainerFilePath {
				// copy callHostFilePath before taking address; range vars have same address for every iteration
				value := callHostFilePath
				if nameAsString, ok := name.(string); ok {
					outputs[strings.TrimSuffix(strings.TrimPrefix(nameAsString, "$("), ")")] = &model.Value{File: &value}
				}
			}
		}
	}
	for callSpecContainerDirPath, name := range containerCallSpec.Dirs {
		if "" == name {
			// skip embedded dirs
			continue
		}

		// add dir outputs
		for callContainerDirPath, callHostDirPath := range containerCall.Dirs {
			if callSpecContainerDirPath == callContainerDirPath {
				// copy callHostDirPath before taking address; range vars have same address for every iteration
				value := callHostDirPath
				outputs[strings.TrimSuffix(strings.TrimPrefix(name, "$("), ")")] = &model.Value{Dir: &value}
			}
		}
	}

	return outputs
}
