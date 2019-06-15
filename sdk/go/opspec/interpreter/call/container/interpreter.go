package container

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"path/filepath"

	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/container/cmd"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/container/dirs"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/container/envvars"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/container/files"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/container/image"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/container/sockets"
	stringPkg "github.com/opctl/sdk-golang/opspec/interpreter/string"
)

type Interpreter interface {
	// Interpret interprets an SCGContainerCall into a DCGContainerCall
	Interpret(
		scope map[string]*model.Value,
		scgContainerCall *model.SCGContainerCall,
		containerID string,
		rootOpID string,
		opHandle model.DataHandle,
	) (*model.DCGContainerCall, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter(
	dataDirPath string,
) Interpreter {
	return _interpreter{
		cmdInterpreter:     cmd.NewInterpreter(),
		dirsInterpreter:    dirs.NewInterpreter(dataDirPath),
		envVarsInterpreter: envvars.NewInterpreter(),
		filesInterpreter:   files.NewInterpreter(dataDirPath),
		imageInterpreter:   image.NewInterpreter(),
		os:                 ios.New(),
		dataDirPath:        dataDirPath,
		stringInterpreter:  stringPkg.NewInterpreter(),
		socketsInterpreter: sockets.NewInterpreter(),
	}
}

type _interpreter struct {
	cmdInterpreter     cmd.Interpreter
	dirsInterpreter    dirs.Interpreter
	envVarsInterpreter envvars.Interpreter
	filesInterpreter   files.Interpreter
	imageInterpreter   image.Interpreter
	os                 ios.IOS
	dataDirPath        string
	stringInterpreter  stringPkg.Interpreter
	socketsInterpreter sockets.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	scgContainerCall *model.SCGContainerCall,
	containerID string,
	rootOpID string,
	opHandle model.DataHandle,
) (*model.DCGContainerCall, error) {

	dcgContainerCall := &model.DCGContainerCall{
		DCGBaseCall: model.DCGBaseCall{
			RootOpID: rootOpID,
			OpHandle: opHandle,
		},
		Dirs:        map[string]string{},
		EnvVars:     map[string]string{},
		Files:       map[string]string{},
		Sockets:     map[string]string{},
		WorkDir:     scgContainerCall.WorkDir,
		ContainerID: containerID,
		Ports:       scgContainerCall.Ports,
	}

	// construct dcg container path
	scratchDirPath := filepath.Join(
		itp.dataDirPath,
		"dcg",
		rootOpID,
		"containers",
		containerID,
		"fs",
	)
	if err := itp.os.MkdirAll(scratchDirPath, 0700); nil != err {
		return nil, err
	}

	// interpret cmd
	var err error
	dcgContainerCall.Cmd, err = itp.cmdInterpreter.Interpret(scope, scgContainerCall.Cmd, opHandle)
	if nil != err {
		return nil, err
	}

	// interpret dirs
	dcgContainerCall.Dirs, err = itp.dirsInterpreter.Interpret(opHandle, scope, scgContainerCall.Dirs, scratchDirPath)
	if nil != err {
		return nil, err
	}

	// interpret envVars
	dcgContainerCall.EnvVars, err = itp.envVarsInterpreter.Interpret(scope, scgContainerCall.EnvVars, opHandle)
	if nil != err {
		return nil, err
	}

	// interpret files
	dcgContainerCall.Files, err = itp.filesInterpreter.Interpret(opHandle, scope, scgContainerCall.Files, scratchDirPath)
	if nil != err {
		return nil, err
	}

	// interpret image
	dcgContainerCall.Image, err = itp.imageInterpreter.Interpret(scope, scgContainerCall.Image, opHandle)
	if nil != err {
		return nil, err
	}

	// interpret name as string
	if nil != scgContainerCall.Name {
		dcgContainerCallName, err := itp.stringInterpreter.Interpret(scope, *scgContainerCall.Name, opHandle)
		if nil != err {
			return nil, err
		}
		dcgContainerCall.Name = dcgContainerCallName.String
	}

	// interpret sockets
	dcgContainerCall.Sockets, err = itp.socketsInterpreter.Interpret(scope, scgContainerCall.Sockets, scratchDirPath)

	return dcgContainerCall, err

}
