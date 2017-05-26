package core

//go:generate counterfeiter -o ./fakeContainerCaller.go --fake-name fakeContainerCaller ./ containerCaller

import (
	"bufio"
	"bytes"
	"github.com/golang-interfaces/iio"
	"github.com/opctl/opctl/util/containerprovider"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/containercall"
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
	containerCall containercall.ContainerCall,
	pubSub pubsub.PubSub,
	dcgNodeRepo dcgNodeRepo,
) containerCaller {

	return _containerCaller{
		containerProvider: containerProvider,
		containerCall:     containerCall,
		pubSub:            pubSub,
		dcgNodeRepo:       dcgNodeRepo,
		io:                iio.New(),
	}

}

type _containerCaller struct {
	containerProvider containerprovider.ContainerProvider
	containerCall     containercall.ContainerCall
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

	dcgContainerCall, err := cc.containerCall.Interpret(
		inboundScope,
		scgContainerCall,
		containerId,
		rootOpId,
		pkgRef,
	)
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
			// ErrClosedPipe is expected if multiple pipe closes
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
			// ErrClosedPipe is expected if multiple pipe closes
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

func (this _containerCaller) streamBinder(
	stream io.Reader,
	bindings map[string]string,
	onBind func(name string, value *string),
) error {
	reader := bufio.NewReader(stream)

	var err error
	for {
		var buffer bytes.Buffer
		var l []byte
		var isPrefix bool

		for {
			// use ReadString NOT Scanner to support long lines
			l, isPrefix, err = reader.ReadLine()
			buffer.Write(l)

			// If we've reached the end of the line, stop reading.
			if !isPrefix {
				break
			}

			// If we're just at the EOF, break
			if err != nil {
				break
			}
		}

		line := buffer.String()
		for boundPrefix, name := range bindings {
			trimmedLine := strings.TrimPrefix(line, boundPrefix)
			if trimmedLine != line {
				// if output trimming had effect we've got a match
				onBind(name, &trimmedLine)
			}
		}

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		return err
	}
	return nil
}

func (this _containerCaller) streamTransmitter(
	stdOutReader io.Reader,
	onChunk func(n int, chunk []byte),
) error {
	chunk := make([]byte, 1024)
	var n int
	var err error

	for {
		// rather than chunking by line, we chunk by time at a rate of 30 FPS (frames per second)
		// why? chunking by line would make TTY behaviors such as line editing behave non-TTY like
		<-time.After(33 * time.Millisecond)
		if n, err = stdOutReader.Read(chunk); n > 0 {
			// always call onChunk if n > 0 to ensure full stream sent; even under error conditions
			onChunk(n, chunk)
		}

		if nil != err {
			break
		}
	}

	if io.EOF == err {
		return nil
	}
	return err
}

func (this _containerCaller) handleStdErr(
	stdErrReader io.Reader,
	dcgContainerCall *model.DCGContainerCall,
	scgContainerCall *model.SCGContainerCall,
) error {
	errChan := make(chan error, 2)
	var wg sync.WaitGroup
	pr, pw := io.Pipe()
	tr := io.TeeReader(stdErrReader, pw)

	wg.Add(1)
	go func() {
		if err := this.streamBinder(
			tr,
			scgContainerCall.StdErr,
			func(name string, value *string) {
				this.pubSub.Publish(&model.Event{
					Timestamp: time.Now().UTC(),
					OutputInitialized: &model.OutputInitializedEvent{
						Name:     name,
						Value:    &model.Data{String: value},
						RootOpId: dcgContainerCall.RootOpId,
						CallId:   dcgContainerCall.ContainerId,
					},
				})
			}); nil != err {
			errChan <- err
		}

		// close Pipe
		pw.Close()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err := this.streamTransmitter(
			pr,
			func(n int, chunk []byte) {
				this.pubSub.Publish(
					&model.Event{
						Timestamp: time.Now().UTC(),
						ContainerStdErrWrittenTo: &model.ContainerStdErrWrittenToEvent{
							Data:        chunk[0:n],
							ContainerId: dcgContainerCall.ContainerId,
							ImageRef:    dcgContainerCall.Image.Ref,
							PkgRef:      dcgContainerCall.PkgRef,
							RootOpId:    dcgContainerCall.RootOpId,
						},
					},
				)
			}); nil != err {
			errChan <- err
		}
		wg.Done()
	}()

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

func (this _containerCaller) handleStdOut(
	stdOutReader io.Reader,
	dcgContainerCall *model.DCGContainerCall,
	scgContainerCall *model.SCGContainerCall,
) error {
	errChan := make(chan error, 2)
	var wg sync.WaitGroup
	pr, pw := io.Pipe()
	tr := io.TeeReader(stdOutReader, pw)

	wg.Add(1)
	go func() {
		if err := this.streamBinder(
			tr,
			scgContainerCall.StdOut,
			func(name string, value *string) {
				this.pubSub.Publish(&model.Event{
					Timestamp: time.Now().UTC(),
					OutputInitialized: &model.OutputInitializedEvent{
						Name:     name,
						Value:    &model.Data{String: value},
						RootOpId: dcgContainerCall.RootOpId,
						CallId:   dcgContainerCall.ContainerId,
					},
				})
			}); nil != err {
			errChan <- err
		}

		// close Pipe
		pw.Close()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err := this.streamTransmitter(
			pr,
			func(n int, chunk []byte) {
				this.pubSub.Publish(
					&model.Event{
						Timestamp: time.Now().UTC(),
						ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenToEvent{
							Data:        chunk[0:n],
							ContainerId: dcgContainerCall.ContainerId,
							ImageRef:    dcgContainerCall.Image.Ref,
							PkgRef:      dcgContainerCall.PkgRef,
							RootOpId:    dcgContainerCall.RootOpId,
						},
					},
				)
			}); nil != err {
			errChan <- err
		}
		wg.Done()
	}()

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
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
