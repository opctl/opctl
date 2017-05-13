package core

//go:generate counterfeiter -o ./fakeContainerCaller.go --fake-name fakeContainerCaller ./ containerCaller

import (
	"bufio"
	"github.com/golang-interfaces/iio"
	"github.com/opctl/opctl/util/containerprovider"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
	"io"
	"strings"
	"sync"
	"time"
)

type containerCaller interface {
	// Executes a container call
	Call(
		inboundScope map[string]*model.Data,
		containerId string,
		scgContainerCall *model.SCGContainerCall,
		pkgRef string,
		rootOpId string,
	) error
}

func newContainerCaller(
	containerProvider containerprovider.ContainerProvider,
	dcgFactory dcgFactory,
	pubSub pubsub.PubSub,
	dcgNodeRepo dcgNodeRepo,
) containerCaller {

	return _containerCaller{
		containerProvider: containerProvider,
		dcgFactory:        dcgFactory,
		pubSub:            pubSub,
		dcgNodeRepo:       dcgNodeRepo,
		io:                iio.New(),
	}

}

type _containerCaller struct {
	containerProvider containerprovider.ContainerProvider
	dcgFactory        dcgFactory
	pubSub            pubsub.PubSub
	dcgNodeRepo       dcgNodeRepo
	io                iio.IIO
}

func (cc _containerCaller) Call(
	inboundScope map[string]*model.Data,
	containerId string,
	scgContainerCall *model.SCGContainerCall,
	pkgRef string,
	rootOpId string,
) error {
	defer func() {
		// defer must be defined before conditional return statements so it always runs

		cc.dcgNodeRepo.DeleteIfExists(containerId)

		cc.containerProvider.DeleteContainerIfExists(containerId)

	}()

	cc.dcgNodeRepo.Add(
		&dcgNodeDescriptor{
			Id:        containerId,
			PkgRef:    pkgRef,
			RootOpId:  rootOpId,
			Container: &dcgContainerDescriptor{},
		},
	)

	dcgContainerCall, err := cc.dcgFactory.Construct(inboundScope, scgContainerCall, containerId, rootOpId, pkgRef)
	if nil != err {
		return err
	}

	cc.txOutputs(dcgContainerCall, scgContainerCall)

	cc.pubSub.Publish(
		&model.Event{
			Timestamp: time.Now().UTC(),
			ContainerStarted: &model.ContainerStartedEvent{
				ContainerId: containerId,
				PkgRef:      pkgRef,
				RootOpId:    rootOpId,
			},
		},
	)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	wg.Add(2)

	containerStdOutReader, containerStdOutWriter := cc.io.Pipe()
	go func() {
		if err := cc.handleStdOut(
			containerStdOutReader,
			dcgContainerCall,
			scgContainerCall,
		); nil != err {
			errChan <- err
		}
		wg.Done()
	}()

	containerStdErrReader, containerStdErrWriter := cc.io.Pipe()
	go func() {
		if err := cc.handleStdErr(
			containerStdErrReader,
			dcgContainerCall,
			scgContainerCall,
		); nil != err {
			errChan <- err
		}
		wg.Done()
	}()

	rawExitCode, err := cc.containerProvider.RunContainer(
		dcgContainerCall,
		cc.pubSub,
		containerStdOutWriter,
		containerStdErrWriter,
	)
	// @TODO: handle no exit code
	var exitCode int64
	if nil != rawExitCode {
		exitCode = *rawExitCode
	}

	wg.Wait()

	cc.pubSub.Publish(
		&model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExitedEvent{
				ContainerId: containerId,
				PkgRef:      pkgRef,
				RootOpId:    rootOpId,
				ExitCode:    exitCode,
			},
		},
	)

	if nil == err && len(errChan) > 0 {
		err = <-errChan
	}
	return err
}

func (this _containerCaller) handleStdErr(
	stdErrReader io.ReadCloser,
	dcgContainerCall *model.DCGContainerCall,
	scgContainerCall *model.SCGContainerCall,
) error {
	defer stdErrReader.Close()
	scanner := bufio.NewScanner(stdErrReader)

	// scan writes until EOF or error
	for scanner.Scan() {

		this.pubSub.Publish(
			&model.Event{
				Timestamp: time.Now().UTC(),
				ContainerStdErrWrittenTo: &model.ContainerStdErrWrittenToEvent{
					Data:        scanner.Bytes(),
					ContainerId: dcgContainerCall.ContainerId,
					ImageRef:    dcgContainerCall.Image.Ref,
					PkgRef:      dcgContainerCall.PkgRef,
					RootOpId:    dcgContainerCall.RootOpId,
				},
			},
		)

		for boundPrefix, name := range scgContainerCall.StdErr {
			rawOutput := scanner.Text()
			trimmedOutput := strings.TrimPrefix(rawOutput, boundPrefix)
			if trimmedOutput != rawOutput {
				// if output trimming had effect we've got a match
				this.pubSub.Publish(&model.Event{
					Timestamp: time.Now().UTC(),
					OutputInitialized: &model.OutputInitializedEvent{
						Name:     name,
						Value:    &model.Data{String: &trimmedOutput},
						RootOpId: dcgContainerCall.RootOpId,
						CallId:   dcgContainerCall.ContainerId,
					},
				})
			}
		}
	}
	return scanner.Err()
}

func (this _containerCaller) handleStdOut(
	stdOutReader io.ReadCloser,
	dcgContainerCall *model.DCGContainerCall,
	scgContainerCall *model.SCGContainerCall,
) error {
	defer stdOutReader.Close()
	scanner := bufio.NewScanner(stdOutReader)

	// scan writes until EOF or error
	for scanner.Scan() {

		this.pubSub.Publish(
			&model.Event{
				Timestamp: time.Now().UTC(),
				ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenToEvent{
					Data:        scanner.Bytes(),
					ContainerId: dcgContainerCall.ContainerId,
					ImageRef:    dcgContainerCall.Image.Ref,
					PkgRef:      dcgContainerCall.PkgRef,
					RootOpId:    dcgContainerCall.RootOpId,
				},
			},
		)
		for boundPrefix, name := range scgContainerCall.StdOut {
			rawOutput := scanner.Text()
			trimmedOutput := strings.TrimPrefix(rawOutput, boundPrefix)
			if trimmedOutput != rawOutput {
				// if output trimming had effect we've got a match
				this.pubSub.Publish(&model.Event{
					Timestamp: time.Now().UTC(),
					OutputInitialized: &model.OutputInitializedEvent{
						Name:     name,
						Value:    &model.Data{String: &trimmedOutput},
						RootOpId: dcgContainerCall.RootOpId,
						CallId:   dcgContainerCall.ContainerId,
					},
				})
			}
		}
	}
	return scanner.Err()
}

func (this _containerCaller) txOutputs(
	dcgContainerCall *model.DCGContainerCall,
	scgContainerCall *model.SCGContainerCall,
) {

	// send socket outputs
	for socketAddr, name := range scgContainerCall.Sockets {
		if "0.0.0.0" == socketAddr {
			this.pubSub.Publish(&model.Event{
				Timestamp: time.Now().UTC(),
				OutputInitialized: &model.OutputInitializedEvent{
					Name:     name,
					Value:    &model.Data{Socket: &dcgContainerCall.ContainerId},
					RootOpId: dcgContainerCall.RootOpId,
					CallId:   dcgContainerCall.ContainerId,
				},
			})
		}
	}

	// send file outputs
	for scgContainerFilePath, name := range scgContainerCall.Files {
		for dcgContainerFilePath, dcgHostFilePath := range dcgContainerCall.Files {
			if scgContainerFilePath == dcgContainerFilePath {
				this.pubSub.Publish(&model.Event{
					Timestamp: time.Now().UTC(),
					OutputInitialized: &model.OutputInitializedEvent{
						Name:     name,
						Value:    &model.Data{File: &dcgHostFilePath},
						RootOpId: dcgContainerCall.RootOpId,
						CallId:   dcgContainerCall.ContainerId,
					},
				})
			}
		}
	}

	// send dir outputs
	for scgContainerDirPath, name := range scgContainerCall.Dirs {
		for dcgContainerDirPath, dcgHostDirPath := range dcgContainerCall.Dirs {
			if scgContainerDirPath == dcgContainerDirPath {
				this.pubSub.Publish(&model.Event{
					Timestamp: time.Now().UTC(),
					OutputInitialized: &model.OutputInitializedEvent{
						Name:     name,
						Value:    &model.Data{Dir: &dcgHostDirPath},
						RootOpId: dcgContainerCall.RootOpId,
						CallId:   dcgContainerCall.ContainerId,
					},
				})
			}
		}
	}

	return
}
