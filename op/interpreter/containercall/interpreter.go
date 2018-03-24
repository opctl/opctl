package containercall

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall/dirs"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall/envvars"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall/files"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall/image"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall/sockets"
	stringPkg "github.com/opspec-io/sdk-golang/op/interpreter/string"
	"path/filepath"
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
	rootFSPath string,
) Interpreter {
	return _interpreter{
		dirsInterpreter:    dirs.NewInterpreter(rootFSPath),
		envVarsInterpreter: envvars.NewInterpreter(),
		filesInterpreter:   files.NewInterpreter(rootFSPath),
		imageInterpreter:   image.NewInterpreter(),
		os:                 ios.New(),
		rootFSPath:         rootFSPath,
		stringInterpreter:  stringPkg.NewInterpreter(),
		socketsInterpreter: sockets.NewInterpreter(),
	}
}

type _interpreter struct {
	dirsInterpreter    dirs.Interpreter
	envVarsInterpreter envvars.Interpreter
	filesInterpreter   files.Interpreter
	imageInterpreter   image.Interpreter
	os                 ios.IOS
	rootFSPath         string
	stringInterpreter  stringPkg.Interpreter
	socketsInterpreter sockets.Interpreter
}

func (cc _interpreter) Interpret(
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
		Name:        scgContainerCall.Name,
		Ports:       scgContainerCall.Ports,
	}

	// construct dcg container path
	scratchDirPath := filepath.Join(
		cc.rootFSPath,
		"dcg",
		rootOpID,
		"containers",
		containerID,
		"fs",
	)
	if err := cc.os.MkdirAll(scratchDirPath, 0700); nil != err {
		return nil, err
	}

	// construct cmd
	for _, cmdEntryExpression := range scgContainerCall.Cmd {
		// interpret each entry as string
		cmdEntry, err := cc.stringInterpreter.Interpret(scope, cmdEntryExpression, opHandle)
		if nil != err {
			return nil, err
		}
		dcgContainerCall.Cmd = append(dcgContainerCall.Cmd, *cmdEntry.String)
	}

	// interpret dirs
	var err error
	dcgContainerCall.Dirs, err = cc.dirsInterpreter.Interpret(opHandle, scope, scgContainerCall.Dirs, scratchDirPath)
	if nil != err {
		return nil, err
	}

	// interpret envVars
	dcgContainerCall.EnvVars, err = cc.envVarsInterpreter.Interpret(scope, scgContainerCall.EnvVars, opHandle)
	if nil != err {
		return nil, err
	}

	// interpret files
	dcgContainerCall.Files, err = cc.filesInterpreter.Interpret(opHandle, scope, scgContainerCall.Files, scratchDirPath)
	if nil != err {
		return nil, err
	}

	// interpret image
	dcgContainerCall.Image, err = cc.imageInterpreter.Interpret(scope, scgContainerCall.Image, opHandle)
	if nil != err {
		return nil, err
	}

	// interpret sockets
	dcgContainerCall.Sockets, err = cc.socketsInterpreter.Interpret(scope, scgContainerCall.Sockets, scratchDirPath)

	return dcgContainerCall, err

}
