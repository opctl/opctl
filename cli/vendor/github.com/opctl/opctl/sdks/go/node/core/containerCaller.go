package core

//go:generate counterfeiter -o ./fakeContainerCaller.go --fake-name fakeContainerCaller ./ containerCaller

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/golang-interfaces/iio"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/types"
	"github.com/opctl/opctl/sdks/go/util/pubsub"
)

type containerCaller interface {
	// Executes a container call
	Call(
		ctx context.Context,
		dcgContainerCall *types.DCGContainerCall,
		inboundScope map[string]*types.Value,
		scgContainerCall *types.SCGContainerCall,
	)
}

func newContainerCaller(
	containerRuntime containerruntime.ContainerRuntime,
	pubSub pubsub.PubSub,
) containerCaller {

	return _containerCaller{
		containerRuntime: containerRuntime,
		pubSub:           pubSub,
		io:               iio.New(),
	}

}

type _containerCaller struct {
	containerRuntime containerruntime.ContainerRuntime
	pubSub           pubsub.PubSub
	io               iio.IIO
}

func (cc _containerCaller) Call(
	ctx context.Context,
	dcgContainerCall *types.DCGContainerCall,
	inboundScope map[string]*types.Value,
	scgContainerCall *types.SCGContainerCall,
) {
	var err error
	outputs := map[string]*types.Value{}
	var exitCode int64

	defer func() {
		event := types.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &types.ContainerExitedEvent{
				ContainerID: dcgContainerCall.ContainerID,
				OpRef:       dcgContainerCall.OpHandle.Ref(),
				RootOpID:    dcgContainerCall.RootOpID,
				ExitCode:    exitCode,
				Outputs:     outputs,
			},
		}

		if nil != err {
			event.ContainerExited.Error = &types.CallEndedEventError{
				Message: err.Error(),
			}
		}

		cc.pubSub.Publish(event)
	}()

	cc.pubSub.Publish(
		types.Event{
			Timestamp: time.Now().UTC(),
			ContainerStarted: &types.ContainerStartedEvent{
				ContainerID: dcgContainerCall.ContainerID,
				OpRef:       dcgContainerCall.OpHandle.Ref(),
				RootOpID:    dcgContainerCall.RootOpID,
			},
		},
	)

	logStdOutPR, logStdOutPW := cc.io.Pipe()
	logStdErrPR, logStdErrPW := cc.io.Pipe()

	// interpret logs
	logChan := make(chan error, 1)
	go func() {
		logChan <- cc.interpretLogs(
			logStdOutPR,
			logStdErrPR,
			dcgContainerCall,
		)
	}()

	outputs = cc.interpretOutputs(
		scgContainerCall,
		dcgContainerCall,
	)

	var rawExitCode *int64
	rawExitCode, err = cc.containerRuntime.RunContainer(
		ctx,
		dcgContainerCall,
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
	dcgContainerCall *types.DCGContainerCall,
) error {
	stdOutLogChan := make(chan error, 1)
	go func() {
		// interpret stdOut
		stdOutLogChan <- readChunks(
			stdOutReader,
			func(chunk []byte) {
				this.pubSub.Publish(
					types.Event{
						Timestamp: time.Now().UTC(),
						ContainerStdOutWrittenTo: &types.ContainerStdOutWrittenToEvent{
							Data:        chunk,
							ContainerID: dcgContainerCall.ContainerID,
							ImageRef:    dcgContainerCall.Image.Ref,
							OpRef:       dcgContainerCall.OpHandle.Ref(),
							RootOpID:    dcgContainerCall.RootOpID,
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
					types.Event{
						Timestamp: time.Now().UTC(),
						ContainerStdErrWrittenTo: &types.ContainerStdErrWrittenToEvent{
							Data:        chunk,
							ContainerID: dcgContainerCall.ContainerID,
							ImageRef:    dcgContainerCall.Image.Ref,
							OpRef:       dcgContainerCall.OpHandle.Ref(),
							RootOpID:    dcgContainerCall.RootOpID,
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
	scgContainerCall *types.SCGContainerCall,
	dcgContainerCall *types.DCGContainerCall,
) map[string]*types.Value {
	outputs := map[string]*types.Value{}

	for socketAddr, name := range scgContainerCall.Sockets {
		// add socket outputs
		if "0.0.0.0" == socketAddr {
			outputs[name] = &types.Value{Socket: &dcgContainerCall.ContainerID}
		}
	}
	for scgContainerFilePath, name := range scgContainerCall.Files {
		// add file outputs
		for dcgContainerFilePath, dcgHostFilePath := range dcgContainerCall.Files {
			if scgContainerFilePath == dcgContainerFilePath {
				// copy dcgHostFilePath before taking address; range vars have same address for every iteration
				value := dcgHostFilePath
				if nameAsString, ok := name.(string); ok {
					outputs[strings.TrimSuffix(strings.TrimPrefix(nameAsString, "$("), ")")] = &types.Value{File: &value}
				}
			}
		}
	}
	for scgContainerDirPath, name := range scgContainerCall.Dirs {
		// add dir outputs
		for dcgContainerDirPath, dcgHostDirPath := range dcgContainerCall.Dirs {
			if scgContainerDirPath == dcgContainerDirPath {
				// copy dcgHostDirPath before taking address; range vars have same address for every iteration
				value := dcgHostDirPath
				outputs[strings.TrimSuffix(strings.TrimPrefix(name, "$("), ")")] = &types.Value{Dir: &value}
			}
		}
	}

	return outputs
}
