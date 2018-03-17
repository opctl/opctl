package core

//go:generate counterfeiter -o ./fakeContainerCaller.go --fake-name fakeContainerCaller ./ containerCaller

import (
	"context"
	"fmt"
	"github.com/golang-interfaces/iio"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/core/containerruntime"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"io"
	"strings"
	"time"
)

type containerCaller interface {
	// Executes a container call
	Call(
		inboundScope map[string]*model.Value,
		containerId string,
		scgContainerCall *model.SCGContainerCall,
		opDirHandle model.DataHandle,
		rootOpId string,
	) error
}

func newContainerCaller(
	containerRuntime containerruntime.ContainerRuntime,
	containerCallInterpreter containercall.Interpreter,
	pubSub pubsub.PubSub,
	dcgNodeRepo dcgNodeRepo,
) containerCaller {

	return _containerCaller{
		containerRuntime: containerRuntime,
		containerCall:    containerCallInterpreter,
		pubSub:           pubSub,
		dcgNodeRepo:      dcgNodeRepo,
		io:               iio.New(),
	}

}

type _containerCaller struct {
	containerRuntime containerruntime.ContainerRuntime
	containerCall    containercall.Interpreter
	pubSub           pubsub.PubSub
	dcgNodeRepo      dcgNodeRepo
	io               iio.IIO
}

func (cc _containerCaller) Call(
	inboundScope map[string]*model.Value,
	containerId string,
	scgContainerCall *model.SCGContainerCall,
	opDirHandle model.DataHandle,
	rootOpId string,
) error {
	defer func() {
		// defer must be defined before conditional return statements so it always runs

		cc.dcgNodeRepo.DeleteIfExists(containerId)

		cc.containerRuntime.DeleteContainerIfExists(containerId)

	}()

	cc.dcgNodeRepo.Add(
		&dcgNodeDescriptor{
			Id:        containerId,
			PkgRef:    opDirHandle.Ref(),
			RootOpId:  rootOpId,
			Container: &dcgContainerDescriptor{},
		},
	)

	dcgContainerCall, err := cc.containerCall.Interpret(
		inboundScope,
		scgContainerCall,
		containerId,
		rootOpId,
		opDirHandle,
	)
	if nil != err {
		return err
	}

	cc.pubSub.Publish(
		model.Event{
			Timestamp: time.Now().UTC(),
			ContainerStarted: &model.ContainerStartedEvent{
				ContainerId: containerId,
				PkgRef:      opDirHandle.Ref(),
				RootOpId:    rootOpId,
			},
		},
	)

	// we need 2 stdOut readers
	stdOutReader, stdOutWriter := cc.io.Pipe()
	logStdOutPR, logStdOutPW := cc.io.Pipe()
	outputsStdOutTR := io.TeeReader(stdOutReader, logStdOutPW)

	// we need 2 stdErr readers
	stdErrReader, stdErrWriter := cc.io.Pipe()
	logStdErrPR, logStdErrPW := cc.io.Pipe()
	outputsStdErrTR := io.TeeReader(stdErrReader, logStdErrPW)

	// interpret logs
	logChan := make(chan error, 1)
	go func() {
		logChan <- cc.interpretLogs(
			logStdOutPR,
			logStdErrPR,
			scgContainerCall,
			dcgContainerCall,
		)
	}()

	// interpret outputs
	outputsChan := make(
		chan struct {
			outputs map[string]*model.Value
			err     error
		}, 1)
	go func() {
		// close pipes
		defer logStdOutPW.Close()
		defer logStdErrPW.Close()

		outputs, err := cc.interpretOutputs(
			outputsStdOutTR,
			outputsStdErrTR,
			scgContainerCall,
			dcgContainerCall,
		)
		outputsChan <- struct {
			outputs map[string]*model.Value
			err     error
		}{
			outputs: outputs,
			err:     err,
		}
	}()

	rawExitCode, err := cc.containerRuntime.RunContainer(
		context.TODO(),
		dcgContainerCall,
		cc.pubSub,
		stdOutWriter,
		stdErrWriter,
	)
	// @TODO: handle no exit code
	var exitCode int64
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

	// wait on outputsChan
	interpretOutputsResult := <-outputsChan
	if nil == err {
		// non-destructively set err
		err = interpretOutputsResult.err
	}

	cc.pubSub.Publish(
		model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExitedEvent{
				ContainerId: containerId,
				PkgRef:      opDirHandle.Ref(),
				RootOpId:    rootOpId,
				ExitCode:    exitCode,
				Outputs:     interpretOutputsResult.outputs,
			},
		},
	)
	return err
}

func (this _containerCaller) interpretLogs(
	stdOutReader io.Reader,
	stdErrReader io.Reader,
	scgContainerCall *model.SCGContainerCall,
	dcgContainerCall *model.DCGContainerCall,
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
						ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenToEvent{
							Data:        chunk,
							ContainerId: dcgContainerCall.ContainerId,
							ImageRef:    dcgContainerCall.Image.Ref,
							PkgRef:      dcgContainerCall.DataHandle.Ref(),
							RootOpId:    dcgContainerCall.RootOpId,
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
						ContainerStdErrWrittenTo: &model.ContainerStdErrWrittenToEvent{
							Data:        chunk,
							ContainerId: dcgContainerCall.ContainerId,
							ImageRef:    dcgContainerCall.Image.Ref,
							PkgRef:      dcgContainerCall.DataHandle.Ref(),
							RootOpId:    dcgContainerCall.RootOpId,
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
	stdOutReader io.Reader,
	stdErrReader io.Reader,
	scgContainerCall *model.SCGContainerCall,
	dcgContainerCall *model.DCGContainerCall,
) (map[string]*model.Value, error) {
	outputs := map[string]*model.Value{}

	for socketAddr, name := range scgContainerCall.Sockets {
		// add socket outputs
		if "0.0.0.0" == socketAddr {
			outputs[name] = &model.Value{Socket: &dcgContainerCall.ContainerId}
		}
	}
	for scgContainerFilePath, name := range scgContainerCall.Files {
		// add file outputs
		for dcgContainerFilePath, dcgHostFilePath := range dcgContainerCall.Files {
			if scgContainerFilePath == dcgContainerFilePath {
				// copy dcgHostFilePath before taking address; range vars have same address for every iteration
				value := dcgHostFilePath
				if nameAsString, ok := name.(string); ok {
					outputs[strings.TrimSuffix(strings.TrimPrefix(nameAsString, "$("), ")")] = &model.Value{File: &value}
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
				outputs[strings.TrimSuffix(strings.TrimPrefix(name, "$("), ")")] = &model.Value{Dir: &value}
			}
		}
	}

	stdErrOutputsChan := make(chan struct {
		outputs map[string]*model.Value
		err     error
	}, 1)
	go func() {
		// add stdErr stdErrOutputs
		stdErrOutputs := map[string]*model.Value{}
		err := bindLines(
			stdErrReader,
			scgContainerCall.StdErr,
			func(name string, value *string) {
				stdErrOutputs[name] = &model.Value{String: value}
			})
		stdErrOutputsChan <- struct {
			outputs map[string]*model.Value
			err     error
		}{outputs: stdErrOutputs, err: err}
	}()

	stdOutOutputsChan := make(chan struct {
		outputs map[string]*model.Value
		err     error
	}, 1)
	go func() {
		// add stdOut stdOutOutputs
		stdOutOutputs := map[string]*model.Value{}
		err := bindLines(
			stdOutReader,
			scgContainerCall.StdOut,
			func(name string, value *string) {
				stdOutOutputs[name] = &model.Value{String: value}
			})
		stdOutOutputsChan <- struct {
			outputs map[string]*model.Value
			err     error
		}{outputs: stdOutOutputs, err: err}
	}()

	// wait for stdErr result
	chanResult := <-stdErrOutputsChan
	if nil != chanResult.err {
		return nil, chanResult.err
	}
	for name, value := range chanResult.outputs {
		outputs[name] = value
	}

	// wait for stdOut result
	chanResult = <-stdOutOutputsChan
	if nil != chanResult.err {
		return nil, chanResult.err
	}
	for name, value := range chanResult.outputs {
		outputs[name] = value
	}

	return outputs, nil
}
